package middleware

// internal/middleware/logger.go
//
// Log ทุก HTTP request ด้วย slog (structured logging)
// structured = log เป็น JSON → ง่ายต่อการ search ใน production

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
)

func NewLogger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		// เรียก handler ถัดไปก่อน
		err := c.Next()

		// หลัง handler ทำงานเสร็จ → log ผลลัพธ์
		duration := time.Since(start)
		status := c.Response().StatusCode()

		// เลือก log level ตาม status code
		switch {
		case status >= 500:
			slog.Error("request completed",
				"method", c.Method(),
				"path", c.Path(),
				"status", status,
				"duration_ms", duration.Milliseconds(),
				"ip", c.IP(),
			)
		case status >= 400:
			slog.Warn("request completed",
				"method", c.Method(),
				"path", c.Path(),
				"status", status,
				"duration_ms", duration.Milliseconds(),
				"ip", c.IP(),
			)
		default:
			slog.Info("request completed",
				"method", c.Method(),
				"path", c.Path(),
				"status", status,
				"duration_ms", duration.Milliseconds(),
				"ip", c.IP(),
			)
		}

		return err
	}
}
