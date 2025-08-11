-- Script สำหรับเพิ่มข้อมูลทดสอบในฐานข้อมูล NurseShift
-- ใช้สำหรับทดสอบ API จริง

-- เชื่อมต่อกับฐานข้อมูล
\c nurseshift;

-- ตั้งค่า search_path
SET search_path TO nurse_shift, public;

-- เพิ่มข้อมูลผู้ใช้ทดสอบ
INSERT INTO users (
    id, email, password_hash, first_name, last_name, phone, role, status, position,
    days_remaining, subscription_expires_at, package_type, max_departments, avatar_url,
    settings, created_at, updated_at
) VALUES 
-- Admin user (password: admin123)
(
    '550e8400-e29b-41d4-a716-446655440001'::uuid,
    'admin@nurseshift.com',
    '$2a$10$OQJnRxKT4dwQD1blpI2lze9/1NPb4XVl.V8Hle6a3p1p7CsIC7I4m',
    'ผู้ดูแล',
    'ระบบ',
    '+66-81-234-5678',
    'admin',
    'active',
    'System Administrator',
    90,
    NOW() + INTERVAL '90 days',
    'enterprise',
    20,
    NULL,
    '{"theme": "dark", "language": "th"}',
    NOW(),
    NOW()
),
-- Regular user (password: user123)
(
    '550e8400-e29b-41d4-a716-446655440002'::uuid,
    'user@nurseshift.com',
    '$2a$10$5G5QXW39SmBuPa9UF0Mft.rAXKMaLi0VeBuvesgjwFslILoE7ej.C',
    'พยาบาล',
    'ทดสอบ',
    '+66-82-345-6789',
    'user',
    'active',
    'หัวหน้าพยาบาล',
    85,
    NOW() + INTERVAL '30 days',
    'trial',
    2,
    NULL,
    '{"theme": "light", "language": "th"}',
    NOW(),
    NOW()
),
-- Test user 3 (password: test123)
(
    '550e8400-e29b-41d4-a716-446655440003'::uuid,
    'test@nurseshift.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    'พยาบาล',
    'ทดสอบ2',
    '+66-83-456-7890',
    'user',
    'active',
    'พยาบาลประจำแผนก',
    60,
    NOW() + INTERVAL '60 days',
    'standard',
    5,
    NULL,
    '{"theme": "light", "language": "th"}',
    NOW(),
    NOW()
);

-- เพิ่มข้อมูลแผนกทดสอบ
INSERT INTO departments (
    id, user_id, name, description, max_nurses, max_assistants, settings, is_active, created_at, updated_at
) VALUES 
(
    '660e8400-e29b-41d4-a716-446655440001'::uuid,
    '550e8400-e29b-41d4-a716-446655440002'::uuid,
    'แผนกผู้ป่วยใน',
    'แผนกผู้ป่วยในสำหรับผู้ป่วยที่ต้องนอนโรงพยาบาล',
    15,
    8,
    '{"shift_start": "08:00", "shift_end": "20:00"}',
    true,
    NOW(),
    NOW()
),
(
    '660e8400-e29b-41d4-a716-446655440002'::uuid,
    '550e8400-e29b-41d4-a716-446655440002'::uuid,
    'แผนกฉุกเฉิน',
    'แผนกฉุกเฉินสำหรับผู้ป่วยที่ต้องได้รับการรักษาเร่งด่วน',
    12,
    6,
    '{"shift_start": "06:00", "shift_end": "18:00"}',
    true,
    NOW(),
    NOW()
);

-- เพิ่มข้อมูลพนักงานในแผนก
INSERT INTO department_staff (
    id, department_id, name, position, phone, email, is_active, created_at, updated_at
) VALUES 
(
    '770e8400-e29b-41d4-a716-446655440001'::uuid,
    '660e8400-e29b-41d4-a716-446655440001'::uuid,
    'สมหญิง ใจดี',
    'พยาบาลประจำแผนก',
    '+66-84-567-8901',
    'somying@hospital.com',
    true,
    NOW(),
    NOW()
),
(
    '770e8400-e29b-41d4-a716-446655440002'::uuid,
    '660e8400-e29b-41d4-a716-446655440001'::uuid,
    'สมชาย รักงาน',
    'พยาบาลช่วยเหลือ',
    '+66-85-678-9012',
    'somchai@hospital.com',
    true,
    NOW(),
    NOW()
),
(
    '770e8400-e29b-41d4-a716-446655440003'::uuid,
    '660e8400-e29b-41d4-a716-446655440002'::uuid,
    'สมศรี ใจเย็น',
    'พยาบาลประจำแผนก',
    '+66-86-789-0123',
    'somsri@hospital.com',
    true,
    NOW(),
    NOW()
);

