'use client'

import { useEffect, useState } from 'react'
import { 
  PlusIcon,
  PencilIcon,
  TrashIcon,
  ClockIcon,
  CalendarIcon,
  CheckIcon
} from '@heroicons/react/24/outline'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'
import { Switch } from '@/components/ui/Switch'
import DashboardLayout from '@/components/layout/DashboardLayout'
import Swal from 'sweetalert2'
import { TimePicker, DatePicker } from 'antd'
import dayjs from 'dayjs'
import 'antd/dist/reset.css'
import settingService from '@/services/settingService'
import { departmentService } from '@/services/departmentService'

// Department options (จะโหลดจาก API และ fallback เป็นว่าง)
type DeptOption = { id: string; name: string; description?: string }

type ShiftItem = { id: string, name: string, startTime: string, endTime: string, nurseCount: number, assistantCount: number, isActive: boolean }

type HolidayItem = { id: string, name: string, startDate: string, endDate: string }

const weekDays = [
  { id: 'monday', name: 'จันทร์', enabled: true },
  { id: 'tuesday', name: 'อังคาร', enabled: true },
  { id: 'wednesday', name: 'พุธ', enabled: true },
  { id: 'thursday', name: 'พฤหัสบดี', enabled: true },
  { id: 'friday', name: 'ศุกร์', enabled: true },
  { id: 'saturday', name: 'เสาร์', enabled: false },
  { id: 'sunday', name: 'อาทิตย์', enabled: false }
]

