'use client'

import { useState, useEffect } from 'react'
import { 
  PlusIcon,
  PencilIcon,
  TrashIcon,
  UserGroupIcon,
  BuildingOfficeIcon
} from '@heroicons/react/24/outline'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card'
import { Button } from '@/components/ui/Button'
import { ButtonGroup } from '@/components/ui/ButtonGroup'
import { Input } from '@/components/ui/Input'
import DashboardLayout from '@/components/layout/DashboardLayout'
import Swal from 'sweetalert2'
import { departmentService, DepartmentWithStats, CreateDepartmentRequest } from '@/services/departmentService'

// Mock data for employees (will be replaced with real API later)
const mockEmployees = [
  { id: '1', firstName: 'สมหญิง', lastName: 'ใจดี', position: 'nurse', departmentId: '1' },
  { id: '2', firstName: 'สมชาย', lastName: 'ขยัน', position: 'assistant', departmentId: '1' },
  { id: '3', firstName: 'สมใส', lastName: 'เก่งกาจ', position: 'nurse', departmentId: '2' },
  { id: '4', firstName: 'สมศรี', lastName: 'มานะ', position: 'nurse', departmentId: '2' },
  { id: '5', firstName: 'สมพร', lastName: 'อุตสาหะ', position: 'assistant', departmentId: '3' }
]