-- เพิ่มข้อมูลกะการทำงาน
INSERT INTO shifts (
    id, department_id, name, type, start_time, end_time, duration_hours, required_nurses, required_assistants, color, is_active, created_at, updated_at
) VALUES 
(
    '880e8400-e29b-41d4-a716-446655440001'::uuid,
    '660e8400-e29b-41d4-a716-446655440001'::uuid,
    'กะเช้า',
    'morning',
    '08:00:00',
    '16:00:00',
    8.00,
    5,
    2,
    '#3B82F6',
    true,
    NOW(),
    NOW()
),
(
    '880e8400-e29b-41d4-a716-446655440002'::uuid,
    '660e8400-e29b-41d4-a716-446655440001'::uuid,
    'กะบ่าย',
    'afternoon',
    '16:00:00',
    '00:00:00',
    8.00,
    4,
    2,
    '#10B981',
    true,
    NOW(),
    NOW()
),
(
    '880e8400-e29b-41d4-a716-446655440003'::uuid,
    '660e8400-e29b-41d4-a716-446655440001'::uuid,
    'กะดึก',
    'night',
    '00:00:00',
    '08:00:00',
    8.00,
    3,
    1,
    '#8B5CF6',
    true,
    NOW(),
    NOW()
);

-- เพิ่มข้อมูลวันทำงาน
INSERT INTO working_days (id, department_id, day_of_week, is_working_day, created_at) VALUES 
('990e8400-e29b-41d4-a716-446655440001'::uuid, '660e8400-e29b-41d4-a716-446655440001'::uuid, 0, true, NOW()),  -- วันอาทิตย์
('990e8400-e29b-41d4-a716-446655440002'::uuid, '660e8400-e29b-41d4-a716-446655440001'::uuid, 1, true, NOW()),  -- วันจันทร์
('990e8400-e29b-41d4-a716-446655440003'::uuid, '660e8400-e29b-41d4-a716-446655440001'::uuid, 2, true, NOW()),  -- วันอังคาร
('990e8400-e29b-41d4-a716-446655440004'::uuid, '660e8400-e29b-41d4-a716-446655440001'::uuid, 3, true, NOW()),  -- วันพุธ
('990e8400-e29b-41d4-a716-446655440005'::uuid, '660e8400-e29b-41d4-a716-446655440001'::uuid, 4, true, NOW()),  -- วันพฤหัสบดี
('990e8400-e29b-41d4-a716-446655440006'::uuid, '660e8400-e29b-41d4-a716-446655440001'::uuid, 5, true, NOW()),  -- วันศุกร์
('990e8400-e29b-41d4-a716-446655440007'::uuid, '660e8400-e29b-41d4-a716-446655440001'::uuid, 6, true, NOW());  -- วันเสาร์

-- เพิ่มข้อมูลวันหยุด
INSERT INTO holidays (
    id, department_id, name, start_date, end_date, is_recurring, recurrence_pattern, created_at, updated_at
) VALUES 
(
    'aa0e8400-e29b-41d4-a716-446655440001'::uuid,
    '660e8400-e29b-41d4-a716-446655440001'::uuid,
    'วันขึ้นปีใหม่',
    '2025-01-01',
    '2025-01-01',
    true,
    '{"type": "annual", "month": 1, "day": 1}',
    NOW(),
    NOW()
),
(
    'aa0e8400-e29b-41d4-a716-446655440002'::uuid,
    '660e8400-e29b-41d4-a716-446655440001'::uuid,
    'วันสงกรานต์',
    '2025-04-13',
    '2025-04-15',
    true,
    '{"type": "annual", "month": 4, "days": [13, 14, 15]}',
    NOW(),
    NOW()
);

-- เพิ่มข้อมูลลำดับความสำคัญ
INSERT INTO scheduling_priorities (
    id, user_id, name, description, priority_order, is_active, config, created_at, updated_at
) VALUES 
(
    'bb0e8400-e29b-41d4-a716-446655440001'::uuid,
    '550e8400-e29b-41d4-a716-446655440002'::uuid,
    'พยาบาลอาวุโส',
    'พยาบาลที่มีประสบการณ์สูงควรได้รับการจัดกะก่อน',
    1,
    true,
    '{"min_years": 5, "bonus_points": 10}',
    NOW(),
    NOW()
),
(
    'bb0e8400-e29b-41d4-a716-446655440002'::uuid,
    '550e8400-e29b-41d4-a716-446655440002'::uuid,
    'พยาบาลใหม่',
    'พยาบาลใหม่ควรได้รับการจัดกะที่เหมาะสม',
    2,
    true,
    '{"max_years": 2, "mentor_required": true}',
    NOW(),
    NOW()
);

-- แสดงข้อมูลที่เพิ่มเข้าไป
SELECT 'Users' as table_name, COUNT(*) as count FROM users
UNION ALL
SELECT 'Departments', COUNT(*) FROM departments
UNION ALL
SELECT 'Department Staff', COUNT(*) FROM department_staff
UNION ALL
SELECT 'Shifts', COUNT(*) FROM shifts
UNION ALL
SELECT 'Working Days', COUNT(*) FROM working_days
UNION ALL
SELECT 'Holidays', COUNT(*) FROM holidays
UNION ALL
SELECT 'Scheduling Priorities', COUNT(*) FROM scheduling_priorities;
