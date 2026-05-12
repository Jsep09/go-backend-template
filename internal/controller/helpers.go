package controller

import (
	"github.com/gofiber/fiber/v3"

	"github.com/Jsep09/go-backend-template/internal/middleware"
)

// mustGetClaims ดึง user claims จาก context
// ถ้าดึงไม่ได้ = ลืม protect route → panic ทันที
func mustGetClaims(c fiber.Ctx) *middleware.SupabaseClaims {
	claims, ok := middleware.GetUserClaims(c)
	if !ok {
		panic("mustGetClaims called on unprotected route")
	}
	return claims
}
