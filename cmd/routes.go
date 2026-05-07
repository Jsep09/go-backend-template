package main

import (
	"github.com/Jsep09/go-backend-template/internal/handlers"
	"github.com/Jsep09/go-backend-template/internal/middleware"
	"github.com/Jsep09/go-backend-template/internal/models"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func registerRoutes(app *fiber.App, db *pgxpool.Pool) {
	// Health check — ไม่ต้อง auth, ใช้ check ว่า service ยัง alive
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(models.Ok(fiber.Map{
			"status":  "ok",
			"service": "go-backend-template",
		}))
	})

	// API v1 group
	v1 := app.Group("/api/v1")

	// Public routes (ไม่ต้อง login)
	// v1.Post("/auth/...", ...)

	// Protected routes (ต้อง JWT)
	protected := v1.Group("", middleware.NewAuthJWT(
		mustGetEnv("SUPABASE_JWT_SECRET"),
	))

	// Example resource
	exampleHandler := handlers.NewExampleHandler(db)
	protected.Get("/examples", exampleHandler.List)
	protected.Get("/examples/:id", exampleHandler.GetByID)
	protected.Post("/examples", exampleHandler.Create)
	protected.Put("/examples/:id", exampleHandler.Update)
	protected.Delete("/examples/:id", exampleHandler.Delete)

	// 404 handler
	app.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(models.Fail("route not found"))
	})
}
