-- Script สำหรับรีเซ็ตฐานข้อมูล NurseShift
-- ใช้สำหรับเชื่อมต่อกับ PostgreSQL ที่ลงในเครื่อง (ไม่ใช่ Docker)

-- ลบฐานข้อมูลเดิม (ถ้ามี)
DROP DATABASE IF EXISTS nurseshift;

-- สร้างฐานข้อมูลใหม่
CREATE DATABASE nurseshift
    WITH 
    OWNER = nuttapong2
    ENCODING = 'UTF8'
    LC_COLLATE = 'C'
    LC_CTYPE = 'C'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

-- เชื่อมต่อกับฐานข้อมูลใหม่
\c nurseshift;

-- ตั้งค่า search_path
SET search_path TO nurse_shift, public;
