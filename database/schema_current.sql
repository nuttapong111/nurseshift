-- NurseShift Management System Database Schema (Current Production Version)
-- PostgreSQL 16.3+
-- This schema reflects the CURRENT database structure after adding department_role

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "citext";

-- Create Schemas
CREATE SCHEMA IF NOT EXISTS nurse_shift;
SET search_path TO nurse_shift, public;

-- ===================================
-- ENUMS
-- ===================================

CREATE TYPE user_role AS ENUM ('admin', 'user');
CREATE TYPE user_status AS ENUM ('active', 'inactive', 'pending', 'suspended');
CREATE TYPE shift_type AS ENUM ('morning', 'afternoon', 'night', 'overtime');
CREATE TYPE leave_status AS ENUM ('pending', 'approved', 'rejected', 'cancelled');
CREATE TYPE leave_type AS ENUM ('sick', 'personal', 'vacation', 'emergency', 'maternity');
CREATE TYPE payment_status AS ENUM ('pending', 'approved', 'rejected', 'expired');
CREATE TYPE package_type AS ENUM ('standard', 'enterprise', 'trial');
CREATE TYPE notification_type AS ENUM ('schedule', 'leave', 'system', 'payment', 'reminder', 'holiday');
CREATE TYPE notification_priority AS ENUM ('low', 'medium', 'high');
-- NEW: Department role enum for staff in departments
CREATE TYPE department_role AS ENUM ('nurse', 'assistant');

-- ===================================
-- CORE TABLES
-- ===================================

-- Users (Admin และ User - หัวหน้าพยาบาล)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email citext UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    role user_role NOT NULL DEFAULT 'user',
    status user_status NOT NULL DEFAULT 'active',
    position VARCHAR(100),
    date_joined DATE NOT NULL DEFAULT CURRENT_DATE,
    avatar_url VARCHAR(500),
    last_login_at TIMESTAMP WITH TIME ZONE,
    days_remaining INTEGER DEFAULT 30,
    subscription_expires_at TIMESTAMP WITH TIME ZONE,
    package_type package_type DEFAULT 'trial',
    max_departments INTEGER DEFAULT 2,
    settings JSONB DEFAULT '{}',
    email_verified BOOLEAN DEFAULT false,
    email_verification_token VARCHAR(255),
    email_verification_expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Departments (แผนกที่หัวหน้าพยาบาลสร้าง)
CREATE TABLE departments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    head_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    max_nurses INTEGER DEFAULT 10,
    max_assistants INTEGER DEFAULT 5,
    settings JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Department Users (ความสัมพันธ์ระหว่าง users และ departments)
-- UPDATED: Added department_role column
CREATE TABLE department_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    department_id UUID NOT NULL REFERENCES departments(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    department_role department_role NOT NULL DEFAULT 'nurse', -- NEW: Role in department
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    assigned_by UUID REFERENCES users(id),
    UNIQUE(department_id, user_id)
);

-- Department Staff (พนักงานในแผนก)
CREATE TABLE department_staff (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    department_id UUID NOT NULL REFERENCES departments(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    position VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Shift Definitions
CREATE TABLE shifts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    department_id UUID NOT NULL REFERENCES departments(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type shift_type NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    duration_hours DECIMAL(4,2) NOT NULL,
    required_nurses INTEGER DEFAULT 1,
    required_assistants INTEGER DEFAULT 0,
    color VARCHAR(7) DEFAULT '#3B82F6', -- Hex color for UI
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Working Days Configuration
CREATE TABLE working_days (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    department_id UUID NOT NULL REFERENCES departments(id) ON DELETE CASCADE,
    day_of_week INTEGER NOT NULL CHECK (day_of_week >= 0 AND day_of_week <= 6), -- 0=Sunday, 6=Saturday
    is_working_day BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(department_id, day_of_week)
);

-- Holidays
CREATE TABLE holidays (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    department_id UUID NOT NULL REFERENCES departments(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_recurring BOOLEAN DEFAULT false,
    recurrence_pattern JSONB, -- For annual holidays
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CHECK (end_date >= start_date)
);

-- ===================================
-- SCHEDULING TABLES
-- ===================================

-- Schedule Assignments
CREATE TABLE schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    department_id UUID NOT NULL REFERENCES departments(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    shift_id UUID NOT NULL REFERENCES shifts(id) ON DELETE CASCADE,
    schedule_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'assigned',
    notes TEXT,
    assigned_by UUID REFERENCES users(id),
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, schedule_date, shift_id) -- Prevent double-booking
);

-- Leave Requests (หัวหน้าเวรกรอกวันที่พนักงานขอหยุดในแต่ละเดือน)
CREATE TABLE leave_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    department_id UUID NOT NULL REFERENCES departments(id) ON DELETE CASCADE,
    leave_type leave_type NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    reason TEXT,
    status leave_status DEFAULT 'pending',
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP WITH TIME ZONE,
    rejection_reason TEXT,
    attachments JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CHECK (end_date >= start_date)
);

-- ===================================
-- NOTIFICATION & AUDIT TABLES
-- ===================================

-- Notifications
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    type notification_type NOT NULL,
    priority notification_priority DEFAULT 'medium',
    is_read BOOLEAN DEFAULT false,
    read_at TIMESTAMP WITH TIME ZONE,
    action_url VARCHAR(500),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Audit Logs
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    table_name VARCHAR(100),
    record_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ===================================
-- PACKAGE & PAYMENT TABLES
-- ===================================

-- Packages
CREATE TABLE packages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    duration_days INTEGER NOT NULL,
    max_departments INTEGER DEFAULT 1,
    max_users INTEGER DEFAULT 5,
    features JSONB DEFAULT '[]',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Payments
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    package_id UUID REFERENCES packages(id) ON DELETE SET NULL,
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'THB',
    status payment_status DEFAULT 'pending',
    payment_method VARCHAR(50),
    transaction_id VARCHAR(255),
    paid_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ===================================
-- PRIORITY & SETTING TABLES
-- ===================================

-- Scheduling Priorities
CREATE TABLE scheduling_priorities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    department_id UUID NOT NULL REFERENCES departments(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    priority_level INTEGER NOT NULL CHECK (priority_level >= 1 AND priority_level <= 10),
    color VARCHAR(7) DEFAULT '#3B82F6',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(department_id, priority_level)
);

-- User Sessions
CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_token VARCHAR(255) UNIQUE NOT NULL,
    ip_address INET,
    user_agent TEXT,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Organizations (สำหรับ multi-tenant ในอนาคต)
CREATE TABLE organizations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    domain VARCHAR(255),
    settings JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ===================================
-- INDEXES
-- ===================================

-- Users indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_organization_id ON users(organization_id);

-- Departments indexes
CREATE INDEX idx_departments_created_by ON departments(created_by);
CREATE INDEX idx_departments_head_user_id ON departments(head_user_id);
CREATE INDEX idx_departments_is_active ON departments(is_active);

-- Department Users indexes
CREATE INDEX idx_department_users_department_id ON department_users(department_id);
CREATE INDEX idx_department_users_user_id ON department_users(user_id);
-- NEW: Index for department_role
CREATE INDEX idx_department_users_department_role ON department_users(department_role);

-- Department Staff indexes
CREATE INDEX idx_department_staff_department_id ON department_staff(department_id);
CREATE INDEX idx_department_staff_is_active ON department_staff(is_active);

-- Shifts indexes
CREATE INDEX idx_shifts_department_id ON shifts(department_id);
CREATE INDEX idx_shifts_type ON shifts(type);
CREATE INDEX idx_shifts_is_active ON shifts(is_active);

-- Schedules indexes
CREATE INDEX idx_schedules_department_id ON schedules(department_id);
CREATE INDEX idx_schedules_user_id ON schedules(user_id);
CREATE INDEX idx_schedules_shift_id ON schedules(shift_id);
CREATE INDEX idx_schedules_date ON schedules(schedule_date);

-- Leave Requests indexes
CREATE INDEX idx_leave_requests_user_id ON leave_requests(user_id);
CREATE INDEX idx_leave_requests_department_id ON leave_requests(department_id);
CREATE INDEX idx_leave_requests_status ON leave_requests(status);
CREATE INDEX idx_leave_requests_date_range ON leave_requests(start_date, end_date);

-- Notifications indexes
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);

-- Audit Logs indexes
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_table_name ON audit_logs(table_name);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);