export default function DepartmentsPage() {
  const [departments, setDepartments] = useState<DepartmentWithStats[]>([])
  const [selectedDepartment, setSelectedDepartment] = useState<string | null>(null)
  const [selectedDepartmentData, setSelectedDepartmentData] = useState<DepartmentWithStats | null>(null)
  const [employees, setEmployees] = useState<{ id: string; firstName: string; lastName: string; position: string; departmentId: string }[]>(mockEmployees)
  const [departmentEmployees, setDepartmentEmployees] = useState<typeof mockEmployees>([])
  const [showAddDepartmentModal, setShowAddDepartmentModal] = useState(false)
  const [showAddEmployeeModal, setShowAddEmployeeModal] = useState(false)
  const [showEditEmployeeModal, setShowEditEmployeeModal] = useState(false)
  const [editingEmployee, setEditingEmployee] = useState<any>(null)
  const [showEditDepartmentModal, setShowEditDepartmentModal] = useState(false)
  const [editingDepartment, setEditingDepartment] = useState<DepartmentWithStats | null>(null)
  const [loading, setLoading] = useState(true)
  
  const [departmentForm, setDepartmentForm] = useState({ 
    name: '', 
    description: '',
    max_nurses: 1000,
    max_assistants: 500
  })
  const [employeeForm, setEmployeeForm] = useState({
    firstName: '',
    lastName: '',
    position: 'nurse' as 'nurse' | 'assistant'
  })

  // Load departments from API
  useEffect(() => {
    const loadDepartments = async () => {
      try {
        setLoading(true)
        const data = await departmentService.getDepartments()
        setDepartments(data)
      } catch (error) {
        console.error('Error loading departments:', error)
        await Swal.fire({
          icon: 'error',
          title: 'เกิดข้อผิดพลาด',
          text: 'ไม่สามารถโหลดข้อมูลแผนกได้',
          confirmButtonColor: '#dc2626'
        })
      } finally {
        setLoading(false)
      }
    }

    loadDepartments()
  }, [])

  // Load department staff when a department is selected
  useEffect(() => {
    console.log('useEffect triggered with selectedDepartment:', selectedDepartment)
    
    const loadDepartmentStaff = async () => {
      if (!selectedDepartment) {
        console.log('No selectedDepartment, setting empty array and returning')
        setDepartmentEmployees([])
        return
      }

      try {
        setLoading(true)
        console.log('Loading staff for department:', selectedDepartment)
        
        const staff = await departmentService.getDepartmentStaff(selectedDepartment)
        console.log('Received staff data:', staff)
        console.log('Staff type:', typeof staff)
        console.log('Staff isArray:', Array.isArray(staff))
        
        if (!staff || !Array.isArray(staff) || staff.length === 0) {
          console.log('No staff found or staff is not an array, setting empty array')
          setDepartmentEmployees([])
          return
        }
        
        // Convert DepartmentStaff to frontend Employee format
        const convertedEmployees = staff.map(s => {
          console.log('Processing staff member:', s)
          return {
            id: s.id, // Keep as UUID string - don't convert to number
            firstName: s.name.split(' ')[0] || s.name,
            lastName: s.name.split(' ').slice(1).join(' ') || '',
            position: s.position as 'nurse' | 'assistant',
            departmentId: s.department_id // Keep as UUID string - don't convert to number
          }
        })
        
        console.log('Converted employees:', convertedEmployees)
        setDepartmentEmployees(convertedEmployees)
      } catch (error) {
        console.error('Error loading department staff:', error)
        console.error('Error details:', {
          message: (error as Error).message,
          stack: (error as Error).stack,
          selectedDepartment
        })
        setDepartmentEmployees([])
        await Swal.fire({
          icon: 'error',
          title: 'เกิดข้อผิดพลาด',
          text: 'ไม่สามารถโหลดข้อมูลพนักงานในแผนกได้',
          confirmButtonColor: '#dc2626'
        })
      } finally {
        setLoading(false)
      }
    }

    loadDepartmentStaff()
  }, [selectedDepartment])

  const handleAddDepartment = async () => {
    if (!departmentForm.name.trim()) {
      await Swal.fire({
        icon: 'warning',
        title: 'กรุณากรอกชื่อแผนก',
        confirmButtonColor: '#2563eb'
      })
      return
    }

    try {
      const createData: CreateDepartmentRequest = {
        name: departmentForm.name,
        description: departmentForm.description || undefined,
        max_nurses: 1000, // Default value - ไม่จำกัดจำนวน
        max_assistants: 500 // Default value - ไม่จำกัดจำนวน
      }

      const newDepartment = await departmentService.createDepartment(createData)
      
      // แสดงข้อความสำเร็จก่อน
      await Swal.fire({
        icon: 'success',
        title: 'เพิ่มแผนกสำเร็จ!',
        text: `แผนก "${departmentForm.name}" ถูกเพิ่มแล้ว`,
        confirmButtonColor: '#2563eb'
      })

      // ลองดึงรายการแผนกใหม่ (ถ้าล้มเหลวก็ไม่เป็นไร)
      try {
        const updatedDepartments = await departmentService.getDepartments()
        setDepartments(updatedDepartments)
      } catch (refreshError) {
        console.warn('Failed to refresh departments list:', refreshError)
        // เพิ่มแผนกใหม่เข้าไปในรายการปัจจุบัน (แปลง Department เป็น DepartmentWithStats)
        const newDepartmentWithStats = {
          ...newDepartment,
          total_employees: 0,
          active_employees: 0,
          nurse_count: 0,
          assistant_count: 0
        }
        setDepartments(prev => [newDepartmentWithStats, ...prev])
      }
      
      setDepartmentForm({ name: '', description: '', max_nurses: 1000, max_assistants: 500 })
      setShowAddDepartmentModal(false)

    } catch (error) {
      console.error('Error creating department:', error)
      await Swal.fire({
        icon: 'error',
        title: 'เกิดข้อผิดพลาด',
        text: 'ไม่สามารถเพิ่มแผนกได้',
        confirmButtonColor: '#dc2626'
      })
    }
  }

  const handleEditDepartment = async () => {
    if (!editingDepartment || !departmentForm.name.trim()) {
      await Swal.fire({
        icon: 'warning',
        title: 'กรุณากรอกชื่อแผนก',
        confirmButtonColor: '#2563eb'
      })
      return
    }

    try {
      const updateData = {
        name: departmentForm.name,
        description: departmentForm.description || undefined,
        max_nurses: departmentForm.max_nurses,
        max_assistants: departmentForm.max_assistants
      }

      await departmentService.updateDepartment(editingDepartment.id, updateData)
      
      // แสดงข้อความสำเร็จ
      await Swal.fire({
        icon: 'success',
        title: 'แก้ไขแผนกสำเร็จ!',
        text: `แผนก "${departmentForm.name}" ถูกอัปเดตแล้ว`,
        confirmButtonColor: '#2563eb'
      })

      // รีเฟรชรายการแผนก
      try {
        const updatedDepartments = await departmentService.getDepartments()
        setDepartments(updatedDepartments)
      } catch (refreshError) {
        console.warn('Failed to refresh departments list:', refreshError)
      }
      
      // รีเซ็ตฟอร์มและปิด modal
      setDepartmentForm({ name: '', description: '', max_nurses: 1000, max_assistants: 500 })
      setShowEditDepartmentModal(false)
      setEditingDepartment(null)

    } catch (error) {
      console.error('Error updating department:', error)
      await Swal.fire({
        icon: 'error',
        title: 'เกิดข้อผิดพลาด',
        text: 'ไม่สามารถแก้ไขแผนกได้',
        confirmButtonColor: '#dc2626'
      })
    }
  }

  const openEditModal = (department: DepartmentWithStats) => {
    setEditingDepartment(department)
    setDepartmentForm({
      name: department.name,
      description: department.description || '',
      max_nurses: department.max_nurses,
      max_assistants: department.max_assistants
    })
    setShowEditDepartmentModal(true)
  }

  const handleDeleteDepartment = async (departmentId: string, departmentName: string) => {
    const result = await Swal.fire({
      title: 'ลบแผนก?',
      text: `คุณต้องการลบแผนก "${departmentName}" หรือไม่?`,
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#dc2626',
      cancelButtonColor: '#6b7280',
      confirmButtonText: 'ลบแผนก',
      cancelButtonText: 'ยกเลิก'
    })

    if (result.isConfirmed) {
      try {
        await departmentService.deleteDepartment(departmentId)
        
        // Refresh departments list
        const updatedDepartments = await departmentService.getDepartments()
        setDepartments(updatedDepartments)
        
        setEmployees(employees.filter(e => e.departmentId !== departmentId))
        if (selectedDepartment === departmentId) {
          setSelectedDepartment(null)
        }

        await Swal.fire({
          icon: 'success',
          title: 'ลบแผนกสำเร็จ!',
          confirmButtonColor: '#2563eb'
        })
      } catch (error) {
        console.error('Error deleting department:', error)
        await Swal.fire({
          icon: 'error',
          title: 'เกิดข้อผิดพลาด',
          text: 'ไม่สามารถลบแผนกได้',
          confirmButtonColor: '#dc2626'
        })
      }
    }
  }

  const handleAddEmployee = async () => {
    if (!employeeForm.firstName.trim() || !employeeForm.lastName.trim()) {
      await Swal.fire({
        icon: 'warning',
        title: 'กรุณากรอกข้อมูลให้ครบ',
        confirmButtonColor: '#2563eb'
      })
      return
    }

    if (!selectedDepartment) {
      await Swal.fire({
        icon: 'warning',
        title: 'กรุณาเลือกแผนกก่อน',
        confirmButtonColor: '#2563eb'
      })
      return
    }

    try {
      // เรียกใช้ API เพื่อเพิ่มพนักงาน
      const newStaff = await departmentService.addDepartmentStaff(selectedDepartment, {
        first_name: employeeForm.firstName,
        last_name: employeeForm.lastName,
        position: employeeForm.position,
        phone: '',
        email: ''
      })

      // แสดงข้อความสำเร็จ
      await Swal.fire({
        icon: 'success',
        title: 'เพิ่มพนักงานสำเร็จ!',
        confirmButtonColor: '#2563eb'
      })

      // รีเฟรชรายการพนักงานในแผนก
      try {
        console.log('Refreshing staff list for department:', selectedDepartment)
        const updatedStaff = await departmentService.getDepartmentStaff(selectedDepartment)
        console.log('Updated staff data:', updatedStaff)
        
        // Convert DepartmentStaff to frontend Employee format
        const convertedEmployees = updatedStaff.map(s => ({
          id: s.id,
          firstName: s.name.split(' ')[0] || s.name,
          lastName: s.name.split(' ').slice(1).join(' ') || '',
          position: s.position as 'nurse' | 'assistant',
          departmentId: s.department_id
        }))
        
        console.log('Converted updated employees:', convertedEmployees)
        setDepartmentEmployees(convertedEmployees)
        
        // Refresh departments list to get updated stats
        try {
          const updatedDepartments = await departmentService.getDepartments()
          setDepartments(updatedDepartments)
        } catch (refreshDeptError) {
          console.warn('Failed to refresh departments list:', refreshDeptError)
        }
      } catch (refreshError) {
        console.warn('Failed to refresh staff list:', refreshError)
        // เพิ่มพนักงานใหม่เข้าไปในรายการปัจจุบัน
        if (selectedDepartment) {
          const newEmployee = {
            id: newStaff.id,
            firstName: employeeForm.firstName,
            lastName: employeeForm.lastName,
            position: employeeForm.position,
            departmentId: selectedDepartment
          }
          setDepartmentEmployees(prev => [...prev, newEmployee])
        }
      }

      // รีเซ็ตฟอร์มและปิด modal
      setEmployeeForm({ firstName: '', lastName: '', position: 'nurse' })
      setShowAddEmployeeModal(false)

    } catch (error) {
      console.error('Error adding employee:', error)
      await Swal.fire({
        icon: 'error',
        title: 'เกิดข้อผิดพลาด',
        text: 'ไม่สามารถเพิ่มพนักงานได้',
        confirmButtonColor: '#dc2626'
      })
    }
  }

  const handleDeleteEmployee = async (employeeId: string, employeeName: string) => {
    const result = await Swal.fire({
      title: 'ลบพนักงาน?',
      text: `คุณต้องการลบ "${employeeName}" หรือไม่?`,
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#dc2626',
      cancelButtonColor: '#6b7280',
      confirmButtonText: 'ลบพนักงาน',
      cancelButtonText: 'ยกเลิก'
    })

    if (result.isConfirmed) {
      try {
        // Find employee to get department info
        const employee = departmentEmployees.find(e => e.id === employeeId)
        if (!employee) {
          throw new Error('ไม่พบข้อมูลพนักงาน')
        }

        // Call API to delete employee from backend
        if (!selectedDepartment) {
          throw new Error('ไม่พบข้อมูลแผนก')
        }
        await departmentService.deleteDepartmentStaff(selectedDepartment, employeeId)
        
        // Remove from local state
        setDepartmentEmployees(prev => prev.filter(e => e.id !== employeeId))
        
        // Refresh departments list to get updated stats
        try {
          const updatedDepartments = await departmentService.getDepartments()
          setDepartments(updatedDepartments)
        } catch (refreshError) {
          console.warn('Failed to refresh departments list:', refreshError)
        }

        await Swal.fire({
          icon: 'success',
          title: 'ลบพนักงานสำเร็จ!',
          confirmButtonColor: '#2563eb'
        })
      } catch (error) {
        console.error('Error deleting employee:', error)
        await Swal.fire({
          icon: 'error',
          title: 'เกิดข้อผิดพลาด',
          text: 'ไม่สามารถลบพนักงานได้',
          confirmButtonColor: '#dc2626'
        })
      }
    }
  }

  const handleEditEmployee = async (employee: any) => {
    // Set editing employee
    setEditingEmployee(employee)
    
    // Set employee form with current employee data
    setEmployeeForm({
      firstName: employee.firstName,
      lastName: employee.lastName,
      position: employee.position
    })
    
    // Show edit employee modal
    setShowEditEmployeeModal(true)
  }

  const handleUpdateEmployee = async () => {
    if (!editingEmployee || !selectedDepartment) {
      await Swal.fire({
        icon: 'error',
        title: 'เกิดข้อผิดพลาด',
        text: 'ไม่พบข้อมูลพนักงานหรือแผนก',
        confirmButtonColor: '#dc2626'
      })
      return
    }

    try {
      // TODO: Call API to update employee
      // await departmentService.updateDepartmentStaff(editingEmployee.id, employeeForm)
      
      // Update local state
      setDepartmentEmployees(prev => 
        prev.map(emp => 
          emp.id === editingEmployee.id 
            ? { ...emp, ...employeeForm }
            : emp
        )
      )

      // Show success message
      await Swal.fire({
        icon: 'success',
        title: 'แก้ไขพนักงานสำเร็จ!',
        text: `พนักงาน "${employeeForm.firstName} ${employeeForm.lastName}" ถูกอัปเดตแล้ว`,
        confirmButtonColor: '#2563eb'
      })

      // Reset form and close modal
      setEmployeeForm({ firstName: '', lastName: '', position: 'nurse' })
      setShowEditEmployeeModal(false)
      setEditingEmployee(null)

    } catch (error) {
      console.error('Error updating employee:', error)
      await Swal.fire({
        icon: 'error',
        title: 'เกิดข้อผิดพลาด',
        text: 'ไม่สามารถแก้ไขพนักงานได้',
        confirmButtonColor: '#dc2626'
      })
    }
  }

  return (
    <DashboardLayout>
      <div className="space-y-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">จัดการแผนกและพนักงาน</h1>
          <p className="text-gray-600 mt-2">เพิ่ม แก้ไข และลบข้อมูลแผนกและพนักงาน</p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Departments List */}
          <Card className="border-0 shadow-md">
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle className="flex items-center">
                    <BuildingOfficeIcon className="w-5 h-5 mr-2" />
                    รายการแผนก
                  </CardTitle>
                  <CardDescription>แผนกทั้งหมดที่คุณจัดการ</CardDescription>
                </div>
                <ButtonGroup direction="horizontal" spacing="normal">
                  <Button onClick={() => setShowAddDepartmentModal(true)} size="sm">
                    <PlusIcon className="w-4 h-4 mr-2" />
                    เพิ่มแผนก
                  </Button>
                </ButtonGroup>
              </div>
            </CardHeader>
            <CardContent>
              {loading ? (
                <div className="text-center py-8">
                  <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
                  <p className="text-gray-500 mt-2">กำลังโหลดข้อมูล...</p>
                </div>
              ) : departments.length > 0 ? (
                <div className="space-y-3">
                  {departments.map((department) => (
                    <div
                      key={department.id}
                      className={`p-4 border rounded-lg cursor-pointer transition-colors ${
                        selectedDepartment === department.id
                          ? 'border-blue-500 bg-blue-50'
                          : 'border-gray-200 hover:border-gray-300'
                      }`}
                      onClick={() => {
                        console.log('Department clicked:', department.id)
                        setSelectedDepartment(department.id)
                      }}
                    >
                      <div className="flex items-center justify-between">
                        <div className="flex-1">
                          <h3 className="font-medium text-gray-900">{department.name}</h3>
                          <p className="text-sm text-gray-500">
                            {department.total_employees} พนักงาน
                          </p>
                          <p className="text-xs text-gray-400">
                            พยาบาล: {department.nurse_count} | ผู้ช่วย: {department.assistant_count}
                          </p>
                        </div>
                        <div className="flex items-center space-x-2">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={(e) => {
                              e.stopPropagation()
                              openEditModal(department)
                            }}
                          >
                            <PencilIcon className="w-4 h-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={(e) => {
                              e.stopPropagation()
                              handleDeleteDepartment(department.id, department.name)
                            }}
                            className="text-red-600 hover:text-red-700"
                          >
                            <TrashIcon className="w-4 h-4" />
                          </Button>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="text-center py-12">
                  <div className="mx-auto w-24 h-24 bg-gray-100 rounded-full flex items-center justify-center mb-4">
                    <BuildingOfficeIcon className="w-12 h-12 text-gray-400" />
                  </div>
                  <h3 className="text-lg font-medium text-gray-900 mb-2">ยังไม่มีแผนกในระบบ</h3>
                  <p className="text-gray-500 mb-4">เริ่มต้นด้วยการสร้างแผนกแรกของคุณ</p>
                  <Button onClick={() => setShowAddDepartmentModal(true)}>
                    <PlusIcon className="w-4 h-4 mr-2" />
                    สร้างแผนกแรก
                  </Button>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Employees List */}
          <Card className="border-0 shadow-md">
            <CardHeader>
              <div className="flex items-center justify-between">
                <div className="flex items-center">
                  <UserGroupIcon className="w-5 h-5 mr-2" />
                  <div>
                    <CardTitle>
                      พนักงานในแผนก
                      {selectedDepartmentData && `: ${selectedDepartmentData.name}`}
                    </CardTitle>
                    <CardDescription>
                      {selectedDepartment ? 'จัดการพนักงานในแผนกที่เลือก' : 'กรุณาเลือกแผนกก่อน'}
                    </CardDescription>
                  </div>
                </div>
                {selectedDepartment && (
                  <Button
                    size="sm"
                    onClick={() => setShowAddEmployeeModal(true)}
                  >
                    <PlusIcon className="w-4 h-4 mr-1" />
                    เพิ่มพนักงาน
                  </Button>
                )}
              </div>
            </CardHeader>
            <CardContent>
              {selectedDepartment ? (
                <div className="space-y-3">
                  {(() => {
                    console.log('Rendering with selectedDepartment:', selectedDepartment, 'departmentEmployees:', departmentEmployees)
                    return null
                  })()}
                  {departmentEmployees.length > 0 ? (
                    departmentEmployees.map((employee) => (
                      <div
                        key={employee.id}
                        className="p-3 border border-gray-200 rounded-lg"
                      >
                        <div className="flex items-center justify-between">
                          <div>
                            <h4 className="font-medium text-gray-900">
                              {employee.firstName} {employee.lastName}
                            </h4>
                            <p className="text-sm text-gray-500">
                              {employee.position === 'nurse' ? 'พยาบาล' : 'ผู้ช่วยพยาบาล'}
                            </p>
                          </div>
                          <div className="flex items-center space-x-2">
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => handleEditEmployee(employee)}
                            >
                              <PencilIcon className="w-4 h-4" />
                            </Button>
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => 
                                handleDeleteEmployee(employee.id, `${employee.firstName} ${employee.lastName}`)
                              }
                              className="text-red-600 hover:text-red-700"
                            >
                              <TrashIcon className="w-4 h-4" />
                            </Button>
                          </div>
                        </div>
                      </div>
                    ))
                  ) : (
                    <div className="text-center py-8 text-gray-500">
                      ยังไม่มีพนักงานในแผนกนี้
                    </div>
                  )}
                </div>
              ) : (
                <div className="text-center py-8 text-gray-500">
                  กรุณาเลือกแผนกเพื่อดูรายการพนักงาน
                </div>
              )}
            </CardContent>
          </Card>
        </div>

        {/* Add Department Modal */}
        {showAddDepartmentModal && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg p-6 w-full max-w-md">
              <h3 className="text-lg font-medium mb-4">เพิ่มแผนกใหม่</h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    ชื่อแผนก *
                  </label>
                  <Input
                    value={departmentForm.name}
                    onChange={(e) => setDepartmentForm({ ...departmentForm, name: e.target.value })}
                    placeholder="กรอกชื่อแผนก"
                    className="text-gray-900 bg-white border-gray-300 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    คำอธิบาย
                  </label>
                  <textarea
                    value={departmentForm.description}
                    onChange={(e) => setDepartmentForm({ ...departmentForm, description: e.target.value })}
                    placeholder="กรอกคำอธิบายแผนก"
                    className="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-gray-900 bg-white"
                    rows={3}
                  />
                </div>

                {/* Hidden fields for future use */}
                <div className="hidden">
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        จำนวนพยาบาลสูงสุด *
                      </label>
                      <Input
                        type="number"
                        min="1"
                        value={departmentForm.max_nurses || 1000}
                        onChange={(e) => setDepartmentForm({ ...departmentForm, max_nurses: parseInt(e.target.value) || 1000 })}
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        จำนวนผู้ช่วยสูงสุด *
                      </label>
                      <Input
                        type="number"
                        min="1"
                        value={departmentForm.max_assistants || 500}
                        onChange={(e) => setDepartmentForm({ ...departmentForm, max_assistants: parseInt(e.target.value) || 500 })}
                      />
                    </div>
                  </div>
                </div>

                <div className="flex justify-end space-x-3">
                  <Button
                    variant="outline"
                    onClick={() => {
                      setShowAddDepartmentModal(false)
                      setDepartmentForm({ name: '', description: '', max_nurses: 1000, max_assistants: 500 })
                    }}
                  >
                    ยกเลิก
                  </Button>
                  <Button onClick={handleAddDepartment}>
                    เพิ่มแผนก
                  </Button>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Edit Department Modal */}
        {showEditDepartmentModal && editingDepartment && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg p-6 w-full max-w-md">
              <h3 className="text-lg font-medium mb-4">แก้ไขแผนก: {editingDepartment.name}</h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    ชื่อแผนก *
                  </label>
                  <Input
                    value={departmentForm.name}
                    onChange={(e) => setDepartmentForm({ ...departmentForm, name: e.target.value })}
                    placeholder="กรอกชื่อแผนก"
                    className="text-gray-900 bg-white border-gray-300 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    คำอธิบาย
                  </label>
                  <textarea
                    value={departmentForm.description}
                    onChange={(e) => setDepartmentForm({ ...departmentForm, description: e.target.value })}
                    placeholder="กรอกคำอธิบายแผนก"
                    className="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-gray-900 bg-white"
                    rows={3}
                  />
                </div>

                {/* Hidden fields for future use */}
                <div className="hidden">
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        จำนวนพยาบาลสูงสุด *
                      </label>
                      <Input
                        type="number"
                        min="1"
                        value={departmentForm.max_nurses || 1000}
                        onChange={(e) => setDepartmentForm({ ...departmentForm, max_nurses: parseInt(e.target.value) || 1000 })}
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        จำนวนผู้ช่วยสูงสุด *
                      </label>
                      <Input
                        type="number"
                        min="1"
                        value={departmentForm.max_assistants || 500}
                        onChange={(e) => setDepartmentForm({ ...departmentForm, max_assistants: parseInt(e.target.value) || 500 })}
                      />
                    </div>
                  </div>
                </div>

                <div className="flex justify-end space-x-3">
                  <Button
                    variant="outline"
                    onClick={() => {
                      setShowEditDepartmentModal(false)
                      setEditingDepartment(null)
                      setDepartmentForm({ name: '', description: '', max_nurses: 1000, max_assistants: 500 })
                    }}
                  >
                    ยกเลิก
                  </Button>
                  <Button onClick={handleEditDepartment}>
                    บันทึกการแก้ไข
                  </Button>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Add Employee Modal */}
        {showAddEmployeeModal && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg p-6 w-full max-w-md">
              <h3 className="text-lg font-medium mb-4">เพิ่มพนักงานใหม่</h3>
              <div className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      ชื่อ
                    </label>
                    <Input
                      value={employeeForm.firstName}
                      onChange={(e) => setEmployeeForm({ ...employeeForm, firstName: e.target.value })}
                      placeholder="ชื่อ"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      นามสกุล
                    </label>
                    <Input
                      value={employeeForm.lastName}
                      onChange={(e) => setEmployeeForm({ ...employeeForm, lastName: e.target.value })}
                      placeholder="นามสกุล"
                    />
                  </div>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    ตำแหน่ง
                  </label>
                  <select
                    value={employeeForm.position}
                    onChange={(e) => setEmployeeForm({ ...employeeForm, position: e.target.value as 'nurse' | 'assistant' })}
                    className="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="nurse">พยาบาล</option>
                    <option value="assistant">ผู้ช่วยพยาบาล</option>
                  </select>
                </div>
                <div className="flex justify-end space-x-3">
                  <Button
                    variant="outline"
                    onClick={() => {
                      setShowAddEmployeeModal(false)
                      setEmployeeForm({ firstName: '', lastName: '', position: 'nurse' })
                    }}
                  >
                    ยกเลิก
                  </Button>
                  <Button onClick={handleAddEmployee}>
                    เพิ่มพนักงาน
                  </Button>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Edit Employee Modal */}
        {showEditEmployeeModal && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg p-6 w-full max-w-md">
              <h3 className="text-lg font-medium mb-4">แก้ไขพนักงาน</h3>
              <div className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      ชื่อ
                    </label>
                    <Input
                      value={employeeForm.firstName}
                      onChange={(e) => setEmployeeForm({ ...employeeForm, firstName: e.target.value })}
                      placeholder="ชื่อ"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      นามสกุล
                    </label>
                    <Input
                      value={employeeForm.lastName}
                      onChange={(e) => setEmployeeForm({ ...employeeForm, lastName: e.target.value })}
                      placeholder="นามสกุล"
                    />
                  </div>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    ตำแหน่ง
                  </label>
                  <select
                    value={employeeForm.position}
                    onChange={(e) => setEmployeeForm({ ...employeeForm, position: e.target.value as 'nurse' | 'assistant' })}
                    className="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="nurse">พยาบาล</option>
                    <option value="assistant">ผู้ช่วยพยาบาล</option>
                  </select>
                </div>
                <div className="flex justify-end space-x-3">
                  <Button
                    variant="outline"
                    onClick={() => {
                      setShowEditEmployeeModal(false)
                      setEditingEmployee(null)
                      setEmployeeForm({ firstName: '', lastName: '', position: 'nurse' })
                    }}
                  >
                    ยกเลิก
                  </Button>
                  <Button onClick={handleUpdateEmployee}>
                    บันทึกการแก้ไข
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
