package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/timeout"
)

// NewTimeout จำกัดเวลาแต่ละ request ไม่เกิน d
// ถ้า handler ทำงานเกิน d → return fiber.ErrRequestTimeout (408)
// ซึ่ง globalErrorHandler จะรับไปตอบ 408 ให้อัตโนมัติ
func NewTimeout(d time.Duration) fiber.Handler {
	return timeout.New(func(c fiber.Ctx) error {
		return c.Next()
	}, timeout.Config{
		Timeout: d,
	})
}
