// internal/models/response.go
//
// ทำไมต้องมีไฟล์นี้?
// เพื่อให้ทุก endpoint ตอบ response format เดียวกันเสมอ
// { "data": ..., "error": null }
// { "data": null, "error": "message" }
// ถ้าไม่มี standard นี้ แต่ละ handler จะ response ต่างกัน debug ยากมาก
package models

type APIResponse struct {
	Data  any     `json:"data"`
	Error *string `json:"error"`
}

func Ok(data any) APIResponse {
	return APIResponse{
		Data:  data,
		Error: nil,
	}
}

func Fail(msg string) APIResponse {
	return APIResponse{
		Data:  nil,
		Error: &msg,
	}
}
