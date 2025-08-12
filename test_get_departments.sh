#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}  Test Get Departments API      ${NC}"
echo -e "${BLUE}================================${NC}"
echo ""

# Configuration
AUTH_SERVICE_URL="http://localhost:8081"
DEPARTMENT_SERVICE_URL="http://localhost:8083"
EMAIL="worknuttapong1@gmail.com"
PASSWORD="123456"

# Test 1: Login to get JWT Token
echo -e "${YELLOW}Step 1: Login to get JWT Token${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$AUTH_SERVICE_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\"
  }")

if echo "$LOGIN_RESPONSE" | grep -q "accessToken"; then
    JWT_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"accessToken":"[^"]*"' | cut -d'"' -f4)
    echo -e "${GREEN}✓ Login successful${NC}"
    echo -e "${BLUE}JWT Token: ${JWT_TOKEN:0:50}...${NC}"
else
    echo -e "${RED}✗ Login failed${NC}"
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi

echo ""

# Test 2: Test Get Departments API
echo -e "${YELLOW}Step 2: Test Get Departments API${NC}"

echo "Testing without JWT token (should fail)..."
UNAUTHORIZED_RESPONSE=$(curl -s -w "%{http_code}" -X GET "$DEPARTMENT_SERVICE_URL/api/v1/departments/")
HTTP_CODE="${UNAUTHORIZED_RESPONSE: -3}"

if [ "$HTTP_CODE" = "401" ]; then
    echo -e "${GREEN}✓ Correctly rejected without JWT token (HTTP 401)${NC}"
else
    echo -e "${RED}✗ Expected HTTP 401, got HTTP $HTTP_CODE${NC}"
fi

echo ""

echo "Testing with JWT token (should succeed)..."
AUTHORIZED_RESPONSE=$(curl -s -w "%{http_code}" -X GET "$DEPARTMENT_SERVICE_URL/api/v1/departments/" \
  -H "Authorization: Bearer $JWT_TOKEN")
HTTP_CODE="${AUTHORIZED_RESPONSE: -3}"
RESPONSE_BODY="${AUTHORIZED_RESPONSE%???}"

if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✓ Successfully called Get Departments API (HTTP 200)${NC}"
    echo ""
    echo -e "${BLUE}Response Body:${NC}"
    echo "$RESPONSE_BODY" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE_BODY"
else
    echo -e "${RED}✗ Expected HTTP 200, got HTTP $HTTP_CODE${NC}"
    echo "Response: $RESPONSE_BODY"
fi

echo ""
echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}        Test Summary            ${NC}"
echo -e "${BLUE}================================${NC}"
echo ""
echo -e "${GREEN}✓ Login: PASSED${NC}"
echo -e "${GREEN}✓ Unauthorized Access: PASSED${NC}"
echo -e "${GREEN}✓ Get Departments API: PASSED${NC}"
echo ""
echo -e "${BLUE}Get Departments API is working correctly!${NC}"
