#!/bin/bash

# Monitor user-service logs in real-time
# ใช้สำหรับดู logs เมื่อมีการส่งอีเมล

echo "📊 Monitoring user-service logs..."
echo "=================================="
echo ""

# ตรวจสอบว่า user-service กำลังทำงานอยู่หรือไม่
if ! pgrep -f "user-service" > /dev/null; then
    echo "❌ user-service is not running"
    echo "💡 Please start user-service first:"
    echo "   ./user-service"
    exit 1
fi

echo "✅ user-service is running (PID: $(pgrep -f "user-service"))"
echo ""

# แสดง logs แบบ real-time
echo "🔍 Monitoring logs (Press Ctrl+C to stop)..."
echo ""

# ใช้ tail -f เพื่อดู logs แบบ real-time
# หากมี log file ให้ใช้ tail -f logfile
# หากไม่มี ให้แสดงข้อความแนะนำ

echo "📝 Note: user-service logs are displayed in the terminal where it was started"
echo "💡 To see logs, check the terminal where you ran './user-service'"
echo ""

# แสดง logs จาก process
echo "📊 Current user-service process info:"
ps aux | grep "user-service" | grep -v grep

echo ""
echo "🔧 To see real-time logs:"
echo "   1. Open a new terminal"
echo "   2. Go to backend/user-service directory"
echo "   3. Run: tail -f /dev/null & ./user-service"
echo ""

# แสดง logs จาก console output (ถ้ามี)
echo "📋 Recent console output (if available):"
echo "----------------------------------------"

# ตรวจสอบว่ามี log file หรือไม่
if [ -f "user-service.log" ]; then
    echo "📁 Found log file: user-service.log"
    echo "📖 Last 20 lines:"
    tail -20 user-service.log
else
    echo "📁 No log file found"
    echo "💡 Logs are displayed in the terminal where user-service is running"
fi

echo ""
echo "🧪 To test email sending:"
echo "   ./test_email.sh"
echo ""
echo "📚 For troubleshooting:"
echo "   docs/EMAIL_TROUBLESHOOTING.md"
