package service

// ─────────────────────────────────────────
// Output types — ส่งกลับไปให้ Controller
// ─────────────────────────────────────────

type ExampleResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      string `json:"user_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ─────────────────────────────────────────
// Input types — รับมาจาก Controller
// ─────────────────────────────────────────

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
