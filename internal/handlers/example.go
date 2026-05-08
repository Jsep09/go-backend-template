// internal/handlers/example.go
//
// ตัวอย่าง CRUD handler
// ตอนนี้ยังไม่มี sqlc → ใช้ pgx ตรงๆ ก่อน
// หลัง sqlc generate แล้วค่อยเปลี่ยน

package handlers

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Jsep09/go-backend-template/internal/middleware"
	"github.com/Jsep09/go-backend-template/internal/models"
)

// validate instance — สร้างครั้งเดียว reuse ได้
var validate = validator.New()

// ─────────────────────────────────────────
// Handler struct
// ─────────────────────────────────────────

type ExampleHandler struct {
	db *pgxpool.Pool
}

func NewExampleHandler(db *pgxpool.Pool) *ExampleHandler {
	return &ExampleHandler{db: db}
}

// ─────────────────────────────────────────
// Request / Response types
// ─────────────────────────────────────────

type CreateExampleRequest struct {
	Name        string `json:"name"         validate:"required,min=1,max=100"`
	Description string `json:"description"  validate:"max=500"`
}

type ExampleResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      string `json:"user_id"`
	CreatedAt   string `json:"created_at"`
}

// ─────────────────────────────────────────
// List GET /api/v1/examples
// ─────────────────────────────────────────

func (h *ExampleHandler) List(c fiber.Ctx) error {
	claims, ok := middleware.GetUserClaims(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			models.Fail("unauthorized"),
		)
	}

	rows, err := h.db.Query(c.Context(),
		`SELECT id, name, description, user_id, created_at
         FROM examples
         WHERE user_id = $1
         ORDER BY created_at DESC`,
		claims.Sub, // parameterized — ป้องกัน SQL injection
	)
	if err != nil {
		slog.Error("failed to query examples", "error", err, "user_id", claims.Sub)
		return c.Status(fiber.StatusInternalServerError).JSON(
			models.Fail("failed to fetch examples"),
		)
	}
	defer rows.Close()

	// สร้าง slice ว่างก่อน (ไม่ใช้ nil เพราะ JSON จะเป็น null แทน [])
	results := make([]ExampleResponse, 0)

	for rows.Next() {
		var ex ExampleResponse
		if err := rows.Scan(
			&ex.ID,
			&ex.Name,
			&ex.Description,
			&ex.UserID,
			&ex.CreatedAt,
		); err != nil {
			slog.Error("failed to scan example row", "error", err)
			continue
		}
		results = append(results, ex)
	}

	return c.JSON(models.Ok(results))
}

// ─────────────────────────────────────────
// GetByID GET /api/v1/examples/:id
// ─────────────────────────────────────────

func (h *ExampleHandler) GetByID(c fiber.Ctx) error {
	claims, ok := middleware.GetUserClaims(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Fail("unauthorized"))
	}

	id := c.Params("id")

	var ex ExampleResponse
	err := h.db.QueryRow(c.Context(),
		`SELECT id, name, description, user_id, created_at
         FROM examples
         WHERE id = $1 AND user_id = $2`, // เช็ค user_id ด้วย — ป้องกัน user ขโมยข้อมูลกัน
		id, claims.Sub,
	).Scan(&ex.ID, &ex.Name, &ex.Description, &ex.UserID, &ex.CreatedAt)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			models.Fail("example not found"),
		)
	}

	return c.JSON(models.Ok(ex))
}

// ─────────────────────────────────────────
// Create POST /api/v1/examples
// ─────────────────────────────────────────

func (h *ExampleHandler) Create(c fiber.Ctx) error {
	claims, ok := middleware.GetUserClaims(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Fail("unauthorized"))
	}

	// 1. Parse body
	var req CreateExampleRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			models.Fail("invalid request body"),
		)
	}

	// 2. Validate
	if err := validate.Struct(req); err != nil {
		// แปลง validation errors ให้อ่านง่าย
		errs := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(
			models.Fail(errs[0].Translate(nil)),
		)
	}

	// 3. Insert
	var ex ExampleResponse
	err := h.db.QueryRow(c.Context(),
		`INSERT INTO examples (name, description, user_id)
         VALUES ($1, $2, $3)
         RETURNING id, name, description, user_id, created_at`,
		req.Name, req.Description, claims.Sub,
	).Scan(&ex.ID, &ex.Name, &ex.Description, &ex.UserID, &ex.CreatedAt)

	if err != nil {
		slog.Error("failed to create example", "error", err, "user_id", claims.Sub)
		return c.Status(fiber.StatusInternalServerError).JSON(
			models.Fail("failed to create example"),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(models.Ok(ex))
}

// ─────────────────────────────────────────
// Update PUT /api/v1/examples/:id
// ─────────────────────────────────────────

func (h *ExampleHandler) Update(c fiber.Ctx) error {
	claims, ok := middleware.GetUserClaims(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Fail("unauthorized"))
	}

	id := c.Params("id")

	var req CreateExampleRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			models.Fail("invalid request body"),
		)
	}

	if err := validate.Struct(req); err != nil {
		errs := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(
			models.Fail(errs[0].Translate(nil)),
		)
	}

	var ex ExampleResponse
	err := h.db.QueryRow(c.Context(),
		`UPDATE examples
         SET name = $1, description = $2
         WHERE id = $3 AND user_id = $4
         RETURNING id, name, description, user_id, created_at`,
		req.Name, req.Description, id, claims.Sub,
	).Scan(&ex.ID, &ex.Name, &ex.Description, &ex.UserID, &ex.CreatedAt)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			models.Fail("example not found or not yours"),
		)
	}

	return c.JSON(models.Ok(ex))
}

// ─────────────────────────────────────────
// Delete DELETE /api/v1/examples/:id
// ─────────────────────────────────────────

func (h *ExampleHandler) Delete(c fiber.Ctx) error {
	claims, ok := middleware.GetUserClaims(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Fail("unauthorized"))
	}

	id := c.Params("id")

	result, err := h.db.Exec(c.Context(),
		`DELETE FROM examples
         WHERE id = $1 AND user_id = $2`,
		id, claims.Sub,
	)

	if err != nil || result.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(
			models.Fail("example not found or not yours"),
		)
	}

	return c.JSON(models.Ok(fiber.Map{"deleted": true}))
}
