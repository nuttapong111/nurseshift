import dayjs from 'dayjs'

const SETTING_SERVICE_URL = process.env.NEXT_PUBLIC_SETTING_SERVICE_URL || 'http://localhost:8085'

const getAuthToken = (): string | null => {
  if (typeof window !== 'undefined') {
    return localStorage.getItem('token')
  }
  return null
}

const request = async <T>(url: string, options: RequestInit = {}): Promise<T> => {
  const token = getAuthToken()
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  }
  if (token) headers['Authorization'] = `Bearer ${token}`

  const res = await fetch(url, { ...options, headers })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) {
    throw new Error(data?.message || `HTTP ${res.status}`)
  }
  return data as T
}

export type WorkingDayPayload = { idOrName: 'sunday'|'monday'|'tuesday'|'wednesday'|'thursday'|'friday'|'saturday'; enabled: boolean }

export const settingService = {
  async getSettings(departmentId: string): Promise<{
    workingDays: Array<{ dayOfWeek: number; isWorkingDay: boolean }>
    shifts: Array<{ id: string; name?: string; type?: string; startTime: string; endTime: string; requiredNurses: number; requiredAssistants: number; isActive: boolean }>
    holidays: Array<{ id: string; name: string; startDate: string; endDate: string }>
  }> {
    const url = `${SETTING_SERVICE_URL}/api/v1/settings?departmentId=${encodeURIComponent(departmentId)}`
    const res = await request<{ status: string; data: any }>(url)
    const d = res.data
    return {
      workingDays: (d?.workingDays || []).map((w: any) => ({
        dayOfWeek: w.dayOfWeek ?? w.DayOfWeek ?? w.day_of_week ?? 0,
        isWorkingDay: !!(w.isWorkingDay ?? w.IsWorkingDay ?? w.is_working_day)
      })),
      shifts: (d?.shifts || []).map((s: any) => ({
        id: s.id ?? s.ID ?? s.Id,
        name: s.name ?? s.Name,
        type: s.type ?? s.Type,
        startTime: s.startTime ?? s.StartTime ?? s.start_time ?? '07:00',
        endTime: s.endTime ?? s.EndTime ?? s.end_time ?? '15:00',
        requiredNurses: s.requiredNurses ?? s.RequiredNurses ?? s.required_nurses ?? 1,
        requiredAssistants: s.requiredAssistants ?? s.RequiredAssistants ?? s.required_assistants ?? 0,
        isActive: !!(s.isActive ?? s.IsActive ?? s.is_active),
      })),
      holidays: (d?.holidays || []).map((h: any) => ({
        id: h.id ?? h.ID ?? h.Id,
        name: h.name ?? h.Name,
        startDate: dayjs(h.startDate ?? h.StartDate ?? h.start_date).format('YYYY-MM-DD'),
        endDate: dayjs(h.endDate ?? h.EndDate ?? h.end_date).format('YYYY-MM-DD'),
      })),
    }
  },

  async updateWorkingDays(departmentId: string, workingDays: WorkingDayPayload[]): Promise<void> {
    const url = `${SETTING_SERVICE_URL}/api/v1/settings?departmentId=${encodeURIComponent(departmentId)}`
    await request(url, {
      method: 'PUT',
      body: JSON.stringify({ workingDays }),
    })
  },

  async createShift(input: {
    departmentId: string
    name: string
    type: string
    startTime: string
    endTime: string
    nurseCount: number
    assistantCount: number
    color?: string
    isActive?: boolean
  }): Promise<string> {
    const url = `${SETTING_SERVICE_URL}/api/v1/settings/shifts`
    const res = await request<{ status: string; id: string }>(url, {
      method: 'POST',
      body: JSON.stringify(input),
    })
    return res.id
  },

  async updateShift(shiftId: string, input: {
    name: string
    type: string
    startTime: string
    endTime: string
    nurseCount: number
    assistantCount: number
    color?: string
  }): Promise<void> {
    const url = `${SETTING_SERVICE_URL}/api/v1/settings/shifts/${encodeURIComponent(shiftId)}`
    await request(url, { method: 'PUT', body: JSON.stringify(input) })
  },

  async toggleShift(shiftId: string, enabled: boolean): Promise<void> {
    const url = `${SETTING_SERVICE_URL}/api/v1/settings/shifts/${encodeURIComponent(shiftId)}/toggle`
    await request(url, { method: 'PATCH', body: JSON.stringify({ enabled }) })
  },

  async deleteShift(shiftId: string): Promise<void> {
    const url = `${SETTING_SERVICE_URL}/api/v1/settings/shifts/${encodeURIComponent(shiftId)}`
    await request(url, { method: 'DELETE' })
  },

  async createHoliday(input: { departmentId: string; name: string; startDate: string; endDate: string; isRecurring?: boolean }): Promise<string> {
    const url = `${SETTING_SERVICE_URL}/api/v1/settings/holidays`
    const res = await request<{ status: string; id: string }>(url, { method: 'POST', body: JSON.stringify(input) })
    return res.id
  },

  async deleteHoliday(holidayId: string): Promise<void> {
    const url = `${SETTING_SERVICE_URL}/api/v1/settings/holidays/${encodeURIComponent(holidayId)}`
    await request(url, { method: 'DELETE' })
  },

  async updateHoliday(holidayId: string, input: { name: string; startDate: string; endDate: string; isRecurring?: boolean }): Promise<void> {
    const url = `${SETTING_SERVICE_URL}/api/v1/settings/holidays/${encodeURIComponent(holidayId)}`
    await request(url, { method: 'PUT', body: JSON.stringify(input) })
  },
}

export default settingService


