-- examples table สำหรับ template นี้
CREATE TABLE IF NOT EXISTS examples (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT        NOT NULL,
    description TEXT        NOT NULL DEFAULT '',
    user_id     UUID        NOT NULL,   -- อ้างอิง Supabase auth.users
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- index เพื่อให้ query WHERE user_id = $1 เร็ว
CREATE INDEX IF NOT EXISTS idx_examples_user_id ON examples(user_id);

-- auto update updated_at เมื่อมีการแก้ไข
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