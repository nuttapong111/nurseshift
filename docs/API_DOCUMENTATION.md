# 📚 **NurseShift API Documentation**

## **🎯 Overview**
NurseShift is a microservices-based nurse shift management system built with Go (Fiber framework) and clean architecture principles. This document provides comprehensive API documentation for all services.

## **🏗️ Architecture**
- **Frontend**: Next.js + TypeScript + Tailwind CSS
- **Backend**: Go (Fiber) Microservices with Clean Architecture
- **Database**: PostgreSQL with proper schema design
- **Authentication**: JWT-based authentication
- **Documentation**: Swagger/OpenAPI 3.0

## **🔧 Services Overview**

| Service | Port | Status | Description |
|---------|------|--------|-------------|
| **Auth Service** | 8081 | ✅ Complete | User authentication, JWT management |
| **User Service** | 8082 | ✅ Complete | User profile and management |
| **Department Service** | 8083 | ✅ Complete | Department and employee management |
| **Schedule Service** | 8084 | ✅ Complete | Shift scheduling and management |
| **Setting Service** | 8085 | ✅ Complete | System settings and configuration |
| **Priority Service** | 8086 | ✅ Complete | Scheduling priority management |
| **Notification Service** | 8087 | ✅ Complete | Notifications and alerts |
| **Package Service** | 8088 | ✅ Complete | Membership package management |
| **Payment Service** | 8089 | ✅ Complete | Payment processing and history |

## **📖 Interactive Documentation**

### **🔗 Swagger UI Access**
Each service provides interactive Swagger UI documentation:

- **🔐 Auth Service**: http://localhost:8081/swagger/
- **👤 User Service**: http://localhost:8082/swagger/
- **🏢 Department Service**: http://localhost:8083/swagger/
- **📅 Schedule Service**: http://localhost:8084/swagger/
- **⚙️ Setting Service**: http://localhost:8085/swagger/
- **🎯 Priority Service**: http://localhost:8086/swagger/
- **📢 Notification Service**: http://localhost:8087/swagger/
- **💰 Package Service**: http://localhost:8088/swagger/
- **💳 Payment Service**: http://localhost:8089/swagger/

## **🔒 Authentication**

### **Bearer Token Format**
```
Authorization: Bearer <JWT_TOKEN>
```

### **Login Endpoint**
```http
POST http://localhost:8081/api/v1/auth/login
Content-Type: application/json

{
  "email": "admin@thephyathai.com",
  "password": "password123"
}
```

### **Response Format**
```json
{
  "status": "success",
  "message": "เข้าสู่ระบบสำเร็จ",
  "data": {
    "user": { ... },
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresIn": 86400
  }
}
```

## **📋 Core API Endpoints**

### **🔐 Auth Service (8081)**
```
POST   /api/v1/auth/login      - User login
POST   /api/v1/auth/register   - User registration  
POST   /api/v1/auth/refresh    - Refresh token
GET    /api/v1/auth/me         - Get current user
GET    /health                 - Health check
```

### **👤 User Service (8082)**
```
GET    /api/v1/users/profile   - Get user profile
PUT    /api/v1/users/profile   - Update profile
POST   /api/v1/users/avatar    - Upload avatar
GET    /api/v1/users           - List users (with filters)
GET    /api/v1/users/search    - Search users
GET    /api/v1/users/stats     - User statistics
GET    /api/v1/users/:id       - Get specific user
GET    /health                 - Health check
```

### **🏢 Department Service (8083)**
```
GET    /api/v1/departments             - List departments
POST   /api/v1/departments             - Create department
GET    /api/v1/departments/stats       - Department statistics
GET    /api/v1/departments/:id         - Get department details
PUT    /api/v1/departments/:id         - Update department
DELETE /api/v1/departments/:id         - Delete department
GET    /api/v1/departments/:id/employees - Department employees
GET    /health                         - Health check
```

### **📅 Schedule Service (8084)**
```
GET    /api/v1/schedules       - Get schedules
POST   /api/v1/schedules       - Create schedule
GET    /health                 - Health check
```

