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

	// ─────────────────────────────────
	// API v1
	// ─────────────────────────────────
	v1 := app.Group("/api/v1")

	// Protected — ต้อง JWT
	auth := v1.Group("", middleware.NewAuthJWT(mustGetEnv("SUPABASE_JWT_SECRET")))

	// me — ดู user info จาก token
	auth.Get("/me", func(c fiber.Ctx) error {
		claims, ok := middleware.GetUserClaims(c)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(
				models.Fail("unauthorized"),
			)
		}
		return c.JSON(models.Ok(fiber.Map{
			"user_id": claims.Sub,
			"role":    claims.Role,
		}))
	})

	// Example resource (ต้องการ DB)
	example := handlers.NewExampleHandler(db)
	auth.Get("/examples", example.List)
	auth.Post("/examples", example.Create)
	auth.Get("/examples/:id", example.GetByID)
	auth.Put("/examples/:id", example.Update)
	auth.Delete("/examples/:id", example.Delete)

	// 404 handler
	app.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(models.Fail("route not found"))
	})
}
