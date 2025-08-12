--
-- PostgreSQL database dump
--

-- Dumped from database version 15.13 (Homebrew)
-- Dumped by pg_dump version 15.13 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: nurse_shift; Type: SCHEMA; Schema: -; Owner: nuttapong2
--

CREATE SCHEMA nurse_shift;


ALTER SCHEMA nurse_shift OWNER TO nuttapong2;

--
-- Name: citext; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;


--
-- Name: EXTENSION citext; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION citext IS 'data type for case-insensitive character strings';


--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: leave_status; Type: TYPE; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TYPE nurse_shift.leave_status AS ENUM (
    'pending',
    'approved',
    'rejected',
    'cancelled'
);


ALTER TYPE nurse_shift.leave_status OWNER TO nuttapong2;

--
-- Name: leave_type; Type: TYPE; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TYPE nurse_shift.leave_type AS ENUM (
    'sick',
    'personal',
    'vacation',
    'emergency',
    'maternity'
);


ALTER TYPE nurse_shift.leave_type OWNER TO nuttapong2;

--
-- Name: notification_priority; Type: TYPE; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TYPE nurse_shift.notification_priority AS ENUM (
    'low',
    'medium',
    'high'
);


ALTER TYPE nurse_shift.notification_priority OWNER TO nuttapong2;

--
-- Name: notification_type; Type: TYPE; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TYPE nurse_shift.notification_type AS ENUM (
    'schedule',
    'leave',
    'system',
    'payment',
    'reminder',
    'holiday'
);


ALTER TYPE nurse_shift.notification_type OWNER TO nuttapong2;

--
-- Name: package_type; Type: TYPE; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TYPE nurse_shift.package_type AS ENUM (
    'standard',
    'enterprise',
    'trial'
);


ALTER TYPE nurse_shift.package_type OWNER TO nuttapong2;

--
-- Name: payment_status; Type: TYPE; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TYPE nurse_shift.payment_status AS ENUM (
    'pending',
    'approved',
    'rejected',
    'expired'
);


ALTER TYPE nurse_shift.payment_status OWNER TO nuttapong2;

--
-- Name: shift_type; Type: TYPE; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TYPE nurse_shift.shift_type AS ENUM (
    'morning',
    'afternoon',
    'night',
    'overtime'
);


ALTER TYPE nurse_shift.shift_type OWNER TO nuttapong2;

--
-- Name: user_role; Type: TYPE; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TYPE nurse_shift.user_role AS ENUM (
    'admin',
    'user'
);


ALTER TYPE nurse_shift.user_role OWNER TO nuttapong2;

--
-- Name: user_status; Type: TYPE; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TYPE nurse_shift.user_status AS ENUM (
    'active',
    'inactive',
    'pending',
    'suspended'
);


ALTER TYPE nurse_shift.user_status OWNER TO nuttapong2;

--
-- Name: create_audit_log(); Type: FUNCTION; Schema: nurse_shift; Owner: nuttapong2
--

CREATE FUNCTION nurse_shift.create_audit_log() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    INSERT INTO audit_logs (
        user_id,
        action,
        resource_type,
        resource_id,
        old_data,
        new_data
    ) VALUES (
        COALESCE(NEW.updated_by, OLD.updated_by, current_setting('app.current_user_id', true)::UUID),
        TG_OP,
        TG_TABLE_NAME,
        COALESCE(NEW.id, OLD.id),
        CASE WHEN TG_OP = 'DELETE' THEN row_to_json(OLD) ELSE NULL END,
        CASE WHEN TG_OP IN ('INSERT', 'UPDATE') THEN row_to_json(NEW) ELSE NULL END
    );
    RETURN COALESCE(NEW, OLD);
END;
$$;


ALTER FUNCTION nurse_shift.create_audit_log() OWNER TO nuttapong2;

--
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: nurse_shift; Owner: nuttapong2
--

CREATE FUNCTION nurse_shift.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;


ALTER FUNCTION nurse_shift.update_updated_at_column() OWNER TO nuttapong2;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: audit_logs; Type: TABLE; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TABLE nurse_shift.audit_logs (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid,
    action character varying(100) NOT NULL,
    resource_type character varying(100) NOT NULL,
    resource_id uuid,
    old_data jsonb,
    new_data jsonb,
    ip_address inet,
    user_agent text,
    created_at timestamp with time zone DEFAULT now()
);


ALTER TABLE nurse_shift.audit_logs OWNER TO nuttapong2;

