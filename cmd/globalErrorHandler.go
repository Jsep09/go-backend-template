package main

import (
	"log/slog"

	"github.com/Jsep09/go-backend-template/internal/models"
	"github.com/gofiber/fiber/v3"
)

func globalErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	// ถ้า error มาจาก fiber เอง (เช่น 404, 405) ให้ใช้ status code นั้น
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	slog.Error("request error",
		"method", c.Method(),
		"path", c.Path(),
		"status", code,
		"error", err.Error(),
	)

	return c.Status(code).JSON(models.Fail(err.Error()))
}