### **⚙️ Setting Service (8085)**
```
GET    /api/v1/settings        - Get system settings
PUT    /api/v1/settings        - Update settings
GET    /health                 - Health check
```

### **🎯 Priority Service (8086)**
```
GET    /api/v1/priorities          - Get priorities
PUT    /api/v1/priorities/:id      - Update priority
PUT    /api/v1/priorities/:id/setting - Update priority settings
GET    /health                     - Health check
```

### **📢 Notification Service (8087)**
```
GET    /api/v1/notifications           - Get notifications
GET    /api/v1/notifications/stats     - Notification statistics
POST   /api/v1/notifications           - Create notification
PUT    /api/v1/notifications/:id/read  - Mark as read
PUT    /api/v1/notifications/read-all  - Mark all as read
DELETE /api/v1/notifications/:id       - Delete notification
GET    /health                         - Health check
```

### **💰 Package Service (8088)**
```
GET    /api/v1/packages            - Get packages
GET    /api/v1/packages/current    - Current user package
GET    /api/v1/packages/stats      - Package statistics
GET    /api/v1/packages/:id        - Get package details
POST   /api/v1/packages/order      - Create package order
PUT    /api/v1/packages/settings   - Update package settings
GET    /health                     - Health check
```

### **💳 Payment Service (8089)**
```
GET    /api/v1/payments            - Payment history
POST   /api/v1/payments            - Create payment
GET    /api/v1/payments/stats      - Payment statistics
GET    /api/v1/payments/:id        - Payment details
PUT    /api/v1/payments/:id        - Update payment (resubmit)
PUT    /api/v1/payments/:id/approve - Approve payment (admin)
PUT    /api/v1/payments/:id/reject  - Reject payment (admin)
GET    /health                     - Health check
```

## **📊 Response Format Standards**

### **Success Response**
```json
{
  "status": "success",
  "message": "ดำเนินการสำเร็จ",
  "data": { ... }
}
```

### **Error Response**
```json
{
  "status": "error", 
  "message": "เกิดข้อผิดพลาด",
  "error": "technical details (optional)"
}
```

### **Paginated Response**
```json
{
  "status": "success",
  "message": "ดึงข้อมูลสำเร็จ",
  "data": {
    "items": [...],
    "total": 100,
    "page": 1,
    "limit": 10,
    "totalPages": 10
  }
}
```

## **🚀 Quick Start**

### **1. Start All Services**
```bash
# Start all microservices
./scripts/start-all-services.sh start

# Check service status
./scripts/start-all-services.sh status

# Stop all services
./scripts/start-all-services.sh stop
```

### **2. Test Authentication**
```bash
# Login
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@thephyathai.com", "password": "password123"}'

# Use token for authenticated requests
curl -X GET http://localhost:8082/api/v1/users/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### **3. Health Checks**
```bash
# Check all services health
./scripts/start-all-services.sh health

# Individual health checks
curl http://localhost:8081/health  # Auth Service
curl http://localhost:8082/health  # User Service
# ... etc for all services
```

## **🔍 Development Notes**

### **Mock Data**
All services currently use mock data for rapid development and testing. Real database integration can be implemented in the next phase.

### **Security Features**
- ✅ JWT Authentication with refresh tokens
- ✅ Password hashing with bcrypt
- ✅ CORS configuration
- ✅ Request validation
- ✅ Error handling

### **Architecture Benefits**
- **Scalability**: Each service can be scaled independently
- **Maintainability**: Clean architecture with clear separation of concerns
- **Testability**: Easy to unit test individual components
- **Documentation**: Comprehensive Swagger documentation
- **Type Safety**: Strong typing with Go structs

## **📞 Support**

For questions or issues:
- **Team**: NurseShift Development Team
- **Email**: support@nurseshift.com
- **Documentation**: Check individual service Swagger UIs for detailed endpoint documentation

---

*Last Updated: August 2025*
*Version: 1.0.0*


