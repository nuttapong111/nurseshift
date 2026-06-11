# 🏥 NurseShift - ระบบจัดการเวรพยาบาล

ระบบจัดการเวรพยาบาลแบบครบวงจร ที่ออกแบบมาเพื่อช่วยให้หัวหน้าพยาบาลสามารถจัดการตารางเวร พนักงาน และการทำงานของแผนกได้อย่างมีประสิทธิภาพ

## ✨ คุณสมบัติหลัก

### 🔐 ระบบจัดการผู้ใช้
- การสมัครสมาชิกและเข้าสู่ระบบ
- ระบบจัดการสิทธิ์ (Admin/User)
- การรีเซ็ตรหัสผ่านผ่านอีเมล
- การจัดการโปรไฟล์ผู้ใช้

### 🏢 ระบบจัดการแผนก
- สร้างและจัดการแผนกต่างๆ
- กำหนดจำนวนพนักงานสูงสุด
- การตั้งค่าเฉพาะแผนก

### 👥 ระบบจัดการพนักงาน
- เพิ่ม/ลบ/แก้ไขข้อมูลพนักงาน
- จัดกลุ่มพนักงานตามแผนก
- ติดตามสถานะการทำงาน

### 📅 ระบบจัดการตารางเวร
- สร้างตารางเวรอัตโนมัติ
- จัดการกะการทำงาน (เช้า/บ่าย/ดึก)
- กำหนดวันทำงานและวันหยุด
- ระบบจัดลำดับความสำคัญ

### 📧 ระบบแจ้งเตือน
- แจ้งเตือนแบบเรียลไทม์
- การแจ้งเตือนผ่านอีเมล
- ระบบการแจ้งเตือนหลายระดับ

### 💰 ระบบจัดการแพ็คเกจ
- แพ็คเกจทดลองใช้ (30 วัน)
- แพ็คเกจมาตรฐาน (990 บาท/เดือน)
- แพ็คเกจระดับองค์กร (2,990 บาท/3 เดือน)

## 🏗️ สถาปัตยกรรมระบบ

### Backend (Microservices)
- **Auth Service** - ระบบจัดการผู้ใช้และการยืนยันตัวตน
- **Department Service** - ระบบจัดการแผนก
- **Schedule Service** - ระบบจัดการตารางเวร
- **User Service** - ระบบจัดการข้อมูลผู้ใช้
- **Notification Service** - ระบบแจ้งเตือน
- **Payment Service** - ระบบจัดการการชำระเงิน
- **Package Service** - ระบบจัดการแพ็คเกจ
- **Priority Service** - ระบบจัดการลำดับความสำคัญ
- **Setting Service** - ระบบจัดการการตั้งค่า
- **Employee Leave Service** - ระบบจัดการการลาของพนักงาน

### Frontend
- **Next.js 14** - React Framework
- **TypeScript** - Type Safety
- **Tailwind CSS** - Styling
- **Responsive Design** - รองรับทุกอุปกรณ์

### Database
- **PostgreSQL** - ฐานข้อมูลหลัก
- **Redis** - Cache และ Session Management

## 🚀 การติดตั้งและใช้งาน

### ความต้องการของระบบ
- Node.js 18+
- Go 1.21+
- PostgreSQL 14+
- Redis 6+

### การติดตั้ง Backend

1. **Clone โปรเจค**
```bash
git clone https://github.com/your-username/nurseshift-final.git
cd nurseshift-final
```

2. **ติดตั้ง Dependencies**
```bash
# Auth Service
cd backend/auth-service
go mod download
```

3. **ตั้งค่าฐานข้อมูล**
```bash
# รัน schema
psql -U your-username -d postgres -f database/schema.sql

# เพิ่มข้อมูลทดสอบ
psql -U your-username -d nurseshift -f scripts/seed_test_data.sql
```

4. **ตั้งค่า Environment Variables**
```bash
# แก้ไข backend/auth-service/config.env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your-username
DB_PASSWORD=your-password
DB_NAME=nurseshift
```

5. **รัน Service**
```bash
cd backend/auth-service
go run cmd/server/main.go
```

### การติดตั้ง Frontend

1. **ติดตั้ง Dependencies**
```bash
cd frontend
npm install
```

2. **รัน Development Server**
```bash
npm run dev
```

## 📧 การตั้งค่าอีเมล

สำหรับการใช้งานระบบรีเซ็ตรหัสผ่าน ต้องตั้งค่าการส่งอีเมล:

### ใช้ Gmail (แนะนำสำหรับการใช้งานจริง)
1. เปิดใช้งาน 2-Factor Authentication
2. สร้าง App Password
3. อัปเดต config.env:
```env
EMAIL_PROVIDER=gmail
EMAIL_FROM_EMAIL=your-email@gmail.com
EMAIL_FROM_PASSWORD=your-16-digit-app-password
```

### ใช้ Mock Service (สำหรับการพัฒนา)
```env
EMAIL_PROVIDER=mock
```

ดูรายละเอียดเพิ่มเติมได้ที่ [docs/GMAIL_SETUP.md](docs/GMAIL_SETUP.md)

## 🧪 การทดสอบ

### API Testing
```bash
# Health Check
curl http://localhost:8081/health

# Login
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@nurseshift.com","password":"admin123"}'
```

### Test Credentials
- **Admin**: admin@nurseshift.com / admin123
- **User**: user@nurseshift.com / user123
- **Test**: test@nurseshift.com / test123

## 📚 API Documentation

- **Swagger UI**: http://localhost:8081/swagger/
- **Health Check**: GET /health
- **Authentication**: POST /api/v1/auth/*
- **User Management**: GET/POST/PUT/DELETE /api/v1/users/*

## 🔧 การพัฒนา

### โครงสร้างโปรเจค
```
nurseshift_final/
├── backend/           # Microservices (Go)
├── frontend/          # Next.js Frontend
├── database/          # Database schemas
├── docs/             # Documentation
├── scripts/          # Utility scripts
└── deployment/       # Deployment configs
```

### การเพิ่ม Microservice ใหม่
```bash
# ใช้ script ที่มีอยู่
./scripts/create-microservice.sh service-name
```

## 🚀 การ Deploy

### Docker
```bash
docker-compose up -d
```

### Kubernetes
```bash
kubectl apply -f deployment/
```

## 🤝 การมีส่วนร่วม

1. Fork โปรเจค
2. สร้าง Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit การเปลี่ยนแปลง (`git commit -m 'Add some AmazingFeature'`)
4. Push ไปยัง Branch (`git push origin feature/AmazingFeature`)
5. เปิด Pull Request


**NurseShift** - ระบบจัดการเวรพยาบาลที่ทันสมัยและใช้งานง่าย 🏥✨
