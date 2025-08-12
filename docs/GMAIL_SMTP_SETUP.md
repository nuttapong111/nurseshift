# การตั้งค่า Gmail SMTP สำหรับการส่งอีเมล

## ขั้นตอนการตั้งค่า

### 1. เปิดใช้งาน 2-Factor Authentication
1. ไปที่ [Google Account Settings](https://myaccount.google.com/)
2. เลือก "Security" จากเมนูด้านซ้าย
3. เปิดใช้งาน "2-Step Verification"

### 2. สร้าง App Password
1. หลังจากเปิด 2-Factor Authentication แล้ว ให้ไปที่ "App passwords"
2. เลือก "Mail" และ "Other (Custom name)"
3. ตั้งชื่อ เช่น "NurseShift Email Service"
4. คัดลอก App Password ที่ได้ (16 ตัวอักษร)

### 3. แก้ไขไฟล์ config.env
แก้ไขไฟล์ `backend/user-service/config.env`:

```env
# SMTP Configuration for Email Service
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-16-digit-app-password
SMTP_FROM=your-email@gmail.com
```

**หมายเหตุสำคัญ:**
- `SMTP_PASSWORD` ต้องเป็น App Password ไม่ใช่รหัสผ่านปกติของ Gmail
- App Password มี 16 ตัวอักษร
- อย่าใส่เครื่องหมายขีด (-) ใน App Password

### 4. ตัวอย่างการตั้งค่า
```env
SMTP_USERNAME=nurseshift@gmail.com
SMTP_PASSWORD=abcd efgh ijkl mnop
SMTP_FROM=nurseshift@gmail.com
```

### 5. รีสตาร์ท Service
หลังจากแก้ไขไฟล์ config.env แล้ว ให้รีสตาร์ท user-service:

```bash
cd backend/user-service
pkill -f "user-service"
./user-service
```

## การทดสอบ

### 1. ทดสอบการส่งอีเมล
```bash
curl -X POST http://localhost:8082/api/v1/users/send-verification-email \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com"}'
```

### 2. ตรวจสอบ Logs
หากการตั้งค่าถูกต้อง จะเห็น log:
```
📧 Sending verification email to: test@example.com
```

หากมีปัญหา จะเห็น error message ที่เกี่ยวข้องกับ SMTP

## การแก้ไขปัญหาที่พบบ่อย

### 1. Authentication Failed
- ตรวจสอบว่าใช้ App Password ไม่ใช่รหัสผ่านปกติ
- ตรวจสอบว่าเปิด 2-Factor Authentication แล้ว

### 2. Connection Refused
- ตรวจสอบว่า SMTP_HOST และ SMTP_PORT ถูกต้อง
- ตรวจสอบการเชื่อมต่ออินเทอร์เน็ต

### 3. Email Not Received
- ตรวจสอบ Spam/Junk folder
- ตรวจสอบว่า SMTP_FROM ถูกต้อง
- ตรวจสอบ logs ของ service

## ความปลอดภัย

### 1. App Password
- App Password มีสิทธิ์จำกัดเฉพาะการส่งอีเมล
- สามารถยกเลิกได้โดยไม่กระทบกับบัญชีหลัก
- ควรใช้เฉพาะในระบบที่เชื่อถือได้

### 2. Environment Variables
- อย่า commit ไฟล์ config.env ที่มี credentials จริง
- ใช้ .env.local หรือ environment variables ของระบบ
- ตรวจสอบ .gitignore ว่ามี .env files

## ข้อมูลเพิ่มเติม

- [Gmail SMTP Settings](https://support.google.com/mail/answer/7126229)
- [App Passwords](https://support.google.com/accounts/answer/185833)
- [SMTP Configuration](https://en.wikipedia.org/wiki/Simple_Mail_Transfer_Protocol)
