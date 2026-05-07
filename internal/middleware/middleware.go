// internal/middleware/cors.go
//
// CORS = Cross-Origin Resource Sharing
// ป้องกันไม่ให้ browser เรียก API จาก domain ที่ไม่ได้รับอนุญาต
// เช่น ถ้า frontend อยู่ที่ evil.com จะเรียก API เราไม่ได้

package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func NewCORS(allowedOrigins []string) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins, // v3 รับ []string ตรงๆ ได้เลย
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	})
}
