package main

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"

	_ "github.com/Jsep09/go-backend-template/docs" // ← trigger init() ให้ swagger handler รู้จัก spec
	"github.com/Jsep09/go-backend-template/internal/middleware"
)

func main() {
	// 2. โหลด .env ก่อน เพราะ logger ต้องการ APP_ENV จาก .env
	if err := godotenv.Load(); err != nil {
		slog.Info("no .env file found, using system environment variables")
	}

	// 3. อ่าน Config จาก environment
	cfg := loadConfig()

	// 1. Setup Logger — dev ใช้ text (อ่านง่าย), prod ใช้ JSON (machine-readable)
	var logHandler slog.Handler
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	if cfg.AppEnv == "production" {
		logHandler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, opts)
	}
	slog.SetDefault(slog.New(logHandler))

	// 4. เชื่อมต่อ Database

	db, err := connectDB(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close() // ปิด connection pool ตอน app ปิด
	slog.Info("database connected successfully")

	// 5. สร้าง Fiber App
	app := fiber.New(fiber.Config{
		// ซ่อน Fiber version จาก response header (security)
		ServerHeader: cfg.AppName,

		// Error handler กลาง — จัดการ error ที่ handler
		ErrorHandler: globalErrorHandler,
	})

	// 6. Register Middleware (ลำดับสำคัญมาก!)
	app.Use(middleware.NewCORS(cfg.AllowedOrigins)) // CORS
	app.Use(middleware.NewRateLimiter(              // Rate limit หลัง CORS
		cfg.RateLimitMax,
		cfg.RateLimitWindow,
	))
	app.Use(middleware.NewLogger()) // Request logger

	// 7. Register Routes
	registerRoutes(app, db, cfg.AppEnv)

	// 8. Start Server + Graceful Shutdown
	startServer(app, cfg.AppPort)
}
