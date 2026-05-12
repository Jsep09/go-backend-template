-- ─────────────────────────────────────────
-- Migration: create_examples
-- Description: สร้าง examples table พร้อม index และ updated_at trigger
-- ─────────────────────────────────────────

-- UP ──────────────────────────────────────

CREATE TABLE IF NOT EXISTS examples (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT        NOT NULL,
    description TEXT        NOT NULL DEFAULT '',
    user_id     UUID        NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- index เพื่อให้ query WHERE user_id = $1 เร็ว
CREATE INDEX IF NOT EXISTS idx_examples_user_id ON examples(user_id);

-- function + trigger สำหรับ auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER examples_updated_at
    BEFORE UPDATE ON examples
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

-- DOWN (rollback) ──────────────────────────
-- DROP TRIGGER IF EXISTS examples_updated_at ON examples;
-- DROP FUNCTION IF EXISTS update_updated_at();
-- DROP TABLE IF EXISTS examples;
