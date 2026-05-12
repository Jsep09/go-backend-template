package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"

	"github.com/Jsep09/go-backend-template/internal/models"
)

type SupabaseClaims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

const UserClaimsKey = "userClaims"

func NewAuthJWT(jwtSecret string) fiber.Handler {
	secret := []byte(jwtSecret)

	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		lower := strings.ToLower(authHeader)
		if !strings.HasPrefix(lower, "bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(
				models.Fail("missing or malformed token"),
			)
		}

		// trim จาก original (authHeader) ไม่ใช่ lower
		// เพราะ token ข้างหลัง "Bearer " เป็น case-sensitive
		tokenStr := authHeader[len("Bearer "):] // ตัด 7 ตัวแรกออก ("Bearer " = 7 chars)
		claims := &SupabaseClaims{}

		_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return secret, nil
		})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(
				models.Fail("invalid or expired token"),
			)
		}

		c.Locals(UserClaimsKey, claims)
		return c.Next()
	}
}

func GetUserClaims(c fiber.Ctx) (*SupabaseClaims, bool) {
	claims, ok := c.Locals(UserClaimsKey).(*SupabaseClaims)
	return claims, ok
}
