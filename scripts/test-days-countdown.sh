#!/bin/bash

# NurseShift Days Countdown Test Script
# This script tests the days remaining countdown system

PROJECT_ROOT="/Volumes/NO NAME/project/nurseshift_final"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧪 Testing NurseShift Days Countdown System...${NC}"
echo "=================================="

# Test user email
TEST_EMAIL="worknuttapong1@gmail.com"

echo -e "${BLUE}📊 Current user status:${NC}"
psql -d nurseshift -c "
SELECT 
    email,
    days_remaining,
    subscription_expires_at,
    created_at,
    package_type,
    status
FROM nurse_shift.users 
WHERE email = '$TEST_EMAIL';
"

echo -e "\n${BLUE}🔄 Running days remaining update function...${NC}"
psql -d nurseshift -c "SELECT nurse_shift.update_user_days_remaining();"

echo -e "\n${BLUE}📊 Updated user status:${NC}"
psql -d nurseshift -c "
SELECT 
    email,
    days_remaining,
    subscription_expires_at,
    created_at,
    package_type,
    status
FROM nurse_shift.users 
WHERE email = '$TEST_EMAIL';
"

echo -e "\n${BLUE}📈 All users subscription status:${NC}"
psql -d nurseshift -c "
SELECT 
    email,
    days_remaining,
    subscription_expires_at,
    package_type,
    status,
    CASE 
        WHEN days_remaining > 7 THEN 'safe'
        WHEN days_remaining > 0 THEN 'warning'
        ELSE 'expired'
    END as days_status
FROM nurse_shift.users 
WHERE status IN ('active', 'inactive', 'pending', 'suspended')
ORDER BY days_remaining ASC;
"

echo -e "\n${GREEN}✅ Days countdown test completed!${NC}"
echo -e "${BLUE}💡 The system should now show the correct remaining days${NC}"
echo -e "${BLUE}💡 Cron service will update this automatically every day at midnight UTC${NC}"
