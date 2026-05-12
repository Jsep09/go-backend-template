package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/Jsep09/go-backend-template/internal/db/generated"
	"github.com/Jsep09/go-backend-template/internal/middleware"
)

// mustGetClaims ดึง user claims จาก context
// ถ้าดึงไม่ได้ = bug ของโปรแกรมเมอร์ (ลืม protect route) → panic ทันที
func mustGetClaims(c fiber.Ctx) *middleware.SupabaseClaims {
	claims, ok := middleware.GetUserClaims(c)
	if !ok {
		panic("mustGetClaims called on unprotected route")
	}
	return claims
}

// parseUUID แปลง string → pgtype.UUID
// ใช้ร่วมกันได้ทุก service ไม่ต้องเขียนซ้ำ
func parseUUID(s string) (pgtype.UUID, error) {
	var id pgtype.UUID
	err := id.Scan(s)
	return id, err
}

// toExampleResponse แปลง db.Example → ExampleResponse
// ใช้ร่วมกันได้ทุก method ไม่ต้องเขียนซ้ำ
func toExampleResponse(ex db.Example) ExampleResponse {
	return ExampleResponse{
		ID:          uuid.UUID(ex.ID.Bytes).String(),
		Name:        ex.Name,
		Description: ex.Description,
		UserID:      uuid.UUID(ex.UserID.Bytes).String(),
		CreatedAt:   ex.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:   ex.UpdatedAt.Time.Format(time.RFC3339),
	}
}
