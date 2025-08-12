#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}  Test Department Staff API     ${NC}"
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

# Test 2: Get Department Staff
echo -e "${YELLOW}Step 2: Test Get Department Staff API${NC}"

DEPARTMENT_ID="d8d95a30-5ce3-4b39-a7b9-1fcf7ddcab7d"
echo "Getting staff for department: $DEPARTMENT_ID"

STAFF_RESPONSE=$(curl -s -w "%{http_code}" -X GET "$DEPARTMENT_SERVICE_URL/api/v1/departments/$DEPARTMENT_ID/staff" \
  -H "Authorization: Bearer $JWT_TOKEN")

HTTP_CODE="${STAFF_RESPONSE: -3}"
RESPONSE_BODY="${STAFF_RESPONSE%???}"

if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✓ Successfully retrieved department staff (HTTP $HTTP_CODE)${NC}"
    echo ""
    echo "Response Body:"
    echo "$RESPONSE_BODY" | jq .
else
    echo -e "${RED}✗ Failed to retrieve department staff (HTTP $HTTP_CODE)${NC}"
    echo "Response: $RESPONSE_BODY"
fi

echo ""

# Test 3: Add Staff Member
echo -e "${YELLOW}Step 3: Test Add Staff Member API${NC}"

echo "Adding test staff member..."
ADD_STAFF_RESPONSE=$(curl -s -w "%{http_code}" -X POST "$DEPARTMENT_SERVICE_URL/api/v1/departments/$DEPARTMENT_ID/staff" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "first_name": "ทดสอบ",
    "last_name": "พนักงาน",
    "position": "nurse",
    "phone": "0812345678",
    "email": "test@example.com",
    "department_role": "nurse"
  }')

HTTP_CODE="${ADD_STAFF_RESPONSE: -3}"
RESPONSE_BODY="${ADD_STAFF_RESPONSE%???}"

if [ "$HTTP_CODE" = "201" ]; then
    echo -e "${GREEN}✓ Successfully added staff member (HTTP $HTTP_CODE)${NC}"
    echo ""
    echo "Response Body:"
    echo "$RESPONSE_BODY" | jq .
    
    # Extract staff ID for verification
    STAFF_ID=$(echo "$RESPONSE_BODY" | jq -r '.data.id')
    echo ""
    echo "Created Staff ID: $STAFF_ID"
else
    echo -e "${RED}✗ Failed to add staff member (HTTP $HTTP_CODE)${NC}"
    echo "Response: $RESPONSE_BODY"
fi

echo ""

# Test 4: Verify Staff Member
echo -e "${YELLOW}Step 4: Verify Added Staff Member${NC}"

if [ -n "$STAFF_ID" ] && [ "$STAFF_ID" != "null" ]; then
    echo "Verifying staff member: $STAFF_ID"
    
    VERIFY_RESPONSE=$(curl -s -w "%{http_code}" -X GET "$DEPARTMENT_SERVICE_URL/api/v1/departments/$DEPARTMENT_ID/staff" \
      -H "Authorization: Bearer $JWT_TOKEN")
    
    HTTP_CODE="${VERIFY_RESPONSE: -3}"
    RESPONSE_BODY="${VERIFY_RESPONSE%???}"
    
    if [ "$HTTP_CODE" = "200" ]; then
        echo -e "${GREEN}✓ Successfully verified staff member (HTTP $HTTP_CODE)${NC}"
        echo ""
        echo "Updated Staff List:"
        echo "$RESPONSE_BODY" | jq .
    else
        echo -e "${RED}✗ Failed to verify staff member (HTTP $HTTP_CODE)${NC}"
        echo "Response: $RESPONSE_BODY"
    fi
else
    echo -e "${YELLOW}⚠️  Skipping verification - no staff ID available${NC}"
fi

echo ""
echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}        Test Summary            ${NC}"
echo -e "${BLUE}================================${NC}"

if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "201" ]; then
    echo -e "${GREEN}✓ Department Staff API is working correctly!${NC}"
else
    echo -e "${RED}✗ Department Staff API has issues${NC}"
fi
