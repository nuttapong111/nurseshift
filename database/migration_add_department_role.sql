-- Migration Script: Add Department Role Support
-- Version: 1.1.0
-- Date: 2025-08-12
-- Description: Add department_role enum and column to department_users table
-- Safe for production - includes rollback support

-- ===================================
-- FORWARD MIGRATION
-- ===================================

-- Step 1: Create department_role enum if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'department_role') THEN
        CREATE TYPE nurse_shift.department_role AS ENUM ('nurse', 'assistant');
        RAISE NOTICE 'Created department_role enum';
    ELSE
        RAISE NOTICE 'department_role enum already exists';
    END IF;
END $$;

-- Step 2: Add department_role column if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'nurse_shift' 
        AND table_name = 'department_users' 
        AND column_name = 'department_role'
    ) THEN
        ALTER TABLE nurse_shift.department_users 
        ADD COLUMN department_role nurse_shift.department_role NOT NULL DEFAULT 'nurse';
        RAISE NOTICE 'Added department_role column';
    ELSE
        RAISE NOTICE 'department_role column already exists';
    END IF;
END $$;

-- Step 3: Add index for department_role if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes 
        WHERE schemaname = 'nurse_shift' 
        AND tablename = 'department_users' 
        AND indexname = 'idx_department_users_department_role'
    ) THEN
        CREATE INDEX idx_department_users_department_role ON nurse_shift.department_users(department_role);
        RAISE NOTICE 'Added department_role index';
    ELSE
        RAISE NOTICE 'department_role index already exists';
    END IF;
END $$;

-- Step 4: Add comment if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_description 
        WHERE objoid = (
            SELECT oid FROM pg_class 
            WHERE relname = 'department_users' 
            AND relnamespace = (SELECT oid FROM pg_namespace WHERE nspname = 'nurse_shift')
        ) 
        AND objsubid = (
            SELECT attnum FROM pg_attribute 
            WHERE attname = 'department_role' 
            AND attrelid = (
                SELECT oid FROM pg_class 
                WHERE relname = 'department_users' 
                AND relnamespace = (SELECT oid FROM pg_namespace WHERE nspname = 'nurse_shift')
            )
        )
    ) THEN
        COMMENT ON COLUMN nurse_shift.department_users.department_role IS 'Role ของผู้ใช้ในแผนก: nurse (พยาบาล) หรือ assistant (ผู้ช่วยพยาบาล)';
        RAISE NOTICE 'Added department_role comment';
    ELSE
        RAISE NOTICE 'department_role comment already exists';
    END IF;
END $$;

-- ===================================
-- VERIFICATION
-- ===================================

-- Verify the migration was successful
SELECT 
    'Migration Status' as check_type,
    CASE 
        WHEN EXISTS (SELECT 1 FROM pg_type WHERE typname = 'department_role') 
        THEN 'PASS' ELSE 'FAIL' 
    END as enum_exists,
    CASE 
        WHEN EXISTS (
            SELECT 1 FROM information_schema.columns 
            WHERE table_schema = 'nurse_shift' 
            AND table_name = 'department_users' 
            AND column_name = 'department_role'
        ) 
        THEN 'PASS' ELSE 'FAIL' 
    END as column_exists,
    CASE 
        WHEN EXISTS (
            SELECT 1 FROM pg_indexes 
            WHERE schemaname = 'nurse_shift' 
            AND tablename = 'department_users' 
            AND indexname = 'idx_department_users_department_role'
        ) 
        THEN 'PASS' ELSE 'FAIL' 
    END as index_exists;

-- ===================================
-- ROLLBACK MIGRATION (if needed)
-- ===================================

/*
-- To rollback this migration, uncomment and run the following:

-- Remove index
DROP INDEX IF EXISTS nurse_shift.idx_department_users_department_role;

-- Remove column
ALTER TABLE nurse_shift.department_users DROP COLUMN IF EXISTS department_role;

-- Remove enum (only if no other tables use it)
-- DROP TYPE IF EXISTS nurse_shift.department_role;

-- Note: Be careful with DROP TYPE as it might be used by other tables
*/
