'use client'

import { useState, useEffect } from 'react'
import { 
  CalendarDaysIcon,
  UserGroupIcon,
  ClockIcon,
  BuildingOfficeIcon,
  PlusIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
  DocumentArrowDownIcon,
  AdjustmentsHorizontalIcon,
  ExclamationTriangleIcon,
  CheckIcon,
  MinusIcon,
  PencilIcon,
  EyeIcon,
  XMarkIcon,
  UserIcon,
  MinusCircleIcon,
  PlusCircleIcon
} from '@heroicons/react/24/outline'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card'
import { Button } from '@/components/ui/Button'
import { ButtonGroup } from '@/components/ui/ButtonGroup'
import dynamic from 'next/dynamic'
const ButtonToolbar = dynamic(() => import('./ButtonToolbar'), { ssr: false })
import { Input } from '@/components/ui/Input'
import DashboardLayout from '@/components/layout/DashboardLayout'
import Swal from 'sweetalert2'
import { scheduleService } from '@/services/scheduleService'
import { departmentService } from '@/services/departmentService'
import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'
import * as XLSX from 'xlsx'
import { saveAs } from 'file-saver'
import html2canvas from 'html2canvas'

// Types
interface Employee {
  id: string
  name: string
  position: string
  department: string
  shiftCounts: Record<string, number>
}

interface ShiftType {
  id: string // backend shift id
  name: string // backend shift name
  startTime: string
  endTime: string
  color: string // tailwind class by type
}

interface DaySchedule {
  date: string
  shifts: {
    [shiftId: string]: {
      name?: string
      startTime?: string
      endTime?: string
      color?: string
      nurses: Employee[]
      assistants: Employee[]
      requiredNurses: number
      requiredAssistants: number
    }
  }
}

interface DepartmentOption { id: string; name: string }

type Row = { staffId: string; name: string; position: 'พยาบาล'|'ผู้ช่วยพยาบาล'; perShift: Record<string, number>; total: number }

// Default shift types (fallback)
const defaultShiftTypes: ShiftType[] = [
  { id: 'morning', name: 'เวรเช้า', startTime: '07:00', endTime: '15:00', color: 'bg-yellow-100 text-yellow-800' },
  { id: 'afternoon', name: 'เวรบ่าย', startTime: '15:00', endTime: '23:00', color: 'bg-blue-100 text-blue-800' },
  { id: 'night', name: 'เวรดึก', startTime: '23:00', endTime: '07:00', color: 'bg-purple-100 text-purple-800' }
]

const mockDepartments: { id: string; name: string; employees: Employee[] }[] = [
  {
    id: 'emergency',
    name: 'แผนกฉุกเฉิน',
    employees: [
      { id: '1', name: 'สมใจ ใจดี', position: 'พยาบาล', department: 'แผนกฉุกเฉิน', shiftCounts: { morning: 8, afternoon: 6, night: 4, total: 18 } },
      { id: '2', name: 'วิไล รักษ์ดี', position: 'พยาบาล', department: 'แผนกฉุกเฉิน', shiftCounts: { morning: 7, afternoon: 5, night: 6, total: 18 } },
      { id: '3', name: 'มณี ใสใส', position: 'ผู้ช่วยพยาบาล', department: 'แผนกฉุกเฉิน', shiftCounts: { morning: 6, afternoon: 7, night: 5, total: 18 } },
      { id: '4', name: 'สุดา ขยันขัน', position: 'ผู้ช่วยพยาบาล', department: 'แผนกฉุกเฉิน', shiftCounts: { morning: 5, afternoon: 8, night: 5, total: 18 } }
    ]
  },
  {
    id: 'internal',
    name: 'แผนกอายุรกรรม',
    employees: [
      { id: '5', name: 'ประไพ แจ่มใส', position: 'พยาบาล', department: 'แผนกอายุรกรรม', shiftCounts: { morning: 9, afternoon: 5, night: 4, total: 18 } },
      { id: '6', name: 'นิรมล สุขใจ', position: 'พยาบาล', department: 'แผนกอายุรกรรม', shiftCounts: { morning: 6, afternoon: 7, night: 5, total: 18 } },
      { id: '7', name: 'บุปผา รื่นรมย์', position: 'ผู้ช่วยพยาบาล', department: 'แผนกอายุรกรรม', shiftCounts: { morning: 5, afternoon: 6, night: 7, total: 18 } }
    ]
  }
]

const monthNames = [
  'มกราคม', 'กุมภาพันธ์', 'มีนาคม', 'เมษายน', 'พฤษภาคม', 'มิถุนายน',
  'กรกฎาคม', 'สิงหาคม', 'กันยายน', 'ตุลาคม', 'พฤศจิกายน', 'ธันวาคม'
]

