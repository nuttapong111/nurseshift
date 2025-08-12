#!/bin/bash

echo "🧪 Testing Password Reset Flow"
echo "================================"

# Test 1: Forgot Password
echo "1. Testing Forgot Password API..."
RESPONSE=$(curl -s -X POST http://localhost:8081/api/v1/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email": "test@nurseshift.com"}')

echo "Response: $RESPONSE"
echo ""

# Test 2: Try to reset password with invalid token
echo "2. Testing Reset Password API with invalid token..."
RESPONSE=$(curl -s -X POST http://localhost:8081/api/v1/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{"token": "123456", "newPassword": "newpassword123"}')

echo "Response: $RESPONSE"
echo ""

# Test 3: Try to reset password with valid token format but not stored
echo "3. Testing Reset Password API with valid token format..."
RESPONSE=$(curl -s -X POST http://localhost:8081/api/v1/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{"token": "999999", "newPassword": "newpassword123"}')

echo "Response: $RESPONSE"
echo ""

echo "✅ Password Reset Flow Test Completed!"
