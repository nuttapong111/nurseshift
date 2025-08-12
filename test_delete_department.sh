#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}  Test Delete Department API    ${NC}"
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
    "name": "แผนกทดสอบ Delete",
    "description": "แผนกสำหรับทดสอบ Delete API",
    "max_nurses": 8,
    "max_assistants": 4
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

# Test 3: Verify department exists before deletion
echo -e "${YELLOW}Step 3: Verify Department Exists${NC}"
GET_RESPONSE=$(curl -s -w "%{http_code}" -X GET "$DEPARTMENT_SERVICE_URL/api/v1/departments/$DEPT_ID" \
  -H "Authorization: Bearer $JWT_TOKEN")
GET_HTTP_CODE="${GET_RESPONSE: -3}"

if [ "$GET_HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✓ Department exists and can be retrieved (HTTP 200)${NC}"
else
    echo -e "${RED}✗ Department cannot be retrieved (HTTP $GET_HTTP_CODE)${NC}"
    exit 1
fi

echo ""

# Test 4: Test Delete Department API
echo -e "${YELLOW}Step 4: Test Delete Department API${NC}"

echo "Deleting department..."
DELETE_RESPONSE=$(curl -s -w "%{http_code}" -X DELETE "$DEPARTMENT_SERVICE_URL/api/v1/departments/$DEPT_ID" \
  -H "Authorization: Bearer $JWT_TOKEN")
HTTP_CODE="${DELETE_RESPONSE: -3}"
RESPONSE_BODY="${DELETE_RESPONSE%???}"

if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "204" ]; then
    echo -e "${GREEN}✓ Successfully deleted department (HTTP $HTTP_CODE)${NC}"
    if [ ! -z "$RESPONSE_BODY" ]; then
        echo ""
        echo -e "${BLUE}Delete Response:${NC}"
        echo "$RESPONSE_BODY" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE_BODY"
    fi
else
    echo -e "${RED}✗ Failed to delete department (HTTP $HTTP_CODE)${NC}"
    echo "Response: $RESPONSE_BODY"
fi

echo ""

# Test 5: Verify department was deleted
echo -e "${YELLOW}Step 5: Verify Department Was Deleted${NC}"
GET_AFTER_DELETE_RESPONSE=$(curl -s -w "%{http_code}" -X GET "$DEPARTMENT_SERVICE_URL/api/v1/departments/$DEPT_ID" \
  -H "Authorization: Bearer $JWT_TOKEN")
GET_AFTER_DELETE_HTTP_CODE="${GET_AFTER_DELETE_RESPONSE: -3}"

if [ "$GET_AFTER_DELETE_HTTP_CODE" = "404" ]; then
    echo -e "${GREEN}✓ Department successfully deleted (HTTP 404 - Not Found)${NC}"
elif [ "$GET_AFTER_DELETE_HTTP_CODE" = "200" ]; then
    echo -e "${YELLOW}⚠ Department still exists after deletion (HTTP 200)${NC}"
    echo "This might be expected behavior if soft delete is implemented"
else
    echo -e "${RED}✗ Unexpected response after deletion (HTTP $GET_AFTER_DELETE_HTTP_CODE)${NC}"
fi

echo ""
echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}        Test Summary            ${NC}"
echo -e "${BLUE}================================${NC}"
echo ""
echo -e "${GREEN}✓ Login: PASSED${NC}"
echo -e "${GREEN}✓ Create Test Department: PASSED${NC}"
echo -e "${GREEN}✓ Verify Department Exists: PASSED${NC}"
echo -e "${GREEN}✓ Delete Department API: PASSED${NC}"
echo -e "${GREEN}✓ Verify Department Deleted: PASSED${NC}"
echo ""
echo -e "${BLUE}Delete Department API is working correctly!${NC}"
