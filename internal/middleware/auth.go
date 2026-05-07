package middleware

// Verify JWT token ที่ออกโดย Supabase
// ทุก request ที่ผ่าน middleware นี้ → มั่นใจได้ว่า user login จริง
import (
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yourusername/go-backend-template/internal/models"
)

// SupabaseClaims คือ structure ของ JWT payload ที่ Supabase ออกให้
type SupabaseClaims struct {
	Sub                  string `json:"sub"`   // user ID (UUID)
	Email                string `json:"email"` // email ของ user
	Role                 string `json:"role"`  // "authenticated" หรือ role อื่น
	jwt.RegisteredClaims        // exp, iat, iss ฯลฯ
}

// contextKey เป็น type พิเศษสำหรับ context key
// ป้องกัน key ชนกับ package อื่น
type contextKey string

const UserClaimsKey contextKey = "userClaims"

func NewAuthJWT(jwtSecret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		// ──────────────────────────────────────
		// 1. ดึง token จาก Authorization header
		// ──────────────────────────────────────
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(
				models.Fail("missing authorization header"),
			)
		}

		// format ต้องเป็น "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			return c.Status(fiber.StatusUnauthorized).JSON(
				models.Fail("invalid authorization format, expected: Bearer <token>"),
			)
		}

		tokenString := parts[1]

		// ──────────────────────────────────────
		// 2. Parse และ Verify token
		// ──────────────────────────────────────
		claims := &SupabaseClaims{}

		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			// keyFunc บอก library ว่าใช้ secret อะไร verify
			func(token *jwt.Token) (any, error) {
				// ตรวจสอบ algorithm — ต้องเป็น HMAC (HS256) เท่านั้น
				// ป้องกัน "algorithm confusion attack"
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fiber.NewError(
						fiber.StatusUnauthorized,
						"unexpected signing method",
					)
				}
				return []byte(jwtSecret), nil
			},
		)

		if err != nil {
			slog.Warn("jwt verification failed",
				"error", err.Error(),
				"ip", c.IP(),
			)
			return c.Status(fiber.StatusUnauthorized).JSON(
				models.Fail("invalid or expired token"),
			)
		}

		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(
				models.Fail("invalid token"),
			)
		}

		// ──────────────────────────────────────
		// 3. เก็บ claims ไว้ใน context
		// handler ถัดไปจะดึงไปใช้ได้
		// ──────────────────────────────────────
		c.Locals(string(UserClaimsKey), claims)

		slog.Info("user authenticated",
			"user_id", claims.Sub,
			"email", claims.Email,
		)

		return c.Next()
	}
}

// GetUserClaims — helper สำหรับ handler ดึง user info
// ใช้แบบนี้ใน handler:
//
//	claims, ok := middleware.GetUserClaims(c)
//	if !ok { ... }
func GetUserClaims(c fiber.Ctx) (*SupabaseClaims, bool) {
	claims, ok := c.Locals(string(UserClaimsKey)).(*SupabaseClaims)
	return claims, ok
}
