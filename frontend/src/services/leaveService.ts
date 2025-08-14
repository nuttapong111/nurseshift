import type { ApiResponse } from '@/types'

const LEAVE_SERVICE_URL = 'http://localhost:8090'

export const getAuthToken = (): string | null => {
  if (typeof window !== 'undefined') {
    return localStorage.getItem('token')
  }
  return null
}

const handleApiResponse = async <T>(response: Response): Promise<T> => {
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}))
    throw new Error(errorData.message || `HTTP error! status: ${response.status}`)
  }
  return response.json()
}

export const leaveService = {
  async list(params?: { month?: string; employeeId?: string; departmentId?: string }) {
    const token = getAuthToken()
    if (!token) throw new Error('No authentication token found')

    const search = new URLSearchParams()
    if (params?.month) search.append('month', params.month)
    if (params?.employeeId) search.append('employeeId', params.employeeId)
    if (params?.departmentId) search.append('departmentId', params.departmentId)

    const res = await fetch(`${LEAVE_SERVICE_URL}/api/v1/leaves?${search}`, {
      headers: { Authorization: `Bearer ${token}` },
    })
    const data = await handleApiResponse<ApiResponse<any>>(res)
    return data.data
  },

  async create(body: {
    employeeId: string
    employeeName: string
    departmentId: string
    departmentName: string
    date: string
    reason?: string
  }) {
    const token = getAuthToken()
    if (!token) throw new Error('No authentication token found')
    const res = await fetch(`${LEAVE_SERVICE_URL}/api/v1/leaves`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
      body: JSON.stringify(body),
    })
    const data = await handleApiResponse<ApiResponse<any>>(res)
    return data.data
  },

  async update(id: string, body: { date?: string; reason?: string }) {
    const token = getAuthToken()
    if (!token) throw new Error('No authentication token found')
    const res = await fetch(`${LEAVE_SERVICE_URL}/api/v1/leaves/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
      body: JSON.stringify(body),
    })
    const data = await handleApiResponse<ApiResponse<any>>(res)
    return data.data
  },

  async remove(id: string) {
    const token = getAuthToken()
    if (!token) throw new Error('No authentication token found')
    const res = await fetch(`${LEAVE_SERVICE_URL}/api/v1/leaves/${id}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` },
    })
    const data = await handleApiResponse<ApiResponse<any>>(res)
    return data.data
  },

  async toggle(id: string) {
    const token = getAuthToken()
    if (!token) throw new Error('No authentication token found')
    const res = await fetch(`${LEAVE_SERVICE_URL}/api/v1/leaves/${id}/toggle`, {
      method: 'PUT',
      headers: { Authorization: `Bearer ${token}` },
    })
    const data = await handleApiResponse<ApiResponse<any>>(res)
    return data.data
  },

  async health() {
    const res = await fetch(`${LEAVE_SERVICE_URL}/health`)
    return handleApiResponse<any>(res)
  },
}

export default leaveService
