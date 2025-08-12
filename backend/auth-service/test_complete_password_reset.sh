#!/bin/bash

echo "🧪 Complete Password Reset Flow Test"
echo "====================================="

# Test 1: Forgot Password with non-existent email
echo "1. Testing Forgot Password with non-existent email..."
RESPONSE=$(curl -s -X POST http://localhost:8081/api/v1/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email": "nonexistent@example.com"}')

echo "Response: $RESPONSE"
echo ""

# Test 2: Forgot Password with existing email
echo "2. Testing Forgot Password with existing email..."
RESPONSE=$(curl -s -X POST http://localhost:8081/api/v1/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@nurseshift.com"}')

echo "Response: $RESPONSE"
echo ""

# Test 3: Reset Password with invalid token
echo "3. Testing Reset Password with invalid token..."
RESPONSE=$(curl -s -X POST http://localhost:8081/api/v1/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{"token": "123456", "newPassword": "newpassword123"}')

echo "Response: $RESPONSE"
echo ""

# Test 4: Reset Password with valid token format but not stored
echo "4. Testing Reset Password with valid token format..."
RESPONSE=$(curl -s -X POST http://localhost:8081/api/v1/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{"token": "999999", "newPassword": "newpassword123"}')

echo "Response: $RESPONSE"
echo ""

# Test 5: Test with Gmail email (if configured)
echo "5. Testing Forgot Password with Gmail email..."
RESPONSE=$(curl -s -X POST http://localhost:8081/api/v1/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email": "nurseshift.user@gmail.com"}')

echo "Response: $RESPONSE"
echo ""

echo "✅ Complete Password Reset Flow Test Completed!"
echo ""
echo "📋 Test Summary:"
echo "- Non-existent email: Should return generic message"
echo "- Existing email: Should return success message"
echo "- Invalid token: Should return error message"
echo "- Gmail email: Should work if user exists in system"
