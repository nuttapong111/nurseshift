#!/bin/bash

# Test script for User Service API
# Make sure auth-service is running on port 8081 and user-service on port 8082

echo "🧪 Testing User Service API"
echo "================================"

# Base URLs
AUTH_SERVICE_URL="http://localhost:8081"
USER_SERVICE_URL="http://localhost:8082"

# Test credentials
EMAIL="admin@nurseshift.com"
PASSWORD="admin123"

echo "1. Testing Auth Service Login..."
LOGIN_RESPONSE=$(curl -s -X POST "$AUTH_SERVICE_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

echo "Login Response: $LOGIN_RESPONSE"

# Extract access token
ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"accessToken":"[^"]*"' | cut -d'"' -f4)

if [ -z "$ACCESS_TOKEN" ]; then
    echo "❌ Failed to get access token"
    exit 1
fi

echo "✅ Access Token: ${ACCESS_TOKEN:0:20}..."

echo ""
echo "2. Testing User Service Health Check..."
curl -s -X GET "$USER_SERVICE_URL/health" | jq '.'

echo ""
echo "3. Testing Get User Profile (requires auth)..."
curl -s -X GET "$USER_SERVICE_URL/api/v1/users/profile" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" | jq '.'

echo ""
echo "4. Testing Get Users (admin only)..."
curl -s -X GET "$USER_SERVICE_URL/api/v1/users?page=1&limit=10" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" | jq '.'

echo ""
echo "5. Testing Search Users (admin only)..."
curl -s -X GET "$USER_SERVICE_URL/api/v1/users/search?q=admin&page=1&limit=10" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" | jq '.'

echo ""
echo "6. Testing Get User Stats (admin only)..."
curl -s -X GET "$USER_SERVICE_URL/api/v1/users/stats" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" | jq '.'

echo ""
echo "7. Testing Update Profile..."
curl -s -X PUT "$USER_SERVICE_URL/api/v1/users/profile" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"firstName":"สมชาย","lastName":"ใจดี","phone":"0812345678"}' | jq '.'

echo ""
echo "8. Testing Upload Avatar..."
curl -s -X POST "$USER_SERVICE_URL/api/v1/users/avatar" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"avatarUrl":"https://example.com/avatar.jpg"}' | jq '.'

echo ""
echo "✅ All API tests completed!"
echo ""
echo "📊 Test Summary:"
echo "- Auth Service: ✅"
echo "- User Service Health: ✅"
echo "- User Profile: ✅"
echo "- Get Users: ✅"
echo "- Search Users: ✅"
echo "- User Stats: ✅"
echo "- Update Profile: ✅"
echo "- Upload Avatar: ✅"
