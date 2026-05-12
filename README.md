# go-backend-template

Production-ready Go backend template สำหรับต่อกับ Supabase พร้อมใช้งานได้ทันที

**Stack:** Go + Fiber v3 · Supabase (Auth + DB) · sqlc · pgx v5 · Swagger

---

## Features

- **JWT Auth** — verify Supabase token อัตโนมัติทุก protected route
- **Clean Architecture** — Controller → Service → DB แยกชัด
- **Type-safe DB** — sqlc generate Go code จาก SQL ให้อัตโนมัติ
- **Swagger UI** — เปิดที่ `/swagger` สำหรับ dev, ปิดอัตโนมัติใน production
- **Pagination** — มาตรฐานเดียวกันทุก list endpoint
- **RLS Ready** — migration มี Row Level Security พร้อมแล้ว
- **Rate Limit + Timeout** — ป้องกัน abuse และ request ค้าง
- **Structured Logging** — text สำหรับ dev, JSON สำหรับ production

---

## ต้องมีก่อน

- [Go 1.22+](https://go.dev/dl/)
- [Supabase account](https://supabase.com) + project
- [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html) — `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
- [swag](https://github.com/swaggo/swag) — `go install github.com/swaggo/swag/cmd/swag@latest`
- [supabase CLI](https://supabase.com/docs/guides/cli) — สำหรับ push migration

---

## เริ่มต้นใช้งาน

**1. Clone และติดตั้ง dependencies**
```bash
git clone https://github.com/Jsep09/go-backend-template.git
cd go-backend-template
go mod tidy
```

**2. ตั้งค่า environment**
```bash
cp .env.example .env
```

แก้ไข `.env`:
```env
DATABASE_URL=postgresql://postgres:[PASSWORD]@db.[PROJECT-REF].supabase.co:5432/postgres
SUPABASE_JWT_SECRET=your-jwt-secret   # Supabase Dashboard → Settings → API → JWT Secret
```

**3. Push migration ขึ้น Supabase**
```bash
make migrate
```

**4. รัน server**
```bash
make run
```

เปิด http://localhost:3000/swagger เพื่อดู API docs

---

## คำสั่งที่ใช้บ่อย

| คำสั่ง | ทำอะไร |
|---|---|
| `make run` | รัน server (development) |
| `make build` | build binary ไว้ที่ `bin/app` |
| `make test` | รัน unit tests |
| `make test-cover` | รัน tests + ดู coverage |
| `make docs` | regenerate Swagger docs |
| `make sqlc` | regenerate DB code จาก SQL |
| `make migrate` | push migration ขึ้น Supabase |

---

## โครงสร้าง Project

```
cmd/                        ← entry point
├── main.go                 ← boot sequence
├── config.go               ← อ่าน environment variables
├── routes.go               ← register routes ทั้งหมด
├── database.go             ← connection pool
├── startServer.go          ← graceful shutdown
└── globalErrorHandler.go   ← จัดการ error response กลาง

internal/
├── controller/             ← HTTP layer (รับ request, ส่ง response)
│   ├── example.go          ← handlers + Swagger annotations
│   ├── types.go            ← request structs
│   └── helpers.go          ← mustGetClaims, parsePagination
├── service/                ← Business logic (ไม่มี fiber)
│   ├── example.go          ← business logic
│   ├── interface.go        ← service interfaces
│   ├── types.go            ← input/output structs
│   └── helpers.go          ← utility functions
├── middleware/             ← JWT, CORS, Rate limit, Timeout, Logger
├── models/                 ← response format มาตรฐาน
└── db/
    ├── queries/            ← SQL files (แก้ที่นี่)
    └── generated/          ← sqlc output (ห้ามแก้มือ)

supabase/
└── migrations/             ← database migrations
```

---

## วิธีเพิ่ม Feature ใหม่ (เช่น `tags`)

### 1. สร้าง Migration

```bash
supabase migration new create_tags
```

เขียน SQL ใน `supabase/migrations/xxx_create_tags.sql`:

```sql
-- UP
CREATE TABLE tags (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT        NOT NULL,
    user_id    UUID        NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE tags ENABLE ROW LEVEL SECURITY;

CREATE POLICY "users can manage own tags"
    ON tags FOR ALL
    USING (auth.uid() = user_id);

-- DOWN
-- DROP TABLE IF EXISTS tags;
```

```bash
make migrate
```

### 2. เขียน SQL Queries

สร้าง `internal/db/queries/tags.sql`:

```sql
-- name: ListTags :many
SELECT * FROM tags WHERE user_id = $1 ORDER BY created_at DESC;

-- name: CreateTag :one
INSERT INTO tags (name, user_id) VALUES ($1, $2)
RETURNING *;
```

```bash
make sqlc   # generate Go code จาก SQL
```

### 3. สร้าง Service

สร้าง `internal/service/tags.go` โดย copy pattern จาก `example.go`:

```go
type TagService struct {
    queries db.Querier
}

func NewTagService(queries db.Querier) *TagService {
    return &TagService{queries: queries}
}

func (s *TagService) List(ctx context.Context, userID string) ([]TagResponse, error) {
    // business logic
}
```

เพิ่ม interface ใน `internal/service/interface.go`:

```go
type TagServiceInterface interface {
    List(ctx context.Context, userID string) ([]TagResponse, error)
    Create(ctx context.Context, input CreateTagInput) (TagResponse, error)
}
```

### 4. สร้าง Controller

สร้าง `internal/controller/tags.go`:

```go
type TagController struct {
    svc service.TagServiceInterface
}

// List godoc
// @Summary  ดึงรายการ tags
// @Tags     tags
// @Security BearerAuth
// @Router   /tags [get]
func (h *TagController) List(c fiber.Ctx) error {
    claims := mustGetClaims(c)
    items, err := h.svc.List(c.Context(), claims.Sub)
    // ...
    return c.JSON(models.Ok(items))
}
```

### 5. Register Routes

เพิ่มใน `cmd/routes.go`:

```go
tagSvc  := service.NewTagService(queries)
tagCtrl := controller.NewTagController(tagSvc)
auth.Get("/tags", tagCtrl.List)
auth.Post("/tags", tagCtrl.Create)
```

### 6. Regenerate Swagger Docs

```bash
make docs
```

---

## Response Format

ทุก endpoint ตอบ format เดียวกัน:

```json
// สำเร็จ
{ "data": { ... }, "error": null }

// ผิดพลาด
{ "data": null, "error": "message" }

// List พร้อม pagination
{
  "data": {
    "items": [...],
    "meta": {
      "page": 1,
      "limit": 20,
      "total": 100,
      "total_pages": 5,
      "has_next": true,
      "has_prev": false
    }
  },
  "error": null
}
```

---

## Authentication

ทุก protected route ต้องส่ง Supabase JWT token:

```
Authorization: Bearer <supabase-access-token>
```

ดึง token จาก Supabase client ฝั่ง frontend:

```ts
const { data: { session } } = await supabase.auth.getSession()
const token = session?.access_token
```

---

## Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `DATABASE_URL` | ✅ | — | Supabase connection string |
| `SUPABASE_JWT_SECRET` | ✅ | — | JWT secret จาก Supabase Dashboard |
| `APP_ENV` | | `development` | `development` หรือ `production` |
| `APP_PORT` | | `3000` | port ที่ server listen |
| `APP_NAME` | | `go-backend-template` | ชื่อ app (ใช้ใน header) |
| `ALLOWED_ORIGINS` | | `http://localhost:5173` | CORS origins คั่นด้วย `,` |
| `RATE_LIMIT_MAX` | | `60` | requests ต่อ window |
| `RATE_LIMIT_WINDOW` | | `1` | window size (นาที) |

---

## Migration Workflow

```bash
# สร้าง migration ใหม่
supabase migration new <ชื่อ>

# ดูสถานะ
supabase migration list

# push ขึ้น Supabase
make migrate
```

> **กฎสำคัญ:** migration ที่ push ไปแล้วห้ามแก้ไข ให้สร้างไฟล์ใหม่เสมอ

---

## License

MIT
