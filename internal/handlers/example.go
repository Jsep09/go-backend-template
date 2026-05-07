package handlers

import (
	"github.com/Jsep09/go-backend-template/internal/models"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExampleHandler struct {
	db *pgxpool.Pool
}

func NewExampleHandler(db *pgxpool.Pool) *ExampleHandler {
	return &ExampleHandler{db: db}
}

func (h *ExampleHandler) List(c fiber.Ctx) error {
	return c.JSON(models.Ok([]any{}))
}

func (h *ExampleHandler) GetByID(c fiber.Ctx) error {
	return c.JSON(models.Ok(fiber.Map{"id": c.Params("id")}))
}

func (h *ExampleHandler) Create(c fiber.Ctx) error {
	return c.Status(fiber.StatusCreated).JSON(models.Ok(nil))
}

func (h *ExampleHandler) Update(c fiber.Ctx) error {
	return c.JSON(models.Ok(nil))
}

func (h *ExampleHandler) Delete(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}
