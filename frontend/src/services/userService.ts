import type { 
  User, 
  ApiResponse, 
  SendVerificationEmailResponse, 
  VerifyEmailResponse, 
  CheckEmailVerificationResponse 
} from '@/types'

import { normalizeBaseUrl } from '@/lib/utils'
const USER_SERVICE_URL = normalizeBaseUrl(process.env.NEXT_PUBLIC_USER_SERVICE_URL, 'http://localhost:8082')

// Helper function to get auth token
export const getAuthToken = (): string | null => {
  if (typeof window !== 'undefined') {
    return localStorage.getItem('token')
  }
  return null
}

// Helper function to handle API responses
const handleApiResponse = async <T>(response: Response): Promise<T> => {
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}))
    throw new Error(errorData.message || `HTTP error! status: ${response.status}`)
  }
  return response.json()
}

// User Service API functions
export const userService = {
  // Get user profile
  async getProfile(): Promise<User> {
    const token = getAuthToken()
    if (!token) {
      throw new Error('No authentication token found')
    }

    const response = await fetch(`${USER_SERVICE_URL}/api/v1/users/profile`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    })

    const data = await handleApiResponse<ApiResponse<User>>(response)
    return data.data!
  },

  // Update user profile
  async updateProfile(profileData: Partial<User>): Promise<User> {
    const token = getAuthToken()
    if (!token) {
      throw new Error('No authentication token found')
    }

    const response = await fetch(`${USER_SERVICE_URL}/api/v1/users/profile`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(profileData),
    })

    const data = await handleApiResponse<ApiResponse<User>>(response)
    return data.data!
  },

  // Upload avatar
  async uploadAvatar(avatarUrl: string): Promise<void> {
    const token = getAuthToken()
    if (!token) {
      throw new Error('No authentication token found')
    }

    const response = await fetch(`${USER_SERVICE_URL}/api/v1/users/avatar`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ avatarUrl }),
    })

    await handleApiResponse<ApiResponse<void>>(response)
  },

  // Get all users (admin only)
  async getUsers(params?: {
    page?: number
    limit?: number
    role?: string
    status?: string
  }): Promise<{
    users: User[]
    total: number
    page: number
    limit: number
    totalPages: number
  }> {
    const token = getAuthToken()
    if (!token) {
      throw new Error('No authentication token found')
    }

    const searchParams = new URLSearchParams()
    if (params?.page) searchParams.append('page', params.page.toString())
    if (params?.limit) searchParams.append('limit', params.limit.toString())
    if (params?.role) searchParams.append('role', params.role)
    if (params?.status) searchParams.append('status', params.status)

    const response = await fetch(`${USER_SERVICE_URL}/api/v1/users?${searchParams}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    })

    const data = await handleApiResponse<ApiResponse<{
      users: User[]
      total: number
      page: number
      limit: number
      totalPages: number
    }>>(response)
    return data.data!
  },

  // Search users (admin only)
  async searchUsers(params: {
    query: string
    page?: number
    limit?: number
    role?: string
    status?: string
  }): Promise<{
    users: User[]
    total: number
    page: number
    limit: number
    totalPages: number
  }> {
    const token = getAuthToken()
    if (!token) {
      throw new Error('No authentication token found')
    }

    const searchParams = new URLSearchParams()
    searchParams.append('q', params.query)
    if (params?.page) searchParams.append('page', params.page.toString())
    if (params?.limit) searchParams.append('limit', params.limit.toString())
    if (params?.role) searchParams.append('role', params.role)
    if (params?.status) searchParams.append('status', params.status)

    const response = await fetch(`${USER_SERVICE_URL}/api/v1/users/search?${searchParams}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    })

    const data = await handleApiResponse<ApiResponse<{
      users: User[]
      total: number
      page: number
      limit: number
      totalPages: number
    }>>(response)
    return data.data!
  },

  // Get user statistics (admin only)
  async getUserStats(): Promise<{
    totalUsers: number
    activeUsers: number
    inactiveUsers: number
    adminCount: number
    userCount: number
  }> {
    const token = getAuthToken()
    if (!token) {
      throw new Error('No authentication token found')
    }

    const response = await fetch(`${USER_SERVICE_URL}/api/v1/users/stats`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    })

    const data = await handleApiResponse<ApiResponse<{
      totalUsers: number
      activeUsers: number
      inactiveUsers: number
      adminCount: number
      userCount: number
    }>>(response)
    return data.data!
  },

  // Get specific user by ID
  async getUserById(userId: string): Promise<User> {
    const token = getAuthToken()
    if (!token) {
      throw new Error('No authentication token found')
    }

    const response = await fetch(`${USER_SERVICE_URL}/api/v1/users/${userId}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    })

    const data = await handleApiResponse<ApiResponse<User>>(response)
    return data.data!
  },

  // Health check
  async healthCheck(): Promise<{ status: string; service: string; timestamp: string }> {
    const response = await fetch(`${USER_SERVICE_URL}/health`)
    return handleApiResponse<{ status: string; service: string; timestamp: string }>(response)
  },

  // Email Verification APIs
  async sendVerificationEmail(email: string): Promise<SendVerificationEmailResponse> {
    const response = await fetch(`${USER_SERVICE_URL}/api/v1/users/send-verification-email`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email }),
    })

    return handleApiResponse<SendVerificationEmailResponse>(response)
  },

  async verifyEmail(token: string): Promise<VerifyEmailResponse> {
    const response = await fetch(`${USER_SERVICE_URL}/api/v1/users/verify-email`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ token }),
    })

    return handleApiResponse<VerifyEmailResponse>(response)
  },

  async checkEmailVerification(email: string): Promise<CheckEmailVerificationResponse> {
    const response = await fetch(`${USER_SERVICE_URL}/api/v1/users/check-email-verification/${encodeURIComponent(email)}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    })

    return handleApiResponse<CheckEmailVerificationResponse>(response)
  }
}

export default userService
