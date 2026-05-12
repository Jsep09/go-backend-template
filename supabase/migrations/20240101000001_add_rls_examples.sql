-- ─────────────────────────────────────────
-- Migration: add_rls_examples
-- Description: เปิด RLS และสร้าง policies สำหรับ examples table
--
-- Why แยกไฟล์จาก create_examples?
-- schema (DDL) กับ security (RLS) ควรแยกกัน
-- ถ้าวันหลังต้อง adjust policy ก็สร้าง migration ใหม่ได้ชัดเจน
--
-- Note: Go backend ใช้ service role → bypass RLS ได้เลย
--       RLS ป้องกัน direct access ผ่าน Supabase client (anon/user key)
-- ─────────────────────────────────────────

-- UP ──────────────────────────────────────

ALTER TABLE examples ENABLE ROW LEVEL SECURITY;

-- SELECT: user เห็นแค่ row ของตัวเอง
CREATE POLICY "users can view own examples"
    ON examples FOR SELECT
    USING (auth.uid() = user_id);

-- INSERT: user insert ได้เฉพาะ row ที่ user_id ตรงกับตัวเอง
--         WITH CHECK ป้องกัน user A สร้างข้อมูลแทน user B
CREATE POLICY "users can insert own examples"
    ON examples FOR INSERT
    WITH CHECK (auth.uid() = user_id);

-- UPDATE: user แก้ได้เฉพาะ row ของตัวเอง
CREATE POLICY "users can update own examples"
    ON examples FOR UPDATE
    USING (auth.uid() = user_id);

-- DELETE: user ลบได้เฉพาะ row ของตัวเอง
CREATE POLICY "users can delete own examples"
    ON examples FOR DELETE
    USING (auth.uid() = user_id);

-- DOWN (rollback) ──────────────────────────
-- DROP POLICY "users can view own examples"   ON examples;
-- DROP POLICY "users can insert own examples" ON examples;
-- DROP POLICY "users can update own examples" ON examples;
-- DROP POLICY "users can delete own examples" ON examples;
-- ALTER TABLE examples DISABLE ROW LEVEL SECURITY;
