export interface User {
  id: string
  email: string
  firstName: string
  lastName: string
  phone?: string
  role: 'user' | 'admin'
  status: 'active' | 'inactive' | 'pending' | 'suspended'
  position?: string
  remainingDays: number
  subscriptionExpiresAt?: string
  packageType: 'standard' | 'enterprise' | 'trial'
  maxDepartments: number
  avatarUrl?: string
  settings?: string
  lastLoginAt?: string
  emailVerified: boolean
  emailVerificationToken?: string
  emailVerificationExpiresAt?: string
  createdAt: string
  updatedAt: string
}

export interface Department {
  id: string
  name: string
  userId: string
  createdAt: Date
  updatedAt: Date
}

export interface Employee {
  id: string
  firstName: string
  lastName: string
  position: 'nurse' | 'assistant'
  departmentId: string
  createdAt: Date
  updatedAt: Date
}

export interface Shift {
  id: string
  name: string
  startTime: string
  endTime: string
  nurseCount: number
  assistantCount: number
  isActive: boolean
  departmentId: string
  createdAt: Date
  updatedAt: Date
}

export interface Schedule {
  id: string
  date: Date
  shiftId: string
  employeeId: string
  departmentId: string
  createdAt: Date
  updatedAt: Date
}

export interface Holiday {
  id: string
  name: string
  date: Date
  departmentId: string
  createdAt: Date
  updatedAt: Date
}

export interface EmployeeLeave {
  id: string
  employeeId: string
  date: Date
  reason?: string
  departmentId: string
  createdAt: Date
  updatedAt: Date
}

export interface Priority {
  id: string
  name: string
  order: number
  isActive: boolean
  userId: string
  createdAt: Date
  updatedAt: Date
}

export interface Package {
  id: string
  name: string
  price: number
  durationDays: number
  description: string
  isActive: boolean
  createdAt: Date
  updatedAt: Date
}

export interface Payment {
  id: string
  userId: string
  packageId: string
  amount: number
  status: 'pending' | 'approved' | 'rejected'
  evidence?: string
  createdAt: Date
  updatedAt: Date
}

export interface DashboardStats {
  totalNurses: number
  totalAssistants: number
  totalShifts: number
  totalDepartments: number
}

export interface AuthResponse {
  status: string
  message: string
  accessToken: string
  refreshToken: string
  expiresAt: Date
  user: User
}

export interface ApiResponse<T> {
  success: boolean
  message: string
  data?: T
}

// Password Reset Types
export interface ForgotPasswordRequest {
  email: string
}

export interface ForgotPasswordResponse {
  status: string
  message: string
}

export interface ResetPasswordRequest {
  token: string
  newPassword: string
}

export interface ResetPasswordResponse {
  status: string
  message: string
}

// Email Verification Types
export interface SendVerificationEmailRequest {
  email: string
}

export interface SendVerificationEmailResponse {
  status: string
  message: string
}

export interface VerifyEmailRequest {
  token: string
}

export interface VerifyEmailResponse {
  status: string
  message: string
}

export interface CheckEmailVerificationRequest {
  email: string
}

export interface CheckEmailVerificationResponse {
  status: string
  message: string
  isVerified: boolean
}
