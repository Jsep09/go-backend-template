package service

// ทำไม test อยู่ใน package เดียวกับ service?
// เพราะต้องการเข้าถึง internal helpers เช่น parseUUID, toExampleResponse
// ถ้าใช้ package service_test จะเห็นแค่ exported symbols

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/Jsep09/go-backend-template/internal/db/generated"
)

// ─────────────────────────────────────────
// Mock DB — implement db.Querier ทั้งหมด
// แต่ละ test จะ override เฉพาะ method ที่ต้องการ
// ─────────────────────────────────────────

type mockQuerier struct {
	listFn    func(ctx context.Context, arg db.ListExamplesParams) ([]db.Example, error)
	countFn   func(ctx context.Context, userID pgtype.UUID) (int64, error)
	getFn     func(ctx context.Context, arg db.GetExampleParams) (db.Example, error)
	createFn  func(ctx context.Context, arg db.CreateExampleParams) (db.Example, error)
	updateFn  func(ctx context.Context, arg db.UpdateExampleParams) (db.Example, error)
	deleteFn  func(ctx context.Context, arg db.DeleteExampleParams) error
}

func (m *mockQuerier) ListExamples(ctx context.Context, arg db.ListExamplesParams) ([]db.Example, error) {
	return m.listFn(ctx, arg)
}
func (m *mockQuerier) CountExamples(ctx context.Context, userID pgtype.UUID) (int64, error) {
	return m.countFn(ctx, userID)
}
func (m *mockQuerier) GetExample(ctx context.Context, arg db.GetExampleParams) (db.Example, error) {
	return m.getFn(ctx, arg)
}
func (m *mockQuerier) CreateExample(ctx context.Context, arg db.CreateExampleParams) (db.Example, error) {
	return m.createFn(ctx, arg)
}
func (m *mockQuerier) UpdateExample(ctx context.Context, arg db.UpdateExampleParams) (db.Example, error) {
	return m.updateFn(ctx, arg)
}
func (m *mockQuerier) DeleteExample(ctx context.Context, arg db.DeleteExampleParams) error {
	return m.deleteFn(ctx, arg)
}

// ─────────────────────────────────────────
// Helpers สำหรับสร้าง test data
// ─────────────────────────────────────────

func newTestUUID(t *testing.T) pgtype.UUID {
	t.Helper()
	uid, _ := parseUUID(uuid.NewString())
	return uid
}

func newFakeExample(id, userID pgtype.UUID, name string) db.Example {
	return db.Example{
		ID:          id,
		UserID:      userID,
		Name:        name,
		Description: "test description",
		CreatedAt:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
}

// ─────────────────────────────────────────
// Tests — GetByID
// ─────────────────────────────────────────

func TestGetByID_Success(t *testing.T) {
	exID := newTestUUID(t)
	uid := newTestUUID(t)
	fake := newFakeExample(exID, uid, "hello")

	svc := NewExampleService(&mockQuerier{
		getFn: func(_ context.Context, arg db.GetExampleParams) (db.Example, error) {
			return fake, nil // จำลองว่า DB เจอ record
		},
	})

	result, err := svc.GetByID(context.Background(),
		uuid.UUID(exID.Bytes).String(),
		uuid.UUID(uid.Bytes).String(),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Name != "hello" {
		t.Errorf("expected name=hello, got %s", result.Name)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	exID := newTestUUID(t)
	uid := newTestUUID(t)

	svc := NewExampleService(&mockQuerier{
		getFn: func(_ context.Context, _ db.GetExampleParams) (db.Example, error) {
			return db.Example{}, pgx.ErrNoRows // จำลองว่า DB ไม่เจอ
		},
	})

	_, err := svc.GetByID(context.Background(),
		uuid.UUID(exID.Bytes).String(),
		uuid.UUID(uid.Bytes).String(),
	)

	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestGetByID_InvalidID(t *testing.T) {
	uid := newTestUUID(t)

	svc := NewExampleService(&mockQuerier{}) // getFn ไม่ถูกเรียกเลย

	_, err := svc.GetByID(context.Background(), "not-a-uuid", uuid.UUID(uid.Bytes).String())

	if err != ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

// ─────────────────────────────────────────
// Tests — Create
// ─────────────────────────────────────────

func TestCreate_Success(t *testing.T) {
	uid := newTestUUID(t)
	exID := newTestUUID(t)

	svc := NewExampleService(&mockQuerier{
		createFn: func(_ context.Context, arg db.CreateExampleParams) (db.Example, error) {
			// เช็คว่า service ส่งข้อมูลถูกต้องไปให้ DB
			if arg.Name != "test name" {
				t.Errorf("expected name=test name, got %s", arg.Name)
			}
			return newFakeExample(exID, uid, arg.Name), nil
		},
	})

	result, err := svc.Create(context.Background(), CreateExampleInput{
		Name:        "test name",
		Description: "test desc",
		UserID:      uuid.UUID(uid.Bytes).String(),
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Name != "test name" {
		t.Errorf("expected name=test name, got %s", result.Name)
	}
}

func TestCreate_InvalidUserID(t *testing.T) {
	svc := NewExampleService(&mockQuerier{})

	_, err := svc.Create(context.Background(), CreateExampleInput{
		Name:   "test",
		UserID: "bad-uuid",
	})

	if err != ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

// ─────────────────────────────────────────
// Tests — List
// ─────────────────────────────────────────

func TestList_Success(t *testing.T) {
	uid := newTestUUID(t)
	exID := newTestUUID(t)

	svc := NewExampleService(&mockQuerier{
		listFn: func(_ context.Context, arg db.ListExamplesParams) ([]db.Example, error) {
			// เช็คว่า pagination ถูกส่งไปถูกต้อง
			if arg.Limit != 10 {
				t.Errorf("expected limit=10, got %d", arg.Limit)
			}
			if arg.Offset != 10 { // page=2, limit=10 → offset=10
				t.Errorf("expected offset=10, got %d", arg.Offset)
			}
			return []db.Example{newFakeExample(exID, uid, "item1")}, nil
		},
		countFn: func(_ context.Context, _ pgtype.UUID) (int64, error) {
			return 25, nil // จำลองว่ามี 25 records ทั้งหมด
		},
	})

	items, total, err := svc.List(context.Background(), ListExamplesInput{
		UserID: uuid.UUID(uid.Bytes).String(),
		Page:   2,
		Limit:  10,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 25 {
		t.Errorf("expected total=25, got %d", total)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}
}

// ─────────────────────────────────────────
// Tests — Delete
// ─────────────────────────────────────────

func TestDelete_Success(t *testing.T) {
	exID := newTestUUID(t)
	uid := newTestUUID(t)
	called := false

	svc := NewExampleService(&mockQuerier{
		deleteFn: func(_ context.Context, _ db.DeleteExampleParams) error {
			called = true
			return nil
		},
	})

	err := svc.Delete(context.Background(),
		uuid.UUID(exID.Bytes).String(),
		uuid.UUID(uid.Bytes).String(),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !called {
		t.Error("expected deleteFn to be called")
	}
}
