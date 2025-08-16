export interface PriorityItem {
  id: string
  name: string
  description?: string
  order: number
  isActive: boolean
  hasSettings?: boolean
  settingType?: string
  settingLabel?: string
  settingUnit?: string
  settingValue?: number
}

const BASE_URL = process.env.NEXT_PUBLIC_PRIORITY_SERVICE_URL || process.env.NEXT_PUBLIC_PRIORITY_API_URL || 'http://localhost:8086'

async function headers(): Promise<HeadersInit> {
  if (typeof window === 'undefined') return { 'Content-Type': 'application/json' }
  const token = localStorage.getItem('token')
  return {
    'Content-Type': 'application/json',
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
  }
}

export const priorityService = {
  async list(departmentId: string): Promise<{ priorities: PriorityItem[] }> {
    const qs = new URLSearchParams({ departmentId })
    const res = await fetch(`${BASE_URL}/api/v1/priorities?${qs.toString()}`, {
      method: 'GET',
      headers: await headers(),
    })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const body = await res.json()
    return body?.data || { priorities: [] }
  },

  async update(priorityId: string, payload: Partial<{ isActive: boolean; order: number }>): Promise<void> {
    const res = await fetch(`${BASE_URL}/api/v1/priorities/${priorityId}`, {
      method: 'PUT',
      headers: await headers(),
      body: JSON.stringify(payload),
    })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
  },

  async swap(aId: string, bId: string): Promise<void> {
    const res = await fetch(`${BASE_URL}/api/v1/priorities/swap`, {
      method: 'POST',
      headers: await headers(),
      body: JSON.stringify({ aId, bId }),
    })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
  },
}


