package main

import (
	"time"

	dbgen "github.com/Jsep09/go-backend-template/internal/db/generated"
	"github.com/Jsep09/go-backend-template/internal/controller"
	"github.com/Jsep09/go-backend-template/internal/middleware"
	"github.com/Jsep09/go-backend-template/internal/models"
	"github.com/Jsep09/go-backend-template/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func registerRoutes(app *fiber.App, pool *pgxpool.Pool, appEnv string) {
	// สร้าง queries จาก pool — ส่งให้ทุก service ใช้ร่วมกัน
	queries := dbgen.New(pool)

	// Swagger — เปิดเฉพาะ development เท่านั้น
	// production ปิดเพื่อไม่ expose API spec ให้คนภายนอก
	if appEnv != "production" {
		app.Get("/swagger", func(c fiber.Ctx) error {
			c.Set("Content-Type", "text/html")
			return c.SendString(`<!DOCTYPE html>
<html>
  <head>
    <title>go-backend-template API</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.onload = () => {
        SwaggerUIBundle({
          url: "/swagger/doc.json",
          dom_id: "#swagger-ui",
          presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
          layout: "BaseLayout",
          persistAuthorization: true,
        })
      }
    </script>
  </body>
</html>`)
		})

		// Swagger spec — orval ชี้มาที่ URL นี้เพื่อ generate TypeScript hooks
		app.Get("/swagger/doc.json", func(c fiber.Ctx) error {
			return c.SendFile("./docs/swagger.json")
		})
	}

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
	v1 := app.Group("/api/v1", middleware.NewTimeout(30*time.Second))

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

	// Example resource — wire service → controller
	exampleSvc := service.NewExampleService(queries)
	exampleCtrl := controller.NewExampleController(exampleSvc)
	auth.Get("/examples", exampleCtrl.List)
	auth.Post("/examples", exampleCtrl.Create)
	auth.Get("/examples/:id", exampleCtrl.GetByID)
	auth.Put("/examples/:id", exampleCtrl.Update)
	auth.Delete("/examples/:id", exampleCtrl.Delete)

	// 404 handler
	app.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(models.Fail("route not found"))
	})
}
