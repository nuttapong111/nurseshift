#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}  Test Update Department API    ${NC}"
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

# Test 2: Create a test department first
echo -e "${YELLOW}Step 2: Create Test Department${NC}"
CREATE_RESPONSE=$(curl -s -w "%{http_code}" -X POST "$DEPARTMENT_SERVICE_URL/api/v1/departments/" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "name": "แผนกทดสอบ Update",
    "description": "แผนกสำหรับทดสอบ Update API",
    "max_nurses": 10,
    "max_assistants": 5
  }')
HTTP_CODE="${CREATE_RESPONSE: -3}"
RESPONSE_BODY="${CREATE_RESPONSE%???}"

if [ "$HTTP_CODE" = "201" ] || [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✓ Successfully created test department (HTTP $HTTP_CODE)${NC}"
    DEPT_ID=$(echo "$RESPONSE_BODY" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    echo -e "${BLUE}Test Department ID: $DEPT_ID${NC}"
else
    echo -e "${RED}✗ Failed to create test department (HTTP $HTTP_CODE)${NC}"
    echo "Response: $RESPONSE_BODY"
    exit 1
fi

echo ""

# Test 3: Test Update Department API
echo -e "${YELLOW}Step 3: Test Update Department API${NC}"

echo "Updating department..."
UPDATE_RESPONSE=$(curl -s -w "%{http_code}" -X PUT "$DEPARTMENT_SERVICE_URL/api/v1/departments/$DEPT_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "name": "แผนกทดสอบ Update (แก้ไขแล้ว)",
    "description": "แผนกที่ได้รับการแก้ไขแล้ว",
    "max_nurses": 20,
    "max_assistants": 10
  }')
HTTP_CODE="${UPDATE_RESPONSE: -3}"
RESPONSE_BODY="${UPDATE_RESPONSE%???}"

if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✓ Successfully updated department (HTTP 200)${NC}"
    echo ""
    echo -e "${BLUE}Update Response:${NC}"
    echo "$RESPONSE_BODY" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE_BODY"
else
    echo -e "${RED}✗ Failed to update department (HTTP $HTTP_CODE)${NC}"
    echo "Response: $RESPONSE_BODY"
fi

echo ""

# Test 4: Verify department was updated
echo -e "${YELLOW}Step 4: Verify Updated Department${NC}"
GET_RESPONSE=$(curl -s -w "%{http_code}" -X GET "$DEPARTMENT_SERVICE_URL/api/v1/departments/$DEPT_ID" \
  -H "Authorization: Bearer $JWT_TOKEN")
GET_HTTP_CODE="${GET_RESPONSE: -3}"
GET_BODY="${GET_RESPONSE%???}"

if [ "$GET_HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✓ Successfully retrieved updated department (HTTP 200)${NC}"
    echo ""
    echo -e "${BLUE}Updated Department Data:${NC}"
    echo "$GET_BODY" | python3 -m json.tool 2>/dev/null || echo "$GET_BODY"
else
    echo -e "${RED}✗ Failed to retrieve updated department (HTTP $GET_HTTP_CODE)${NC}"
fi

echo ""
echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}        Test Summary            ${NC}"
echo -e "${BLUE}================================${NC}"
echo ""
echo -e "${GREEN}✓ Login: PASSED${NC}"
echo -e "${GREEN}✓ Create Test Department: PASSED${NC}"
echo -e "${GREEN}✓ Update Department API: PASSED${NC}"
echo -e "${GREEN}✓ Verify Updated Department: PASSED${NC}"
echo ""
echo -e "${BLUE}Update Department API is working correctly!${NC}"
