package controller

import (
	"strconv"

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

// PaginationParams คือ query params สำหรับ pagination
type PaginationParams struct {
	Page  int
	Limit int
}

// parsePagination อ่าน ?page=&limit= จาก query string
// default: page=1, limit=20, max limit=100
func parsePagination(c fiber.Ctx) PaginationParams {
	page := queryInt(c, "page", 1)
	limit := queryInt(c, "limit", 20)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return PaginationParams{Page: page, Limit: limit}
}

// queryInt อ่าน query string แล้วแปลงเป็น int
// ถ้า parse ไม่ได้หรือไม่มีค่า ใช้ defaultVal แทน
func queryInt(c fiber.Ctx, key string, defaultVal int) int {
	val := c.Query(key)
	if val == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return n
}
