# การแก้ไขปัญหาการส่งอีเมล (Email Troubleshooting Guide)

## ปัญหาที่พบบ่อย

### 1. อีเมลไม่ถูกส่ง (Email not sent)

#### อาการ
- API ส่งคืน success แต่ไม่ได้รับอีเมล
- ใน console แสดง "⚠️ SMTP not configured"
- ไม่มี error message

#### สาเหตุ
- SMTP credentials ไม่ถูกตั้งค่าใน `config.env`
- ใช้ placeholder values แทนค่าจริง

#### วิธีแก้ไข
1. ตรวจสอบไฟล์ `backend/user-service/config.env`
2. แก้ไข SMTP settings ให้ถูกต้อง:

```env
SMTP_USERNAME=your-real-email@gmail.com
SMTP_PASSWORD=your-16-digit-app-password
SMTP_FROM=your-real-email@gmail.com
```

3. รีสตาร์ท user-service:
```bash
cd backend/user-service
pkill -f "user-service"
./user-service
```

### 2. Authentication Failed

#### อาการ
- Error: "failed to authenticate"
- "535 Authentication failed"

#### สาเหตุ
- ใช้รหัสผ่านปกติของ Gmail แทน App Password
- ไม่ได้เปิด 2-Factor Authentication
- App Password หมดอายุ

#### วิธีแก้ไข
1. เปิด 2-Factor Authentication ใน Google Account
2. สร้าง App Password ใหม่
3. ใช้ App Password ใน `SMTP_PASSWORD`
4. รีสตาร์ท service

### 3. Connection Refused

#### อาการ
- Error: "connection refused"
- "failed to connect to SMTP server"

#### สาเหตุ
- SMTP_HOST หรือ SMTP_PORT ไม่ถูกต้อง
- Firewall หรือ network blocking
- Gmail SMTP server ไม่สามารถเข้าถึงได้

#### วิธีแก้ไข
1. ตรวจสอบ SMTP settings:
```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
```

2. ตรวจสอบการเชื่อมต่ออินเทอร์เน็ต
3. ตรวจสอบ firewall settings

### 4. Email Not Received

#### อาการ
- API ส่งคืน success
- ไม่มี error ใน logs
- อีเมลไม่ปรากฏใน inbox

#### สาเหตุ
- อีเมลไปตกใน Spam/Junk folder
- SMTP_FROM ไม่ถูกต้อง
- Gmail filtering

#### วิธีแก้ไข
1. ตรวจสอบ Spam/Junk folder
2. ตรวจสอบ SMTP_FROM ว่าตรงกับ SMTP_USERNAME
3. เพิ่ม email address ใน contacts
4. ตรวจสอบ Gmail filters

## การตรวจสอบและ Debug

### 1. ตรวจสอบ SMTP Configuration

รันสคริปต์ทดสอบ:
```bash
cd backend/user-service
./test_email.sh
```

### 2. ตรวจสอบ Logs

ดู logs ของ user-service:
```bash
# ใน terminal ที่รัน user-service
# หรือดู console output
```

### 3. ทดสอบ SMTP Connection

ทดสอบการเชื่อมต่อ SMTP:
```bash
telnet smtp.gmail.com 587
```

### 4. ตรวจสอบ Environment Variables

ตรวจสอบว่า environment variables ถูกอ่านหรือไม่:
```bash
cd backend/user-service
grep -E "SMTP_" config.env
```

## การตั้งค่า SMTP ที่ถูกต้อง

### 1. Gmail SMTP Settings
```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-16-digit-app-password
SMTP_FROM=your-email@gmail.com
```

### 2. การสร้าง App Password
1. เปิด [Google Account Settings](https://myaccount.google.com/)
2. ไปที่ Security > 2-Step Verification
3. เลือก App passwords
4. สร้าง App Password สำหรับ "Mail"
5. คัดลอก 16 ตัวอักษร

### 3. ตัวอย่างการตั้งค่า
```env
SMTP_USERNAME=nurseshift@gmail.com
SMTP_PASSWORD=abcd efgh ijkl mnop
SMTP_FROM=nurseshift@gmail.com
```

## การทดสอบ

### 1. ทดสอบการส่งอีเมล
```bash
curl -X POST http://localhost:8082/api/v1/users/send-verification-email \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com"}'
```

### 2. ตรวจสอบสถานะ
```bash
curl http://localhost:8082/api/v1/users/check-email-verification/test@example.com
```

### 3. ตรวจสอบ Logs
```bash
# ดู console output ของ user-service
# หรือใช้ log file ถ้ามี
```

## การแก้ไขปัญหาแบบ Step-by-Step

### Step 1: ตรวจสอบ Configuration
```bash
cd backend/user-service
cat config.env | grep SMTP
```

### Step 2: ตรวจสอบ Service Status
```bash
curl http://localhost:8082/health
```

### Step 3: ทดสอบการส่งอีเมล
```bash
./test_email.sh
```

### Step 4: ตรวจสอบ Logs
ดู console output ของ user-service

### Step 5: แก้ไข Configuration
แก้ไขไฟล์ `config.env` ตามที่แนะนำ

### Step 6: รีสตาร์ท Service
```bash
pkill -f "user-service"
./user-service
```

### Step 7: ทดสอบอีกครั้ง
```bash
./test_email.sh
```

## การป้องกันปัญหาในอนาคต

### 1. ใช้ Environment Variables
```bash
export SMTP_USERNAME="your-email@gmail.com"
export SMTP_PASSWORD="your-app-password"
export SMTP_FROM="your-email@gmail.com"
```

### 2. ตรวจสอบ Configuration ก่อน Deploy
```bash
./test_email.sh
```

### 3. ใช้ Configuration Validation
เพิ่มการตรวจสอบ SMTP configuration ใน startup

### 4. Logging ที่ดี
เพิ่ม detailed logging สำหรับ SMTP operations

## ข้อมูลเพิ่มเติม

- [Gmail SMTP Setup Guide](docs/GMAIL_SMTP_SETUP.md)
- [Email Verification System](docs/email-verification.md)
- [API Documentation](backend/user-service/API_DOCUMENTATION.md)

## ติดต่อ Support

หากยังมีปัญหา กรุณา:
1. ตรวจสอบ logs ของ user-service
2. รัน `./test_email.sh` และแจ้งผลลัพธ์
3. ตรวจสอบ SMTP configuration
4. แจ้ง error message ที่ได้รับ
