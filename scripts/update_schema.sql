-- อัปเดต Schema สำหรับระบบการยืนยันอีเมล
-- ไฟล์: scripts/update_schema.sql

BEGIN;

-- เพิ่มคอลัมน์ใหม่ในตาราง users (ใน schema nurse_shift)
ALTER TABLE nurse_shift.users ADD COLUMN IF NOT EXISTS email_verified BOOLEAN DEFAULT FALSE;
ALTER TABLE nurse_shift.users ADD COLUMN IF NOT EXISTS email_verification_token VARCHAR(255);
ALTER TABLE nurse_shift.users ADD COLUMN IF NOT EXISTS email_verification_expires_at TIMESTAMP;

-- อัปเดตข้อมูลที่มีอยู่ (ให้ผู้ใช้ที่มีอยู่มี email_verified = true)
UPDATE nurse_shift.users SET email_verified = true WHERE email_verified IS NULL;

-- สร้าง index สำหรับประสิทธิภาพ
CREATE INDEX IF NOT EXISTS idx_users_email_verified ON nurse_shift.users(email_verified);
CREATE INDEX IF NOT EXISTS idx_users_verification_token ON nurse_shift.users(email_verification_token);

-- แสดงข้อความยืนยัน
DO $$
BEGIN
    RAISE NOTICE 'Schema updated successfully!';
    RAISE NOTICE 'Added columns: email_verified, email_verification_token, email_verification_expires_at';
    RAISE NOTICE 'Created indexes for better performance';
END $$;

COMMIT;
