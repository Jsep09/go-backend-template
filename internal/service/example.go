package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"

	db "github.com/Jsep09/go-backend-template/internal/db/generated"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
)

// ─────────────────────────────────────────
// Service struct
// ─────────────────────────────────────────

type ExampleService struct {
	queries *db.Queries
}

func NewExampleService(queries *db.Queries) *ExampleService {
	return &ExampleService{queries: queries}
}

// ─────────────────────────────────────────
// Types — plain Go ไม่มี fiber
// ─────────────────────────────────────────

type ExampleResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      string `json:"user_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type CreateExampleInput struct {
	Name        string
	Description string
	UserID      string
}

type ListExamplesInput struct {
	UserID string
	Page   int
	Limit  int
}

// ─────────────────────────────────────────
// Methods — รับ plain string, return plain types
// ─────────────────────────────────────────

func (s *ExampleService) List(ctx context.Context, input ListExamplesInput) ([]ExampleResponse, int64, error) {
	uid, err := parseUUID(input.UserID)
	if err != nil {
		return nil, 0, ErrInvalidInput
	}

	// query items พร้อม pagination
	examples, err := s.queries.ListExamples(ctx, db.ListExamplesParams{
		UserID: uid,
		Limit:  int32(input.Limit),
		Offset: int32((input.Page - 1) * input.Limit),
	})
	if err != nil {
		slog.Error("failed to list examples", "error", err, "user_id", input.UserID)
		return nil, 0, err
	}

	// query total count สำหรับคำนวณ total_pages
	total, err := s.queries.CountExamples(ctx, uid)
	if err != nil {
		slog.Error("failed to count examples", "error", err, "user_id", input.UserID)
		return nil, 0, err
	}

	results := make([]ExampleResponse, 0, len(examples))
	for _, ex := range examples {
		results = append(results, toExampleResponse(ex))
	}
	return results, total, nil
}

func (s *ExampleService) GetByID(ctx context.Context, id, userID string) (ExampleResponse, error) {
	exID, err := parseUUID(id)
	if err != nil {
		return ExampleResponse{}, ErrInvalidInput
	}

	uid, err := parseUUID(userID)
	if err != nil {
		return ExampleResponse{}, ErrInvalidInput
	}

	ex, err := s.queries.GetExample(ctx, db.GetExampleParams{
		ID:     exID,
		UserID: uid,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ExampleResponse{}, ErrNotFound // แปลงเป็น sentinel error
		}
		slog.Error("failed to get example", "error", err)
		return ExampleResponse{}, err
	}

	return toExampleResponse(ex), nil
}

func (s *ExampleService) Create(ctx context.Context, input CreateExampleInput) (ExampleResponse, error) {
	uid, err := parseUUID(input.UserID)
	if err != nil {
		return ExampleResponse{}, ErrInvalidInput
	}

	ex, err := s.queries.CreateExample(ctx, db.CreateExampleParams{
		Name:        input.Name,
		Description: input.Description,
		UserID:      uid,
	})
	if err != nil {
		slog.Error("failed to create example", "error", err, "user_id", input.UserID)
		return ExampleResponse{}, err
	}

	return toExampleResponse(ex), nil
}

func (s *ExampleService) Update(ctx context.Context, id string, input CreateExampleInput) (ExampleResponse, error) {
	exID, err := parseUUID(id)
	if err != nil {
		return ExampleResponse{}, ErrInvalidInput
	}

	uid, err := parseUUID(input.UserID)
	if err != nil {
		return ExampleResponse{}, ErrInvalidInput
	}

	ex, err := s.queries.UpdateExample(ctx, db.UpdateExampleParams{
		Name:        input.Name,
		Description: input.Description,
		ID:          exID,
		UserID:      uid,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ExampleResponse{}, ErrNotFound
		}
		slog.Error("failed to update example", "error", err)
		return ExampleResponse{}, err
	}

	return toExampleResponse(ex), nil
}

func (s *ExampleService) Delete(ctx context.Context, id, userID string) error {
	exID, err := parseUUID(id)
	if err != nil {
		return ErrInvalidInput
	}

	uid, err := parseUUID(userID)
	if err != nil {
		return ErrInvalidInput
	}

	if err := s.queries.DeleteExample(ctx, db.DeleteExampleParams{
		ID:     exID,
		UserID: uid,
	}); err != nil {
		slog.Error("failed to delete example", "error", err)
		return err
	}

	return nil
}
