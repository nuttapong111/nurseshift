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
  EyeIcon
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

// Types
interface Employee {
  id: number
  name: string
  position: 'พยาบาล' | 'ผู้ช่วยพยาบาล'
  department: string
  shiftCounts: {
    morning: number
    afternoon: number
    night: number
    total: number
  }
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
      { id: 1, name: 'สมใจ ใจดี', position: 'พยาบาล', department: 'แผนกฉุกเฉิน', shiftCounts: { morning: 8, afternoon: 6, night: 4, total: 18 } },
      { id: 2, name: 'วิไล รักษ์ดี', position: 'พยาบาล', department: 'แผนกฉุกเฉิน', shiftCounts: { morning: 7, afternoon: 5, night: 6, total: 18 } },
      { id: 3, name: 'มณี ใสใส', position: 'ผู้ช่วยพยาบาล', department: 'แผนกฉุกเฉิน', shiftCounts: { morning: 6, afternoon: 7, night: 5, total: 18 } },
      { id: 4, name: 'สุดา ขยันขัน', position: 'ผู้ช่วยพยาบาล', department: 'แผนกฉุกเฉิน', shiftCounts: { morning: 5, afternoon: 8, night: 5, total: 18 } }
    ]
  },
  {
    id: 'internal',
    name: 'แผนกอายุรกรรม',
    employees: [
      { id: 5, name: 'ประไพ แจ่มใส', position: 'พยาบาล', department: 'แผนกอายุรกรรม', shiftCounts: { morning: 9, afternoon: 5, night: 4, total: 18 } },
      { id: 6, name: 'นิรมล สุขใจ', position: 'พยาบาล', department: 'แผนกอายุรกรรม', shiftCounts: { morning: 6, afternoon: 7, night: 5, total: 18 } },
      { id: 7, name: 'บุปผา รื่นรมย์', position: 'ผู้ช่วยพยาบาล', department: 'แผนกอายุรกรรม', shiftCounts: { morning: 5, afternoon: 6, night: 7, total: 18 } }
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

  useEffect(() => { setMounted(true) }, [])
  const [showReduceModal, setShowReduceModal] = useState(false)
  const [reduceForm, setReduceForm] = useState({
    date: '',
    shift: '',
    nursesToReduce: 0,
    assistantsToReduce: 0
  })
  const [employees, setEmployees] = useState<Employee[]>([])
  const [rawAssignments, setRawAssignments] = useState<any[]>([])

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

  // โหลดจำนวนเวรทั้งหมดของเดือน (ทุกแผนก)
  useEffect(() => {
    const loadMonthlyShifts = async () => {
      try {
        const month = `${currentDate.getFullYear()}-${String(currentDate.getMonth() + 1).padStart(2,'0')}`
        const all = await scheduleService.list({ month })
        setStats(prev => ({ ...prev, totalShifts: (all || []).length }))
      } catch (_e) {
        // ignore
      }
    }
    loadMonthlyShifts()
  }, [currentDate])

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
        const month = `${currentYear}-${String(currentMonth + 1).padStart(2,'0')}`
        const [items, metaDays] = await Promise.all([
          scheduleService.list({ departmentId: selectedDepartment, month }),
          scheduleService.calendarMeta(selectedDepartment, month)
        ])
        const metaObj: {[k:string]: {isWorking:boolean; isHoliday:boolean}} = {}
        ;(metaDays || []).forEach((d:any) => { metaObj[d.date] = { isWorking: !!d.isWorking, isHoliday: !!d.isHoliday } })
        setCalendarMeta(metaObj)
        setRawAssignments(items)

        // Build calendar using required counts from shifts
        const byShiftId: Record<string, { requiredNurses: number; requiredAssistants: number; name: string; startTime: string; endTime: string; color: string }> = {}
        types.forEach(t => {
          const s = shifts.find(x => x.id === t.id)!
          byShiftId[t.id] = { requiredNurses: s?.requiredNurse || 0, requiredAssistants: s?.requiredAsst || 0, name: t.name, startTime: t.startTime, endTime: t.endTime, color: t.color }
        })

        const daysInMonth = new Date(currentYear, currentMonth + 1, 0).getDate()
        const newSchedules: {[key: string]: DaySchedule} = {}
        for (let day = 1; day <= daysInMonth; day++) {
          const dateKey = `${currentYear}-${String(currentMonth + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`
          newSchedules[dateKey] = {
            date: dateKey,
            shifts: {
              // initialize per backend shift id
              ...types.reduce((acc, st) => {
                const meta = byShiftId[st.id]
                acc[st.id] = { name: meta?.name, startTime: meta?.startTime, endTime: meta?.endTime, color: meta?.color, nurses: [], assistants: [], requiredNurses: meta?.requiredNurses || 0, requiredAssistants: meta?.requiredAssistants || 0 }
                return acc
              }, {} as any)
            }
          }
        }

        // Fill counts by shiftId per date
        const counts: Record<string, Record<string, { n: number; a: number }>> = {}
        const names: Record<string, Record<string, { nurses: string[]; assistants: string[] }>> = {}
        for (const rec of items as any[]) {
          if (!byShiftId[rec.shiftId]) continue
          const dateKey = rec.scheduleDate
          const sid = rec.shiftId
          if (!counts[dateKey]) counts[dateKey] = {}
          if (!counts[dateKey][sid]) counts[dateKey][sid] = { n:0, a:0 }
          if (!names[dateKey]) names[dateKey] = {}
          if (!names[dateKey][sid]) names[dateKey][sid] = { nurses: [], assistants: [] }
          if (rec.departmentRole === 'nurse') { counts[dateKey][sid].n += 1; names[dateKey][sid].nurses.push(rec.userName || '-') }
          else if (rec.departmentRole === 'assistant') { counts[dateKey][sid].a += 1; names[dateKey][sid].assistants.push(rec.userName || '-') }
        }
        const mkDummy = (n: number): Employee[] => Array.from({ length: n }).map((_, i) => ({ id: i+1, name: '-', position: 'พยาบาล', department: '', shiftCounts: { morning:0, afternoon:0, night:0, total:0 } }))
        Object.entries(counts).forEach(([dateKey, v]) => {
          const ds = newSchedules[dateKey]
          if (!ds) return
          const mkNamed = (arr: string[], position: 'พยาบาล'|'ผู้ช่วยพยาบาล'): Employee[] => arr.map((name, i) => ({ id: i+1, name, position, department: '', shiftCounts: { morning:0, afternoon:0, night:0, total:0 } }))
          Object.entries(v).forEach(([sid, cnt]) => {
            const nm = names[dateKey]?.[sid]
            ds.shifts[sid] = ds.shifts[sid] || { name: byShiftId[sid]?.name, startTime: byShiftId[sid]?.startTime, endTime: byShiftId[sid]?.endTime, color: byShiftId[sid]?.color, nurses: [], assistants: [], requiredNurses: byShiftId[sid]?.requiredNurses || 0, requiredAssistants: byShiftId[sid]?.requiredAssistants || 0 }
            const nurseNames = (nm?.nurses && nm.nurses.length > 0) ? nm.nurses : Array.from({length: cnt.n}).map(() => '-')
            const asstNames = (nm?.assistants && nm.assistants.length > 0) ? nm.assistants : Array.from({length: cnt.a}).map(() => '-')
            ds.shifts[sid].nurses = mkNamed(nurseNames, 'พยาบาล')
            ds.shifts[sid].assistants = mkNamed(asstNames, 'ผู้ช่วยพยาบาล')
          })
        })

        setSchedules(newSchedules)
      } catch (_err) {
        // fallback to mock if API fails
        const dept = mockDepartments.find(d => d.id === selectedDepartment)
        setEmployees(dept ? [...dept.employees] : [])
    generateMockSchedule()
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

  const handleAutoGenerateBackend = async () => {
    if (!selectedDepartment) return
    try {
      const res = await scheduleService.autoGenerate(selectedDepartment, formatMonth(currentDate))
      await Swal.fire({ icon: 'success', title: 'สร้างตารางเวรอัตโนมัติ (Backend) สำเร็จ', text: `บันทึก ${res.inserted} รายการ`, confirmButtonColor: '#2563eb' })
    } catch (e: any) {
      await Swal.fire({ icon: 'error', title: 'ไม่สามารถสร้างตารางเวรได้', text: e?.message || 'เกิดข้อผิดพลาด', confirmButtonColor: '#2563eb' })
    }
  }

  const handleAutoGenerateAI = async () => {
    if (!selectedDepartment) return
    try {
      const res = await scheduleService.aiGenerate(selectedDepartment, formatMonth(currentDate))
      await Swal.fire({ icon: 'success', title: 'สร้างตารางเวรด้วย AI สำเร็จ', text: `บันทึก ${res.inserted} รายการ`, confirmButtonColor: '#2563eb' })
    } catch (e: any) {
      await Swal.fire({ icon: 'error', title: 'ไม่สามารถสร้างตารางเวรด้วย AI ได้', text: e?.message || 'เกิดข้อผิดพลาด', confirmButtonColor: '#2563eb' })
    }
  }

  // Build dynamic summary columns by backend shifts (use shiftTypes directly)
  const calculateShiftSummary = () => {
    const shiftIdToName: Record<string, string> = {}
    shiftTypes.forEach(st => { shiftIdToName[st.id] = st.name })
    const shiftIds = Object.keys(shiftIdToName)

    type Row = { name: string; position: 'พยาบาล'|'ผู้ช่วยพยาบาล'; perShift: Record<string, number>; total: number }
    const map: Record<string, Row> = {}
    rawAssignments.forEach((a: any) => {
      const name = a.userName || '-'
      const role = a.departmentRole === 'assistant' ? 'ผู้ช่วยพยาบาล' : 'พยาบาล'
      if (!map[name]) map[name] = { name, position: role, perShift: {}, total: 0 }
      map[name].perShift[a.shiftId] = (map[name].perShift[a.shiftId] || 0) + 1
      map[name].total++
    })

    // transform rows to include all shift columns
    const rows = Object.values(map).map(r => {
      shiftIds.forEach(sid => { if (!(sid in r.perShift)) r.perShift[sid] = 0 })
      return r
    }).sort((a,b)=>a.name.localeCompare(b.name))

    return { rows, shiftIds, shiftIdToName }
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

    // Update employee shift counts
    const updatedEmployees = employees.map(emp => {
      const removedNurse = removedNurses.find(n => n.id === emp.id)
      const removedAssistant = removedAssistants.find(a => a.id === emp.id)
      
      if (removedNurse || removedAssistant) {
        return {
          ...emp,
          shiftCounts: {
            ...emp.shiftCounts,
            [reduceForm.shift]: emp.shiftCounts[reduceForm.shift as keyof typeof emp.shiftCounts] - 1,
            total: emp.shiftCounts.total - 1
          }
        }
      }
      return emp
    })

    setSchedules(updatedSchedules)
    setEmployees(updatedEmployees)
    setShowReduceModal(false)
    setReduceForm({ date: '', shift: '', nursesToReduce: 0, assistantsToReduce: 0 })

    await Swal.fire({
      icon: 'success',
      title: 'ปรับลดจำนวนพนักงานสำเร็จ!',
      html: `
        <div class="text-left">
          <p><strong>พนักงานที่ถูกลดออก:</strong></p>
          ${removedNurses.length > 0 ? `<p>พยาบาล: ${removedNurses.map(n => n.name).join(', ')}</p>` : ''}
          ${removedAssistants.length > 0 ? `<p>ผู้ช่วยพยาบาล: ${removedAssistants.map(a => a.name).join(', ')}</p>` : ''}
          <p class="mt-2"><strong>เหตุผล:</strong> เลือกจากคนที่มีจำนวนเวรมากที่สุด</p>
        </div>
      `,
      confirmButtonColor: '#2563eb'
    })
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
                <Button onClick={handleReduceStaff} className="w-full sm:w-auto">ปรับลดพนักงาน</Button>
                <Button onClick={handleAutoGenerateBackend} className="w-full sm:w-auto">สร้างอัตโนมัติ (Backend)</Button>
                <Button variant="outline" onClick={handleAutoGenerateAI} className="w-full sm:w-auto">สร้างอัตโนมัติ (AI)</Button>
              </div>
            </div>
          </div>
        </div>

        {/* Shift Summary Table - dynamic columns by backend shifts */}
        <Card className="border-0 shadow-md">
          <CardHeader>
            <CardTitle className="flex items-center">
              <UserGroupIcon className="w-5 h-5 mr-2" />
              สรุปจำนวนเวรของพนักงาน
            </CardTitle>
            <CardDescription>แสดงจำนวนเวรตามชื่อเวรที่มาจากระบบหลังบ้าน</CardDescription>
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
                <Button variant="outline" size="sm">
                  <DocumentArrowDownIcon className="w-4 h-4 mr-1" />
                  ส่งออก
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
                                  {(shiftData.startTime || shiftData.endTime) && (
                                    <span className="text-[10px] text-gray-600">{shiftData.startTime} - {shiftData.endTime}</span>
                                  )}
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
