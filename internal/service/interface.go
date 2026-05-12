package service

import "context"

// ExampleServiceInterface กำหนด contract ที่ Controller ใช้
// ทำให้ test ได้โดยไม่ต้องต่อ database จริง
type ExampleServiceInterface interface {
	List(ctx context.Context, input ListExamplesInput) ([]ExampleResponse, int64, error)
	GetByID(ctx context.Context, id, userID string) (ExampleResponse, error)
	Create(ctx context.Context, input CreateExampleInput) (ExampleResponse, error)
	Update(ctx context.Context, id string, input CreateExampleInput) (ExampleResponse, error)
	Delete(ctx context.Context, id, userID string) error
}
