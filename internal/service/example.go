package service

import (
	"errors"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"

	"github.com/Jsep09/go-backend-template/internal/models"

	db "github.com/Jsep09/go-backend-template/internal/db/generated"
)

var validate = validator.New()

// ─────────────────────────────────────────
// Handler struct
// ─────────────────────────────────────────

type ExampleHandler struct {
	queries *db.Queries
}

func NewExampleHandler(queries *db.Queries) *ExampleHandler {
	return &ExampleHandler{queries: queries}
}

// ─────────────────────────────────────────
// Request / Response types
// ─────────────────────────────────────────

type CreateExampleRequest struct {
	Name        string `json:"name"        validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=500"`
}

type ExampleResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      string `json:"user_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ─────────────────────────────────────────
// List GET /api/v1/examples
// ─────────────────────────────────────────

func (h *ExampleHandler) List(c fiber.Ctx) error {
	claims := mustGetClaims(c)

	userID, err := parseUUID(claims.Sub)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Fail("invalid user id"))
	}

	examples, err := h.queries.ListExamples(c.Context(), userID)
	if err != nil {
		slog.Error("failed to list examples", "error", err, "user_id", claims.Sub)
		return c.Status(fiber.StatusInternalServerError).JSON(models.Fail("failed to fetch examples"))
	}

	results := make([]ExampleResponse, 0, len(examples))
	for _, ex := range examples {
		results = append(results, toExampleResponse(ex))
	}

	return c.JSON(models.Ok(results))
}

// ─────────────────────────────────────────
// GetByID GET /api/v1/examples/:id
// ─────────────────────────────────────────

func (h *ExampleHandler) GetByID(c fiber.Ctx) error {
	claims := mustGetClaims(c)

	exID, err := parseUUID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Fail("invalid id"))
	}

	userID, err := parseUUID(claims.Sub)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Fail("invalid user id"))
	}

	ex, err := h.queries.GetExample(c.Context(), db.GetExampleParams{
		ID:     exID,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(models.Fail("example not found"))
		}
		slog.Error("failed to get example", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.Fail("failed to get example"))
	}

	return c.JSON(models.Ok(toExampleResponse(ex)))
}

// ─────────────────────────────────────────
// Create POST /api/v1/examples
// ─────────────────────────────────────────

func (h *ExampleHandler) Create(c fiber.Ctx) error {
	claims := mustGetClaims(c)

	var req CreateExampleRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Fail("invalid request body"))
	}

	if err := validate.Struct(req); err != nil {
		errs := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(models.Fail(errs[0].Translate(nil)))
	}

	userID, err := parseUUID(claims.Sub)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Fail("invalid user id"))
	}

	ex, err := h.queries.CreateExample(c.Context(), db.CreateExampleParams{
		Name:        req.Name,
		Description: req.Description,
		UserID:      userID,
	})
	if err != nil {
		slog.Error("failed to create example", "error", err, "user_id", claims.Sub)
		return c.Status(fiber.StatusInternalServerError).JSON(models.Fail("failed to create example"))
	}

	return c.Status(fiber.StatusCreated).JSON(models.Ok(toExampleResponse(ex)))
}

// ─────────────────────────────────────────
// Update PUT /api/v1/examples/:id
// ─────────────────────────────────────────

func (h *ExampleHandler) Update(c fiber.Ctx) error {
	claims := mustGetClaims(c)

	exID, err := parseUUID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Fail("invalid id"))
	}

	userID, err := parseUUID(claims.Sub)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Fail("invalid user id"))
	}

	var req CreateExampleRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Fail("invalid request body"))
	}

	if err := validate.Struct(req); err != nil {
		errs := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(models.Fail(errs[0].Translate(nil)))
	}

	ex, err := h.queries.UpdateExample(c.Context(), db.UpdateExampleParams{
		Name:        req.Name,
		Description: req.Description,
		ID:          exID,
		UserID:      userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(models.Fail("example not found or not yours"))
		}
		slog.Error("failed to update example", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.Fail("failed to update example"))
	}

	return c.JSON(models.Ok(toExampleResponse(ex)))
}

// ─────────────────────────────────────────
// Delete DELETE /api/v1/examples/:id
// ─────────────────────────────────────────

func (h *ExampleHandler) Delete(c fiber.Ctx) error {
	claims := mustGetClaims(c)

	exID, err := parseUUID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Fail("invalid id"))
	}

	userID, err := parseUUID(claims.Sub)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Fail("invalid user id"))
	}

	err = h.queries.DeleteExample(c.Context(), db.DeleteExampleParams{
		ID:     exID,
		UserID: userID,
	})
	if err != nil {
		slog.Error("failed to delete example", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.Fail("failed to delete example"))
	}

	return c.JSON(models.Ok(fiber.Map{"deleted": true}))
}
