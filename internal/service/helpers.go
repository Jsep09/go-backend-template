package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/Jsep09/go-backend-template/internal/db/generated"
)

// parseUUID แปลง string → pgtype.UUID (ใช้ภายใน service เท่านั้น)
func parseUUID(s string) (pgtype.UUID, error) {
	var id pgtype.UUID
	err := id.Scan(s)
	return id, err
}

// toExampleResponse แปลง db.Example → ExampleResponse
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
