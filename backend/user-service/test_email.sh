#!/bin/bash

# Test Email Verification API
# ใช้สำหรับทดสอบการส่งอีเมลยืนยัน

echo "🧪 Testing Email Verification API"
echo "=================================="

# ตรวจสอบว่า user-service กำลังทำงานอยู่หรือไม่
if ! curl -s http://localhost:8082/health > /dev/null 2>&1; then
    echo "❌ user-service is not running on port 8082"
    echo "💡 Please start user-service first:"
    echo "   cd backend/user-service && ./user-service"
    exit 1
fi

echo "✅ user-service is running on port 8082"
echo ""

# ทดสอบการส่งอีเมลยืนยัน
echo "📧 Testing Send Verification Email API..."
echo "----------------------------------------"

# อีเมลสำหรับทดสอบ (แก้ไขตามต้องการ)
TEST_EMAIL="test@example.com"

echo "📝 Sending verification email to: $TEST_EMAIL"

RESPONSE=$(curl -s -X POST http://localhost:8082/api/v1/users/send-verification-email \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$TEST_EMAIL\"}")

echo "📤 Response:"
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"

echo ""
echo "🔍 Check user-service logs for email sending details:"
echo "   - If SMTP is configured: Look for '✅ Verification email sent successfully'"
echo "   - If SMTP is not configured: Look for '⚠️  SMTP not configured'"
echo ""

# ตรวจสอบ SMTP configuration
echo "🔧 Checking SMTP Configuration..."
echo "--------------------------------"

if [ -f "config.env" ]; then
    echo "📁 config.env found"
    
    if grep -q "SMTP_USERNAME" config.env; then
        SMTP_USER=$(grep "SMTP_USERNAME" config.env | cut -d'=' -f2)
        if [ "$SMTP_USER" != "your-email@gmail.com" ] && [ -n "$SMTP_USER" ]; then
            echo "✅ SMTP_USERNAME is configured: $SMTP_USER"
        else
            echo "❌ SMTP_USERNAME is not properly configured"
        fi
    else
        echo "❌ SMTP_USERNAME not found in config.env"
    fi
    
    if grep -q "SMTP_PASSWORD" config.env; then
        SMTP_PASS=$(grep "SMTP_PASSWORD" config.env | cut -d'=' -f2)
        if [ "$SMTP_PASS" != "your-app-password" ] && [ -n "$SMTP_PASS" ]; then
            echo "✅ SMTP_PASSWORD is configured"
        else
            echo "❌ SMTP_PASSWORD is not properly configured"
        fi
    else
        echo "❌ SMTP_PASSWORD not found in config.env"
    fi
    
    if grep -q "SMTP_FROM" config.env; then
        SMTP_FROM=$(grep "SMTP_FROM" config.env | cut -d'=' -f2)
        if [ "$SMTP_FROM" != "your-email@gmail.com" ] && [ -n "$SMTP_FROM" ]; then
            echo "✅ SMTP_FROM is configured: $SMTP_FROM"
        else
            echo "❌ SMTP_FROM is not properly configured"
        fi
    else
        echo "❌ SMTP_FROM not found in config.env"
    fi
else
    echo "❌ config.env not found"
fi

echo ""
echo "📚 For SMTP setup instructions, see: docs/GMAIL_SMTP_SETUP.md"
echo ""

# ทดสอบการตรวจสอบสถานะการยืนยันอีเมล
echo "🔍 Testing Check Email Verification Status API..."
echo "------------------------------------------------"

echo "📝 Checking verification status for: $TEST_EMAIL"

STATUS_RESPONSE=$(curl -s "http://localhost:8082/api/v1/users/check-email-verification/$TEST_EMAIL")

echo "📤 Status Response:"
echo "$STATUS_RESPONSE" | jq '.' 2>/dev/null || echo "$STATUS_RESPONSE"

echo ""
echo "✨ Test completed!"
echo ""
echo "💡 Next steps:"
echo "   1. Configure SMTP settings in config.env if not done"
echo "   2. Restart user-service after configuration changes"
echo "   3. Check email inbox for verification email"
echo "   4. Use the verification token to verify email"
