'use client'

import { useState, useEffect } from 'react'
import { 
  PlusIcon,
  PencilIcon,
  TrashIcon,
  CalendarIcon,
  UserIcon
} from '@heroicons/react/24/outline'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card'
import { Button } from '@/components/ui/Button'
import { ButtonGroup } from '@/components/ui/ButtonGroup'
import { Input } from '@/components/ui/Input'
import { Switch } from '@/components/ui/Switch'
import DashboardLayout from '@/components/layout/DashboardLayout'
import Swal from 'sweetalert2'
import { DatePicker, Select } from 'antd'
import dayjs from 'dayjs'
import 'antd/dist/reset.css'
import { leaveService } from '@/services/leaveService'
import { departmentService } from '@/services/departmentService'

// ลบ mockDepartments และใช้ข้อมูลจริงจาก API

export default function EmployeeLeavesPage() {
  const [employees, setEmployees] = useState<any[]>([])
  const [leaves, setLeaves] = useState<any[]>([])
  const [departments, setDepartments] = useState<Array<{ id: string; name: string }>>([])
  const [selectedMonth, setSelectedMonth] = useState(dayjs())
  const [selectedEmployee, setSelectedEmployee] = useState('')
  const [selectedDepartment, setSelectedDepartment] = useState('')
  
  const [showAddLeaveModal, setShowAddLeaveModal] = useState(false)
  const [editLeaveModal, setEditLeaveModal] = useState<{ open: boolean; item: any | null; reason: string }>({ open: false, item: null, reason: '' })
  
  const [leaveForm, setLeaveForm] = useState({
    departmentId: '',
    employeeId: '',
    startDate: null as dayjs.Dayjs | null,
    endDate: null as dayjs.Dayjs | null,
    reason: ''
  })

  // โหลดแผนกจาก API จริง
  useEffect(() => {
    const loadDepartments = async () => {
      try {
        const list = await departmentService.getDepartments()
        setDepartments(list.map((d: any) => ({ id: d.id, name: d.name })))
      } catch (e) {
        console.error('โหลดแผนกไม่สำเร็จ', e)
        setDepartments([])
      }
    }
    loadDepartments()
  }, [])

  // โหลดพนักงานตามแผนกที่เลือก (ถ้ามี) ทั้งในฟิลเตอร์หรือในฟอร์ม
  useEffect(() => {
    const depId = leaveForm.departmentId || selectedDepartment
    if (!depId) {
      setEmployees([])
      return
    }
    const loadEmployees = async () => {
      try {
        const staff = await departmentService.getDepartmentStaff(depId)
        const deptName = departments.find(d => d.id === depId)?.name || ''
        const mapped = (staff || []).map((s: any) => {
          const name = s.name || ''
          const parts = name.split(' ')
          const firstName = parts[0] || name
          const lastName = parts.slice(1).join(' ')
          return {
            id: String(s.id),
            firstName,
            lastName,
            position: s.position,
            departmentId: s.department_id,
            departmentName: deptName,
          }
        })
        setEmployees(mapped)
      } catch (e) {
        console.error('โหลดพนักงานไม่สำเร็จ', e)
        setEmployees([])
      }
    }
    loadEmployees()
  }, [leaveForm.departmentId, selectedDepartment, departments])

  // Load leaves from API (filter by month on server; allow optional department/employee filters)
  useEffect(() => {
    const load = async () => {
      try {
        const month = selectedMonth ? selectedMonth.format('YYYY-MM') : undefined
        const data = await leaveService.list({ 
          month,
          departmentId: selectedDepartment || undefined,
          employeeId: selectedEmployee || undefined,
        })
        const items = (data || []).map((l: any) => ({
          id: l.id as string,
          ids: [l.id as string],
          employeeId: String(l.userId),
          startDate: l.startDate,
          endDate: l.endDate ?? l.startDate,
          reason: l.reason || '-',
          isActive: l.status !== 'cancelled',
          status: l.status,
          employeeName: l.userName,
          departmentId: l.departmentId,
          departmentName: l.departmentName,
        }))

        items.sort((a: any, b: any) => {
          if (a.employeeId !== b.employeeId) return a.employeeId.localeCompare(b.employeeId)
          if (a.departmentId !== b.departmentId) return a.departmentId.localeCompare(b.departmentId)
          return String(a.startDate).localeCompare(String(b.startDate))
        })

        const grouped: any[] = []
        for (const it of items) {
          const last = grouped[grouped.length - 1]
          const contiguous = last && last.employeeId === it.employeeId && last.departmentId === it.departmentId && last.reason === it.reason && last.isActive === it.isActive && dayjs(it.startDate).diff(dayjs(last.endDate), 'day') === 1
          if (contiguous) {
            last.endDate = it.endDate
            last.ids.push(it.id)
          } else {
            grouped.push({ ...it })
          }
        }
        setLeaves(grouped)
      } catch (e) {
        console.error('Failed to load leaves', e)
        setLeaves([])
      }
    }
    load()
  }, [selectedMonth, selectedDepartment, selectedEmployee])

  const handleToggleLeave = async (leaveItem: any) => {
    try {
      await Promise.all((leaveItem.ids || [leaveItem.id]).map((id: string) => leaveService.toggle(String(id))))
      setLeaves(prev => prev.map((l: any) => l.id === leaveItem.id ? { ...l, isActive: !leaveItem.isActive } : l))
      await Swal.fire({ icon: 'success', title: leaveItem.isActive ? 'ปิดการใช้งานวันหยุดแล้ว' : 'เปิดการใช้งานวันหยุดแล้ว', timer: 1200, showConfirmButton: false })
    } catch (e: any) {
      console.error(e)
      await Swal.fire({ icon: 'error', title: 'สลับสถานะไม่สำเร็จ', text: e.message || 'เกิดข้อผิดพลาด' })
    }
  }

  const handleAddLeave = async () => {
    if (!leaveForm.departmentId || !leaveForm.employeeId || !leaveForm.startDate || !leaveForm.endDate) {
      await Swal.fire({ icon: 'warning', title: 'กรุณากรอกข้อมูลให้ครบ', text: 'กรุณาเลือกแผนก พนักงาน และวันที่เริ่มต้น-สิ้นสุด', confirmButtonColor: '#2563eb' })
      return
    }

    if (leaveForm.startDate.isAfter(leaveForm.endDate)) {
      await Swal.fire({ icon: 'warning', title: 'วันที่ไม่ถูกต้อง', text: 'วันที่เริ่มต้นต้องไม่เกินวันที่สิ้นสุด', confirmButtonColor: '#2563eb' })
      return
    }

    const employee = employees.find(emp => String(emp.id) === leaveForm.employeeId)
    if (!employee) return

    const startDate = leaveForm.startDate.startOf('day')
    const endDate = leaveForm.endDate.startOf('day')
    const daysDiff = endDate.diff(startDate, 'day') + 1

    try {
      const tasks: Promise<any>[] = []
      for (let i = 0; i < daysDiff; i++) {
        const currentDate = startDate.add(i, 'day').format('YYYY-MM-DD')
        tasks.push(
          leaveService.create({
            employeeId: leaveForm.employeeId,
            employeeName: `${employee.firstName} ${employee.lastName}`,
            departmentId: leaveForm.departmentId,
            departmentName: employee.departmentName,
            date: currentDate,
            reason: leaveForm.reason || undefined,
          })
        )
      }
      await Promise.all(tasks)

      // Reload list with current filters (use the same mapping + grouping as initial load)
      const month = selectedMonth ? selectedMonth.format('YYYY-MM') : undefined
      const data = await leaveService.list({ 
        month, 
        departmentId: selectedDepartment || undefined, 
        employeeId: selectedEmployee || undefined,
      })
      const items = (data || []).map((l: any) => ({
        id: l.id as string,
        ids: [l.id as string],
        employeeId: String(l.userId),
        startDate: l.startDate,
        endDate: l.endDate ?? l.startDate,
        reason: l.reason || '-',
        isActive: l.status !== 'cancelled',
        status: l.status,
        employeeName: l.userName,
        departmentId: l.departmentId,
        departmentName: l.departmentName,
      }))

      items.sort((a: any, b: any) => {
        if (a.employeeId !== b.employeeId) return a.employeeId.localeCompare(b.employeeId)
        if (a.departmentId !== b.departmentId) return a.departmentId.localeCompare(b.departmentId)
        return String(a.startDate).localeCompare(String(b.startDate))
      })

      const grouped: any[] = []
      for (const it of items) {
        const last = grouped[grouped.length - 1]
        const contiguous = last && last.employeeId === it.employeeId && last.departmentId === it.departmentId && last.reason === it.reason && last.isActive === it.isActive && dayjs(it.startDate).diff(dayjs(last.endDate), 'day') === 1
        if (contiguous) {
          last.endDate = it.endDate
          last.ids.push(it.id)
        } else {
          grouped.push({ ...it })
        }
      }
      setLeaves(grouped)

      setLeaveForm({ departmentId: '', employeeId: '', startDate: null, endDate: null, reason: '' })
      setShowAddLeaveModal(false)
      await Swal.fire({ icon: 'success', title: `เพิ่มวันหยุดพนักงานสำเร็จ! (${daysDiff} วัน)`, confirmButtonColor: '#2563eb' })
    } catch (e: any) {
      console.error(e)
      await Swal.fire({ icon: 'error', title: 'บันทึกไม่สำเร็จ', text: e.message || 'เกิดข้อผิดพลาด' })
    }
  }

  const handleDeleteLeave = async (leaveItem: any) => {
    const result = await Swal.fire({
      title: 'ลบวันหยุด?',
      text: `คุณต้องการลบวันหยุดของ "${leaveItem.employeeName}" ช่วงวันที่ ${new Date(leaveItem.startDate).toLocaleDateString('th-TH')} - ${new Date(leaveItem.endDate).toLocaleDateString('th-TH')} หรือไม่?`,
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#dc2626',
      cancelButtonColor: '#6b7280',
      confirmButtonText: 'ลบวันหยุด',
      cancelButtonText: 'ยกเลิก'
    })

    if (result.isConfirmed) {
      try {
        const ids: string[] = leaveItem.ids || [leaveItem.id]
        await Promise.all(ids.map((id) => leaveService.remove(String(id))))
        setLeaves(leaves.filter(l => l.id !== leaveItem.id))
        await Swal.fire({ icon: 'success', title: 'ลบวันหยุดสำเร็จ!', confirmButtonColor: '#2563eb' })
      } catch (e: any) {
        console.error(e)
        await Swal.fire({ icon: 'error', title: 'ลบไม่สำเร็จ', text: e.message || 'เกิดข้อผิดพลาด' })
      }
    }
  }

  const filteredEmployeesForForm = employees.filter(emp => 
    leaveForm.departmentId ? emp.departmentId === leaveForm.departmentId : true
  )

  const filteredLeaves = leaves.filter(leave => {
    const matchMonth = true
    const matchEmployee = selectedEmployee ? String(leave.employeeId) === selectedEmployee : true
    const matchDepartment = selectedDepartment ? leave.departmentId === selectedDepartment : true
    return matchMonth && matchEmployee && matchDepartment
  })

  const filteredEmployeesForFilter = employees.filter(emp => 
    selectedDepartment ? emp.departmentId === selectedDepartment : true
  )

  return (
    <DashboardLayout>
      <div className="space-y-8">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">จัดการวันหยุดพนักงาน</h1>
            <p className="text-gray-600 mt-2">เพิ่ม แก้ไข และจัดการวันหยุดของพนักงานในแผนก</p>
          </div>
          <Button onClick={() => setShowAddLeaveModal(true)}>
            <PlusIcon className="w-4 h-4 mr-2" />
            เพิ่มวันหยุด
          </Button>
        </div>

        {/* Filters */}
        <Card className="border-0 shadow-md">
          <CardHeader>
            <CardTitle>ตัวกรองข้อมูล</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  เลือกแผนก
                </label>
                <select
                  value={selectedDepartment}
                  onChange={(e) => {
                    setSelectedDepartment(e.target.value)
                    setSelectedEmployee('')
                  }}
                  className="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 h-10"
                >
                  <option value="">ทุกแผนก</option>
                  {departments.map((dept) => (
                    <option key={dept.id} value={dept.id}>
                      {dept.name}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  เลือกเดือน
                </label>
                <DatePicker
                  picker="month"
                  value={selectedMonth}
                  onChange={(date) => setSelectedMonth(date!)}
                  className="w-full h-10"
                  placeholder="เลือกเดือน"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  เลือกพนักงาน
                </label>
                <select
                  value={selectedEmployee}
                  onChange={(e) => setSelectedEmployee(e.target.value)}
                  className="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 h-10"
                  disabled={!selectedDepartment}
                >
                  <option value="">{selectedDepartment ? 'ทุกคน' : 'กรุณาเลือกแผนกก่อน'}</option>
                  {filteredEmployeesForFilter.map((employee) => (
                    <option key={employee.id} value={employee.id}>
                      {employee.firstName} {employee.lastName} ({employee.position})
                    </option>
                  ))}
                </select>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Leaves List */}
        <Card className="border-0 shadow-md">
          <CardHeader>
            <CardTitle className="flex items-center">
              <CalendarIcon className="w-5 h-5 mr-2" />
              รายการวันหยุดพนักงาน ({filteredLeaves.length} รายการ)
            </CardTitle>
            <CardDescription>
              วันหยุดของพนักงานในเดือน {selectedMonth ? selectedMonth.format('MMMM YYYY') : 'ทั้งหมด'}
            </CardDescription>
          </CardHeader>
          <CardContent>
            {filteredLeaves.length > 0 ? (
              <div className="space-y-4">
                {filteredLeaves.map((leave) => (
                  <div key={leave.id} className="p-4 border border-gray-200 rounded-lg">
                    <div className="flex items-center justify-between">
                      <div className="flex-1">
                        <div className="flex items-center space-x-3 mb-2">
                          <UserIcon className="w-5 h-5 text-gray-400" />
                          <h3 className="font-medium text-gray-900">{leave.employeeName}</h3>
                        </div>
                        <div className="ml-8 space-y-1">
                          <p className="text-sm text-gray-600">
                            <strong>แผนก:</strong> {leave.departmentName}
                          </p>
                          <p className="text-sm text-gray-600">
                            <strong>วันที่:</strong> {new Date(leave.startDate).toLocaleDateString('th-TH')} {dayjs(leave.endDate).isAfter(dayjs(leave.startDate),'day') ? `- ${new Date(leave.endDate).toLocaleDateString('th-TH')}` : ''}
                          </p>
                          <p className="text-sm text-gray-600">
                            <strong>เหตุผล:</strong> {leave.reason}
                          </p>
                          <p className="text-sm">
                            <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                              leave.isActive 
                                ? 'bg-green-100 text-green-800' 
                                : 'bg-red-100 text-red-800'
                            }`}>
                              {leave.isActive ? 'ใช้งาน' : 'ไม่ใช้งาน'}
                            </span>
                          </p>
                        </div>
                      </div>
                      <div className="flex items-center space-x-3">
                          <div className="flex items-center space-x-2">
                          <span className="text-sm text-gray-600">
                            {leave.isActive ? 'เปิดใช้งาน' : 'ปิดใช้งาน'}
                          </span>
                          <Switch
                              checked={leave.isActive}
                              onChange={() => handleToggleLeave(leave)}
                          />
                        </div>
                        <ButtonGroup direction="horizontal" spacing="tight">
                          <Button 
                            variant="ghost" 
                            size="sm"
                            onClick={() => setEditLeaveModal({ open: true, item: leave, reason: leave.reason || '' })}
                          >
                            <PencilIcon className="w-4 h-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => handleDeleteLeave(leave)}
                            className="text-red-600 hover:text-red-700"
                          >
                            <TrashIcon className="w-4 h-4" />
                          </Button>
                        </ButtonGroup>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-8 text-gray-500">
                ไม่มีข้อมูลวันหยุดพนักงานในช่วงที่เลือก
              </div>
            )}
          </CardContent>
        </Card>

        {/* Add Leave Modal */}
        {showAddLeaveModal && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-md shadow-xl">
              <h3 className="text-lg font-medium mb-4 text-gray-900 dark:text-white">เพิ่มวันหยุดพนักงาน</h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    เลือกแผนก
                  </label>
                  <select
                    value={leaveForm.departmentId}
                    onChange={(e) => setLeaveForm({ 
                      ...leaveForm, 
                      departmentId: e.target.value,
                      employeeId: ''
                    })}
                    className="w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 h-10 bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                  >
                    <option value="">เลือกแผนก</option>
                    {departments.map((dept) => (
                      <option key={dept.id} value={dept.id}>
                        {dept.name}
                      </option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    เลือกพนักงาน
                  </label>
                  <select
                    value={leaveForm.employeeId}
                    onChange={(e) => setLeaveForm({ ...leaveForm, employeeId: e.target.value })}
                    className="w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 h-10 bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                    disabled={!leaveForm.departmentId}
                  >
                    <option value="">
                      {!leaveForm.departmentId ? 'กรุณาเลือกแผนกก่อน' : 'เลือกพนักงาน'}
                    </option>
                    {filteredEmployeesForForm.map((employee) => (
                      <option key={employee.id} value={employee.id}>
                        {employee.firstName} {employee.lastName} ({employee.position})
                      </option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    วันที่เริ่มต้น
                  </label>
                  <DatePicker
                    value={leaveForm.startDate}
                    onChange={(date) => {
                      // หากเลือกวันเริ่มใหม่ แล้ววันสิ้นสุดก่อนหน้าอยู่ก่อนวันเริ่ม ให้รีเซ็ต/ขยับวันสิ้นสุด
                      if (date && leaveForm.endDate && date.isAfter(leaveForm.endDate, 'day')) {
                        setLeaveForm({ ...leaveForm, startDate: date, endDate: date })
                      } else {
                        setLeaveForm({ ...leaveForm, startDate: date })
                      }
                    }}
                    className="w-full h-10"
                    placeholder="เลือกวันที่เริ่มต้น"
                    disabledDate={(current) => {
                      if (!current) return false
                      if (leaveForm.endDate) {
                        return current.isAfter(leaveForm.endDate, 'day')
                      }
                      return false
                    }}
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    วันที่สิ้นสุด
                  </label>
                  <DatePicker
                    value={leaveForm.endDate}
                    onChange={(date) => setLeaveForm({ ...leaveForm, endDate: date })}
                    className="w-full h-10"
                    placeholder="เลือกวันที่สิ้นสุด"
                    disabled={!leaveForm.startDate}
                    disabledDate={(current) => {
                      if (!current) return false
                      if (leaveForm.startDate) {
                        return current.isBefore(leaveForm.startDate, 'day')
                      }
                      return false
                    }}
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    เหตุผล (ไม่บังคับ)
                  </label>
                  <Input
                    value={leaveForm.reason}
                    onChange={(e) => setLeaveForm({ ...leaveForm, reason: e.target.value })}
                    placeholder="เช่น ลาป่วย, ลากิจ, ลาคลอด"
                    className="dark:bg-gray-700 dark:border-gray-600 dark:text-white h-10"
                  />
                </div>
                <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-md p-3">
                  <p className="text-sm text-blue-800 dark:text-blue-200">
                    <strong>คำแนะนำ:</strong> เลือกแผนกก่อน แล้วค่อยเลือกพนักงานในแผนกนั้น
                  </p>
                </div>
                <div className="flex justify-end mt-6">
                  <Button
                    variant="outline"
                    onClick={() => {
                      setShowAddLeaveModal(false)
                      setLeaveForm({ departmentId: '', employeeId: '', startDate: null, endDate: null, reason: '' })
                    }}
                    className="px-6 py-2 mr-16"
                  >
                    ยกเลิก
                  </Button>
                  <Button 
                    onClick={handleAddLeave}
                    className="px-6 py-2"
                  >
                    เพิ่มวันหยุด
                  </Button>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Edit Leave Modal */}
        {editLeaveModal.open && editLeaveModal.item && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-md shadow-xl">
              <h3 className="text-lg font-medium mb-4 text-gray-900 dark:text-white">แก้ไขวันหยุดพนักงาน</h3>
              <div className="space-y-4">
                <div>
                  <p className="text-sm text-gray-600 dark:text-gray-300">
                    พนักงาน: {editLeaveModal.item.employeeName}
                  </p>
                  <p className="text-sm text-gray-600 dark:text-gray-300">
                    ช่วงวันที่: {new Date(editLeaveModal.item.startDate).toLocaleDateString('th-TH')} {dayjs(editLeaveModal.item.endDate).isAfter(dayjs(editLeaveModal.item.startDate),'day') ? `- ${new Date(editLeaveModal.item.endDate).toLocaleDateString('th-TH')}` : ''}
                  </p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">เหตุผล</label>
                  <Input
                    value={editLeaveModal.reason}
                    onChange={(e) => setEditLeaveModal(prev => ({ ...prev, reason: e.target.value }))}
                    placeholder="เช่น ลาป่วย, ลากิจ, ลาคลอด"
                    className="dark:bg-gray-700 dark:border-gray-600 dark:text-white h-10"
                  />
                </div>
                <div className="flex justify-end mt-6">
                  <Button
                    variant="outline"
                    onClick={() => setEditLeaveModal({ open: false, item: null, reason: '' })}
                    className="px-6 py-2 mr-16"
                  >
                    ยกเลิก
                  </Button>
                  <Button
                    onClick={async () => {
                      try {
                        const ids: string[] = editLeaveModal.item.ids || [editLeaveModal.item.id]
                        await Promise.all(ids.map((id) => leaveService.update(id, { reason: editLeaveModal.reason || undefined })))

                        // Reload list with current filters
                        const month = selectedMonth ? selectedMonth.format('YYYY-MM') : undefined
                        const data = await leaveService.list({ 
                          month, 
                          departmentId: selectedDepartment || undefined, 
                          employeeId: selectedEmployee || undefined,
                        })
                        const items = (data || []).map((l: any) => ({
                          id: l.id as string,
                          ids: [l.id as string],
                          employeeId: String(l.userId),
                          startDate: l.startDate,
                          endDate: l.endDate ?? l.startDate,
                          reason: l.reason || '-',
                          isActive: l.status !== 'cancelled',
                          status: l.status,
                          employeeName: l.userName,
                          departmentId: l.departmentId,
                          departmentName: l.departmentName,
                        }))
                        items.sort((a: any, b: any) => {
                          if (a.employeeId !== b.employeeId) return a.employeeId.localeCompare(b.employeeId)
                          if (a.departmentId !== b.departmentId) return a.departmentId.localeCompare(b.departmentId)
                          return String(a.startDate).localeCompare(String(b.startDate))
                        })
                        const grouped: any[] = []
                        for (const it of items) {
                          const last = grouped[grouped.length - 1]
                          const contiguous = last && last.employeeId === it.employeeId && last.departmentId === it.departmentId && last.reason === it.reason && last.isActive === it.isActive && dayjs(it.startDate).diff(dayjs(last.endDate), 'day') === 1
                          if (contiguous) {
                            last.endDate = it.endDate
                            last.ids.push(it.id)
                          } else {
                            grouped.push({ ...it })
                          }
                        }
                        setLeaves(grouped)

                        setEditLeaveModal({ open: false, item: null, reason: '' })
                        await Swal.fire({ icon: 'success', title: 'แก้ไขวันหยุดสำเร็จ', timer: 1200, showConfirmButton: false })
                      } catch (e: any) {
                        console.error(e)
                        await Swal.fire({ icon: 'error', title: 'บันทึกไม่สำเร็จ', text: e.message || 'เกิดข้อผิดพลาด' })
                      }
                    }}
                    className="px-6 py-2"
                  >
                    บันทึก
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