export default function DepartmentSettingsPage() {
  const [departments, setDepartments] = useState<DeptOption[]>([])
  const [selectedDepartment, setSelectedDepartment] = useState<DeptOption | null>(null)
  const [shifts, setShifts] = useState<ShiftItem[]>([])
  const [holidays, setHolidays] = useState<HolidayItem[]>([])
  const [workingDays, setWorkingDays] = useState(weekDays)
  
  const [showAddShiftModal, setShowAddShiftModal] = useState(false)
  const [showAddHolidayModal, setShowAddHolidayModal] = useState(false)
  const [editingHoliday, setEditingHoliday] = useState<HolidayItem | null>(null)
  
  const [shiftForm, setShiftForm] = useState({
    name: '',
    startTime: '07:00',
    endTime: '15:00',
    nurseCount: 1,
    assistantCount: 1
  })
  const [editingShift, setEditingShift] = useState<ShiftItem | null>(null)
  
  const [holidayForm, setHolidayForm] = useState({
    name: '',
    startDate: '',
    endDate: ''
  })

  const dayIdToIndex: Record<string, number> = { sunday: 0, monday: 1, tuesday: 2, wednesday: 3, thursday: 4, friday: 5, saturday: 6 }
  const indexToDayId = ['sunday','monday','tuesday','wednesday','thursday','friday','saturday']

  const loadSettings = async (departmentId: string) => {
    try {
      const data = await settingService.getSettings(departmentId)
      // Working days
      const mappedDays = weekDays.map(d => ({ ...d, enabled: false }))
      for (const w of data.workingDays) {
        const id = indexToDayId[w.dayOfWeek]
        const idx = mappedDays.findIndex(d => d.id === id)
        if (idx >= 0) mappedDays[idx].enabled = Boolean(w.isWorkingDay)
      }
      setWorkingDays(mappedDays)
      // Shifts
      setShifts(data.shifts.map(s => ({
        id: s.id,
        name: s.name || '',
        startTime: s.startTime,
        endTime: s.endTime,
        nurseCount: s.requiredNurses,
        assistantCount: s.requiredAssistants,
        isActive: s.isActive
      })))
      // Holidays
      setHolidays(data.holidays.map(h => ({ id: h.id, name: h.name, startDate: h.startDate, endDate: h.endDate })))
    } catch (e) {
      console.error('loadSettings error', e)
    }
  }

  // Load departments then settings
  useEffect(() => {
    (async () => {
      try {
        const deps = await departmentService.getDepartments()
        const options: DeptOption[] = deps.map(d => ({ id: d.id, name: d.name, description: d.description }))
        setDepartments(options)
        const first = options[0]
        if (first) {
          setSelectedDepartment(first)
          await loadSettings(first.id)
        }
      } catch (e) { console.error('load departments failed', e) }
    })()
  }, [])

  const handleToggleShift = async (shiftId: any, currentStatus: boolean) => {
    setShifts(shifts.map(shift => 
      shift.id === shiftId 
        ? { ...shift, isActive: !currentStatus }
        : shift
    ))
    
    try {
      await settingService.toggleShift(String(shiftId), !currentStatus)
    } catch (e) { console.error(e) }

    await Swal.fire({
      icon: 'success',
      title: currentStatus ? 'ปิดการใช้งานเวรแล้ว' : 'เปิดการใช้งานเวรแล้ว',
      timer: 1500,
      showConfirmButton: false
    })
  }

  const handleToggleWorkingDay = async (dayId: string) => {
    const updated = workingDays.map(day => 
      day.id === dayId 
        ? { ...day, enabled: !day.enabled }
        : day)
    setWorkingDays(updated)

    // ส่งขึ้น API
    try {
      if (selectedDepartment) {
        const payload = updated.map(d => ({ idOrName: d.id as any, enabled: d.enabled }))
        await settingService.updateWorkingDays(selectedDepartment.id, payload)
      }
    } catch (e) { console.error(e) }
  }

  const handleAddOrEditShift = async () => {
    if (!shiftForm.name || !shiftForm.startTime || !shiftForm.endTime) {
      await Swal.fire({
        icon: 'warning',
        title: 'กรุณากรอกข้อมูลให้ครบ',
        confirmButtonColor: '#2563eb'
      })
      return
    }

    if (editingShift) {
      try {
        await settingService.updateShift(String(editingShift.id), {
          name: shiftForm.name,
          type: 'custom',
          startTime: shiftForm.startTime,
          endTime: shiftForm.endTime,
          nurseCount: shiftForm.nurseCount,
          assistantCount: shiftForm.assistantCount,
          color: '#3B82F6',
        })
        setShifts(shifts.map(s => String(s.id) === String(editingShift.id) ? { ...s, name: shiftForm.name, startTime: shiftForm.startTime, endTime: shiftForm.endTime, nurseCount: shiftForm.nurseCount, assistantCount: shiftForm.assistantCount } : s))
        await Swal.fire({ icon: 'success', title: 'บันทึกเวรสำเร็จ!', confirmButtonColor: '#2563eb' })
      } catch (e) { console.error(e) }
    } else {
      let newShiftId: string | null = null
      if (selectedDepartment) {
        try {
          newShiftId = await settingService.createShift({
            departmentId: selectedDepartment.id,
            name: shiftForm.name,
            type: 'custom',
            startTime: shiftForm.startTime,
            endTime: shiftForm.endTime,
            nurseCount: shiftForm.nurseCount,
            assistantCount: shiftForm.assistantCount,
            color: '#3B82F6',
            isActive: true,
          })
        } catch (e) { console.error(e) }
      }
      const newShift = {
        id: newShiftId || String(Date.now()),
        name: shiftForm.name,
        startTime: shiftForm.startTime,
        endTime: shiftForm.endTime,
        nurseCount: shiftForm.nurseCount,
        assistantCount: shiftForm.assistantCount,
        isActive: true
      }
      setShifts([...shifts, newShift])
      await Swal.fire({ icon: 'success', title: 'เพิ่มเวรสำเร็จ!', confirmButtonColor: '#2563eb' })
    }
    setShiftForm({ name: '', startTime: '07:00', endTime: '15:00', nurseCount: 1, assistantCount: 1 })
    setEditingShift(null)
    setShowAddShiftModal(false)
  }

  const handleAddOrEditHoliday = async () => {
    if (!holidayForm.name || !holidayForm.startDate || !holidayForm.endDate) {
      await Swal.fire({
        icon: 'warning',
        title: 'กรุณากรอกข้อมูลให้ครบ',
        confirmButtonColor: '#2563eb'
      })
      return
    }

    // Check if end date is before start date
    if (new Date(holidayForm.endDate) < new Date(holidayForm.startDate)) {
      await Swal.fire({
        icon: 'warning',
        title: 'วันที่สิ้นสุดต้องไม่น้อยกว่าวันที่เริ่มต้น',
        confirmButtonColor: '#2563eb'
      })
      return
    }

    if (editingHoliday) {
      try {
        await settingService.updateHoliday(String(editingHoliday.id), {
          name: holidayForm.name,
          startDate: holidayForm.startDate,
          endDate: holidayForm.endDate,
          isRecurring: false,
        })
        setHolidays(holidays.map(h => String(h.id) === String(editingHoliday.id) ? { id: editingHoliday.id, name: holidayForm.name, startDate: holidayForm.startDate, endDate: holidayForm.endDate } : h))
        await Swal.fire({ icon: 'success', title: 'บันทึกวันหยุดสำเร็จ!', confirmButtonColor: '#2563eb' })
      } catch (e) { console.error(e) }
    } else {
      let newId: string | null = null
      if (selectedDepartment) {
        try {
          newId = await settingService.createHoliday({
            departmentId: selectedDepartment.id,
            name: holidayForm.name,
            startDate: holidayForm.startDate,
            endDate: holidayForm.endDate,
            isRecurring: false,
          })
        } catch (e) { console.error(e) }
      }
      const newHoliday = {
        id: newId || String(Date.now()),
        name: holidayForm.name,
        startDate: holidayForm.startDate,
        endDate: holidayForm.endDate
      }
      setHolidays([...holidays, newHoliday])
      await Swal.fire({ icon: 'success', title: 'เพิ่มวันหยุดสำเร็จ!', confirmButtonColor: '#2563eb' })
    }
    setHolidayForm({ name: '', startDate: '', endDate: '' })
    setEditingHoliday(null)
    setShowAddHolidayModal(false)
  }

  const handleDeleteShift = async (shiftId: string, shiftName: string) => {
    const result = await Swal.fire({
      title: 'ลบเวร?',
      text: `คุณต้องการลบเวร "${shiftName}" หรือไม่?`,
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#dc2626',
      cancelButtonColor: '#6b7280',
      confirmButtonText: 'ลบเวร',
      cancelButtonText: 'ยกเลิก'
    })

    if (result.isConfirmed) {
      try { await settingService.deleteShift(String(shiftId)) } catch (e) { console.error(e) }
      setShifts(shifts.filter(s => String(s.id) !== String(shiftId)))
      await Swal.fire({
        icon: 'success',
        title: 'ลบเวรสำเร็จ!',
        confirmButtonColor: '#2563eb'
      })
    }
  }

  const handleDeleteHoliday = async (holidayId: number, holidayName: string) => {
    const result = await Swal.fire({
      title: 'ลบวันหยุด?',
      text: `คุณต้องการลบวันหยุด "${holidayName}" หรือไม่?`,
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#dc2626',
      cancelButtonColor: '#6b7280',
      confirmButtonText: 'ลบวันหยุด',
      cancelButtonText: 'ยกเลิก'
    })

    if (result.isConfirmed) {
      try { await settingService.deleteHoliday(String(holidayId)) } catch (e) { console.error(e) }
      setHolidays(holidays.filter(h => h.id !== String(holidayId)))
      await Swal.fire({
        icon: 'success',
        title: 'ลบวันหยุดสำเร็จ!',
        confirmButtonColor: '#2563eb'
      })
    }
  }

  return (
    <DashboardLayout>
      <div className="space-y-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">ตั้งค่าแผนก</h1>
          <p className="text-gray-600 mt-2">จัดการเวร วันหยุด และวันทำงานของแผนก</p>
          
          {/* Department Selection */}
          <div className="mt-6">
            <label htmlFor="department-select" className="block text-sm font-medium text-gray-700 mb-2">
              เลือกแผนก
            </label>
              <select
              id="department-select"
                value={selectedDepartment?.id ?? ''}
              onChange={(e) => {
                const dept = departments.find(d => d.id === e.target.value)
                if (dept) { setSelectedDepartment(dept); loadSettings(dept.id) }
              }}
              className="block w-full max-w-md px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
            >
              <option value="" disabled>เลือกแผนก</option>
              {departments.map((dept) => (
                <option key={dept.id} value={dept.id}>
                  {dept.name}
                </option>
              ))}
            </select>
            {selectedDepartment ? (
              <p className="mt-2 text-sm text-gray-500">
                {selectedDepartment.description}
              </p>
            ) : null}
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Shifts Management */}
          <Card className="border-0 shadow-md">
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle className="flex items-center">
                    <ClockIcon className="w-5 h-5 mr-2" />
                    การตั้งค่าเวร - {selectedDepartment?.name ?? ''}
                  </CardTitle>
                  <CardDescription>สร้างและจัดการเวรในแผนก</CardDescription>
                </div>
                <Button onClick={() => setShowAddShiftModal(true)} size="sm">
                  <PlusIcon className="w-4 h-4 mr-2" />
                  เพิ่มเวร
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {shifts.map((shift) => (
                  <div key={shift.id} className="p-4 border border-gray-200 rounded-lg">
                    <div className="flex items-center justify-between mb-3">
                      <div className="flex-1">
                        <h3 className="font-medium text-gray-900">{shift.name || 'เวร'}</h3>
                        <p className="text-sm text-gray-500">
                          {shift.startTime} - {shift.endTime}
                        </p>
                        <p className="text-sm text-gray-500">
                          พยาบาล {shift.nurseCount} คน | ผู้ช่วย {shift.assistantCount} คน
                        </p>
                      </div>
                      <div className="flex items-center space-x-2">
                        <div className="flex items-center space-x-2">
                          <span className="text-sm text-gray-600">
                            {shift.isActive ? 'เปิดใช้งาน' : 'ปิดใช้งาน'}
                          </span>
                          <Switch
                            checked={shift.isActive}
                            onChange={() => handleToggleShift(shift.id, shift.isActive)}
                          />
                        </div>
                        <Button variant="ghost" size="sm" onClick={() => { setEditingShift(shift); setShiftForm({ name: shift.name || '', startTime: shift.startTime, endTime: shift.endTime, nurseCount: shift.nurseCount, assistantCount: shift.assistantCount }); setShowAddShiftModal(true) }}>
                          <PencilIcon className="w-4 h-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleDeleteShift(shift.id, shift.name)}
                          className="text-red-600 hover:text-red-700"
                        >
                          <TrashIcon className="w-4 h-4" />
                        </Button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Working Days */}
          <Card className="border-0 shadow-md">
            <CardHeader>
              <CardTitle className="flex items-center">
                <CheckIcon className="w-5 h-5 mr-2" />
                วันทำงาน - {selectedDepartment?.name ?? ''}
              </CardTitle>
              <CardDescription>เลือกวันที่ทำงานในแต่ละสัปดาห์</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {workingDays.map((day) => (
                  <div key={day.id} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                    <span className="font-medium text-gray-900">{day.name}</span>
                    <div className="flex items-center space-x-2">
                      <span className="text-sm text-gray-600">
                        {day.enabled ? 'ทำงาน' : 'หยุด'}
                      </span>
                      <Switch
                        checked={day.enabled}
                        onChange={() => handleToggleWorkingDay(day.id)}
                      />
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Holidays Management */}
        <Card className="border-0 shadow-md">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="flex items-center">
                  <CalendarIcon className="w-5 h-5 mr-2" />
                  วันหยุดประจำปี
                </CardTitle>
                <CardDescription>กำหนดวันหยุดของแผนกในปีนี้</CardDescription>
              </div>
              <Button onClick={() => setShowAddHolidayModal(true)} size="sm">
                <PlusIcon className="w-4 h-4 mr-2" />
                เพิ่มวันหยุด
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {holidays.map((holiday) => (
                <div key={holiday.id} className="p-4 border border-gray-200 rounded-lg">
                  <div className="flex items-center justify-between">
                                      <div>
                    <h3 className="font-medium text-gray-900">{holiday.name}</h3>
                    <p className="text-sm text-gray-500">
                      {holiday.startDate === holiday.endDate 
                        ? new Date(holiday.startDate).toLocaleDateString('th-TH')
                        : `${new Date(holiday.startDate).toLocaleDateString('th-TH')} - ${new Date(holiday.endDate).toLocaleDateString('th-TH')}`
                      }
                    </p>
                    <p className="text-xs text-gray-400">
                      {holiday.startDate === holiday.endDate 
                        ? '1 วัน'
                        : `${Math.ceil((new Date(holiday.endDate).getTime() - new Date(holiday.startDate).getTime()) / (1000 * 3600 * 24)) + 1} วัน`
                      }
                    </p>
                  </div>
                    <div className="flex items-center space-x-2">
                      <Button variant="ghost" size="sm" onClick={() => { setEditingHoliday(holiday); setHolidayForm({ name: holiday.name, startDate: holiday.startDate, endDate: holiday.endDate }); setShowAddHolidayModal(true) }}>
                        <PencilIcon className="w-4 h-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleDeleteHoliday(Number(holiday.id), holiday.name)}
                        className="text-red-600 hover:text-red-700"
                      >
                        <TrashIcon className="w-4 h-4" />
                      </Button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Add/Edit Shift Modal */}
        {showAddShiftModal && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg p-6 w-full max-w-md">
              <h3 className="text-lg font-medium mb-4">{editingShift ? 'แก้ไขเวร' : 'เพิ่มเวรใหม่'}</h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">ชื่อเวร</label>
                  <Input
                    value={shiftForm.name}
                    onChange={(e) => setShiftForm({ ...shiftForm, name: e.target.value })}
                    placeholder="เช่น เวรเช้า"
                  />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">เวลาเริ่มต้น</label>
                    <TimePicker
                      value={dayjs(shiftForm.startTime, 'HH:mm')}
                      onChange={(time) => setShiftForm({ ...shiftForm, startTime: time?.format('HH:mm') || '07:00' })}
                      format="HH:mm"
                      className="w-full"
                      placeholder="เลือกเวลาเริ่มต้น"
                      showNow={false}
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">เวลาสิ้นสุด</label>
                    <TimePicker
                      value={dayjs(shiftForm.endTime, 'HH:mm')}
                      onChange={(time) => setShiftForm({ ...shiftForm, endTime: time?.format('HH:mm') || '15:00' })}
                      format="HH:mm"
                      className="w-full"
                      placeholder="เลือกเวลาสิ้นสุด"
                      showNow={false}
                    />
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">จำนวนพยาบาล</label>
                    <Input
                      type="number"
                      min="1"
                      value={shiftForm.nurseCount}
                      onChange={(e) => setShiftForm({ ...shiftForm, nurseCount: parseInt(e.target.value) })}
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">จำนวนผู้ช่วย</label>
                    <Input
                      type="number"
                      min="1"
                      value={shiftForm.assistantCount}
                      onChange={(e) => setShiftForm({ ...shiftForm, assistantCount: parseInt(e.target.value) })}
                    />
                  </div>
                </div>
                <div className="flex justify-end space-x-3">
                  <Button
                    variant="outline"
                    onClick={() => {
                      setShowAddShiftModal(false)
                      setShiftForm({ name: '', startTime: '07:00', endTime: '15:00', nurseCount: 1, assistantCount: 1 })
                    }}
                  >
                    ยกเลิก
                  </Button>
                  <Button onClick={handleAddOrEditShift}>{editingShift ? 'บันทึก' : 'เพิ่มเวร'}</Button>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Add/Edit Holiday Modal */}
        {showAddHolidayModal && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg p-6 w-full max-w-md">
              <h3 className="text-lg font-medium mb-4">{editingHoliday ? 'แก้ไขวันหยุด' : 'เพิ่มวันหยุดใหม่'}</h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">ชื่อวันหยุด</label>
                  <Input
                    value={holidayForm.name}
                    onChange={(e) => setHolidayForm({ ...holidayForm, name: e.target.value })}
                    placeholder="เช่น วันสงกรานต์"
                  />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">วันที่เริ่มต้น</label>
                    <DatePicker
                      value={holidayForm.startDate ? dayjs(holidayForm.startDate) : null}
                      onChange={(date) => {
                        const newStartDate = date?.format('YYYY-MM-DD') || ''
                        setHolidayForm({ 
                          ...holidayForm, 
                          startDate: newStartDate,
                          // Auto-set end date if not set yet
                          endDate: holidayForm.endDate || newStartDate
                        })
                      }}
                      className="w-full"
                      placeholder="เลือกวันที่เริ่มต้น"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">วันที่สิ้นสุด</label>
                    <DatePicker
                      value={holidayForm.endDate ? dayjs(holidayForm.endDate) : null}
                      onChange={(date) => setHolidayForm({ ...holidayForm, endDate: date?.format('YYYY-MM-DD') || '' })}
                      disabledDate={(date) => date && holidayForm.startDate ? date.isBefore(dayjs(holidayForm.startDate)) : false}
                      className="w-full"
                      placeholder="เลือกวันที่สิ้นสุด"
                    />
                  </div>
                </div>
                <div className="bg-blue-50 border border-blue-200 rounded-md p-3">
                  <p className="text-sm text-blue-800">
                    <strong>คำแนะนำ:</strong> หากเป็นวันหยุดเพียง 1 วัน ให้เลือกวันที่เดียวกันทั้งวันเริ่มต้นและสิ้นสุด
                  </p>
                </div>
                <div className="flex justify-end space-x-3">
                  <Button
                    variant="outline"
                    onClick={() => {
                      setShowAddHolidayModal(false)
                      setHolidayForm({ name: '', startDate: '', endDate: '' })
                    }}
                  >
                    ยกเลิก
                  </Button>
                  <Button onClick={handleAddOrEditHoliday}>{editingHoliday ? 'บันทึก' : 'เพิ่มวันหยุด'}</Button>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </DashboardLayout>
  )
}
