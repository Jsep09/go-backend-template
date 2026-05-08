package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"

	"github.com/Jsep09/go-backend-template/internal/models"
)

func NewRateLimiter(max int, window time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: window,

		KeyGenerator: func(c fiber.Ctx) string {
			return c.IP()
		},

		LimitReached: func(c fiber.Ctx) error {
			slog.Warn("rate limit exceeded", "ip", c.IP(), "path", c.Path())
			return c.Status(fiber.StatusTooManyRequests).JSON(
				models.Fail("too many requests, please try again later"),
			)
		},

		Next: func(c fiber.Ctx) bool {
			return c.Method() == fiber.MethodOptions
		},
	})
}
