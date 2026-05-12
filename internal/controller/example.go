package controller

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"

	"github.com/Jsep09/go-backend-template/internal/models"
	"github.com/Jsep09/go-backend-template/internal/service"
)

var validate = validator.New()

// ─────────────────────────────────────────
// Controller struct
// ─────────────────────────────────────────

type ExampleController struct {
	svc *service.ExampleService
}

func NewExampleController(svc *service.ExampleService) *ExampleController {
	return &ExampleController{svc: svc}
}

// ─────────────────────────────────────────
// Request types — HTTP input เท่านั้น
// ─────────────────────────────────────────

type CreateExampleRequest struct {
	Name        string `json:"name"        validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=500"`
}

// ─────────────────────────────────────────
// List GET /api/v1/examples
// ─────────────────────────────────────────

// List godoc
// @Summary      ดึงรายการ examples ทั้งหมดของ user
// @Tags         examples
// @Security     BearerAuth
// @Produce      json
// @Param        page   query     int  false  "หน้าที่ต้องการ (default: 1)"
// @Param        limit  query     int  false  "จำนวนต่อหน้า (default: 20, max: 100)"
// @Success      200  {object}  models.APIResponse{data=models.PaginatedResponse[service.ExampleResponse]}
// @Failure      401  {object}  models.APIResponse
// @Failure      500  {object}  models.APIResponse
// @Router       /examples [get]
func (h *ExampleController) List(c fiber.Ctx) error {
	claims := mustGetClaims(c)
	pagination := parsePagination(c)

	items, total, err := h.svc.List(c.Context(), service.ListExamplesInput{
		UserID: claims.Sub,
		Page:   pagination.Page,
		Limit:  pagination.Limit,
	})
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			return c.Status(fiber.StatusBadRequest).JSON(models.Fail("invalid request"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.Fail("failed to fetch examples"))
	}

	return c.JSON(models.Ok(models.Paginate(items, total, pagination.Page, pagination.Limit)))
}

// ─────────────────────────────────────────
// GetByID GET /api/v1/examples/:id
// ─────────────────────────────────────────

// GetByID godoc
// @Summary      ดึง example ตาม ID
// @Tags         examples
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Example ID (UUID)"
// @Success      200  {object}  models.APIResponse{data=service.ExampleResponse}
// @Failure      400  {object}  models.APIResponse
// @Failure      404  {object}  models.APIResponse
// @Failure      500  {object}  models.APIResponse
// @Router       /examples/{id} [get]
func (h *ExampleController) GetByID(c fiber.Ctx) error {
	claims := mustGetClaims(c)

	result, err := h.svc.GetByID(c.Context(), c.Params("id"), claims.Sub)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			return c.Status(fiber.StatusBadRequest).JSON(models.Fail("invalid id"))
		}
		if errors.Is(err, service.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.Fail("example not found"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.Fail("failed to get example"))
	}

	return c.JSON(models.Ok(result))
}

// ─────────────────────────────────────────
// Create POST /api/v1/examples
// ─────────────────────────────────────────

// Create godoc
// @Summary      สร้าง example ใหม่
// @Tags         examples
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      CreateExampleRequest  true  "ข้อมูล example"
// @Success      201   {object}  models.APIResponse{data=service.ExampleResponse}
// @Failure      400   {object}  models.APIResponse
// @Failure      500   {object}  models.APIResponse
// @Router       /examples [post]
func (h *ExampleController) Create(c fiber.Ctx) error {
	claims := mustGetClaims(c)

	var req CreateExampleRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Fail("invalid request body"))
	}

	if err := validate.Struct(req); err != nil {
		errs := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(models.Fail(errs[0].Translate(nil)))
	}

	result, err := h.svc.Create(c.Context(), service.CreateExampleInput{
		Name:        req.Name,
		Description: req.Description,
		UserID:      claims.Sub,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Fail("failed to create example"))
	}

	return c.Status(fiber.StatusCreated).JSON(models.Ok(result))
}

// ─────────────────────────────────────────
// Update PUT /api/v1/examples/:id
// ─────────────────────────────────────────

// Update godoc
// @Summary      แก้ไข example
// @Tags         examples
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      string                true  "Example ID (UUID)"
// @Param        body  body      CreateExampleRequest  true  "ข้อมูลที่แก้ไข"
// @Success      200   {object}  models.APIResponse{data=service.ExampleResponse}
// @Failure      400   {object}  models.APIResponse
// @Failure      404   {object}  models.APIResponse
// @Failure      500   {object}  models.APIResponse
// @Router       /examples/{id} [put]
func (h *ExampleController) Update(c fiber.Ctx) error {
	claims := mustGetClaims(c)

	var req CreateExampleRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Fail("invalid request body"))
	}

	if err := validate.Struct(req); err != nil {
		errs := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(models.Fail(errs[0].Translate(nil)))
	}

	result, err := h.svc.Update(c.Context(), c.Params("id"), service.CreateExampleInput{
		Name:        req.Name,
		Description: req.Description,
		UserID:      claims.Sub,
	})
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.Fail("example not found or not yours"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.Fail("failed to update example"))
	}

	return c.JSON(models.Ok(result))
}

// ─────────────────────────────────────────
// Delete DELETE /api/v1/examples/:id
// ─────────────────────────────────────────

// Delete godoc
// @Summary      ลบ example
// @Tags         examples
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Example ID (UUID)"
// @Success      200  {object}  models.APIResponse
// @Failure      400  {object}  models.APIResponse
// @Failure      500  {object}  models.APIResponse
// @Router       /examples/{id} [delete]
func (h *ExampleController) Delete(c fiber.Ctx) error {
	claims := mustGetClaims(c)

	if err := h.svc.Delete(c.Context(), c.Params("id"), claims.Sub); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Fail("failed to delete example"))
	}

	return c.JSON(models.Ok(fiber.Map{"deleted": true}))
}
