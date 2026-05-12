package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"

	db "github.com/Jsep09/go-backend-template/internal/db/generated"
)

// ErrNotFound — sentinel error ให้ controller แปลงเป็น 404
var ErrNotFound = errors.New("not found")

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

// ─────────────────────────────────────────
// Methods — รับ plain string, return plain types
// ─────────────────────────────────────────

func (s *ExampleService) List(ctx context.Context, userID string) ([]ExampleResponse, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return nil, err
	}

	examples, err := s.queries.ListExamples(ctx, uid)
	if err != nil {
		slog.Error("failed to list examples", "error", err, "user_id", userID)
		return nil, err
	}

	results := make([]ExampleResponse, 0, len(examples))
	for _, ex := range examples {
		results = append(results, toExampleResponse(ex))
	}
	return results, nil
}

func (s *ExampleService) GetByID(ctx context.Context, id, userID string) (ExampleResponse, error) {
	exID, err := parseUUID(id)
	if err != nil {
		return ExampleResponse{}, err
	}

	uid, err := parseUUID(userID)
	if err != nil {
		return ExampleResponse{}, err
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
		return ExampleResponse{}, err
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
		return ExampleResponse{}, err
	}

	uid, err := parseUUID(input.UserID)
	if err != nil {
		return ExampleResponse{}, err
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
		return err
	}

	uid, err := parseUUID(userID)
	if err != nil {
		return err
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
