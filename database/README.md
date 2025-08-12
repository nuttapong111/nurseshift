# NurseShift Database Management

## 📁 ไฟล์ในโฟลเดอร์นี้

### Schema Files
- **`schema_current.sql`** - Schema ปัจจุบันที่ตรงกับฐานข้อมูล production
- **`schema.sql`** - Schema หลักของระบบ
- **`schema_old.sql`** - Schema เก่า (เก็บไว้เป็น reference)

### Migration Files
- **`migration_add_department_role.sql`** - Migration สำหรับเพิ่ม department_role support
- **`add_department_role_enum.sql`** - Script สำหรับเพิ่ม enum และ column (development)

### Data Files
- **`seed.sql`** - ข้อมูลเริ่มต้นสำหรับ development
- **`seed_test_data.sql`** - ข้อมูลทดสอบ

### Docker
- **`docker-compose.yml`** - สำหรับรัน PostgreSQL ใน Docker

## 🚀 การใช้งานใน Production

### 1. Deploy Schema ใหม่
```bash
# รัน schema ปัจจุบัน
psql -h [HOST] -U [USER] -d [DATABASE] -f schema_current.sql
```

### 2. รัน Migration (ถ้าจำเป็น)
```bash
# รัน migration script
psql -h [HOST] -U [USER] -d [DATABASE] -f migration_add_department_role.sql
```

### 3. ตรวจสอบสถานะ
```bash
# ตรวจสอบว่า migration สำเร็จ
psql -h [HOST] -U [USER] -d [DATABASE] -c "
SELECT 
    'Migration Status' as check_type,
    CASE WHEN EXISTS (SELECT 1 FROM pg_type WHERE typname = 'department_role') THEN 'PASS' ELSE 'FAIL' END as enum_exists,
    CASE WHEN EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'nurse_shift' AND table_name = 'department_users' AND column_name = 'department_role') THEN 'PASS' ELSE 'FAIL' END as column_exists;
"
```

## 🔄 การอัปเดต Schema

### เมื่อมีการเปลี่ยนแปลงฐานข้อมูล:

1. **อัปเดตฐานข้อมูล** (development)
2. **อัปเดต `schema_current.sql`** ให้ตรงกับฐานข้อมูล
3. **สร้าง migration script** (ถ้าจำเป็น)
4. **ทดสอบ migration** ใน development
5. **Deploy ไป production**

### ตัวอย่างการอัปเดต schema_current.sql:
```bash
# Export schema ปัจจุบันจากฐานข้อมูล
pg_dump -h localhost -U postgres -d nurseshift --schema-only --no-owner --no-privileges > schema_current.sql
```

## ⚠️ ข้อควรระวัง

### Production Deployment:
- **สำรองฐานข้อมูล** ก่อนรัน migration
- **ทดสอบ migration** ใน staging environment ก่อน
- **รัน migration** ในช่วง off-peak hours
- **มี rollback plan** เตรียมไว้

### Schema Changes:
- **ไม่ลบ column** ที่มีข้อมูลอยู่
- **เพิ่ม column ใหม่** ควรมี default value
- **สร้าง index** สำหรับ column ที่ใช้ query บ่อย
- **อัปเดต comments** ให้ชัดเจน

## 📊 สถานะปัจจุบัน

### ✅ สิ่งที่เพิ่มแล้ว:
- `department_role` enum (`'nurse'`, `'assistant'`)
- `department_role` column ในตาราง `department_users`
- Index สำหรับ `department_role`
- Comments สำหรับ column และ table

### 🔍 การใช้งาน:
- **ตาราง `users`**: role เป็น `admin`, `user` (ระบบ authentication)
- **ตาราง `department_users`**: `department_role` เป็น `nurse`, `assistant` (role ในแผนก)
- **การแยก role**: ไม่สับสนระหว่าง role ในระบบกับ role ในแผนก

## 🛠️ การแก้ไขปัญหา

### ถ้า migration ล้มเหลว:
```bash
# ตรวจสอบ error logs
psql -h [HOST] -U [USER] -d [DATABASE] -c "\l+"

# ตรวจสอบ table structure
psql -h [HOST] -U [USER] -d [DATABASE] -c "\d+ nurse_shift.department_users"

# ตรวจสอบ enum types
psql -h [HOST] -U [USER] -d [DATABASE] -c "\dT+ nurse_shift.*"
```

### Rollback (ถ้าจำเป็น):
```bash
# เปิดไฟล์ migration script และ uncomment rollback section
# จากนั้นรัน rollback commands
```
