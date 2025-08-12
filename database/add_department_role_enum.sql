-- เพิ่ม enum type สำหรับ role ในแผนก
CREATE TYPE nurse_shift.department_role AS ENUM ('nurse', 'assistant');

-- เพิ่ม column department_role ในตาราง department_users
ALTER TABLE nurse_shift.department_users 
ADD COLUMN department_role nurse_shift.department_role NOT NULL DEFAULT 'nurse';

-- เพิ่ม index สำหรับ department_role
CREATE INDEX idx_department_users_department_role ON nurse_shift.department_users(department_role);

-- อัปเดตข้อมูลที่มีอยู่ให้เป็น 'nurse' เป็นค่าเริ่มต้น
UPDATE nurse_shift.department_users 
SET department_role = 'nurse' 
WHERE department_role IS NULL;

-- เพิ่ม comment สำหรับ column
COMMENT ON COLUMN nurse_shift.department_users.department_role IS 'Role ของผู้ใช้ในแผนก: nurse (พยาบาล) หรือ assistant (ผู้ช่วยพยาบาล)';
