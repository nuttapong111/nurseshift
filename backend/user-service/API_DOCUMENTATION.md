# User Service API Documentation

## Overview
User Service เป็น microservice ที่จัดการข้อมูลผู้ใช้ในระบบ NurseShift โดยใช้ JWT Token สำหรับ authentication และ authorization

## Base URL
```
http://localhost:8082
```

## Authentication
ทุก API ที่ต้องการ authentication ต้องส่ง JWT Token ใน header:
```
Authorization: Bearer <jwt_token>
```

## API Endpoints

### 1. Health Check
**GET** `/health`

ตรวจสอบสถานะของ service

**Response:**
```json
{
  "status": "ok",
  "service": "user-service",
  "timestamp": "2025-08-11T16:29:21.393663+07:00"
}
```

### 2. Get User Profile
**GET** `/api/v1/users/profile`

ดึงข้อมูลโปรไฟล์ของผู้ใช้ที่ login อยู่

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Response:**
```json
{
  "status": "success",
  "message": "ดึงข้อมูลโปรไฟล์สำเร็จ",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "email": "admin@nurseshift.com",
    "firstName": "ผู้ดูแล",
    "lastName": "ระบบ",
    "phone": "+66-81-234-5678",
    "role": "admin",
    "status": "active",
    "position": "System Administrator",
    "remainingDays": 90,
    "subscriptionExpiresAt": "2025-11-09T10:17:49.136878+07:00",
    "packageType": "enterprise",
    "maxDepartments": 20,
    "avatarUrl": null,
    "settings": "{\"theme\": \"dark\", \"language\": \"th\"}",
    "lastLoginAt": "2025-08-11T16:29:24.814913+07:00",
    "createdAt": "2025-08-11T10:17:49.136878+07:00",
    "updatedAt": "2025-08-11T16:29:24.815158+07:00"
  }
}
```

### 3. Update User Profile
**PUT** `/api/v1/users/profile`

อัปเดตข้อมูลโปรไฟล์ของผู้ใช้

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "firstName": "ชื่อใหม่",
  "lastName": "นามสกุลใหม่",
  "phone": "เบอร์โทรใหม่",
  "position": "ตำแหน่งใหม่"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "อัปเดตโปรไฟล์สำเร็จ",
  "data": {
    // ข้อมูลผู้ใช้ที่อัปเดตแล้ว
  }
}
```

### 4. Get Users (Admin Only)
**GET** `/api/v1/users`

ดึงรายการผู้ใช้ทั้งหมด (เฉพาะ admin)

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Query Parameters:**
- `page` (optional): หน้าข้อมูล (default: 1)
- `limit` (optional): จำนวนข้อมูลต่อหน้า (default: 10)
- `role` (optional): กรองตาม role (admin, user)
- `status` (optional): กรองตาม status (active, inactive, pending, suspended)

**Response:**
```json
{
  "status": "success",
  "message": "ดึงข้อมูลผู้ใช้สำเร็จ",
  "data": {
    "users": [
      // รายการผู้ใช้
    ],
    "total": 4,
    "page": 1,
    "limit": 10,
    "totalPages": 1
  }
}
```

### 5. Search Users (Admin Only)
**GET** `/api/v1/users/search`

ค้นหาผู้ใช้ตามคำค้นหา (เฉพาะ admin)

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Query Parameters:**
- `q` (required): คำค้นหา
- `page` (optional): หน้าข้อมูล (default: 1)
- `limit` (optional): จำนวนข้อมูลต่อหน้า (default: 10)
- `role` (optional): กรองตาม role
- `status` (optional): กรองตาม status

**Response:**
```json
{
  "status": "success",
  "message": "ค้นหาผู้ใช้สำเร็จ",
  "data": {
    "users": [
      // รายการผู้ใช้ที่ค้นพบ
    ],
    "total": 2,
    "page": 1,
    "limit": 10,
    "totalPages": 1
  }
}
```

### 6. Get User Statistics (Admin Only)
**GET** `/api/v1/users/stats`

ดึงสถิติข้อมูลผู้ใช้ (เฉพาะ admin)

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Response:**
```json
{
  "status": "success",
  "message": "ดึงสถิติผู้ใช้สำเร็จ",
  "data": {
    "totalUsers": 4,
    "activeUsers": 4,
    "inactiveUsers": 0,
    "adminCount": 1,
    "userCount": 3
  }
}
```

### 7. Get Specific User
**GET** `/api/v1/users/{id}`

ดึงข้อมูลผู้ใช้ตาม ID

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Path Parameters:**
- `id`: UUID ของผู้ใช้

**Response:**
```json
{
  "status": "success",
  "message": "ดึงข้อมูลผู้ใช้สำเร็จ",
  "data": {
    // ข้อมูลผู้ใช้
  }
}
```

### 8. Upload Avatar
**POST** `/api/v1/users/avatar`

อัปเดตรูปโปรไฟล์

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "avatarUrl": "https://example.com/avatar.jpg"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "อัปเดตรูปโปรไฟล์สำเร็จ"
}
```

## Error Responses

### 401 Unauthorized
```json
{
  "error": "Authorization header required",
  "message": "กรุณาเข้าสู่ระบบ"
}
```

### 403 Forbidden
```json
{
  "error": "Insufficient permissions",
  "message": "คุณไม่มีสิทธิ์เข้าถึงข้อมูลนี้"
}
```

### 500 Internal Server Error
```json
{
  "status": "error",
  "message": "ไม่สามารถดึงข้อมูลโปรไฟล์ได้",
  "error": "error details"
}
```

## Data Models

### User Entity
```typescript
interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  phone?: string;
  role: 'admin' | 'user';
  status: 'active' | 'inactive' | 'pending' | 'suspended';
  position?: string;
  remainingDays: number;
  subscriptionExpiresAt?: Date;
  packageType: 'standard' | 'enterprise' | 'trial';
  maxDepartments: number;
  avatarUrl?: string;
  settings?: string;
  lastLoginAt?: Date;
  createdAt: Date;
  updatedAt: Date;
}
```

## Testing

### Run Test Script
```bash
./test_user_service_api.sh
```

### Manual Testing with curl
```bash
# Login to get token
curl -X POST "http://localhost:8081/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@nurseshift.com","password":"admin123"}'

# Use token to access protected endpoints
curl -X GET "http://localhost:8082/api/v1/users/profile" \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json"
```

## Notes
- Service ใช้ PostgreSQL database
- JWT Token ต้องได้จาก Auth Service
- Admin สามารถเข้าถึงข้อมูลผู้ใช้ทั้งหมดได้
- User ปกติสามารถเข้าถึงเฉพาะข้อมูลของตัวเองได้
- ข้อมูลทั้งหมดเป็นภาษาไทย
