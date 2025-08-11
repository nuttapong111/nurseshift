# การตั้งค่า Gmail สำหรับส่งอีเมลใน NurseShift

## ขั้นตอนการตั้งค่า Gmail

### 1. เปิดใช้งาน 2-Factor Authentication
1. ไปที่ [Google Account Settings](https://myaccount.google.com/)
2. เลือก "Security" (ความปลอดภัย)
3. เปิดใช้งาน "2-Step Verification" (การยืนยันตัวตน 2 ขั้นตอน)

### 2. สร้าง App Password
1. หลังจากเปิดใช้งาน 2-Factor Authentication แล้ว
2. ไปที่ "App passwords" (รหัสผ่านแอป)
3. เลือก "Mail" และ "Other (Custom name)"
4. ตั้งชื่อเป็น "NurseShift"
5. คัดลอกรหัสผ่าน 16 หลักที่ได้

### 3. อัปเดต Config.env
แก้ไขไฟล์ `backend/auth-service/config.env`:

```env
# Email Configuration
EMAIL_PROVIDER=gmail
EMAIL_FROM_EMAIL=your-email@gmail.com
EMAIL_FROM_PASSWORD=your-16-digit-app-password
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
```

**หมายเหตุสำคัญ:**
- ใช้ App Password ไม่ใช่รหัสผ่าน Gmail ปกติ
- App Password มี 16 หลัก
- อย่าแชร์ App Password กับผู้อื่น

### 4. ทดสอบการส่งอีเมล

#### ใช้ Mock Service (สำหรับการพัฒนา)
```env
EMAIL_PROVIDER=mock
```

#### ใช้ Gmail Service (สำหรับการใช้งานจริง)
```env
EMAIL_PROVIDER=gmail
EMAIL_FROM_EMAIL=your-email@gmail.com
EMAIL_FROM_PASSWORD=your-16-digit-app-password
```

## การใช้งาน API

### Forgot Password
```bash
POST /api/v1/auth/forgot-password
{
  "email": "user@example.com"
}
```

### Reset Password
```bash
POST /api/v1/auth/reset-password
{
  "token": "123456",
  "newPassword": "newpassword123"
}
```

## การแก้ไขปัญหา

### ปัญหาที่พบบ่อย

1. **Authentication failed**
   - ตรวจสอบว่าใช้ App Password ไม่ใช่รหัสผ่าน Gmail ปกติ
   - ตรวจสอบว่าเปิดใช้งาน 2-Factor Authentication แล้ว

2. **Connection timeout**
   - ตรวจสอบ firewall และ network settings
   - ตรวจสอบว่า port 587 ไม่ถูกบล็อก

3. **Rate limiting**
   - Gmail มีการจำกัดจำนวนอีเมลที่ส่งต่อวัน
   - สำหรับการใช้งานจริง ควรพิจารณาใช้บริการอีเมลอื่น เช่น SendGrid, Mailgun

## ความปลอดภัย

- เก็บ App Password ไว้เป็นความลับ
- ใช้ environment variables ไม่ใช่ hardcode ในโค้ด
- ตรวจสอบ logs เพื่อดูการใช้งานที่ไม่ปกติ
- พิจารณาใช้ Redis หรือ database สำหรับเก็บ reset tokens แทน in-memory storage

## การพัฒนา

สำหรับการพัฒนาและทดสอบ ใช้ mock service:
```env
EMAIL_PROVIDER=mock
```

Mock service จะแสดงข้อมูลอีเมลใน console แทนการส่งจริง
