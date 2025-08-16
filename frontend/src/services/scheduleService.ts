import { getAuthToken } from './userService'

import { normalizeBaseUrl } from '@/lib/utils'
const API_BASE_URL = normalizeBaseUrl(process.env.NEXT_PUBLIC_SCHEDULE_SERVICE_URL || process.env.NEXT_PUBLIC_SCHEDULE_API_URL, 'http://localhost:8084')

export interface ScheduleItem {
  id: string
  departmentId: string
  userId: string
  shiftId: string
  scheduleDate: string
  status: string
  notes?: string | null
}

export interface ShiftDef {
  id: string
  departmentId: string
  name: string
  type: string
  startTime: string
  endTime: string
  requiredNurse: number
  requiredAsst: number
  color: string
}

export interface CalendarMetaDay {
  date: string
  isWorking: boolean
  isHoliday: boolean
}

class ScheduleService {
  private async headers(): Promise<HeadersInit> {
    const token = getAuthToken()
    if (!token) throw new Error('No authentication token found')
    return { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` }
  }

  async list(params: { departmentId?: string; month?: string }): Promise<ScheduleItem[]> {
    const qs = new URLSearchParams()
    if (params.departmentId) qs.set('departmentId', params.departmentId)
    if (params.month) qs.set('month', params.month)
    const res = await fetch(`${API_BASE_URL}/api/v1/schedules/?${qs.toString()}`, { headers: await this.headers() })
    if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`)
    const body = await res.json()
    return body.data || []
  }

  async getShifts(departmentId: string): Promise<ShiftDef[]> {
    const qs = new URLSearchParams({ departmentId })
    const res = await fetch(`${API_BASE_URL}/api/v1/schedules/shifts?${qs.toString()}`, { headers: await this.headers() })
    if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`)
    const body = await res.json()
    return body.data || []
  }

  async getAvailableStaff(params: { departmentId: string; date: string; shiftId: string }): Promise<Array<{ id: string; name: string; position: 'nurse'|'assistant' }>> {
    const qs = new URLSearchParams(params as any)
    const res = await fetch(`${API_BASE_URL}/api/v1/schedules/available-staff?${qs.toString()}`, { headers: await this.headers() })
    if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`)
    const body = await res.json()
    return body.data || []
  }

  async create(data: { departmentId: string; date: string }): Promise<ScheduleItem> {
    const res = await fetch(`${API_BASE_URL}/api/v1/schedules/`, {
      method: 'POST',
      headers: await this.headers(),
      body: JSON.stringify(data),
    })
    if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`)
    const body = await res.json()
    return body.data
  }

  async update(id: string, data: Partial<{ status: string; notes: string; shiftId: string }>): Promise<void> {
    const res = await fetch(`${API_BASE_URL}/api/v1/schedules/${id}`, {
      method: 'PUT',
      headers: await this.headers(),
      body: JSON.stringify(data),
    })
    if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`)
  }

  async remove(id: string): Promise<void> {
    const res = await fetch(`${API_BASE_URL}/api/v1/schedules/${id}`, {
      method: 'DELETE',
      headers: await this.headers(),
    })
    if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`)
  }

  async editShift(data: {
    departmentId: string
    date: string
    shiftId: string
    addNurses: string[]
    addAssistants: string[]
    removeNurses: string[]
    removeAssistants: string[]
  }): Promise<void> {
    const res = await fetch(`${API_BASE_URL}/api/v1/schedules/edit-shift`, {
      method: 'POST',
      headers: await this.headers(),
      body: JSON.stringify(data),
    })
    if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`)
  }

  async checkShiftOverlap(data: {
    departmentId: string
    date: string
    shiftId: string
    staffId: string
  }): Promise<{ canAssign: boolean; reason?: string }> {
    const res = await fetch(`${API_BASE_URL}/api/v1/schedules/check-overlap`, {
      method: 'POST',
      headers: await this.headers(),
      body: JSON.stringify(data),
    })
    if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`)
    const body = await res.json()
    return body.data
  }

  async autoGenerate(departmentId: string, month: string): Promise<{ inserted: number }> {
    const res = await fetch(`${API_BASE_URL}/api/v1/schedules/auto-generate`, {
      method: 'POST',
      headers: await this.headers(),
      body: JSON.stringify({ departmentId, month }),
    })
    if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`)
    const body = await res.json()
    return body.data
  }

  async aiGenerate(departmentId: string, month: string): Promise<{ inserted: number }> {
    const res = await fetch(`${API_BASE_URL}/api/v1/schedules/ai-generate`, {
      method: 'POST',
      headers: await this.headers(),
      body: JSON.stringify({ departmentId, month }),
    })
    if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`)
    const body = await res.json()
    return body.data
  }

  async calendarMeta(departmentId: string, month: string): Promise<CalendarMetaDay[]> {
    const qs = new URLSearchParams({ departmentId, month })
    const res = await fetch(`${API_BASE_URL}/api/v1/calendar-meta?${qs.toString()}`, {
      headers: await this.headers(),
    })
    if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`)
    const body = await res.json()
    return body.data || []
  }
}

export const scheduleService = new ScheduleService()


