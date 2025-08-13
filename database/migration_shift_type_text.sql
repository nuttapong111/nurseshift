-- Migrate shifts.type from enum shift_type to TEXT to allow free-text values
BEGIN;

-- Only proceed if column exists
ALTER TABLE nurse_shift.shifts
    ALTER COLUMN type TYPE TEXT USING type::text;

-- Optional: drop enum if no longer used anywhere else (safe if not referenced)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type t JOIN pg_namespace n ON n.oid=t.typnamespace 
               WHERE t.typname='shift_type' AND n.nspname='public') THEN
        -- try to drop; ignore if in use
        BEGIN
            DROP TYPE public.shift_type;
        EXCEPTION WHEN dependent_objects_still_exist THEN
            -- leave the enum if still referenced elsewhere
            NULL;
        END;
    END IF;
END$$;

COMMIT;

