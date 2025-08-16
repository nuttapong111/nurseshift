import { getAuthToken } from './userService'
import { normalizeBaseUrl } from '@/lib/utils'

const API_BASE_URL = normalizeBaseUrl(process.env.NEXT_PUBLIC_DEPARTMENT_SERVICE_URL || process.env.NEXT_PUBLIC_DEPARTMENT_API_URL, 'http://localhost:8083')

export interface Department {
  id: string
  name: string
  description?: string
  head_user_id?: string
  max_nurses: number
  max_assistants: number
  settings?: string
  is_active: boolean
  created_by?: string
  created_at: string
  updated_at: string
}

export interface DepartmentWithStats extends Department {
  total_employees: number
  active_employees: number
  nurse_count: number
  assistant_count: number
}

export interface DepartmentStaff {
  id: string
  department_id: string
  name: string
  position: string
  phone?: string
  email?: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface CreateDepartmentRequest {
  name: string
  description?: string
  head_user_id?: string
  max_nurses: number
  max_assistants: number
  settings?: string
}

export interface UpdateDepartmentRequest {
  name?: string
  description?: string
  head_user_id?: string
  max_nurses?: number
  max_assistants?: number
  settings?: string
  is_active?: boolean
}

class DepartmentService {
  private async getHeaders(): Promise<HeadersInit> {
    const token = await getAuthToken()
    return {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    }
  }

  // Get all departments for the authenticated user
  async getDepartments(): Promise<DepartmentWithStats[]> {
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/departments/`, {
        method: 'GET',
        headers: await this.getHeaders(),
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const result = await response.json()
      return result.data || []
    } catch (error) {
      console.error('Error fetching departments:', error)
      throw error
    }
  }

  // Get specific department by ID
  async getDepartment(id: string): Promise<Department> {
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/departments/${id}`, {
        method: 'GET',
        headers: await this.getHeaders(),
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const result = await response.json()
      return result.data
    } catch (error) {
      console.error('Error fetching department:', error)
      throw error
    }
  }

  // Create new department
  async createDepartment(data: CreateDepartmentRequest): Promise<Department> {
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/departments/`, {
        method: 'POST',
        headers: await this.getHeaders(),
        body: JSON.stringify(data),
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`)
      }

      const result = await response.json()
      return result.data
    } catch (error) {
      console.error('Error creating department:', error)
      throw error
    }
  }

  // Update department
  async updateDepartment(id: string, data: UpdateDepartmentRequest): Promise<Department> {
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/departments/${id}`, {
        method: 'PUT',
        headers: await this.getHeaders(),
        body: JSON.stringify(data),
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`)
      }

      const result = await response.json()
      return result.data
    } catch (error) {
      console.error('Error updating department:', error)
      throw error
    }
  }

  // Delete department
  async deleteDepartment(id: string): Promise<void> {
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/departments/${id}`, {
        method: 'DELETE',
        headers: await this.getHeaders(),
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`)
      }
    } catch (error) {
      console.error('Error deleting department:', error)
      throw error
    }
  }

  // Get department staff
  async getDepartmentStaff(departmentId: string): Promise<DepartmentStaff[]> {
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/departments/${departmentId}/staff`, {
        method: 'GET',
        headers: await this.getHeaders(),
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const result = await response.json()
      console.log('getDepartmentStaff raw response:', result)
      console.log('getDepartmentStaff result.data:', result.data)
      console.log('getDepartmentStaff result.data.staff:', result.data?.staff)
      console.log('getDepartmentStaff result.data.staff type:', typeof result.data?.staff)
      console.log('getDepartmentStaff result.data.staff isArray:', Array.isArray(result.data?.staff))
      
      // Check if result.data.staff exists and is an array
      if (!result.data?.staff || !Array.isArray(result.data.staff)) {
        console.warn('getDepartmentStaff: result.data.staff is not an array, returning empty array')
        return []
      }
      
      return result.data.staff
    } catch (error) {
      console.error('Error fetching department staff:', error)
      throw error
    }
  }

  // Get department statistics
  async getDepartmentStats(departmentId: string): Promise<any> {
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/departments/${departmentId}/stats`, {
        method: 'GET',
        headers: await this.getHeaders(),
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const result = await response.json()
      return result.data
    } catch (error) {
      console.error('Error fetching department stats:', error)
      throw error
    }
  }

  // Add staff member to department
  async addDepartmentStaff(departmentId: string, data: {
    first_name: string
    last_name: string
    position: 'nurse' | 'assistant'
    phone?: string
    email?: string
  }): Promise<DepartmentStaff> {
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/departments/${departmentId}/staff`, {
        method: 'POST',
        headers: await this.getHeaders(),
        body: JSON.stringify({
          ...data,
          department_role: data.position // position และ department_role ต้องตรงกัน
        }),
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`)
      }

      const result = await response.json()
      return result.data
    } catch (error) {
      console.error('Error adding department staff:', error)
      throw error
    }
  }

  // Delete department staff member
  async deleteDepartmentStaff(departmentId: string, staffId: string): Promise<void> {
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/departments/${departmentId}/staff/${staffId}`, {
        method: 'DELETE',
        headers: await this.getHeaders(),
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`)
      }
    } catch (error) {
      console.error('Error deleting department staff:', error)
      throw error
    }
  }
}

export const departmentService = new DepartmentService()
