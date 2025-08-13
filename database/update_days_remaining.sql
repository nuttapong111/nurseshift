-- Function to update user days remaining
-- This function calculates and updates the remaining days for all active users

CREATE OR REPLACE FUNCTION nurse_shift.update_user_days_remaining()
RETURNS INTEGER AS $$
DECLARE
    updated_count INTEGER := 0;
    expired_count INTEGER := 0;
BEGIN
    -- Update days remaining for users with subscription_expires_at
    UPDATE nurse_shift.users 
    SET 
        days_remaining = GREATEST(0, 
            EXTRACT(EPOCH FROM (subscription_expires_at - NOW())) / 86400
        ),
        updated_at = NOW()
    WHERE 
        subscription_expires_at IS NOT NULL 
        AND subscription_expires_at > NOW()
        AND status = 'active';
    
    GET DIAGNOSTICS updated_count = ROW_COUNT;
    
    -- Suspend expired users
    UPDATE nurse_shift.users 
    SET 
        status = 'suspended',
        days_remaining = 0,
        updated_at = NOW()
    WHERE 
        subscription_expires_at IS NOT NULL 
        AND subscription_expires_at <= NOW()
        AND status = 'active';
    
    GET DIAGNOSTICS expired_count = ROW_COUNT;
    
    -- Log the update
    INSERT INTO nurse_shift.audit_logs (
        user_id, action, resource_type, resource_id, 
        old_data, new_data, ip_address, user_agent, created_at
    ) VALUES (
        NULL, 'UPDATE_DAYS_REMAINING', 'SYSTEM', NULL,
        jsonb_build_object('updated_users', updated_count, 'expired_users', expired_count),
        jsonb_build_object('timestamp', NOW()),
        '127.0.0.1', 'cron-service', NOW()
    );
    
    RETURN updated_count + expired_count;
END;
$$ LANGUAGE plpgsql;

-- Function to get user subscription status
CREATE OR REPLACE FUNCTION nurse_shift.get_user_subscription_status(user_email TEXT)
RETURNS TABLE(
    email TEXT,
    days_remaining INTEGER,
    subscription_expires_at TIMESTAMP WITH TIME ZONE,
    package_type TEXT,
    status TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        u.email::TEXT,
        u.days_remaining,
        u.subscription_expires_at,
        u.package_type::TEXT,
        u.status::TEXT
    FROM nurse_shift.users u
    WHERE u.email = user_email;
END;
$$ LANGUAGE plpgsql;

-- Function to manually update a specific user's days remaining
CREATE OR REPLACE FUNCTION nurse_shift.update_specific_user_days(user_email TEXT)
RETURNS INTEGER AS $$
DECLARE
    updated_count INTEGER := 0;
BEGIN
    UPDATE nurse_shift.users 
    SET 
        days_remaining = GREATEST(0, 
            EXTRACT(EPOCH FROM (subscription_expires_at - NOW())) / 86400
        ),
        updated_at = NOW()
    WHERE 
        email = user_email
        AND subscription_expires_at IS NOT NULL 
        AND subscription_expires_at > NOW()
        AND status = 'active';
    
    GET DIAGNOSTICS updated_count = ROW_COUNT;
    
    RETURN updated_count;
END;
$$ LANGUAGE plpgsql;

-- Create a view for monitoring user subscription status
CREATE OR REPLACE VIEW nurse_shift.user_subscription_status AS
SELECT 
    u.id,
    u.email,
    u.first_name,
    u.last_name,
    u.days_remaining,
    u.subscription_expires_at,
    u.package_type,
    u.status,
    u.created_at,
    u.updated_at,
    CASE 
        WHEN u.subscription_expires_at IS NULL THEN 'unlimited'
        WHEN u.subscription_expires_at > NOW() THEN 'active'
        ELSE 'expired'
    END as subscription_status,
    CASE 
        WHEN u.days_remaining > 7 THEN 'safe'
        WHEN u.days_remaining > 0 THEN 'warning'
        ELSE 'expired'
    END as days_status
FROM nurse_shift.users u
    WHERE u.status IN ('active', 'inactive', 'pending', 'suspended')
ORDER BY u.days_remaining ASC, u.subscription_expires_at ASC;

-- Grant permissions
GRANT EXECUTE ON FUNCTION nurse_shift.update_user_days_remaining() TO nuttapong2;
GRANT EXECUTE ON FUNCTION nurse_shift.get_user_subscription_status(TEXT) TO nuttapong2;
GRANT EXECUTE ON FUNCTION nurse_shift.update_specific_user_days(TEXT) TO nuttapong2;
GRANT SELECT ON nurse_shift.user_subscription_status TO nuttapong2;
