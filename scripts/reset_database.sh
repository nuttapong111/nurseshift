#!/bin/bash

# Script สำหรับรีเซ็ตฐานข้อมูล NurseShift และรัน schema ใหม่
# ใช้สำหรับเชื่อมต่อกับ PostgreSQL ที่ลงในเครื่อง (ไม่ใช่ Docker)

echo "🚀 เริ่มต้นรีเซ็ตฐานข้อมูล NurseShift..."

# ตั้งค่าตัวแปร
DB_USER="nuttapong2"
DB_NAME="nurseshift"
SCHEMA_FILE="database/schema.sql"
RESET_SCRIPT="scripts/reset_database.sql"

# ตรวจสอบว่าไฟล์ schema.sql มีอยู่หรือไม่
if [ ! -f "$SCHEMA_FILE" ]; then
    echo "❌ ไม่พบไฟล์ $SCHEMA_FILE"
    exit 1
fi

# ตรวจสอบว่าไฟล์ reset script มีอยู่หรือไม่
if [ ! -f "$RESET_SCRIPT" ]; then
    echo "❌ ไม่พบไฟล์ $RESET_SCRIPT"
    exit 1
fi

echo "📋 ข้อมูลการเชื่อมต่อ:"
echo "   - User: $DB_USER"
echo "   - Database: $DB_NAME"
echo "   - Schema File: $SCHEMA_FILE"
echo ""

# ขั้นตอนที่ 1: รีเซ็ตฐานข้อมูล
echo "🔄 ขั้นตอนที่ 1: รีเซ็ตฐานข้อมูล..."
psql -U "$DB_USER" -d postgres -f "$RESET_SCRIPT"

if [ $? -eq 0 ]; then
    echo "✅ รีเซ็ตฐานข้อมูลสำเร็จ"
else
    echo "❌ รีเซ็ตฐานข้อมูลล้มเหลว"
    exit 1
fi

# ขั้นตอนที่ 2: รัน schema ใหม่
echo "🔄 ขั้นตอนที่ 2: รัน schema ใหม่..."
psql -U "$DB_USER" -d "$DB_NAME" -f "$SCHEMA_FILE"

if [ $? -eq 0 ]; then
    echo "✅ รัน schema ใหม่สำเร็จ"
else
    echo "❌ รัน schema ใหม่ล้มเหลว"
    exit 1
fi

echo ""
echo "🎉 รีเซ็ตฐานข้อมูลและรัน schema ใหม่เสร็จสิ้น!"
echo "📊 ฐานข้อมูล $DB_NAME พร้อมใช้งานแล้ว"
echo ""
echo "💡 คำแนะนำ:"
echo "   - ตรวจสอบการเชื่อมต่อโดยรัน: psql -U $DB_USER -d $DB_NAME"
echo "   - ดูตารางทั้งหมด: \dt nurse_shift.*"
echo "   - ออกจาก psql: \q"