--
-- Name: department_staff; Type: TABLE; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TABLE nurse_shift.department_staff (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    department_id uuid NOT NULL,
    name character varying(100) NOT NULL,
    "position" character varying(100) NOT NULL,
    phone character varying(20),
    email character varying(255),
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE nurse_shift.department_staff OWNER TO nuttapong2;

--
-- Name: departments; Type: TABLE; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TABLE nurse_shift.departments (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    max_nurses integer DEFAULT 10,
    max_assistants integer DEFAULT 5,
    settings jsonb DEFAULT '{}'::jsonb,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE nurse_shift.departments OWNER TO nuttapong2;

--
-- Name: holidays; Type: TABLE; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TABLE nurse_shift.holidays (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    department_id uuid NOT NULL,
    name character varying(255) NOT NULL,
    start_date date NOT NULL,
    end_date date NOT NULL,
    is_recurring boolean DEFAULT false,
    recurrence_pattern jsonb,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT holidays_check CHECK ((end_date >= start_date))
);


ALTER TABLE nurse_shift.holidays OWNER TO nuttapong2;

--
-- Name: leave_requests; Type: TABLE; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TABLE nurse_shift.leave_requests (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    department_id uuid NOT NULL,
    staff_id uuid NOT NULL,
    start_date date NOT NULL,
    end_date date NOT NULL,
    leave_type nurse_shift.leave_type NOT NULL,
    reason text,
    status nurse_shift.leave_status DEFAULT 'pending'::nurse_shift.leave_status,
    requested_by uuid NOT NULL,
    approved_by uuid,
    approved_at timestamp with time zone,
    rejection_reason text,
    notes text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE nurse_shift.leave_requests OWNER TO nuttapong2;

--
-- Name: notifications; Type: TABLE; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TABLE nurse_shift.notifications (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid,
    type nurse_shift.notification_type NOT NULL,
    priority nurse_shift.notification_priority DEFAULT 'medium'::nurse_shift.notification_priority,
    title character varying(255) NOT NULL,
    message text NOT NULL,
    action_url character varying(500),
    is_read boolean DEFAULT false,
    read_at timestamp with time zone,
    data jsonb,
    created_at timestamp with time zone DEFAULT now()
);


ALTER TABLE nurse_shift.notifications OWNER TO nuttapong2;

--
-- Name: packages; Type: TABLE; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TABLE nurse_shift.packages (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(255) NOT NULL,
    type nurse_shift.package_type NOT NULL,
    price numeric(10,2) NOT NULL,
    duration_days integer NOT NULL,
    max_departments integer,
    features jsonb DEFAULT '[]'::jsonb NOT NULL,
    is_active boolean DEFAULT true,
    is_popular boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE nurse_shift.packages OWNER TO nuttapong2;

--
-- Name: payments; Type: TABLE; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TABLE nurse_shift.payments (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    package_id uuid NOT NULL,
    amount numeric(10,2) NOT NULL,
    status nurse_shift.payment_status DEFAULT 'pending'::nurse_shift.payment_status,
    payment_date date NOT NULL,
    evidence_url character varying(500),
    approved_by uuid,
    approved_at timestamp with time zone,
    rejection_reason text,
    extended_days integer,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE nurse_shift.payments OWNER TO nuttapong2;

--
-- Name: schedules; Type: TABLE; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TABLE nurse_shift.schedules (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    department_id uuid NOT NULL,
    staff_id uuid NOT NULL,
    shift_id uuid NOT NULL,
    schedule_date date NOT NULL,
    status character varying(20) DEFAULT 'assigned'::character varying,
    notes text,
    assigned_by uuid,
    assigned_at timestamp with time zone DEFAULT now(),
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE nurse_shift.schedules OWNER TO nuttapong2;

--
-- Name: scheduling_priorities; Type: TABLE; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TABLE nurse_shift.scheduling_priorities (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    priority_order integer NOT NULL,
    is_active boolean DEFAULT true,
    config jsonb DEFAULT '{}'::jsonb,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE nurse_shift.scheduling_priorities OWNER TO nuttapong2;

--
-- Name: shifts; Type: TABLE; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TABLE nurse_shift.shifts (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    department_id uuid NOT NULL,
    name character varying(100) NOT NULL,
    type nurse_shift.shift_type NOT NULL,
    start_time time without time zone NOT NULL,
    end_time time without time zone NOT NULL,
    duration_hours numeric(4,2) NOT NULL,
    required_nurses integer DEFAULT 1,
    required_assistants integer DEFAULT 0,
    color character varying(7) DEFAULT '#3B82F6'::character varying,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE nurse_shift.shifts OWNER TO nuttapong2;

--
-- Name: user_sessions; Type: TABLE; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TABLE nurse_shift.user_sessions (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    token_hash character varying(255) NOT NULL,
    ip_address inet,
    user_agent text,
    created_at timestamp with time zone DEFAULT now(),
    expires_at timestamp with time zone NOT NULL,
    revoked_at timestamp with time zone
);


ALTER TABLE nurse_shift.user_sessions OWNER TO nuttapong2;

--
-- Name: users; Type: TABLE; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TABLE nurse_shift.users (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    email public.citext NOT NULL,
    password_hash character varying(255) NOT NULL,
    first_name character varying(100) NOT NULL,
    last_name character varying(100) NOT NULL,
    phone character varying(20),
    role nurse_shift.user_role DEFAULT 'user'::nurse_shift.user_role NOT NULL,
    status nurse_shift.user_status DEFAULT 'active'::nurse_shift.user_status NOT NULL,
    "position" character varying(100),
    days_remaining integer DEFAULT 30,
    subscription_expires_at timestamp with time zone,
    package_type nurse_shift.package_type DEFAULT 'trial'::nurse_shift.package_type,
    max_departments integer DEFAULT 2,
    avatar_url character varying(500),
    settings jsonb DEFAULT '{}'::jsonb,
    last_login_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE nurse_shift.users OWNER TO nuttapong2;

--
-- Name: working_days; Type: TABLE; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TABLE nurse_shift.working_days (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    department_id uuid NOT NULL,
    day_of_week integer NOT NULL,
    is_working_day boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT working_days_day_of_week_check CHECK (((day_of_week >= 0) AND (day_of_week <= 6)))
);


ALTER TABLE nurse_shift.working_days OWNER TO nuttapong2;

--
-- Data for Name: audit_logs; Type: TABLE DATA; Schema: nurse_shift; Owner: nuttapong2
--

COPY nurse_shift.audit_logs (id, user_id, action, resource_type, resource_id, old_data, new_data, ip_address, user_agent, created_at) FROM stdin;
\.


--
-- Data for Name: department_staff; Type: TABLE DATA; Schema: nurse_shift; Owner: nuttapong2
--

COPY nurse_shift.department_staff (id, department_id, name, "position", phone, email, is_active, created_at, updated_at) FROM stdin;
770e8400-e29b-41d4-a716-446655440001	660e8400-e29b-41d4-a716-446655440001	สมหญิง ใจดี	พยาบาลประจำแผนก	+66-84-567-8901	somying@hospital.com	t	2025-08-11 10:17:49.146501+07	2025-08-11 10:17:49.146501+07
770e8400-e29b-41d4-a716-446655440002	660e8400-e29b-41d4-a716-446655440001	สมชาย รักงาน	พยาบาลช่วยเหลือ	+66-85-678-9012	somchai@hospital.com	t	2025-08-11 10:17:49.146501+07	2025-08-11 10:17:49.146501+07
770e8400-e29b-41d4-a716-446655440003	660e8400-e29b-41d4-a716-446655440002	สมศรี ใจเย็น	พยาบาลประจำแผนก	+66-86-789-0123	somsri@hospital.com	t	2025-08-11 10:17:49.146501+07	2025-08-11 10:17:49.146501+07
\.


--
-- Data for Name: departments; Type: TABLE DATA; Schema: nurse_shift; Owner: nuttapong2
--

COPY nurse_shift.departments (id, user_id, name, description, max_nurses, max_assistants, settings, is_active, created_at, updated_at) FROM stdin;
660e8400-e29b-41d4-a716-446655440001	550e8400-e29b-41d4-a716-446655440002	แผนกผู้ป่วยใน	แผนกผู้ป่วยในสำหรับผู้ป่วยที่ต้องนอนโรงพยาบาล	15	8	{"shift_end": "20:00", "shift_start": "08:00"}	t	2025-08-11 10:17:49.144456+07	2025-08-11 10:17:49.144456+07
660e8400-e29b-41d4-a716-446655440002	550e8400-e29b-41d4-a716-446655440002	แผนกฉุกเฉิน	แผนกฉุกเฉินสำหรับผู้ป่วยที่ต้องได้รับการรักษาเร่งด่วน	12	6	{"shift_end": "18:00", "shift_start": "06:00"}	t	2025-08-11 10:17:49.144456+07	2025-08-11 10:17:49.144456+07
\.


--
-- Data for Name: holidays; Type: TABLE DATA; Schema: nurse_shift; Owner: nuttapong2
--

COPY nurse_shift.holidays (id, department_id, name, start_date, end_date, is_recurring, recurrence_pattern, created_at, updated_at) FROM stdin;
aa0e8400-e29b-41d4-a716-446655440001	660e8400-e29b-41d4-a716-446655440001	วันขึ้นปีใหม่	2025-01-01	2025-01-01	t	{"day": 1, "type": "annual", "month": 1}	2025-08-11 10:17:49.150497+07	2025-08-11 10:17:49.150497+07
aa0e8400-e29b-41d4-a716-446655440002	660e8400-e29b-41d4-a716-446655440001	วันสงกรานต์	2025-04-13	2025-04-15	t	{"days": [13, 14, 15], "type": "annual", "month": 4}	2025-08-11 10:17:49.150497+07	2025-08-11 10:17:49.150497+07
\.


--
-- Data for Name: leave_requests; Type: TABLE DATA; Schema: nurse_shift; Owner: nuttapong2
--

COPY nurse_shift.leave_requests (id, department_id, staff_id, start_date, end_date, leave_type, reason, status, requested_by, approved_by, approved_at, rejection_reason, notes, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: notifications; Type: TABLE DATA; Schema: nurse_shift; Owner: nuttapong2
--

COPY nurse_shift.notifications (id, user_id, type, priority, title, message, action_url, is_read, read_at, data, created_at) FROM stdin;
\.


--
-- Data for Name: packages; Type: TABLE DATA; Schema: nurse_shift; Owner: nuttapong2
--

COPY nurse_shift.packages (id, name, type, price, duration_days, max_departments, features, is_active, is_popular, created_at, updated_at) FROM stdin;
9de9f65d-f735-45cd-ac4c-e7cb5215ca7f	แพ็คเกจทดลองใช้	trial	0.00	30	2	["จัดการแผนกพื้นฐาน", "พนักงานสูงสุด 10 คน", "ตารางเวรแบบง่าย"]	t	f	2025-08-11 10:06:32.945129+07	2025-08-11 10:06:32.945129+07
4561bdee-c817-469d-8523-82c79ea12bf2	แพ็คเกจมาตรฐาน	standard	990.00	30	5	["จัดการหลายแผนก", "พนักงานไม่จำกัด", "ตารางเวรอัตโนมัติ", "การแจ้งเตือนแบบเรียลไทม์", "รายงานและสถิติ"]	t	t	2025-08-11 10:06:32.945129+07	2025-08-11 10:06:32.945129+07
55c32284-d842-4b47-a6d8-6efe12d1179d	แพ็คเกจระดับองค์กร	enterprise	2990.00	90	20	["จัดการหลายแผนกไม่จำกัด", "พนักงานไม่จำกัด", "ตารางเวรอัตโนมัติด้วย AI", "การแจ้งเตือนแบบเรียลไทม์", "รายงานและสถิติขั้นสูง", "การสำรองข้อมูล", "การสนับสนุนลูกค้าแบบพิเศษ"]	t	f	2025-08-11 10:06:32.945129+07	2025-08-11 10:06:32.945129+07
\.


--
-- Data for Name: payments; Type: TABLE DATA; Schema: nurse_shift; Owner: nuttapong2
--

COPY nurse_shift.payments (id, user_id, package_id, amount, status, payment_date, evidence_url, approved_by, approved_at, rejection_reason, extended_days, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: schedules; Type: TABLE DATA; Schema: nurse_shift; Owner: nuttapong2
--

COPY nurse_shift.schedules (id, department_id, staff_id, shift_id, schedule_date, status, notes, assigned_by, assigned_at, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: scheduling_priorities; Type: TABLE DATA; Schema: nurse_shift; Owner: nuttapong2
--

COPY nurse_shift.scheduling_priorities (id, user_id, name, description, priority_order, is_active, config, created_at, updated_at) FROM stdin;
bb0e8400-e29b-41d4-a716-446655440001	550e8400-e29b-41d4-a716-446655440002	พยาบาลอาวุโส	พยาบาลที่มีประสบการณ์สูงควรได้รับการจัดกะก่อน	1	t	{"min_years": 5, "bonus_points": 10}	2025-08-11 10:17:49.151386+07	2025-08-11 10:17:49.151386+07
bb0e8400-e29b-41d4-a716-446655440002	550e8400-e29b-41d4-a716-446655440002	พยาบาลใหม่	พยาบาลใหม่ควรได้รับการจัดกะที่เหมาะสม	2	t	{"max_years": 2, "mentor_required": true}	2025-08-11 10:17:49.151386+07	2025-08-11 10:17:49.151386+07
\.


--
-- Data for Name: shifts; Type: TABLE DATA; Schema: nurse_shift; Owner: nuttapong2
--

COPY nurse_shift.shifts (id, department_id, name, type, start_time, end_time, duration_hours, required_nurses, required_assistants, color, is_active, created_at, updated_at) FROM stdin;
880e8400-e29b-41d4-a716-446655440001	660e8400-e29b-41d4-a716-446655440001	กะเช้า	morning	08:00:00	16:00:00	8.00	5	2	#3B82F6	t	2025-08-11 10:17:49.147812+07	2025-08-11 10:17:49.147812+07
880e8400-e29b-41d4-a716-446655440002	660e8400-e29b-41d4-a716-446655440001	กะบ่าย	afternoon	16:00:00	00:00:00	8.00	4	2	#10B981	t	2025-08-11 10:17:49.147812+07	2025-08-11 10:17:49.147812+07
880e8400-e29b-41d4-a716-446655440003	660e8400-e29b-41d4-a716-446655440001	กะดึก	night	00:00:00	08:00:00	8.00	3	1	#8B5CF6	t	2025-08-11 10:17:49.147812+07	2025-08-11 10:17:49.147812+07
\.


--
-- Data for Name: user_sessions; Type: TABLE DATA; Schema: nurse_shift; Owner: nuttapong2
--

COPY nurse_shift.user_sessions (id, user_id, token_hash, ip_address, user_agent, created_at, expires_at, revoked_at) FROM stdin;
b78fcbfb-a950-4348-903a-9fefbb80cca4	8fe56507-de13-42a6-b2bc-3268699202a6	ae9c83fa24d0bd5b10eb4d31465339b27388d1dcb52807c8d94442bdfa7a9312	127.0.0.1	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36	0001-01-01 06:42:04+06:42:04	2025-08-11 10:52:16.103063+07	2025-08-11 14:48:37.474414+07
4f57bad3-8435-46ee-a3fc-d9ba716c5f8d	8fe56507-de13-42a6-b2bc-3268699202a6	a509bf0fdc15fbec407730d6d28b65d05105d02af19eb0efb03e3748b6f43a7c	127.0.0.1	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36	0001-01-01 06:42:04+06:42:04	2025-08-11 10:52:46.472381+07	2025-08-11 14:48:37.474414+07
088a48db-1d66-4a08-8fc4-2b2adbf47aea	8fe56507-de13-42a6-b2bc-3268699202a6	d3f656f20c3f383c83f2ad3924ab2389b89ea85655231bb8f13118cfce91f463	127.0.0.1	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36	0001-01-01 06:42:04+06:42:04	2025-08-11 14:00:26.569843+07	2025-08-11 14:48:37.474414+07
5960ce1a-17c2-47e5-9a02-20dd35474675	8fe56507-de13-42a6-b2bc-3268699202a6	13950e2af8b79a1fc57f65d820c33c94cac161f7224db6f7c58d6eb5f7be033f	127.0.0.1	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36	0001-01-01 06:42:04+06:42:04	2025-08-11 15:38:29.429185+07	\N
facd7de4-c3e1-48fc-be2f-653ea967a3b2	550e8400-e29b-41d4-a716-446655440001	2e37cc766bc3223942a710419e405b7c14c1fe50812db1538586cb2d63b70950	127.0.0.1	curl/8.7.1	0001-01-01 06:42:04+06:42:04	2025-08-11 16:59:24.811751+07	\N
8ca161e8-61ce-43b1-86fc-601ed280c296	550e8400-e29b-41d4-a716-446655440001	2be0e534fa6bcb21d0924704026fdca121f0ae000d6a07a5e62279b0a0656ee9	127.0.0.1	curl/8.7.1	0001-01-01 06:42:04+06:42:04	2025-08-11 17:10:15.669709+07	\N
09968e53-82d6-4d13-8a32-7ef82e5e7ac1	8fe56507-de13-42a6-b2bc-3268699202a6	9e91218b8b55db448e47ab151f08a401dce3d552ef3200bffc27b5a38fe0fce4	127.0.0.1	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36	0001-01-01 06:42:04+06:42:04	2025-08-11 17:11:17.052073+07	\N
dabe9a33-6a60-4469-8570-0b5550c28aa3	8fe56507-de13-42a6-b2bc-3268699202a6	52e4750ff925ab44971310b01d7a56c17484866a571c369c2c2a52c66d69e05b	127.0.0.1	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36	0001-01-01 06:42:04+06:42:04	2025-08-11 17:13:13.500549+07	\N
e0965f87-4ac1-4337-8715-3e1b79d34e7a	8fe56507-de13-42a6-b2bc-3268699202a6	d8301bbb9ee0cbccf5d15057183244c8f92cb5486068fe602bea5e6549dc2f18	127.0.0.1	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36	0001-01-01 06:42:04+06:42:04	2025-08-11 17:14:01.991089+07	\N
19115bd0-03eb-444d-bfcf-9e9568b8997f	8fe56507-de13-42a6-b2bc-3268699202a6	79a886e6081e4fe36355449677f86efb4b13da9e79dbdc506b271905d8ddf747	127.0.0.1	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36	0001-01-01 06:42:04+06:42:04	2025-08-11 17:17:56.218309+07	\N
79962a3a-4a52-482e-87b5-6a7c3af0e176	8fe56507-de13-42a6-b2bc-3268699202a6	8e66496b4c7918fcb4563068f2ec0b23d772a548bcee4e5134fd8f2ebc562456	127.0.0.1	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36	0001-01-01 06:42:04+06:42:04	2025-08-11 17:18:14.497612+07	\N
caee13ed-ae6a-43bc-9e55-7ab4e0431cd5	8fe56507-de13-42a6-b2bc-3268699202a6	1a88a69d65e3d4e5c021c31a5e0b4e8544f004e099322c6c592d17626d75506b	127.0.0.1	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36	0001-01-01 06:42:04+06:42:04	2025-08-11 17:19:48.19583+07	\N
f272bfd2-1163-4ee0-bada-8efb6cc1693f	550e8400-e29b-41d4-a716-446655440001	9d09ddfc2f4b72c0a949de523c8e8f30dd0cb44a0816703ff09ec744d234f362	127.0.0.1	curl/8.7.1	0001-01-01 06:42:04+06:42:04	2025-08-11 17:27:25.790062+07	\N
1b10378e-3f8c-42ef-8bfa-888f70250411	550e8400-e29b-41d4-a716-446655440001	dd9073de4e9847cc98e4fd2f818bc1331e161b27ec8d1be1324ec440b7b224e7	127.0.0.1	curl/8.7.1	0001-01-01 06:42:04+06:42:04	2025-08-11 17:28:46.629399+07	\N
fdc611aa-fcc2-48d4-903a-1f0cbdda5d85	8fe56507-de13-42a6-b2bc-3268699202a6	95592707f54ac9aa033e72452a91ed19dfd9084cad0c513ca11cae21b4f7f527	127.0.0.1	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36	0001-01-01 06:42:04+06:42:04	2025-08-11 17:33:31.44334+07	\N
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: nurse_shift; Owner: nuttapong2
--

COPY nurse_shift.users (id, email, password_hash, first_name, last_name, phone, role, status, "position", days_remaining, subscription_expires_at, package_type, max_departments, avatar_url, settings, last_login_at, created_at, updated_at) FROM stdin;
550e8400-e29b-41d4-a716-446655440002	user@nurseshift.com	$2a$10$5G5QXW39SmBuPa9UF0Mft.rAXKMaLi0VeBuvesgjwFslILoE7ej.C	พยาบาล	ทดสอบ	+66-82-345-6789	user	active	หัวหน้าพยาบาล	85	2025-09-10 10:17:49.136878+07	trial	2	\N	{"theme": "light", "language": "th"}	\N	2025-08-11 10:17:49.136878+07	2025-08-11 10:17:49.136878+07
550e8400-e29b-41d4-a716-446655440003	test@nurseshift.com	$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi	พยาบาล	ทดสอบ2	+66-83-456-7890	user	active	พยาบาลประจำแผนก	60	2025-10-10 10:17:49.136878+07	standard	5	\N	{"theme": "light", "language": "th"}	\N	2025-08-11 10:17:49.136878+07	2025-08-11 10:17:49.136878+07
550e8400-e29b-41d4-a716-446655440001	admin@nurseshift.com	$2a$10$OQJnRxKT4dwQD1blpI2lze9/1NPb4XVl.V8Hle6a3p1p7CsIC7I4m	ผู้ดูแล	ระบบ	0812345678	admin	active	System Administrator	90	2025-11-09 10:17:49.136878+07	enterprise	20	https://example.com/avatar.jpg	{"theme": "dark", "language": "th"}	2025-08-11 16:58:46.631398+07	2025-08-11 10:17:49.136878+07	2025-08-11 16:58:46.631662+07
8fe56507-de13-42a6-b2bc-3268699202a6	worknuttapong1@gmail.com	$2a$12$09HRIWjtXplsm1bigFmToes4tWhW0blQm24YHhhzLicsYPHDWxl8O	Nuttapong	Silwuti	\N	user	active	หัวหน้าแผนก	30	\N	trial	2	\N	{}	2025-08-11 17:03:31.447959+07	2025-08-11 10:21:59.754391+07	2025-08-11 17:03:31.448229+07
\.


--
-- Data for Name: working_days; Type: TABLE DATA; Schema: nurse_shift; Owner: nuttapong2
--

COPY nurse_shift.working_days (id, department_id, day_of_week, is_working_day, created_at) FROM stdin;
990e8400-e29b-41d4-a716-446655440001	660e8400-e29b-41d4-a716-446655440001	0	t	2025-08-11 10:17:49.148884+07
990e8400-e29b-41d4-a716-446655440002	660e8400-e29b-41d4-a716-446655440001	1	t	2025-08-11 10:17:49.148884+07
990e8400-e29b-41d4-a716-446655440003	660e8400-e29b-41d4-a716-446655440001	2	t	2025-08-11 10:17:49.148884+07
990e8400-e29b-41d4-a716-446655440004	660e8400-e29b-41d4-a716-446655440001	3	t	2025-08-11 10:17:49.148884+07
990e8400-e29b-41d4-a716-446655440005	660e8400-e29b-41d4-a716-446655440001	4	t	2025-08-11 10:17:49.148884+07
990e8400-e29b-41d4-a716-446655440006	660e8400-e29b-41d4-a716-446655440001	5	t	2025-08-11 10:17:49.148884+07
990e8400-e29b-41d4-a716-446655440007	660e8400-e29b-41d4-a716-446655440001	6	t	2025-08-11 10:17:49.148884+07
\.


--
-- Name: audit_logs audit_logs_pkey; Type: CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.audit_logs
    ADD CONSTRAINT audit_logs_pkey PRIMARY KEY (id);


--
-- Name: department_staff department_staff_pkey; Type: CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.department_staff
    ADD CONSTRAINT department_staff_pkey PRIMARY KEY (id);


--
-- Name: departments departments_pkey; Type: CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.departments
    ADD CONSTRAINT departments_pkey PRIMARY KEY (id);


--
-- Name: departments departments_user_id_name_key; Type: CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.departments
    ADD CONSTRAINT departments_user_id_name_key UNIQUE (user_id, name);


--
-- Name: holidays holidays_pkey; Type: CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.holidays
    ADD CONSTRAINT holidays_pkey PRIMARY KEY (id);


--
-- Name: leave_requests leave_requests_pkey; Type: CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.leave_requests
    ADD CONSTRAINT leave_requests_pkey PRIMARY KEY (id);


--
-- Name: notifications notifications_pkey; Type: CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.notifications
    ADD CONSTRAINT notifications_pkey PRIMARY KEY (id);


--
-- Name: packages packages_name_key; Type: CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.packages
    ADD CONSTRAINT packages_name_key UNIQUE (name);


--
-- Name: packages packages_pkey; Type: CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.packages
    ADD CONSTRAINT packages_pkey PRIMARY KEY (id);


--
-- Name: payments payments_pkey; Type: CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.payments
    ADD CONSTRAINT payments_pkey PRIMARY KEY (id);


--
-- Name: schedules schedules_pkey; Type: CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.schedules
    ADD CONSTRAINT schedules_pkey PRIMARY KEY (id);


--
-- Name: schedules schedules_staff_id_schedule_date_shift_id_key; Type: CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.schedules
    ADD CONSTRAINT schedules_staff_id_schedule_date_shift_id_key UNIQUE (staff_id, schedule_date, shift_id);


--
-- Name: scheduling_priorities scheduling_priorities_pkey; Type: CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.scheduling_priorities
    ADD CONSTRAINT scheduling_priorities_pkey PRIMARY KEY (id);


--
-- Name: scheduling_priorities scheduling_priorities_user_id_name_key; Type: CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.scheduling_priorities
    ADD CONSTRAINT scheduling_priorities_user_id_name_key UNIQUE (user_id, name);


--
-- Name: scheduling_priorities scheduling_priorities_user_id_priority_order_key; Type: CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.scheduling_priorities
    ADD CONSTRAINT scheduling_priorities_user_id_priority_order_key UNIQUE (user_id, priority_order);


--
-- Name: shifts shifts_pkey; Type: CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.shifts
    ADD CONSTRAINT shifts_pkey PRIMARY KEY (id);


--
-- Name: user_sessions user_sessions_pkey; Type: CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.user_sessions
    ADD CONSTRAINT user_sessions_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: working_days working_days_department_id_day_of_week_key; Type: CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.working_days
    ADD CONSTRAINT working_days_department_id_day_of_week_key UNIQUE (department_id, day_of_week);


--
-- Name: working_days working_days_pkey; Type: CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.working_days
    ADD CONSTRAINT working_days_pkey PRIMARY KEY (id);


--
-- Name: idx_department_staff_department_id; Type: INDEX; Schema: nurse_shift; Owner: nuttapong2
--

CREATE INDEX idx_department_staff_department_id ON nurse_shift.department_staff USING btree (department_id);


--
-- Name: idx_departments_user_id; Type: INDEX; Schema: nurse_shift; Owner: nuttapong2
--

CREATE INDEX idx_departments_user_id ON nurse_shift.departments USING btree (user_id);


--
-- Name: idx_leave_requests_dates; Type: INDEX; Schema: nurse_shift; Owner: nuttapong2
--

CREATE INDEX idx_leave_requests_dates ON nurse_shift.leave_requests USING btree (start_date, end_date);


--
-- Name: idx_leave_requests_department_id; Type: INDEX; Schema: nurse_shift; Owner: nuttapong2
--

CREATE INDEX idx_leave_requests_department_id ON nurse_shift.leave_requests USING btree (department_id);


--
-- Name: idx_leave_requests_requested_by; Type: INDEX; Schema: nurse_shift; Owner: nuttapong2
--

CREATE INDEX idx_leave_requests_requested_by ON nurse_shift.leave_requests USING btree (requested_by);


--
-- Name: idx_leave_requests_staff_id; Type: INDEX; Schema: nurse_shift; Owner: nuttapong2
--

CREATE INDEX idx_leave_requests_staff_id ON nurse_shift.leave_requests USING btree (staff_id);


--
-- Name: idx_leave_requests_status; Type: INDEX; Schema: nurse_shift; Owner: nuttapong2
--

CREATE INDEX idx_leave_requests_status ON nurse_shift.leave_requests USING btree (status);


--
-- Name: idx_notifications_created_at; Type: INDEX; Schema: nurse_shift; Owner: nuttapong2
--

CREATE INDEX idx_notifications_created_at ON nurse_shift.notifications USING btree (created_at);


--
-- Name: idx_notifications_is_read; Type: INDEX; Schema: nurse_shift; Owner: nuttapong2
--

CREATE INDEX idx_notifications_is_read ON nurse_shift.notifications USING btree (is_read);


--
-- Name: idx_notifications_user_id; Type: INDEX; Schema: nurse_shift; Owner: nuttapong2
--

CREATE INDEX idx_notifications_user_id ON nurse_shift.notifications USING btree (user_id);


--
-- Name: idx_payments_status; Type: INDEX; Schema: nurse_shift; Owner: nuttapong2
--

CREATE INDEX idx_payments_status ON nurse_shift.payments USING btree (status);


--
-- Name: idx_payments_user_id; Type: INDEX; Schema: nurse_shift; Owner: nuttapong2
--

CREATE INDEX idx_payments_user_id ON nurse_shift.payments USING btree (user_id);


--
-- Name: idx_schedules_date; Type: INDEX; Schema: nurse_shift; Owner: nuttapong2
--

CREATE INDEX idx_schedules_date ON nurse_shift.schedules USING btree (schedule_date);


--
-- Name: idx_schedules_department_id; Type: INDEX; Schema: nurse_shift; Owner: nuttapong2
--

CREATE INDEX idx_schedules_department_id ON nurse_shift.schedules USING btree (department_id);


--
-- Name: idx_schedules_shift_id; Type: INDEX; Schema: nurse_shift; Owner: nuttapong2
--

CREATE INDEX idx_schedules_shift_id ON nurse_shift.schedules USING btree (shift_id);


--
-- Name: idx_schedules_staff_id; Type: INDEX; Schema: nurse_shift; Owner: nuttapong2
--

CREATE INDEX idx_schedules_staff_id ON nurse_shift.schedules USING btree (staff_id);


--
-- Name: idx_user_sessions_expires_at; Type: INDEX; Schema: nurse_shift; Owner: nuttapong2
--

CREATE INDEX idx_user_sessions_expires_at ON nurse_shift.user_sessions USING btree (expires_at);


--
-- Name: idx_user_sessions_token_hash; Type: INDEX; Schema: nurse_shift; Owner: nuttapong2
--

CREATE INDEX idx_user_sessions_token_hash ON nurse_shift.user_sessions USING btree (token_hash);


--
-- Name: idx_user_sessions_user_id; Type: INDEX; Schema: nurse_shift; Owner: nuttapong2
--

CREATE INDEX idx_user_sessions_user_id ON nurse_shift.user_sessions USING btree (user_id);


--
-- Name: idx_users_email; Type: INDEX; Schema: nurse_shift; Owner: nuttapong2
--

CREATE INDEX idx_users_email ON nurse_shift.users USING btree (email);


--
-- Name: idx_users_role; Type: INDEX; Schema: nurse_shift; Owner: nuttapong2
--

CREATE INDEX idx_users_role ON nurse_shift.users USING btree (role);


--
-- Name: idx_users_status; Type: INDEX; Schema: nurse_shift; Owner: nuttapong2
--

CREATE INDEX idx_users_status ON nurse_shift.users USING btree (status);


--
-- Name: departments update_departments_updated_at; Type: TRIGGER; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TRIGGER update_departments_updated_at BEFORE UPDATE ON nurse_shift.departments FOR EACH ROW EXECUTE FUNCTION nurse_shift.update_updated_at_column();


--
-- Name: holidays update_holidays_updated_at; Type: TRIGGER; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TRIGGER update_holidays_updated_at BEFORE UPDATE ON nurse_shift.holidays FOR EACH ROW EXECUTE FUNCTION nurse_shift.update_updated_at_column();


--
-- Name: leave_requests update_leave_requests_updated_at; Type: TRIGGER; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TRIGGER update_leave_requests_updated_at BEFORE UPDATE ON nurse_shift.leave_requests FOR EACH ROW EXECUTE FUNCTION nurse_shift.update_updated_at_column();


--
-- Name: packages update_packages_updated_at; Type: TRIGGER; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TRIGGER update_packages_updated_at BEFORE UPDATE ON nurse_shift.packages FOR EACH ROW EXECUTE FUNCTION nurse_shift.update_updated_at_column();


--
-- Name: payments update_payments_updated_at; Type: TRIGGER; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TRIGGER update_payments_updated_at BEFORE UPDATE ON nurse_shift.payments FOR EACH ROW EXECUTE FUNCTION nurse_shift.update_updated_at_column();


--
-- Name: schedules update_schedules_updated_at; Type: TRIGGER; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TRIGGER update_schedules_updated_at BEFORE UPDATE ON nurse_shift.schedules FOR EACH ROW EXECUTE FUNCTION nurse_shift.update_updated_at_column();


--
-- Name: scheduling_priorities update_scheduling_priorities_updated_at; Type: TRIGGER; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TRIGGER update_scheduling_priorities_updated_at BEFORE UPDATE ON nurse_shift.scheduling_priorities FOR EACH ROW EXECUTE FUNCTION nurse_shift.update_updated_at_column();


--
-- Name: shifts update_shifts_updated_at; Type: TRIGGER; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TRIGGER update_shifts_updated_at BEFORE UPDATE ON nurse_shift.shifts FOR EACH ROW EXECUTE FUNCTION nurse_shift.update_updated_at_column();


--
-- Name: users update_users_updated_at; Type: TRIGGER; Schema: nurse_shift; Owner: nuttapong2
--

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON nurse_shift.users FOR EACH ROW EXECUTE FUNCTION nurse_shift.update_updated_at_column();


--
-- Name: audit_logs audit_logs_user_id_fkey; Type: FK CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.audit_logs
    ADD CONSTRAINT audit_logs_user_id_fkey FOREIGN KEY (user_id) REFERENCES nurse_shift.users(id);


--
-- Name: department_staff department_staff_department_id_fkey; Type: FK CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.department_staff
    ADD CONSTRAINT department_staff_department_id_fkey FOREIGN KEY (department_id) REFERENCES nurse_shift.departments(id) ON DELETE CASCADE;


--
-- Name: departments departments_user_id_fkey; Type: FK CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.departments
    ADD CONSTRAINT departments_user_id_fkey FOREIGN KEY (user_id) REFERENCES nurse_shift.users(id) ON DELETE CASCADE;


--
-- Name: holidays holidays_department_id_fkey; Type: FK CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.holidays
    ADD CONSTRAINT holidays_department_id_fkey FOREIGN KEY (department_id) REFERENCES nurse_shift.departments(id) ON DELETE CASCADE;


--
-- Name: leave_requests leave_requests_approved_by_fkey; Type: FK CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.leave_requests
    ADD CONSTRAINT leave_requests_approved_by_fkey FOREIGN KEY (approved_by) REFERENCES nurse_shift.users(id);


--
-- Name: leave_requests leave_requests_department_id_fkey; Type: FK CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.leave_requests
    ADD CONSTRAINT leave_requests_department_id_fkey FOREIGN KEY (department_id) REFERENCES nurse_shift.departments(id) ON DELETE CASCADE;


--
-- Name: leave_requests leave_requests_requested_by_fkey; Type: FK CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.leave_requests
    ADD CONSTRAINT leave_requests_requested_by_fkey FOREIGN KEY (requested_by) REFERENCES nurse_shift.users(id);


--
-- Name: leave_requests leave_requests_staff_id_fkey; Type: FK CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.leave_requests
    ADD CONSTRAINT leave_requests_staff_id_fkey FOREIGN KEY (staff_id) REFERENCES nurse_shift.department_staff(id) ON DELETE CASCADE;


--
-- Name: notifications notifications_user_id_fkey; Type: FK CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.notifications
    ADD CONSTRAINT notifications_user_id_fkey FOREIGN KEY (user_id) REFERENCES nurse_shift.users(id) ON DELETE CASCADE;


--
-- Name: payments payments_approved_by_fkey; Type: FK CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.payments
    ADD CONSTRAINT payments_approved_by_fkey FOREIGN KEY (approved_by) REFERENCES nurse_shift.users(id);


--
-- Name: payments payments_package_id_fkey; Type: FK CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.payments
    ADD CONSTRAINT payments_package_id_fkey FOREIGN KEY (package_id) REFERENCES nurse_shift.packages(id);


--
-- Name: payments payments_user_id_fkey; Type: FK CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.payments
    ADD CONSTRAINT payments_user_id_fkey FOREIGN KEY (user_id) REFERENCES nurse_shift.users(id) ON DELETE CASCADE;


--
-- Name: schedules schedules_assigned_by_fkey; Type: FK CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.schedules
    ADD CONSTRAINT schedules_assigned_by_fkey FOREIGN KEY (assigned_by) REFERENCES nurse_shift.users(id);


--
-- Name: schedules schedules_department_id_fkey; Type: FK CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.schedules
    ADD CONSTRAINT schedules_department_id_fkey FOREIGN KEY (department_id) REFERENCES nurse_shift.departments(id) ON DELETE CASCADE;


--
-- Name: schedules schedules_shift_id_fkey; Type: FK CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.schedules
    ADD CONSTRAINT schedules_shift_id_fkey FOREIGN KEY (shift_id) REFERENCES nurse_shift.shifts(id) ON DELETE CASCADE;


--
-- Name: schedules schedules_staff_id_fkey; Type: FK CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.schedules
    ADD CONSTRAINT schedules_staff_id_fkey FOREIGN KEY (staff_id) REFERENCES nurse_shift.department_staff(id) ON DELETE CASCADE;


--
-- Name: scheduling_priorities scheduling_priorities_user_id_fkey; Type: FK CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.scheduling_priorities
    ADD CONSTRAINT scheduling_priorities_user_id_fkey FOREIGN KEY (user_id) REFERENCES nurse_shift.users(id) ON DELETE CASCADE;


--
-- Name: shifts shifts_department_id_fkey; Type: FK CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.shifts
    ADD CONSTRAINT shifts_department_id_fkey FOREIGN KEY (department_id) REFERENCES nurse_shift.departments(id) ON DELETE CASCADE;


--
-- Name: user_sessions user_sessions_user_id_fkey; Type: FK CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.user_sessions
    ADD CONSTRAINT user_sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES nurse_shift.users(id) ON DELETE CASCADE;


--
-- Name: working_days working_days_department_id_fkey; Type: FK CONSTRAINT; Schema: nurse_shift; Owner: nuttapong2
--

ALTER TABLE ONLY nurse_shift.working_days
    ADD CONSTRAINT working_days_department_id_fkey FOREIGN KEY (department_id) REFERENCES nurse_shift.departments(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

