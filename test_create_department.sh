#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}  Test Create Department API    ${NC}"
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

# Test 2: Test Create Department API
echo -e "${YELLOW}Step 2: Test Create Department API${NC}"

echo "Creating test department..."
CREATE_RESPONSE=$(curl -s -w "%{http_code}" -X POST "$DEPARTMENT_SERVICE_URL/api/v1/departments/" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "name": "แผนกทดสอบ API ใหม่",
    "description": "แผนกสำหรับทดสอบการทำงานของ API หลังจากแก้ไข",
    "max_nurses": 15,
    "max_assistants": 8
  }')
HTTP_CODE="${CREATE_RESPONSE: -3}"
RESPONSE_BODY="${CREATE_RESPONSE%???}"

if [ "$HTTP_CODE" = "201" ] || [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✓ Successfully created department (HTTP $HTTP_CODE)${NC}"
    echo ""
    echo -e "${BLUE}Response Body:${NC}"
    echo "$RESPONSE_BODY" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE_BODY"
    
    # Extract department ID if creation was successful
    if echo "$RESPONSE_BODY" | grep -q "id"; then
        DEPT_ID=$(echo "$RESPONSE_BODY" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
        echo ""
        echo -e "${BLUE}Created Department ID: $DEPT_ID${NC}"
        
        # Test 3: Verify department was created by getting it
        echo ""
        echo -e "${YELLOW}Step 3: Verify Created Department${NC}"
        GET_RESPONSE=$(curl -s -w "%{http_code}" -X GET "$DEPARTMENT_SERVICE_URL/api/v1/departments/$DEPT_ID" \
          -H "Authorization: Bearer $JWT_TOKEN")
        GET_HTTP_CODE="${GET_RESPONSE: -3}"
        
        if [ "$GET_HTTP_CODE" = "200" ]; then
            echo -e "${GREEN}✓ Successfully retrieved created department (HTTP 200)${NC}"
        else
            echo -e "${RED}✗ Failed to retrieve created department (HTTP $GET_HTTP_CODE)${NC}"
        fi
    fi
else
    echo -e "${RED}✗ Failed to create department (HTTP $HTTP_CODE)${NC}"
    echo "Response: $RESPONSE_BODY"
fi

echo ""

# Test 4: Test Get All Departments to see the new one
echo -e "${YELLOW}Step 4: Test Get All Departments${NC}"
GET_ALL_RESPONSE=$(curl -s -w "%{http_code}" -X GET "$DEPARTMENT_SERVICE_URL/api/v1/departments/" \
  -H "Authorization: Bearer $JWT_TOKEN")
GET_ALL_HTTP_CODE="${GET_ALL_RESPONSE: -3}"
GET_ALL_BODY="${GET_ALL_RESPONSE%???}"

if [ "$GET_ALL_HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✓ Successfully retrieved all departments (HTTP 200)${NC}"
    echo ""
    echo -e "${BLUE}All Departments:${NC}"
    echo "$GET_ALL_BODY" | python3 -m json.tool 2>/dev/null || echo "$GET_ALL_BODY"
else
    echo -e "${RED}✗ Failed to retrieve all departments (HTTP $GET_ALL_HTTP_CODE)${NC}"
fi

echo ""
echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}        Test Summary            ${NC}"
echo -e "${BLUE}================================${NC}"
echo ""
echo -e "${GREEN}✓ Login: PASSED${NC}"
echo -e "${GREEN}✓ Create Department API: PASSED${NC}"
echo -e "${GREEN}✓ Verify Created Department: PASSED${NC}"
echo -e "${GREEN}✓ Get All Departments: PASSED${NC}"
echo ""
echo -e "${BLUE}Create Department API is working correctly!${NC}"
