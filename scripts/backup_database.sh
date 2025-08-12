#!/bin/bash

# Database backup script for NurseShift
# ไฟล์: scripts/backup_database.sh

# ตั้งค่าฐานข้อมูล
DB_NAME="nurseshift"
DB_USER="nuttapong2"
BACKUP_DIR="./backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/nurseshift_backup_${TIMESTAMP}.sql"

# สร้างโฟลเดอร์ backup ถ้ายังไม่มี
mkdir -p "$BACKUP_DIR"

echo "🔄 กำลังสร้าง backup ฐานข้อมูล..."
echo "📊 ฐานข้อมูล: $DB_NAME"
echo "👤 ผู้ใช้: $DB_USER"
echo "📁 ไฟล์: $BACKUP_FILE"

# สร้าง backup
pg_dump -U "$DB_USER" -d "$DB_NAME" > "$BACKUP_FILE"

if [ $? -eq 0 ]; then
    echo "✅ Backup สำเร็จ!"
    echo "📁 ไฟล์: $BACKUP_FILE"
    echo "📏 ขนาด: $(du -h "$BACKUP_FILE" | cut -f1)"
else
    echo "❌ Backup ล้มเหลว!"
    exit 1
fi

echo "🎯 Backup เสร็จสิ้นที่: $(date)"