-- Payments indexes
CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_created_at ON payments(created_at);

-- User Sessions indexes
CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_token ON user_sessions(session_token);
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at);

-- ===================================
-- TRIGGERS
-- ===================================

-- Function to update updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply updated_at trigger to relevant tables
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_departments_updated_at BEFORE UPDATE ON departments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_department_staff_updated_at BEFORE UPDATE ON department_staff FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_shifts_updated_at BEFORE UPDATE ON shifts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_holidays_updated_at BEFORE UPDATE ON holidays FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_schedules_updated_at BEFORE UPDATE ON schedules FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_leave_requests_updated_at BEFORE UPDATE ON leave_requests FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_notifications_updated_at BEFORE UPDATE ON notifications FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_packages_updated_at BEFORE UPDATE ON packages FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_payments_updated_at BEFORE UPDATE ON payments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_scheduling_priorities_updated_at BEFORE UPDATE ON scheduling_priorities FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_user_sessions_updated_at BEFORE UPDATE ON user_sessions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_organizations_updated_at BEFORE UPDATE ON organizations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ===================================
-- COMMENTS
-- ===================================

COMMENT ON TABLE users IS 'ผู้ใช้ระบบ (Admin และ User)';
COMMENT ON TABLE departments IS 'แผนกที่หัวหน้าพยาบาลสร้าง';
COMMENT ON TABLE department_users IS 'ความสัมพันธ์ระหว่าง users และ departments พร้อม role ในแผนก';
COMMENT ON TABLE department_staff IS 'พนักงานในแผนก';
COMMENT ON TABLE shifts IS 'กะการทำงาน';
COMMENT ON TABLE working_days IS 'วันทำงานของแผนก';
COMMENT ON TABLE holidays IS 'วันหยุดของแผนก';
COMMENT ON TABLE schedules IS 'ตารางเวร';
COMMENT ON TABLE leave_requests IS 'คำขอหยุดงาน';
COMMENT ON TABLE notifications IS 'การแจ้งเตือน';
COMMENT ON TABLE audit_logs IS 'บันทึกการทำงานของระบบ';
COMMENT ON TABLE packages IS 'แพ็คเกจสมาชิก';
COMMENT ON TABLE payments IS 'การชำระเงิน';
COMMENT ON TABLE scheduling_priorities IS 'ความสำคัญในการจัดตารางเวร';
COMMENT ON TABLE user_sessions IS 'เซสชันผู้ใช้';
COMMENT ON TABLE organizations IS 'องค์กร (สำหรับ multi-tenant)';

COMMENT ON COLUMN department_users.department_role IS 'Role ของผู้ใช้ในแผนก: nurse (พยาบาล) หรือ assistant (ผู้ช่วยพยาบาล)';
COMMENT ON COLUMN users.role IS 'Role ในระบบ: admin หรือ user';
COMMENT ON COLUMN users.status IS 'สถานะผู้ใช้: active, inactive, pending, suspended';
