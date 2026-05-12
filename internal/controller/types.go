package controller

// ─────────────────────────────────────────
// Request types — HTTP input เท่านั้น
// ไม่ควรใช้ type นี้นอก controller layer
// ─────────────────────────────────────────

type CreateExampleRequest struct {
	Name        string `json:"name"        validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=500"`
}