export default function SchedulePage() {
  const [currentDate, setCurrentDate] = useState(new Date())
  const [selectedDepartment, setSelectedDepartment] = useState<string>('')
  const [departments, setDepartments] = useState<DepartmentOption[]>([])
  const [schedules, setSchedules] = useState<{[key: string]: DaySchedule}>({})
  const [shiftTypes, setShiftTypes] = useState<ShiftType[]>(defaultShiftTypes)
  const [calendarMeta, setCalendarMeta] = useState<{[date: string]: { isWorking: boolean; isHoliday: boolean }}>({})
  const [mounted, setMounted] = useState(false)
  const [isGeneratingAI, setIsGeneratingAI] = useState(false)
  const [isGeneratingBackend, setIsGeneratingBackend] = useState(false)
  const [rawAssignments, setRawAssignments] = useState<any[]>([])

  useEffect(() => { setMounted(true) }, [])
  const [showReduceModal, setShowReduceModal] = useState(false)
  const [showEditShiftModal, setShowEditShiftModal] = useState(false)
  const [editingShift, setEditingShift] = useState<{
    date: string
    shiftId: string
    shiftName: string
    startTime: string
    endTime: string
    nurses: Employee[]
    assistants: Employee[]
    requiredNurses: number
    requiredAssistants: number
  } | null>(null)
  const [reduceForm, setReduceForm] = useState({
    date: '',
    shift: '',
    nursesToReduce: 0,
    assistantsToReduce: 0
  })
  const [employees, setEmployees] = useState<Employee[]>([])
  const [availableFromApi, setAvailableFromApi] = useState<Employee[]>([])

  // Stats (จริงจาก API)
  const [stats, setStats] = useState({
    totalNurses: 0,
    totalAssistants: 0,
    totalShifts: 0,
    totalDepartments: 0
  })

  // Load departments (real) once
  useEffect(() => {
    const loadDepartments = async () => {
      try {
        const list = await departmentService.getDepartments()
        const opts = (list || []).map(d => ({ id: d.id, name: d.name }))
        if (opts.length > 0) {
          setDepartments(opts)
          // default select OPD if exists, else first
          const opd = opts.find(o => o.name.toLowerCase().includes('opd'))
          setSelectedDepartment((opd || opts[0]).id)
          // สรุปจำนวนบุคลากรทั้งหมดจากทุกแผนก
          const totalNurses = (list || []).reduce((sum, d:any) => sum + (d.nurse_count || 0), 0)
          const totalAssistants = (list || []).reduce((sum, d:any) => sum + (d.assistant_count || 0), 0)
          setStats(prev => ({
            ...prev,
            totalNurses,
            totalAssistants,
            totalDepartments: opts.length
          }))
        } else {
          // fallback to mock
          setDepartments(mockDepartments.map(d => ({ id: d.id, name: d.name })))
          setSelectedDepartment(mockDepartments[0]?.id || '')
        }
      } catch (_e) {
        // fallback to mock on error
        setDepartments(mockDepartments.map(d => ({ id: d.id, name: d.name })))
        setSelectedDepartment(mockDepartments[0]?.id || '')
      }
    }
    loadDepartments()
  }, [])

  // โหลดจำนวนเวรทั้งหมดของเดือน (ตามแผนกที่เลือก)
  useEffect(() => {
    const loadMonthlyShifts = async () => {
      try {
        const month = `${currentDate.getFullYear()}-${String(currentDate.getMonth() + 1).padStart(2,'0')}`
        const params: { departmentId?: string; month: string } = { month }
        if (selectedDepartment) params.departmentId = selectedDepartment
        const all: any = await scheduleService.list(params)

        let total = 0
        if (Array.isArray(all)) {
          // กรณีเป็น array ของ assignments
          total = all.length
        } else if (all && typeof all === 'object') {
          // กรณีเป็น calendar map: รวมจำนวนพนักงานที่ถูกมอบหมายในทุกวัน/เวร
          Object.values(all as Record<string, any>).forEach((day: any) => {
            const shifts = day?.shifts || {}
            Object.values(shifts).forEach((s: any) => {
              total += (s.nurses?.length || 0) + (s.assistants?.length || 0)
            })
          })
        }
        setStats(prev => ({ ...prev, totalShifts: total }))
      } catch (_e) {
        setStats(prev => ({ ...prev, totalShifts: 0 }))
      }
    }
    loadMonthlyShifts()
  }, [currentDate, selectedDepartment])

  useEffect(() => {
    const load = async () => {
      if (!selectedDepartment) return
      try {
        // Load real shifts
        const shifts = await scheduleService.getShifts(selectedDepartment)
        // Assign distinct colors per shift using a palette (stable by index)
        const palette = [
          'bg-purple-100 text-purple-800',
          'bg-blue-100 text-blue-800',
          'bg-yellow-100 text-yellow-800',
          'bg-green-100 text-green-800',
          'bg-pink-100 text-pink-800',
          'bg-orange-100 text-orange-800',
          'bg-teal-100 text-teal-800',
          'bg-red-100 text-red-800'
        ]
        const types: ShiftType[] = shifts.map((s, idx) => {
          const color = palette[idx % palette.length]
          return { id: s.id, name: s.name, startTime: s.startTime, endTime: s.endTime, color }
        })
        setShiftTypes(types)

        // Load real schedules + calendar meta
        const month = formatMonth(currentDate)
        const [items, metaDays] = await Promise.all([
          scheduleService.list({ departmentId: selectedDepartment, month }),
          scheduleService.calendarMeta(selectedDepartment, month)
        ])
        
        // Build calendar meta map
        const metaMap: {[date: string]: { isWorking: boolean; isHoliday: boolean }} = {}
        ;(metaDays || []).forEach(m => { metaMap[m.date] = { isWorking: m.isWorking, isHoliday: m.isHoliday } })
        setCalendarMeta(metaMap)

        // Case 1: API returns array of assignments
        if (Array.isArray(items)) {
          const normalized = (items || []).map(normalizeAssignment)
          setRawAssignments(normalized)
          buildScheduleStructure(normalized, types)
          return
        }

        // Case 2: API returns calendar map keyed by date
        if (items && typeof items === 'object') {
          const newSchedules: Record<string, DaySchedule> = {}
          const flat: any[] = []
          const byShiftId: Record<string, ShiftType> = {}
          types.forEach(st => { byShiftId[st.id] = st })

          Object.entries(items as Record<string, any>).forEach(([dateKey, dayObj]) => {
            const shiftsObj: Record<string, any> = {}
            const dayShifts = (dayObj as any).shifts || {}
            Object.entries(dayShifts).forEach(([sid, s]: [string, any]) => {
              const nurses: Employee[] = (s.nurses || []).map((p: any) => ({
                id: p.staff_id || p.id,
                name: p.name,
                position: 'พยาบาล',
                department: '',
                shiftCounts: { morning: 0, afternoon: 0, night: 0, total: 0 }
              }))
              const assistants: Employee[] = (s.assistants || []).map((p: any) => ({
                id: p.staff_id || p.id,
                name: p.name,
                position: 'ผู้ช่วยพยาบาล',
                department: '',
                shiftCounts: { morning: 0, afternoon: 0, night: 0, total: 0 }
              }))

              // Push to flat assignments for summary/export compatibility
              nurses.forEach(n => flat.push({ scheduleDate: dateKey, shiftId: sid, departmentRole: 'nurse', userId: n.id, userName: n.name }))
              assistants.forEach(a => flat.push({ scheduleDate: dateKey, shiftId: sid, departmentRole: 'assistant', userId: a.id, userName: a.name }))

              // Use palette class color from shiftTypes (not backend hex) for consistent UI
              const classColor = byShiftId[sid]?.color || 'bg-blue-50 text-blue-800'

              shiftsObj[sid] = {
                name: s.name,
                startTime: s.startTime,
                endTime: s.endTime,
                color: classColor,
                nurses,
                assistants,
                requiredNurses: s.requiredNurses || 0,
                requiredAssistants: s.requiredAssistants || 0
              }
            })

            newSchedules[dateKey] = { date: dateKey, shifts: shiftsObj }
        })

        setSchedules(newSchedules)
          setRawAssignments(flat)
          return
        }

        // Fallback: empty
        setSchedules({})
        setRawAssignments([])
      } catch (error) {
        console.error('Failed to reload schedule data:', error)
      }
    }
    load()
  }, [selectedDepartment, currentDate])

  const generateMockSchedule = () => {
    const year = currentDate.getFullYear()
    const month = currentDate.getMonth()
    const daysInMonth = new Date(year, month + 1, 0).getDate()
    const newSchedules: {[key: string]: DaySchedule} = {}

    for (let day = 1; day <= daysInMonth; day++) {
      const dateKey = `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`
      
      newSchedules[dateKey] = {
        date: dateKey,
        shifts: {
          morning: {
            nurses: employees.filter(e => e.position === 'พยาบาล').slice(0, 2),
            assistants: employees.filter(e => e.position === 'ผู้ช่วยพยาบาล').slice(0, 1),
            requiredNurses: 2,
            requiredAssistants: 1
          },
          afternoon: {
            nurses: employees.filter(e => e.position === 'พยาบาล').slice(1, 3),
            assistants: employees.filter(e => e.position === 'ผู้ช่วยพยาบาล').slice(0, 2),
            requiredNurses: 2,
            requiredAssistants: 2
          },
          night: {
            nurses: employees.filter(e => e.position === 'พยาบาล').slice(0, 1),
            assistants: employees.filter(e => e.position === 'ผู้ช่วยพยาบาล').slice(0, 1),
            requiredNurses: 1,
            requiredAssistants: 1
          }
        }
      }
    }
    
    setSchedules(newSchedules)
  }

  const currentMonth = currentDate.getMonth()
  const currentYear = currentDate.getFullYear()

  const previousMonth = () => {
    setCurrentDate(new Date(currentYear, currentMonth - 1, 1))
  }

  const nextMonth = () => {
    setCurrentDate(new Date(currentYear, currentMonth + 1, 1))
  }

  const formatMonth = (d: Date) => `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`

  const toDateKey = (raw: any): string => {
    try {
      if (!raw) return ''
      if (raw instanceof Date) {
        const y = raw.getFullYear()
        const m = String(raw.getMonth() + 1).padStart(2, '0')
        const d = String(raw.getDate()).padStart(2, '0')
        return `${y}-${m}-${d}`
      }
      const s = String(raw)
      // Extract y, m, d tokens regardless of separators
      const mtx = s.match(/(\d{4})\D?(\d{1,2})\D?(\d{1,2})/)
      if (mtx) {
        const y = Number(mtx[1])
        const m = String(Number(mtx[2])).padStart(2, '0')
        const d = String(Number(mtx[3])).padStart(2, '0')
        return `${y}-${m}-${d}`
      }
      const d2 = new Date(s.replaceAll('/', '-'))
      if (!isNaN(d2.getTime())) return toDateKey(d2)
    } catch (_) {}
    return ''
  }

  const normalizeText = (t?: string) => (t || '').toString().trim().toLowerCase()

  // Normalize raw schedule assignment from API to a consistent shape
  const normalizeAssignment = (a: any) => {
    const scheduleDateRaw = a.scheduleDate || a.schedule_date || a.date || a.schedule_day
    const scheduleDate = toDateKey(scheduleDateRaw)
    const shiftId = String(a.shiftId || a.shift_id || a.shift || a.shiftTypeId || '')
    const shiftName = a.shiftName || a.shift_name || a.shiftTypeName || a.name || ''
    const startTime = a.startTime || a.start_time || ''
    const endTime = a.endTime || a.end_time || ''
    const departmentRoleRaw = a.departmentRole || a.department_role || a.role || a.staffRole
    const departmentRole = (departmentRoleRaw === 'assistant' || departmentRoleRaw === 'ผู้ช่วยพยาบาล') ? 'assistant' : 'nurse'
    const userId = a.staffId || a.staff_id || a.userId || a.user_id || ''
    const name = a.userName || a.staffName || [a.first_name, a.last_name].filter(Boolean).join(' ') || a.name || ''
    return {
      ...a,
      scheduleDate,
      shiftId,
      shiftName,
      startTime,
      endTime,
      departmentRole,
      staffId: userId,
      userId: userId,
      userName: name,
      staffName: name,
      staffRole: departmentRole
    }
  }

  // Function to reload schedule data
  const reloadScheduleData = async () => {
    if (!selectedDepartment) return
    try {
      // Load real shifts
      const shifts = await scheduleService.getShifts(selectedDepartment)
      // Assign distinct colors per shift using a palette (stable by index)
      const palette = [
        'bg-purple-100 text-purple-800',
        'bg-blue-100 text-blue-800',
        'bg-yellow-100 text-yellow-800',
        'bg-green-100 text-green-800',
        'bg-pink-100 text-pink-800',
        'bg-orange-100 text-orange-800',
        'bg-teal-100 text-teal-800',
        'bg-red-100 text-red-800'
      ]
      const types: ShiftType[] = shifts.map((s, idx) => {
        const color = palette[idx % palette.length]
        return { id: s.id, name: s.name, startTime: s.startTime, endTime: s.endTime, color }
      })
      setShiftTypes(types)

      // Load real schedules + calendar meta
      const month = formatMonth(currentDate)
      const [items, metaDays] = await Promise.all([
        scheduleService.list({ departmentId: selectedDepartment, month }),
        scheduleService.calendarMeta(selectedDepartment, month)
      ])
      
      // Build calendar meta map
      const metaMap: {[date: string]: { isWorking: boolean; isHoliday: boolean }} = {}
      ;(metaDays || []).forEach(m => { metaMap[m.date] = { isWorking: m.isWorking, isHoliday: m.isHoliday } })
      setCalendarMeta(metaMap)

      // Case 1: API returns array of assignments
      if (Array.isArray(items)) {
        const normalized = (items || []).map(normalizeAssignment)
        setRawAssignments(normalized)
        buildScheduleStructure(normalized, types)
        return
      }

      // Case 2: API returns calendar map keyed by date
      if (items && typeof items === 'object') {
        const newSchedules: Record<string, DaySchedule> = {}
        const flat: any[] = []
        const byShiftId: Record<string, ShiftType> = {}
        types.forEach(st => { byShiftId[st.id] = st })

        Object.entries(items as Record<string, any>).forEach(([dateKey, dayObj]) => {
          const shiftsObj: Record<string, any> = {}
          const dayShifts = (dayObj as any).shifts || {}
          Object.entries(dayShifts).forEach(([sid, s]: [string, any]) => {
            const nurses: Employee[] = (s.nurses || []).map((p: any) => ({
              id: p.staff_id || p.id,
              name: p.name,
              position: 'พยาบาล',
              department: '',
              shiftCounts: { morning: 0, afternoon: 0, night: 0, total: 0 }
            }))
            const assistants: Employee[] = (s.assistants || []).map((p: any) => ({
              id: p.staff_id || p.id,
              name: p.name,
              position: 'ผู้ช่วยพยาบาล',
              department: '',
              shiftCounts: { morning: 0, afternoon: 0, night: 0, total: 0 }
            }))

            // Push to flat assignments for summary/export compatibility
            nurses.forEach(n => flat.push({ scheduleDate: dateKey, shiftId: sid, departmentRole: 'nurse', userId: n.id, userName: n.name }))
            assistants.forEach(a => flat.push({ scheduleDate: dateKey, shiftId: sid, departmentRole: 'assistant', userId: a.id, userName: a.name }))

            // Use palette class color from shiftTypes (not backend hex) for consistent UI
            const classColor = byShiftId[sid]?.color || 'bg-blue-50 text-blue-800'

            shiftsObj[sid] = {
              name: s.name,
              startTime: s.startTime,
              endTime: s.endTime,
              color: classColor,
              nurses,
              assistants,
              requiredNurses: s.requiredNurses || 0,
              requiredAssistants: s.requiredAssistants || 0
            }
          })

          newSchedules[dateKey] = { date: dateKey, shifts: shiftsObj }
        })

        setSchedules(newSchedules)
        setRawAssignments(flat)
        return
      }

      // Fallback: empty
      setSchedules({})
      setRawAssignments([])
    } catch (error) {
      console.error('Failed to reload schedule data:', error)
    }
  }

  // Function to build schedule structure from raw assignments
  const buildScheduleStructure = (items: any[], types: ShiftType[]) => {
    const year = currentDate.getFullYear()
    const month = currentDate.getMonth()
    const daysInMonth = new Date(year, month + 1, 0).getDate()
    
    const newSchedules: Record<string, DaySchedule> = {}
    const byShiftId: Record<string, ShiftType> = {}
    types.forEach(st => { byShiftId[st.id] = st })

    // Initialize all days
    for (let day = 1; day <= daysInMonth; day++) {
      const dateKey = `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`
      newSchedules[dateKey] = {
        date: dateKey,
        shifts: {
          // initialize per backend shift id
          ...types.reduce((acc, st) => {
            const meta = byShiftId[st.id]
            acc[st.id] = { 
              name: meta?.name, 
              startTime: meta?.startTime, 
              endTime: meta?.endTime, 
              color: meta?.color, 
              nurses: [], 
              assistants: [], 
              requiredNurses: 0, 
              requiredAssistants: 0 
            }
            return acc
          }, {} as Record<string, any>)
        }
      }
    }

    // Populate with normalized assignments
    const normalizedItems = items.map(normalizeAssignment)
    if (process.env.NODE_ENV !== 'production' && normalizedItems.length) {
      console.debug('[schedule] normalized sample:', normalizedItems[0])
    }
    normalizedItems.forEach((item: any) => {
      const schedule = newSchedules[item.scheduleDate]
      if (!schedule) return

      const role = item.departmentRole
      // Try direct id first
      let foundKey: string | null = schedule.shifts[item.shiftId] ? item.shiftId : null
      let shift = foundKey ? schedule.shifts[foundKey] : undefined
      if (!shift) {
        // Try by name contains (normalized)
        const itemName = normalizeText(item.shiftName)
        const byName = Object.entries(schedule.shifts).find(([sid, s]) => {
          const nm = normalizeText((s as any).name)
          return itemName && nm && (nm === itemName || nm.includes(itemName) || itemName.includes(nm))
        })
        if (byName) { foundKey = byName[0]; shift = byName[1] as any }
      }
      if (!shift) {
        // Try by time window
        const byTime = Object.entries(schedule.shifts).find(([sid, s]) => (s as any).startTime === item.startTime && (s as any).endTime === item.endTime)
        if (byTime) { foundKey = byTime[0]; shift = byTime[1] as any }
      }
      if (!shift) {
        // Final fallback: first shift type
        const firstKey = Object.keys(schedule.shifts)[0]
        if (firstKey) { foundKey = firstKey; shift = schedule.shifts[firstKey] }
      }
      if (!shift) return

      const person = {
        id: item.staffId || item.userId || `${item.userName}|${role}`,
        name: item.userName || 'Unknown',
        position: (role === 'assistant' ? 'ผู้ช่วยพยาบาล' : 'พยาบาล') as 'ผู้ช่วยพยาบาล' | 'พยาบาล',
        department: 'Unknown',
        shiftCounts: { morning: 0, afternoon: 0, night: 0, total: 0 }
      }

      if (role === 'assistant') {
        shift.assistants.push(person)
      } else {
        shift.nurses.push(person)
      }

      if (process.env.NODE_ENV !== 'production') {
        console.debug('[schedule] mapped', item.scheduleDate, '->', foundKey, 'name=', (shift as any).name, 'role=', role, 'person=', person.name)
      }
    })

    setSchedules(newSchedules)
  }

  const handleAutoGenerateBackend = async () => {
    if (!selectedDepartment || isGeneratingBackend) return
    setIsGeneratingBackend(true)
    try {
      const res = await scheduleService.autoGenerate(selectedDepartment, formatMonth(currentDate))
      
      // Show a brief loading indicator for data refresh
      const refreshingSwal = Swal.fire({
        title: 'กำลังรีเฟรชข้อมูล...',
        html: '<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>',
        showConfirmButton: false,
        allowOutsideClick: false,
        allowEscapeKey: false
      })
      
      // Refresh data immediately after successful generation
      await reloadScheduleData()
      
      // Close loading and show success message
      Swal.close()
      await Swal.fire({ 
        icon: 'success', 
        title: 'สร้างตารางเวรอัตโนมัติสำเร็จ!', 
        text: `บันทึก ${res.inserted} รายการ และรีเฟรชข้อมูลแล้ว`, 
        confirmButtonColor: '#2563eb',
        timer: 3000,
        timerProgressBar: true
      })
      
    } catch (e: any) {
      await Swal.fire({ icon: 'error', title: 'ไม่สามารถสร้างตารางเวรได้', text: e?.message || 'เกิดข้อผิดพลาด', confirmButtonColor: '#2563eb' })
    } finally {
      setIsGeneratingBackend(false)
    }
  }

  const handleAutoGenerateAI = async () => {
    if (!selectedDepartment || isGeneratingAI) return
    setIsGeneratingAI(true)
    
    // Show loading with progress
    const loadingSwal = Swal.fire({
      title: 'กำลังสร้างตารางเวรด้วย AI',
      html: `
        <div class="flex flex-col items-center space-y-4">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <p class="text-sm text-gray-600">กรุณารอสักครู่... AI กำลังวิเคราะห์และจัดเวรที่เหมาะสม</p>
          <div class="w-full bg-gray-200 rounded-full h-2">
            <div class="bg-blue-600 h-2 rounded-full animate-pulse" style="width: 60%"></div>
          </div>
        </div>
      `,
      allowOutsideClick: false,
      allowEscapeKey: false,
      showConfirmButton: false,
      didOpen: () => {
        Swal.getPopup()?.classList.add('swal2-no-backdrop')
      }
    })

    try {
      const res = await scheduleService.aiGenerate(selectedDepartment, formatMonth(currentDate))
      Swal.close()
      
      // Show a brief loading indicator for data refresh
      const refreshingSwal = Swal.fire({
        title: 'กำลังรีเฟรชข้อมูล...',
        html: '<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>',
        showConfirmButton: false,
        allowOutsideClick: false,
        allowEscapeKey: false
      })
      
      // Refresh data immediately after successful generation
      await reloadScheduleData()
      
      // Close loading and show success message
      Swal.close()
      await Swal.fire({ 
        icon: 'success', 
        title: 'สร้างตารางเวรด้วย AI สำเร็จ!', 
        text: `บันทึก ${res.inserted} รายการ และรีเฟรชข้อมูลแล้ว`, 
        confirmButtonColor: '#2563eb',
        timer: 3000,
        timerProgressBar: true
      })
      
    } catch (e: any) {
      Swal.close()
      await Swal.fire({ icon: 'error', title: 'ไม่สามารถสร้างตารางเวรด้วย AI ได้', text: e?.message || 'เกิดข้อผิดพลาด', confirmButtonColor: '#2563eb' })
    } finally {
      setIsGeneratingAI(false)
    }
  }

  // Calculate total working hours for an employee based on actual shift assignments
  const calculateTotalHours = (employeeName: string): number => {
    let totalHours = 0
    
    // Get summary data to find shift counts for this employee
    const sum = calculateShiftSummary()
    const employeeRow = sum.rows?.find(r => r.name === employeeName)
    
    if (!employeeRow) return 0
    
    // Calculate hours from shift types and counts
    shiftTypes.forEach(shift => {
      const shiftId = shift.id
      const shiftCount = employeeRow.perShift[shiftId] || 0
      
      if (shiftCount > 0) {
        // Calculate hours based on shift time
        const startTime = shift.startTime.split(':')
        const endTime = shift.endTime.split(':')
        const startMinutes = parseInt(startTime[0]) * 60 + parseInt(startTime[1])
        let endMinutes = parseInt(endTime[0]) * 60 + parseInt(endTime[1])
        
        // Handle overnight shifts
        if (endMinutes < startMinutes) {
          endMinutes += 24 * 60
        }
        
        const shiftDuration = (endMinutes - startMinutes) / 60
        totalHours += shiftCount * shiftDuration
      }
    })
    
    return totalHours
  }

  // Add Thai font support to jsPDF
  const addThaiFont = (doc: jsPDF) => {
    try {
      // Try courier first as it has better Unicode support
      doc.setFont('courier')
    } catch (error) {
      console.warn('Could not set Thai-compatible font, falling back to default')
      doc.setFont('helvetica')
    }
  }

  // Alternative export using HTML to Canvas method
  const exportSummaryToImagePDF = async () => {
    try {
      // Show loading
      Swal.fire({
        title: 'กำลังสร้าง PDF...',
        html: '<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>',
        showConfirmButton: false,
        allowOutsideClick: false,
        allowEscapeKey: false
      })

      // Create a temporary table element with Thai content
      const tableData = calculateShiftSummary()
      const rows = tableData.rows || []
      const shiftIds = tableData.shiftIds || []
      const nameMap = tableData.shiftIdToName || {}

      // Create HTML table
      const tableHTML = `
        <div style="font-family: 'Sarabun', 'TH Sarabun PSK', Arial, sans-serif; padding: 20px; background: white;">
          <h2 style="text-align: center; margin-bottom: 10px;">รายงานสรุปเวร</h2>
          <p style="text-align: center; margin-bottom: 20px;">เดือน: ${currentDate.getFullYear()}/${String(currentDate.getMonth() + 1).padStart(2, '0')}</p>
          <table style="width: 100%; border-collapse: collapse; font-size: 12px;">
            <thead>
              <tr style="background-color: #3f51b5; color: white;">
                <th style="border: 1px solid #ddd; padding: 8px; text-align: left;">ชื่อ-นามสกุล</th>
                <th style="border: 1px solid #ddd; padding: 8px; text-align: left;">ตำแหน่ง</th>
                ${shiftIds.map(sid => `<th style="border: 1px solid #ddd; padding: 8px; text-align: center;">${nameMap[sid] || 'เวร'}</th>`).join('')}
                <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">รวม</th>
                <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">ชั่วโมงรวม</th>
              </tr>
            </thead>
            <tbody>
              ${rows.map((r, index) => `
                <tr style="background-color: ${index % 2 === 0 ? '#f5f5f5' : 'white'};">
                  <td style="border: 1px solid #ddd; padding: 8px;">${r.name || ''}</td>
                  <td style="border: 1px solid #ddd; padding: 8px;">${r.position || ''}</td>
                  ${shiftIds.map(sid => `<td style="border: 1px solid #ddd; padding: 8px; text-align: center;">${r.perShift[sid] || 0}</td>`).join('')}
                  <td style="border: 1px solid #ddd; padding: 8px; text-align: center; font-weight: bold;">${r.total || 0}</td>
                  <td style="border: 1px solid #ddd; padding: 8px; text-align: center; color: #2563eb; font-weight: bold;">${calculateTotalHours(r.name || '').toFixed(1)} ชม.</td>
                </tr>
              `).join('')}
            </tbody>
          </table>
        </div>
      `

      // Create temporary div
      const tempDiv = document.createElement('div')
      tempDiv.innerHTML = tableHTML
      tempDiv.style.position = 'absolute'
      tempDiv.style.left = '-9999px'
      tempDiv.style.width = '1200px'
      document.body.appendChild(tempDiv)

      // Convert to canvas
      const canvas = await html2canvas(tempDiv, {
        backgroundColor: 'white',
        scale: 2, // Higher quality
        useCORS: true
      })

      // Remove temp div
      document.body.removeChild(tempDiv)

      // Create PDF with image
      const doc = new jsPDF('portrait', 'mm', 'a4')
      const imgData = canvas.toDataURL('image/png')
      
      const imgWidth = 190 // A4 width minus margins
      const imgHeight = (canvas.height * imgWidth) / canvas.width
      
      doc.addImage(imgData, 'PNG', 10, 10, imgWidth, imgHeight)
      doc.save(`สรุปเวร_${currentDate.getFullYear()}_${String(currentDate.getMonth() + 1).padStart(2, '0')}.pdf`)

      Swal.close()
      
    } catch (error: any) {
      await Swal.fire({
        icon: 'error',
        title: 'ไม่สามารถสร้าง PDF ได้',
        text: error?.message || 'เกิดข้อผิดพลาด',
        confirmButtonColor: '#2563eb'
      })
    }
  }

  // Export functions with Thai support
  const exportSummaryToPDF = () => {
    const doc = new jsPDF()
    
    // Try to add Thai font support
    addThaiFont(doc)
    
    // Use mixed language titles
    doc.setFontSize(16)
    doc.text('รายงานสรุปเวร (Shift Summary Report)', 20, 20)
    doc.setFontSize(12)
    doc.text(`เดือน: ${currentDate.getFullYear()}/${String(currentDate.getMonth() + 1).padStart(2, '0')}`, 20, 30)
    
    // Get data from calculateShiftSummary (same as web display)
    const sum = calculateShiftSummary()
    const rows = sum.rows || []
    const shiftIds: string[] = sum.shiftIds || []
    const nameMap = sum.shiftIdToName || {}
    
    // Prepare headers with Thai language
    const headers = ['ชื่อ-นามสกุล', 'ตำแหน่ง', ...shiftIds.map(sid => nameMap[sid] || 'เวร'), 'รวม', 'ชั่วโมงรวม']
    
    // Prepare data from the same source as web table - keep original Thai data
    const data = rows.map(r => {
      // Calculate hours for this employee
      const totalHours = calculateTotalHours(r.name || '').toFixed(1)
      
      return [
        r.name || '',
        r.position || '',
        ...shiftIds.map(sid => r.perShift[sid] || 0),
        r.total || 0,
        totalHours + ' ชม.'
      ]
    })
    
    autoTable(doc, {
      head: [headers],
      body: data,
      startY: 40,
      styles: {
        font: 'courier', // Use courier for better Unicode support
        fontSize: 10,
        cellPadding: 3,
        lineColor: [200, 200, 200],
        lineWidth: 0.5
      },
      headStyles: {
        fillColor: [63, 81, 181],
        textColor: 255,
        fontStyle: 'bold',
        font: 'courier'
      },
      alternateRowStyles: {
        fillColor: [245, 245, 245]
      },
      margin: { top: 50, left: 10, right: 10 },
      columnStyles: {
        0: { cellWidth: 45 }, // Name column wider for Thai text
        1: { cellWidth: 30 }  // Position column
      }
    })
    
    doc.save(`สรุปเวร_${currentDate.getFullYear()}_${String(currentDate.getMonth() + 1).padStart(2, '0')}.pdf`)
  }

  const exportSummaryToExcel = () => {
    // Create workbook 
    const wb = XLSX.utils.book_new()
    
    // Get data from calculateShiftSummary (same as web display)
    const sum = calculateShiftSummary()
    const rows = sum.rows || []
    const shiftIds: string[] = sum.shiftIds || []
    const nameMap = sum.shiftIdToName || {}
    
    // Prepare headers - keep Thai for Excel (better support)
    const headers = ['ชื่อ-นามสกุล', 'ตำแหน่ง', ...shiftIds.map(sid => nameMap[sid] || 'เวร'), 'รวม', 'ชั่วโมงรวม']
    
    // Prepare data from the same source as web table
    const data = rows.map(r => {
      // Calculate hours for this employee
      const totalHours = parseFloat(calculateTotalHours(r.name || '').toFixed(1))
      
      return [
        r.name || '',
        r.position || '',
        ...shiftIds.map(sid => r.perShift[sid] || 0),
        r.total || 0,
        totalHours
      ]
    })
    
    // Create worksheet
    const ws = XLSX.utils.aoa_to_sheet([headers, ...data])
    
    // Set column widths
    const wscols = [
      { wch: 20 }, // Name
      { wch: 15 }, // Position
      ...shiftIds.map(() => ({ wch: 12 })), // Shift columns
      { wch: 10 }, // Total
      { wch: 12 }  // Total Hours
    ]
    ws['!cols'] = wscols
    
    // Add worksheet to workbook
    XLSX.utils.book_append_sheet(wb, ws, 'สรุปเวร')
    
    // Save file
    const fileName = `สรุปเวร_${currentDate.getFullYear()}_${String(currentDate.getMonth() + 1).padStart(2, '0')}.xlsx`
    XLSX.writeFile(wb, fileName)
  }

  // Calendar PDF with HTML to Canvas method (Thai support)
  const exportCalendarToImagePDF = async () => {
    try {
      // Show loading
      Swal.fire({
        title: 'กำลังสร้าง PDF ตารางเวร...',
        html: '<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>',
        showConfirmButton: false,
        allowOutsideClick: false,
        allowEscapeKey: false
      })

      const year = currentDate.getFullYear()
      const month = currentDate.getMonth()
      const daysInMonth = new Date(year, month + 1, 0).getDate()
      
      // Get data from calculateShiftSummary
      const sum = calculateShiftSummary()
      const rows = sum.rows || []

      // Create HTML table for calendar
      const calendarHTML = `
        <div style="font-family: 'Sarabun', 'TH Sarabun PSK', Arial, sans-serif; padding: 20px; background: white; min-width: 1400px;">
          <h2 style="text-align: center; margin-bottom: 10px;">ตารางเวรประจำเดือน</h2>
          <p style="text-align: center; margin-bottom: 20px;">เดือน: ${year}/${String(month + 1).padStart(2, '0')}</p>
          <table style="width: 100%; border-collapse: collapse; font-size: 10px;">
            <thead>
              <tr style="background-color: #3f51b5; color: white;">
                <th style="border: 1px solid #ddd; padding: 6px; text-align: center; width: 40px;">ลำดับ</th>
                <th style="border: 1px solid #ddd; padding: 6px; text-align: left; width: 150px;">ชื่อ-นามสกุล</th>
                ${Array.from({length: daysInMonth}, (_, i) => `<th style="border: 1px solid #ddd; padding: 4px; text-align: center; width: 35px;">${i + 1}</th>`).join('')}
              </tr>
            </thead>
            <tbody>
              ${rows.map((r, index) => {
                const makeDayCell = (dayIndex: number) => {
                  const day = dayIndex + 1
                  const dateStr = `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`
                  const daySchedule = schedules[dateStr]
                  let shiftsText = ''
                  if (daySchedule?.shifts) {
                    const targetName = (r.name || '').trim()
                    const shiftNames = Object.values(daySchedule.shifts).reduce((acc: string[], shift: any) => {
                      const inNurses = (shift.nurses || []).some((n: any) => (n.name || '').trim() === targetName)
                      const inAsst = (shift.assistants || []).some((a: any) => (a.name || '').trim() === targetName)
                      if (inNurses || inAsst) acc.push(shift.name || 'เวร')
                      return acc
                    }, [])
                    shiftsText = shiftNames.join('<br/>')
                  }
                  return `<td style="border: 1px solid #ddd; padding: 2px; text-align: center; font-size: 8px; line-height: 1.2; vertical-align: top;">${shiftsText}</td>`
                }
                return `
                  <tr style="background-color: ${index % 2 === 0 ? '#f5f5f5' : 'white'};">
                    <td style="border: 1px solid #ddd; padding: 4px; text-align: center;">${index + 1}</td>
                    <td style="border: 1px solid #ddd; padding: 4px; text-align: left;">${r.name || ''}</td>
                    ${Array.from({length: daysInMonth}, (_, dayIndex) => makeDayCell(dayIndex)).join('')}
                  </tr>
                `
              }).join('')}
            </tbody>
          </table>
        </div>
      `

      // Create temporary div
      const tempDiv = document.createElement('div')
      tempDiv.innerHTML = calendarHTML
      tempDiv.style.position = 'absolute'
      tempDiv.style.left = '-9999px'
      tempDiv.style.width = '1400px'
      document.body.appendChild(tempDiv)

      // Convert to canvas
      const canvas = await html2canvas(tempDiv, {
        backgroundColor: 'white',
        scale: 1.5, // Smaller scale for calendar due to width
        useCORS: true,
        allowTaint: true
      })

      // Remove temp div
      document.body.removeChild(tempDiv)

      // Create PDF with image (landscape for calendar)
      const doc = new jsPDF('landscape', 'mm', 'a4')
      const imgData = canvas.toDataURL('image/png')
      
      const pageWidth = 297 // A4 landscape width
      const pageHeight = 210 // A4 landscape height
      const imgWidth = pageWidth - 20 // margins
      const imgHeight = (canvas.height * imgWidth) / canvas.width
      
      // If image is too tall, scale it down
      if (imgHeight > pageHeight - 20) {
        const scaledHeight = pageHeight - 20
        const scaledWidth = (canvas.width * scaledHeight) / canvas.height
        doc.addImage(imgData, 'PNG', 10, 10, scaledWidth, scaledHeight)
      } else {
        doc.addImage(imgData, 'PNG', 10, 10, imgWidth, imgHeight)
      }
      
      doc.save(`ตารางเวร_${year}_${String(month + 1).padStart(2, '0')}.pdf`)

      Swal.close()
      
    } catch (error: any) {
      await Swal.fire({
        icon: 'error',
        title: 'ไม่สามารถสร้าง PDF ตารางเวรได้',
        text: error?.message || 'เกิดข้อผิดพลาด',
        confirmButtonColor: '#2563eb'
      })
    }
  }

  const exportCalendarToPDF = () => {
    const doc = new jsPDF('landscape', 'mm', 'a4')
    
    // Try to add Thai font support
    addThaiFont(doc)
    
    // Title with Thai
    doc.setFontSize(16)
    doc.text('ตารางเวรประจำเดือน (Monthly Shift Schedule)', 20, 20)
    doc.setFontSize(12)
    doc.text(`เดือน: ${currentDate.getFullYear()}/${String(currentDate.getMonth() + 1).padStart(2, '0')}`, 20, 30)
    
    // Get days in month
    const year = currentDate.getFullYear()
    const month = currentDate.getMonth()
    const daysInMonth = new Date(year, month + 1, 0).getDate()
    
    // Get data from calculateShiftSummary to match web display
    const sum = calculateShiftSummary()
    const rows = sum.rows || []
    
    // Headers with Thai
    const headers = ['ลำดับ', 'ชื่อ-นามสกุล', ...Array.from({length: daysInMonth}, (_, i) => (i + 1).toString())]
    
    // Prepare data using the same employees from summary
    const data = rows.map((r, index) => {
      const row = [index + 1, r.name || '']
      
      // Find the actual employee data
      const emp = employees.find(e => e.name === r.name)
      
      for (let day = 1; day <= daysInMonth; day++) {
        const dateStr = `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`
        const daySchedule = schedules[dateStr]
        let shiftsText = ''
        
        if (daySchedule?.shifts && emp) {
          Object.values(daySchedule.shifts).forEach(shift => {
            const foundNurse = shift.nurses.find(n => n.id === emp.id)
            const foundAssistant = shift.assistants.find(a => a.id === emp.id)
            if (foundNurse || foundAssistant) {
              shiftsText += (shift.name?.substring(0, 3) || 'X') + ' '
            }
          })
        }
        row.push(shiftsText.trim())
      }
      return row
    })
    
    autoTable(doc, {
      head: [headers],
      body: data,
      startY: 40,
      styles: {
        font: 'courier', // Use courier for better Unicode support
        fontSize: 6,
        cellPadding: 1,
        halign: 'center',
        valign: 'middle'
      },
      headStyles: {
        fillColor: [63, 81, 181],
        textColor: 255,
        fontSize: 7,
        fontStyle: 'bold',
        font: 'courier'
      },
      columnStyles: {
        0: { cellWidth: 12, halign: 'center' },
        1: { cellWidth: 40, halign: 'left' } // Wider for Thai names
      },
      margin: { left: 5, right: 5 }
    })
    
    doc.save(`ตารางเวร_${currentDate.getFullYear()}_${String(currentDate.getMonth() + 1).padStart(2, '0')}.pdf`)
  }

  const exportCalendarToExcel = () => {
    const year = currentDate.getFullYear()
    const month = currentDate.getMonth()
    const daysInMonth = new Date(year, month + 1, 0).getDate()
    
    // Create workbook
    const wb = XLSX.utils.book_new()
    
    // Get data from calculateShiftSummary to match web display
    const sum = calculateShiftSummary()
    const rows = sum.rows || []
    
    // Headers with Thai support for Excel
    const headers = ['ลำดับ', 'ชื่อ-นามสกุล', ...Array.from({length: daysInMonth}, (_, i) => (i + 1).toString())]
    
    // Prepare data using the same employees from summary
    const data = rows.map((r, index) => {
      const row = [index + 1, r.name || '']
      
      for (let day = 1; day <= daysInMonth; day++) {
        const dateStr = `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`
        const daySchedule = schedules[dateStr]
        let shiftsText = ''
        
        if (daySchedule?.shifts) {
          const targetName = (r.name || '').trim()
          const shiftNames: string[] = []
          Object.values(daySchedule.shifts).forEach((shift: any) => {
            const inNurses = (shift.nurses || []).some((n: any) => (n.name || '').trim() === targetName)
            const inAsst = (shift.assistants || []).some((a: any) => (a.name || '').trim() === targetName)
            if (inNurses || inAsst) {
              shiftNames.push(shift.name || 'เวร')
            }
          })
          shiftsText = shiftNames.join('\n')
        }
        row.push(shiftsText.trim())
      }
      return row
    })
    
    // Create worksheet
    const ws = XLSX.utils.aoa_to_sheet([headers, ...data])
    
    // Set column widths - wider for Thai shift names
    const wscols = [
      { wch: 5 },  // No.
      { wch: 25 }, // Name (wider for Thai names)
      ...Array.from({length: daysInMonth}, () => ({ wch: 15 })) // Days (wider for Thai shift names)
    ]
    ws['!cols'] = wscols
    
    // Set row height for wrap text
    const range = XLSX.utils.decode_range(ws['!ref'] || 'A1')
    for (let R = range.s.r + 1; R <= range.e.r; R++) {
      for (let C = 2; C <= range.e.c; C++) { // Start from day columns
        const cellRef = XLSX.utils.encode_cell({ r: R, c: C })
        if (ws[cellRef]) {
          if (!ws[cellRef].s) ws[cellRef].s = {}
          ws[cellRef].s.alignment = { wrapText: true, vertical: 'top' }
        }
      }
    }
    
    // Add worksheet to workbook
    XLSX.utils.book_append_sheet(wb, ws, 'ตารางเวร')
    
    // Save file
    const fileName = `ตารางเวร_${currentDate.getFullYear()}_${String(currentDate.getMonth() + 1).padStart(2, '0')}.xlsx`
    XLSX.writeFile(wb, fileName)
  }

  // Build dynamic summary columns by backend shifts (use shiftTypes directly)
  const calculateShiftSummary = () => {
    const shiftIdToName: Record<string, string> = {}
    shiftTypes.forEach(st => { shiftIdToName[st.id] = st.name })
    const shiftIds = Object.keys(shiftIdToName)

    // aggregate by staffId (ไม่ใช้ชื่อ เพื่อเลี่ยงปัญหาชื่อซ้ำข้ามตำแหน่ง)
    const arr = Array.isArray(rawAssignments) ? rawAssignments : []
    const map: Record<string, Row> = {}
    arr.forEach((a: any) => {
      const staffId = a.staffId || a.userId || `${a.userName}|${a.departmentRole}`
      const name = a.staffName || a.userName || '-'
      const role = a.staffRole ? (a.staffRole.includes('ช่วย') || a.staffRole === 'assistant' ? 'ผู้ช่วยพยาบาล' : 'พยาบาล')
                               : (a.departmentRole === 'assistant' ? 'ผู้ช่วยพยาบาล' : 'พยาบาล')
      if (!map[staffId]) map[staffId] = { staffId, name, position: role, perShift: {}, total: 0 }
      map[staffId].perShift[a.shiftId] = (map[staffId].perShift[a.shiftId] || 0) + 1
      map[staffId].total++
    })

    const rows = Object.values(map).map(r => {
      shiftIds.forEach(sid => { if (!(sid in r.perShift)) r.perShift[sid] = 0 })
      return r
    }).sort((a,b)=>a.name.localeCompare(b.name))

    return { rows, shiftIds, shiftIdToName }
  }

  const handleEditShift = (date: string, shiftId: string) => {
    const daySchedule = schedules[date]
    if (!daySchedule?.shifts?.[shiftId]) return

    const shift = daySchedule.shifts[shiftId]
    setEditingShift({
      date,
      shiftId,
      shiftName: shift.name || '',
      startTime: shift.startTime || '',
      endTime: shift.endTime || '',
      nurses: shift.nurses || [],
      assistants: shift.assistants || [],
      requiredNurses: shift.requiredNurses || 0,
      requiredAssistants: shift.requiredAssistants || 0
    })
    setShowEditShiftModal(true)
    ;(async () => {
      try {
        const list = await scheduleService.getAvailableStaff({ departmentId: selectedDepartment, date, shiftId })
        setAvailableFromApi(list.map(i => ({ id: i.id, name: i.name, position: i.position, department: '', shiftCounts: {} as any })))
      } catch {
        setAvailableFromApi([])
      }
    })()
  }

  const handleReduceStaff = async () => {
    setShowReduceModal(true)
  }

  const processStaffReduction = async () => {
    if (!reduceForm.date || !reduceForm.shift) {
      await Swal.fire({
        icon: 'warning',
        title: 'กรุณากรอกข้อมูลให้ครบ',
        confirmButtonColor: '#2563eb'
      })
      return
    }

    const schedule = schedules[reduceForm.date]
    if (!schedule) return

    const shift = schedule.shifts[reduceForm.shift]
    
    // Sort staff by total shift count (descending) to remove those with most shifts first
    const sortedNurses = shift.nurses.sort((a, b) => b.shiftCounts.total - a.shiftCounts.total)
    const sortedAssistants = shift.assistants.sort((a, b) => b.shiftCounts.total - a.shiftCounts.total)
    
    const nursesToKeep = Math.max(0, shift.nurses.length - reduceForm.nursesToReduce)
    const assistantsToKeep = Math.max(0, shift.assistants.length - reduceForm.assistantsToReduce)
    
    const newNurses = sortedNurses.slice(0, nursesToKeep)
    const newAssistants = sortedAssistants.slice(0, assistantsToKeep)
    
    const removedNurses = sortedNurses.slice(nursesToKeep)
    const removedAssistants = sortedAssistants.slice(assistantsToKeep)

    // Update schedule
    const updatedSchedules = {
      ...schedules,
      [reduceForm.date]: {
        ...schedule,
        shifts: {
          ...schedule.shifts,
          [reduceForm.shift]: {
            ...shift,
            nurses: newNurses,
            assistants: newAssistants
          }
        }
      }
    }

    try {
      // Show loading
      Swal.fire({
        title: 'กำลังปรับลดพนักงาน...',
        html: '<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>',
        showConfirmButton: false,
        allowOutsideClick: false,
        allowEscapeKey: false
      })
      
      // Find schedule items to remove from API
      const scheduleItems = rawAssignments.filter(item => 
        item.scheduleDate === reduceForm.date && item.shiftId === reduceForm.shift
      )
      
      const nursesToRemove = scheduleItems.filter(item => 
        removedNurses.some(nurse => nurse.id === item.userId)
      )
      const assistantsToRemove = scheduleItems.filter(item => 
        removedAssistants.some(assistant => assistant.id === item.userId)
      )
      
      const itemsToRemove = [...nursesToRemove, ...assistantsToRemove]
      
      // Remove items via API
      for (const item of itemsToRemove) {
        await scheduleService.remove(item.id)
      }
      
      // Reload data to reflect changes
      await reloadScheduleData()
      
      // Show success message
      const removedCount = itemsToRemove.length
    await Swal.fire({
      icon: 'success',
        title: 'ปรับลดพนักงานสำเร็จ',
      html: `
        <div class="text-left">
          <p><strong>พนักงานที่ถูกลดออก:</strong></p>
          ${removedNurses.length > 0 ? `<p>พยาบาล: ${removedNurses.map(n => n.name).join(', ')}</p>` : ''}
          ${removedAssistants.length > 0 ? `<p>ผู้ช่วยพยาบาล: ${removedAssistants.map(a => a.name).join(', ')}</p>` : ''}
            <p class="mt-2"><strong>รวม:</strong> ลดออก ${removedCount} คน และรีเฟรชข้อมูลแล้ว</p>
        </div>
      `,
        confirmButtonColor: '#2563eb',
        timer: 4000,
        timerProgressBar: true
      })
      
      setShowReduceModal(false)
      setReduceForm({ date: '', shift: '', nursesToReduce: 0, assistantsToReduce: 0 })
      
    } catch (error: any) {
      await Swal.fire({
        icon: 'error',
        title: 'ไม่สามารถปรับลดพนักงานได้',
        text: error?.message || 'เกิดข้อผิดพลาด',
      confirmButtonColor: '#2563eb'
    })
    }
  }

  // Generate calendar
  const firstDayOfMonth = new Date(currentYear, currentMonth, 1)
  const lastDayOfMonth = new Date(currentYear, currentMonth + 1, 0)
  const daysInMonth = lastDayOfMonth.getDate()
  const startingDayOfWeek = firstDayOfMonth.getDay()

  const calendar = []
  for (let i = 0; i < startingDayOfWeek; i++) {
    calendar.push(null)
  }
  for (let day = 1; day <= daysInMonth; day++) {
    calendar.push(day)
  }

  if (!mounted) {
    return (
      <DashboardLayout>
        <div className="p-6 animate-pulse text-gray-500">กำลังโหลด...</div>
      </DashboardLayout>
    )
  }

  return (
    <DashboardLayout>
      <div className="space-y-8">
        {/* Header */}
        <div>
          <h1 className="text-3xl font-bold text-gray-900">ตารางเวรประจำเดือน</h1>
          <p className="text-gray-600 mt-2">
            จัดการและติดตามตารางเวรของระบบ
          </p>
        </div>

        {/* Stats Cards - ข้อมูลสรุปรวมทั้งระบบ */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <Card className="border-0 shadow-md hover:shadow-lg transition-shadow">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">จำนวนพยาบาลทั้งหมด</CardTitle>
              <UserGroupIcon className="h-4 w-4 text-blue-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-blue-600">{stats.totalNurses}</div>
              <p className="text-xs text-muted-foreground">ในระบบทั้งหมด</p>
            </CardContent>
          </Card>

          <Card className="border-0 shadow-md hover:shadow-lg transition-shadow">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">จำนวนผู้ช่วยพยาบาลทั้งหมด</CardTitle>
              <UserGroupIcon className="h-4 w-4 text-green-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-green-600">{stats.totalAssistants}</div>
              <p className="text-xs text-muted-foreground">ในระบบทั้งหมด</p>
            </CardContent>
          </Card>

          <Card className="border-0 shadow-md hover:shadow-lg transition-shadow">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">จำนวนเวรที่สร้างทั้งหมด</CardTitle>
              <ClockIcon className="h-4 w-4 text-purple-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-purple-600">{stats.totalShifts}</div>
              <p className="text-xs text-muted-foreground">เวรทั้งหมดในเดือนนี้</p>
            </CardContent>
          </Card>

          <Card className="border-0 shadow-md hover:shadow-lg transition-shadow">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">จำนวนแผนกทั้งหมด</CardTitle>
              <BuildingOfficeIcon className="h-4 w-4 text-orange-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-orange-600">{stats.totalDepartments}</div>
              <p className="text-xs text-muted-foreground">แผนกที่คุณจัดการ</p>
            </CardContent>
          </Card>
        </div>

        {/* Department Selection - ข้อมูลเฉพาะแผนก */}
        <div className="space-y-4">
          <div>
            <h2 className="text-xl font-semibold text-gray-900">
                ข้อมูลแผนก: {departments.find(d => d.id === selectedDepartment)?.name || '—'}
            </h2>
            <p className="text-gray-600 mt-1">
              จัดการและติดตามตารางเวรของแผนกที่เลือก
            </p>
          </div>
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between space-y-4 sm:space-y-0 sm:space-x-4">
            <select
              value={selectedDepartment}
              onChange={(e) => setSelectedDepartment(e.target.value)}
              className="border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 w-full sm:w-auto"
            >
                {departments.map(dept => (
                <option key={dept.id} value={dept.id}>{dept.name}</option>
              ))}
            </select>
            <div className="flex flex-col sm:flex-row space-y-2 sm:space-y-0 sm:space-x-2">
              <div className="flex flex-col sm:flex-row sm:items-center sm:space-x-2 w-full">
                {/* <Button onClick={handleReduceStaff} className="w-full sm:w-auto">ปรับลดพนักงาน</Button>
                <Button 
                  onClick={() => {
                    Swal.fire({
                      title: 'แก้ไขเวร',
                      text: 'ฟีเจอร์นี้กำลังอยู่ในระหว่างการพัฒนา',
                      icon: 'info',
                      confirmButtonColor: '#2563eb'
                    })
                  }} 
                  className="w-full sm:w-auto bg-orange-500 hover:bg-orange-600"
                >
                  <PencilIcon className="w-4 h-4 mr-1" />
                  แก้ไขเวร
                </Button> */}
                <Button 
                  onClick={handleAutoGenerateBackend} 
                  disabled={isGeneratingBackend || isGeneratingAI}
                  className="w-full sm:w-auto"
                >
                  {isGeneratingBackend ? (
                    <>
                      <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                      กำลังสร้าง...
                    </>
                  ) : (
                    'สร้างอัตโนมัติ (Backend)'
                  )}
                </Button>
                <Button 
                  variant="outline" 
                  onClick={handleAutoGenerateAI} 
                  disabled={isGeneratingAI || isGeneratingBackend}
                  className="w-full sm:w-auto"
                >
                  {isGeneratingAI ? (
                    <>
                      <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600 mr-2"></div>
                      AI กำลังวิเคราะห์...
                    </>
                  ) : (
                    'สร้างอัตโนมัติ (AI)'
                  )}
                </Button>
              </div>
            </div>
          </div>
        </div>

        {/* Shift Summary Table - dynamic columns by backend shifts */}
        <Card className="border-0 shadow-md">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
            <CardTitle className="flex items-center">
              <UserGroupIcon className="w-5 h-5 mr-2" />
              สรุปจำนวนเวรของพนักงาน
            </CardTitle>
            <CardDescription>แสดงจำนวนเวรตามชื่อเวรที่มาจากระบบหลังบ้าน</CardDescription>
              </div>
              <div className="flex space-x-2">
                <Button
                  onClick={exportSummaryToImagePDF}
                  size="sm"
                  variant="outline"
                  className="flex items-center"
                  title="ส่งออก PDF รองรับภาษาไทย"
                >
                  <DocumentArrowDownIcon className="w-4 h-4 mr-1" />
                  PDF 
                </Button>
                <Button
                  onClick={exportSummaryToExcel}
                  size="sm"
                  variant="outline"
                  className="flex items-center"
                >
                  <DocumentArrowDownIcon className="w-4 h-4 mr-1" />
                  Excel
                </Button>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <div className="overflow-x-auto">
              {(() => {
                const sum = calculateShiftSummary()
                const rows = sum.rows || []
                const shiftIds: string[] = sum.shiftIds || []
                const nameMap = sum.shiftIdToName || {}
                return (
              <table className="w-full border-collapse">
                <thead>
                  <tr className="border-b border-gray-200">
                    <th className="text-left py-3 px-4 font-medium text-gray-900">ชื่อ-นามสกุล</th>
                    <th className="text-left py-3 px-4 font-medium text-gray-900">ตำแหน่ง</th>
                        {shiftIds.map(sid => (
                          <th key={sid} className="text-center py-3 px-4 font-medium text-gray-900">{nameMap[sid] || 'เวร'}</th>
                        ))}
                    <th className="text-center py-3 px-4 font-medium text-gray-900">รวม</th>
                    <th className="text-center py-3 px-4 font-medium text-gray-900">ชั่วโมงรวม</th>
                  </tr>
                </thead>
                <tbody>
                      {rows.map((r, idx) => (
                        <tr key={idx} className="border-b border-gray-100 hover:bg-gray-50">
                          <td className="py-3 px-4">{r.name}</td>
                        <td className="py-3 px-4">
                            <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${r.position === 'พยาบาล' ? 'bg-blue-100 text-blue-800' : 'bg-green-100 text-green-800'}`}>
                              {r.position}
                          </span>
                        </td>
                          {shiftIds.map(sid => {
                            const color = shiftTypes.find(s => s.id === sid)?.color || 'bg-gray-100 text-gray-800'
                            return (
                              <td key={sid} className="text-center py-3 px-4">
                                <span className={`inline-flex items-center justify-center w-8 h-8 rounded-full text-sm font-medium ${color}`}>
                                  {r.perShift[sid] || 0}
                          </span>
                        </td>
                            )
                          })}
                        <td className="text-center py-3 px-4">
                            <span className="font-semibold text-gray-900">{r.total}</span>
                        </td>
                        <td className="text-center py-3 px-4">
                            <span className="font-semibold text-blue-600">
                              {calculateTotalHours(r.name || '').toFixed(1)} ชม.
                            </span>
                        </td>
                      </tr>
                      ))}
                </tbody>
              </table>
                )
              })()}
            </div>
          </CardContent>
        </Card>

        {/* Calendar */}
        <Card className="border-0 shadow-md">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>ตารางเวรประจำเดือน</CardTitle>
                <CardDescription>คลิกที่วันที่เพื่อดูรายละเอียดและจัดการเวร</CardDescription>
              </div>
              <div className="flex items-center space-x-2">
                <Button
                  onClick={exportCalendarToImagePDF}
                  variant="outline"
                  size="sm"
                  className="flex items-center"
                  title="ส่งออก PDF ตารางเวรรองรับภาษาไทย"
                >
                  <DocumentArrowDownIcon className="w-4 h-4 mr-1" />
                  PDF 
                </Button>
                <Button
                  onClick={exportCalendarToExcel}
                  variant="outline"
                  size="sm"
                  className="flex items-center"
                >
                  <DocumentArrowDownIcon className="w-4 h-4 mr-1" />
                  Excel
                </Button>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            {/* Calendar Header */}
            <div className="flex items-center justify-between mb-6">
              <Button variant="outline" size="sm" onClick={previousMonth}>
                <ChevronLeftIcon className="w-4 h-4" />
              </Button>
              <h3 className="text-xl font-semibold">
                {monthNames[currentMonth]} {currentYear}
              </h3>
              <Button variant="outline" size="sm" onClick={nextMonth}>
                <ChevronRightIcon className="w-4 h-4" />
              </Button>
            </div>

            {/* Calendar Grid */}
            <div className="hidden md:grid grid-cols-7 gap-1 mb-4">
              {/* Day headers */}
              {['อาทิตย์', 'จันทร์', 'อังคาร', 'พุธ', 'พฤหัสบดี', 'ศุกร์', 'เสาร์'].map((day) => (
                <div key={day} className="p-3 text-center text-sm font-medium text-gray-500 bg-gray-50 rounded-md">
                  {day}
                </div>
              ))}
              
              {/* Calendar days */}
              {calendar.map((day, index) => {
                const dateKey = day ? `${currentYear}-${String(currentMonth + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}` : ''
                const daySchedule = dateKey ? schedules[dateKey] : null
                
                const meta = calendarMeta[dateKey]
                const isOff = meta && (!meta.isWorking || meta.isHoliday)
                return (
                  <div
                    key={dateKey || `empty-${index}`}
                    className={`p-2 border border-gray-200 rounded-md min-h-[120px] ${
                      day ? 'bg-white hover:bg-blue-50 cursor-pointer' : 'bg-gray-50'
                    }`}
                  >
                    {day && isOff && (
                      <>
                        <div className="font-semibold text-gray-900 mb-2">{day}</div>
                        <div className="text-xs p-2 rounded bg-red-50 text-red-700 font-medium">วันหยุด</div>
                      </>
                    )}
                    {day && !isOff && daySchedule && (
                      <>
                        <div className="font-semibold text-gray-900 mb-2">{day}</div>
                        <div className="space-y-1">
                          {Object.entries(daySchedule.shifts).map(([shiftId, shiftData]) => {
                            const totalStaff = shiftData.nurses.length + shiftData.assistants.length
                            return (
                              <div key={shiftId} className={`text-xs p-1 rounded ${shiftData.color || ''}`}>
                                <div className="flex items-center justify-between">
                                  <span className="font-medium">{shiftData.name || shiftTypes.find(s=>s.id===shiftId)?.name || 'เวร'}</span>
                                  <div className="flex items-center space-x-2">
                                  {(shiftData.startTime || shiftData.endTime) && (
                                    <span className="text-[10px] text-gray-600">{shiftData.startTime} - {shiftData.endTime}</span>
                                  )}
                                    <button 
                                      onClick={() => handleEditShift(dateKey, shiftId)}
                                      className="p-1 hover:bg-gray-100 rounded"
                                      title="แก้ไขเวร"
                                    >
                                      <PencilIcon className="w-3 h-3 text-gray-600" />
                                    </button>
                                  </div>
                                </div>
                                <div>พยาบาล: {shiftData.nurses.length} คน{shiftData.nurses.length>0 && ` — ${shiftData.nurses.map(n=>n.name).join(', ')}`}</div>
                                <div>ผู้ช่วย: {shiftData.assistants.length} คน{shiftData.assistants.length>0 && ` — ${shiftData.assistants.map(a=>a.name).join(', ')}`}</div>
                                <div className="font-medium">รวม: {totalStaff} คน</div>
                              </div>
                            )
                          })}
                        </div>
                      </>
                    )}
                  </div>
                )
              })}
            </div>

            {/* Mobile Schedule View */}
            <div className="md:hidden space-y-4">
              {calendar.filter(day => day).map((day, index) => {
                const dNum = Number(day)
                const dateKey = `${currentYear}-${String(currentMonth + 1).padStart(2, '0')}-${String(dNum).padStart(2, '0')}`
                const daySchedule = schedules[dateKey]
                const dayName = ['อาทิตย์', 'จันทร์', 'อังคาร', 'พุธ', 'พฤหัสบดี', 'ศุกร์', 'เสาร์'][new Date(currentYear, currentMonth, dNum).getDay()]
                const meta = calendarMeta[dateKey]
                const isOff = meta && (!meta.isWorking || meta.isHoliday)
                
                return (
                  <Card key={index} className="border-0 shadow-sm">
                    <CardHeader className="pb-3">
                      <CardTitle className="text-lg flex items-center justify-between">
                        <span>{dayName} ที่ {dNum}</span>
                        <span className="text-sm font-normal text-gray-500">{dateKey}</span>
                      </CardTitle>
                    </CardHeader>
                    <CardContent>
                      {isOff ? (
                        <div className="text-xs p-2 rounded bg-red-50 text-red-700 font-medium">วันหยุด</div>
                      ) : daySchedule ? (
                        <div className="space-y-3">
                          {Object.entries(daySchedule.shifts).map(([shiftId, shiftData]) => {
                            const shiftType = shiftTypes.find(s => s.id === shiftId)
                            const totalStaff = shiftData.nurses.length + shiftData.assistants.length

                            return (
                              <div key={shiftId} className={`p-3 rounded-lg ${shiftType?.color}`}>
                                <div className="flex items-center justify-between mb-2">
                                  <h4 className="font-medium">{shiftType?.name}</h4>
                                  <span className="text-sm">รวม: {totalStaff} คน</span>
                                </div>
                                <div className="grid grid-cols-2 gap-4 text-sm">
                                  <div>
                                    <span className="font-medium">พยาบาล:</span> {shiftData.nurses.length} คน
                                  </div>
                                  <div>
                                    <span className="font-medium">ผู้ช่วย:</span> {shiftData.assistants.length} คน
                                  </div>
                                </div>
                                <div className="text-xs text-gray-600 mt-1">
                                  {shiftType?.startTime} - {shiftType?.endTime}
                                </div>
                              </div>
                            )
                          })}
                        </div>
                      ) : (
                        <p className="text-gray-500 text-center py-4">ไม่มีตารางเวรในวันนี้</p>
                      )}
                    </CardContent>
                  </Card>
                )
              })}
            </div>

            {/* Legend */}
            <div className="mt-6 flex flex-wrap gap-4 justify-center">
              {shiftTypes.map(shift => (
                <div key={shift.id} className="flex items-center space-x-2">
                  <div className={`w-4 h-4 rounded ${shift.color}`}></div>
                  <span className="text-sm text-gray-600">{shift.name} ({shift.startTime}-{shift.endTime})</span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Edit Shift Modal */}
        {showEditShiftModal && editingShift && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
            <div className="bg-white rounded-xl shadow-2xl p-6 max-w-6xl w-full">
              <div className="flex items-center justify-between mb-6">
                <div>
                  <h3 className="text-2xl font-semibold text-gray-800">แก้ไขเวร {editingShift.shiftName}</h3>
                  <p className="text-gray-600 mt-1">
                    วันที่: {editingShift.date} | เวลา: {editingShift.startTime} - {editingShift.endTime}
                  </p>
                </div>
                <button 
                  onClick={() => setShowEditShiftModal(false)}
                  className="p-2 hover:bg-gray-100 rounded-full transition-colors"
                >
                  <XMarkIcon className="w-6 h-6 text-gray-500" />
                </button>
              </div>

              <div className="grid grid-cols-2 gap-6">
                {/* Left Column - Current Staff */}
                <div className="space-y-6">
                  {/* Current Nurses Card */}
                  <div className="bg-blue-50 rounded-xl p-4 border border-blue-100">
                    <div className="flex items-center justify-between mb-4">
                      <h4 className="text-lg font-medium text-blue-900">
                        พยาบาลที่ขึ้นเวร
                        <span className="ml-2 text-sm text-blue-600">
                          ({editingShift.nurses.length}/{editingShift.requiredNurses} คน)
                        </span>
                      </h4>
                    </div>
                    <div className="space-y-2">
                      {editingShift.nurses.map((nurse, idx) => (
                        <div key={idx} className="flex items-center justify-between bg-white p-3 rounded-lg shadow-sm border border-blue-100">
                          <div className="flex items-center space-x-3">
                            <UserIcon className="w-5 h-5 text-blue-500" />
                            <span className="text-gray-700">{nurse.name}</span>
                          </div>
                          <button 
                            onClick={() => {
                              if (editingShift) {
                                setEditingShift({
                                  ...editingShift,
                                  nurses: editingShift.nurses.filter((_, i) => i !== idx)
                                })
                              }
                            }}
                            className="p-1.5 hover:bg-red-50 rounded-full text-red-500 transition-colors"
                            title="นำออกจากเวร"
                          >
                            <MinusCircleIcon className="w-5 h-5" />
                          </button>
                        </div>
                      ))}
                    </div>
                  </div>

                  {/* Current Assistants Card */}
                  <div className="bg-green-50 rounded-xl p-4 border border-green-100">
                    <div className="flex items-center justify-between mb-4">
                      <h4 className="text-lg font-medium text-green-900">
                        ผู้ช่วยพยาบาลที่ขึ้นเวร
                        <span className="ml-2 text-sm text-green-600">
                          ({editingShift.assistants.length}/{editingShift.requiredAssistants} คน)
                        </span>
                      </h4>
                    </div>
                    <div className="space-y-2">
                      {editingShift.assistants.map((assistant, idx) => (
                        <div key={idx} className="flex items-center justify-between bg-white p-3 rounded-lg shadow-sm border border-green-100">
                          <div className="flex items-center space-x-3">
                            <UserIcon className="w-5 h-5 text-green-500" />
                            <span className="text-gray-700">{assistant.name}</span>
                          </div>
                          <button 
                            onClick={() => {
                              if (editingShift) {
                                setEditingShift({
                                  ...editingShift,
                                  assistants: editingShift.assistants.filter((_, i) => i !== idx)
                                })
                              }
                            }}
                            className="p-1.5 hover:bg-red-50 rounded-full text-red-500 transition-colors"
                            title="นำออกจากเวร"
                          >
                            <MinusCircleIcon className="w-5 h-5" />
                          </button>
                        </div>
                      ))}
                    </div>
                  </div>
                </div>

                {/* Right Column - Available Staff */}
                <div className="space-y-6">
                  {/* Available Nurses Card */}
                  <div className="bg-indigo-50 rounded-xl p-4 border border-indigo-100">
                    <div className="flex items-center justify-between mb-4">
                      <h4 className="text-lg font-medium text-indigo-900">พยาบาลที่สามารถขึ้นเวรได้</h4>
                    </div>
                    <div className="space-y-2">
                      {(availableFromApi.length ? availableFromApi : employees)
                        .filter(emp => emp.position === 'nurse' && 
                          !editingShift.nurses.some(n => String(n.id) === String(emp.id)))
                        .map((nurse, idx) => (
                          <div key={idx} className="flex items-center justify-between bg-white p-3 rounded-lg shadow-sm border border-indigo-100">
                            <div className="flex items-center space-x-3">
                              <UserIcon className="w-5 h-5 text-indigo-500" />
                              <span className="text-gray-700">{nurse.name}</span>
                            </div>
                            <button 
                                                          onClick={() => {
                              if (editingShift) {
                                setEditingShift({
                                  ...editingShift,
                                  nurses: [...editingShift.nurses, nurse]
                                })
                              }
                            }}
                              className="p-1.5 hover:bg-indigo-100 rounded-full text-indigo-600 transition-colors"
                              title="เพิ่มเข้าเวร"
                            >
                              <PlusCircleIcon className="w-5 h-5" />
                            </button>
                          </div>
                        ))}
                    </div>
                  </div>

                  {/* Available Assistants Card */}
                  <div className="bg-purple-50 rounded-xl p-4 border border-purple-100">
                    <div className="flex items-center justify-between mb-4">
                      <h4 className="text-lg font-medium text-purple-900">ผู้ช่วยพยาบาลที่สามารถขึ้นเวรได้</h4>
                    </div>
                    <div className="space-y-2">
                      {(availableFromApi.length ? availableFromApi : employees)
                        .filter(emp => emp.position === 'assistant' && 
                          !editingShift.assistants.some(a => String(a.id) === String(emp.id)))
                        .map((assistant, idx) => (
                          <div key={idx} className="flex items-center justify-between bg-white p-3 rounded-lg shadow-sm border border-purple-100">
                            <div className="flex items-center space-x-3">
                              <UserIcon className="w-5 h-5 text-purple-500" />
                              <span className="text-gray-700">{assistant.name}</span>
                            </div>
                            <button 
                                                          onClick={() => {
                              if (editingShift) {
                                setEditingShift({
                                  ...editingShift,
                                  assistants: [...editingShift.assistants, assistant]
                                })
                              }
                            }}
                              className="p-1.5 hover:bg-purple-100 rounded-full text-purple-600 transition-colors"
                              title="เพิ่มเข้าเวร"
                            >
                              <PlusCircleIcon className="w-5 h-5" />
                            </button>
                          </div>
                        ))}
                    </div>
                  </div>
                </div>
              </div>

              <div className="mt-8 flex justify-end space-x-3">
                <Button 
                  variant="outline" 
                  onClick={() => setShowEditShiftModal(false)}
                  className="px-6"
                >
                  ยกเลิก
                </Button>
                <Button 
                  onClick={async () => {
                    try {
                      // Show loading
                      Swal.fire({
                        title: 'กำลังบันทึกการแก้ไข...',
                        html: '<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>',
                        showConfirmButton: false,
                        allowOutsideClick: false,
                        allowEscapeKey: false
                      })

                      // Get original staff IDs from the shift data
                      const originalShift = schedules[editingShift.date]?.shifts[editingShift.shiftId]
                      
                      if (!originalShift) {
                        throw new Error('ไม่พบข้อมูลเวรเดิม')
                      }

                      const originalNurseIds = originalShift.nurses.map(n => String(n.id))
                      const originalAssistantIds = originalShift.assistants.map(a => String(a.id))

                      // Get current staff IDs from editing state
                      const currentNurseIds = editingShift.nurses.map(n => String(n.id))
                      const currentAssistantIds = editingShift.assistants.map(a => String(a.id))

                      // Calculate differences
                      const addNurses = currentNurseIds.filter(id => !originalNurseIds.includes(id))
                      const removeNurses = originalNurseIds.filter(id => !currentNurseIds.includes(id))
                      const addAssistants = currentAssistantIds.filter(id => !originalAssistantIds.includes(id))
                      const removeAssistants = originalAssistantIds.filter(id => !currentAssistantIds.includes(id))

                      // Check overlap for staff to be added
                      const checkOverlap = async (staffId: string) => {
                        try {
                          const result = await scheduleService.checkShiftOverlap({
                            departmentId: selectedDepartment,
                            date: editingShift.date,
                            shiftId: editingShift.shiftId,
                            staffId
                          })
                          return result.canAssign
                        } catch (error) {
                          console.error('Error checking overlap:', error)
                          return false
                        }
                      }

                      // Filter out staff with overlapping shifts
                      const addableNurses = await Promise.all(
                        addNurses.map(async id => ({
                          id,
                          canAdd: await checkOverlap(id)
                        }))
                      )
                      const addableAssistants = await Promise.all(
                        addAssistants.map(async id => ({
                          id,
                          canAdd: await checkOverlap(id)
                        }))
                      )

                      // Edit shift
                      await scheduleService.editShift({
                        departmentId: selectedDepartment,
                        date: editingShift.date,
                        shiftId: editingShift.shiftId,
                        addNurses: addableNurses.filter(n => n.canAdd).map(n => n.id),
                        addAssistants: addableAssistants.filter(a => a.canAdd).map(a => a.id),
                        removeNurses: removeNurses,
                        removeAssistants: removeAssistants
                      })

                      // Refresh data
                      await reloadScheduleData()

                      // Close modal and show success
                      setShowEditShiftModal(false)
                      Swal.fire({
                        icon: 'success',
                        title: 'แก้ไขเวรสำเร็จ',
                        text: 'ข้อมูลถูกบันทึกและรีเฟรชแล้ว',
                        confirmButtonColor: '#2563eb',
                        timer: 3000,
                        timerProgressBar: true
                      })

                    } catch (error: any) {
                      Swal.fire({
                        icon: 'error',
                        title: 'ไม่สามารถแก้ไขเวรได้',
                        text: error?.message || 'เกิดข้อผิดพลาด',
                        confirmButtonColor: '#2563eb'
                      })
                    }
                  }}
                  className="px-6 bg-blue-600 hover:bg-blue-700"
                >
                  <CheckIcon className="w-5 h-5 mr-2" />
                  บันทึก
                </Button>
              </div>
            </div>
          </div>
        )}

        {/* Reduce Staff Modal */}
        {showReduceModal && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg p-6 w-full max-w-md">
              <h3 className="text-lg font-medium mb-4">ปรับลดจำนวนพนักงาน</h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    วันที่
                  </label>
                  <Input
                    type="date"
                    value={reduceForm.date}
                    onChange={(e) => setReduceForm({ ...reduceForm, date: e.target.value })}
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    เวรที่จะปรับลด
                  </label>
                  <select
                    value={reduceForm.shift}
                    onChange={(e) => setReduceForm({ ...reduceForm, shift: e.target.value })}
                    className="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="">เลือกเวร</option>
                    <option value="morning">เวรเช้า (07:00-15:00)</option>
                    <option value="afternoon">เวรบ่าย (15:00-23:00)</option>
                    <option value="night">เวรดึก (23:00-07:00)</option>
                  </select>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      จำนวนพยาบาลที่จะลด
                    </label>
                    <Input
                      type="number"
                      value={reduceForm.nursesToReduce}
                      onChange={(e) => setReduceForm({ ...reduceForm, nursesToReduce: Number(e.target.value) })}
                      min="0"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      จำนวนผู้ช่วยที่จะลด
                    </label>
                    <Input
                      type="number"
                      value={reduceForm.assistantsToReduce}
                      onChange={(e) => setReduceForm({ ...reduceForm, assistantsToReduce: Number(e.target.value) })}
                      min="0"
                    />
                  </div>
                </div>
                <div className="bg-yellow-50 border border-yellow-200 rounded-md p-3">
                  <p className="text-sm text-yellow-800">
                    <strong>หมายเหตุ:</strong> ระบบจะลดพนักงานที่มีจำนวนเวรมากที่สุดออกก่อน
                    เพื่อความเป็นธรรมในการกระจายภาระงาน
                  </p>
                </div>
                <div className="flex justify-end space-x-3">
                  <Button
                    variant="outline"
                    onClick={() => {
                      setShowReduceModal(false)
                      setReduceForm({ date: '', shift: '', nursesToReduce: 0, assistantsToReduce: 0 })
                    }}
                  >
                    ยกเลิก
                  </Button>
                  <Button onClick={processStaffReduction}>
                    ปรับลดพนักงาน
                  </Button>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </DashboardLayout>
  )
}
