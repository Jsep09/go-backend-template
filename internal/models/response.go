// internal/models/response.go
//
// ทำไมต้องมีไฟล์นี้?
// เพื่อให้ทุก endpoint ตอบ response format เดียวกันเสมอ
// { "data": ..., "error": null }
// { "data": null, "error": "message" }
// ถ้าไม่มี standard นี้ แต่ละ handler จะ response ต่างกัน debug ยากมาก
package models

// ─────────────────────────────────────────
// Standard response
// ─────────────────────────────────────────

type APIResponse struct {
	Data  any     `json:"data"`
	Error *string `json:"error"`
}

func Ok(data any) APIResponse {
	return APIResponse{Data: data, Error: nil}
}

func Fail(msg string) APIResponse {
	return APIResponse{Data: nil, Error: &msg}
}

// ─────────────────────────────────────────
// Paginated response
// ─────────────────────────────────────────

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

type PaginatedResponse[T any] struct {
	Items []T            `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}

// Paginate สร้าง PaginatedResponse พร้อม meta ให้อัตโนมัติ
// ใช้ได้กับทุก type — Paginate(items, total, page, limit)
func Paginate[T any](items []T, total int64, page, limit int) PaginatedResponse[T] {
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}
	return PaginatedResponse[T]{
		Items: items,
		Meta: PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}
}
