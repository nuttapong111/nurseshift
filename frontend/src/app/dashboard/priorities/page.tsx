'use client'

import { useEffect, useState } from 'react'
import { 
  PencilIcon,
  StarIcon,
  ArrowUpIcon,
  ArrowDownIcon
} from '@heroicons/react/24/outline'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card'
import { Button } from '@/components/ui/Button'
import { ButtonGroup } from '@/components/ui/ButtonGroup'
import { Input } from '@/components/ui/Input'
import { Switch } from '@/components/ui/Switch'
import DashboardLayout from '@/components/layout/DashboardLayout'
import Swal from 'sweetalert2'
import { priorityService } from '@/services/priorityService'
import { departmentService } from '@/services/departmentService'

// NOTE: คงโครง UI เดิม แต่ดึงข้อมูลจริงจาก priority-service

// Priority settings configuration
const prioritySettings = {
  maxShiftTypeDifference: { min: 0, max: 5, step: 1 },
  maxConsecutiveNightShifts: { min: 1, max: 5, step: 1 },
  maxConsecutiveShifts: { min: 1, max: 10, step: 1 },
  maxConsecutiveWorkHours: { min: 12, max: 72, step: 6 },
  maxTotalWorkHoursDifference: { min: 8, max: 80, step: 2 },
}

export default function PrioritiesPage() {
  const [departments, setDepartments] = useState<Array<{ id: string; name: string }>>([])
  const [selectedDepartment, setSelectedDepartment] = useState<string>('')
  const [priorities, setPriorities] = useState<any[]>([])

  useEffect(() => {
    const loadDeps = async () => {
      try {
        const list = await departmentService.getDepartments()
        const mapped = (list || []).map((d: any) => ({ id: d.id, name: d.name }))
        setDepartments(mapped)
        if (mapped[0]) setSelectedDepartment(mapped[0].id)
      } catch (e) {
        setDepartments([])
      }
    }
    loadDeps()
  }, [])

  useEffect(() => {
    const load = async () => {
      if (!selectedDepartment) return
      try {
        const data = await priorityService.list(selectedDepartment)
        setPriorities(data.priorities || [])
      } catch (e) {
        setPriorities([])
      }
    }
    load()
  }, [selectedDepartment])

  const handleTogglePriority = async (priorityId: string, currentStatus: boolean) => {
    try {
      await priorityService.update(priorityId, { isActive: !currentStatus })
      setPriorities(prev => prev.map(p => p.id === priorityId ? { ...p, isActive: !currentStatus } : p))
      await Swal.fire({ icon: 'success', title: currentStatus ? 'ปิดการใช้งานความสำคัญแล้ว' : 'เปิดการใช้งานความสำคัญแล้ว', timer: 1200, showConfirmButton: false })
    } catch (e: any) {
      await Swal.fire({ icon: 'error', title: 'อัปเดตไม่สำเร็จ', text: e.message || 'เกิดข้อผิดพลาด' })
    }
  }

  const handleMoveUp = async (priorityId: string) => {
    const currentIndex = priorities.findIndex((p) => p.id === priorityId)
    if (currentIndex <= 0) return
    try {
      await priorityService.swap(priorityId, priorities[currentIndex - 1].id)
      const newPriorities = [...priorities]
      const temp = newPriorities[currentIndex]
      newPriorities[currentIndex] = newPriorities[currentIndex - 1]
      newPriorities[currentIndex - 1] = temp
      // Reassign orders visually
      newPriorities.forEach((p, idx) => (p.order = idx + 1))
      setPriorities(newPriorities)
    } catch (e: any) {
      await Swal.fire({ icon: 'error', title: 'สลับลำดับไม่สำเร็จ', text: e.message || 'เกิดข้อผิดพลาด' })
    }
  }

  const handleMoveDown = async (priorityId: string) => {
    const currentIndex = priorities.findIndex((p) => p.id === priorityId)
    if (currentIndex >= priorities.length - 1) return
    try {
      await priorityService.swap(priorityId, priorities[currentIndex + 1].id)
      const newPriorities = [...priorities]
      const temp = newPriorities[currentIndex]
      newPriorities[currentIndex] = newPriorities[currentIndex + 1]
      newPriorities[currentIndex + 1] = temp
      newPriorities.forEach((p, idx) => (p.order = idx + 1))
      setPriorities(newPriorities)
    } catch (e: any) {
      await Swal.fire({ icon: 'error', title: 'สลับลำดับไม่สำเร็จ', text: e.message || 'เกิดข้อผิดพลาด' })
    }
  }

  const handleSettingChange = async (_priorityId: string, _newValue: number) => {
    try {
      // optimistic update
      setPriorities(priorities.map(priority => 
        priority.id === _priorityId 
          ? { ...priority, settingValue: _newValue }
          : priority
      ))
      // persist to API
      await fetch(`${process.env.NEXT_PUBLIC_PRIORITY_API_URL || 'http://localhost:8086'}/api/v1/priorities/${_priorityId}/setting`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${localStorage.getItem('token')}` },
        body: JSON.stringify({ settingValue: Number(_newValue) })
      })
    } catch (e) {
      // ignore here; toast handled in pencil flow
    }
  }



  // Sort priorities by order
  const sortedPriorities = [...priorities].sort((a, b) => a.order - b.order)

  return (
    <DashboardLayout>
      <div className="space-y-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">จัดการความสำคัญ</h1>
          <p className="text-gray-600 mt-2">ตั้งค่าลำดับความสำคัญในการจัดตารางเวรอัตโนมัติ ระบบมีความสำคัญที่กำหนดไว้แล้วให้ครบถ้วน</p>
        </div>

        {/* Priority Settings Info */}
        <Card className="border-0 shadow-md bg-blue-50 border-blue-200">
          <CardContent className="pt-6">
            <div className="flex items-start space-x-3">
              <StarIcon className="w-6 h-6 text-blue-600 mt-0.5" />
              <div>
                <h3 className="font-medium text-blue-900 mb-2">เกี่ยวกับการตั้งค่าความสำคัญ</h3>
                <p className="text-blue-800 text-sm leading-relaxed">
                  ระบบมีความสำคัญที่กำหนดไว้แล้วสำหรับการจัดตารางเวรอัตโนมัติ โดยจะพิจารณาตามลำดับจากบนลงล่าง 
                  คุณสามารถเปลี่ยนลำดับได้โดยใช้ปุ่มลูกศร เปิด/ปิดการใช้งานแต่ละความสำคัญ และปรับค่าพารามิเตอร์ได้
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Priorities List */}
        <Card className="border-0 shadow-md">
            <CardHeader>
            <CardTitle className="flex items-center">
              <StarIcon className="w-5 h-5 mr-2" />
              รายการความสำคัญ ({sortedPriorities.filter(p => p.isActive).length}/{sortedPriorities.length} กำลังใช้งาน)
            </CardTitle>
            <CardDescription>
              ความสำคัญที่กำหนดไว้ในระบบ - ปรับลำดับด้วยปุ่มลูกศร ลำดับที่ 1 จะมีความสำคัญสูงสุด
            </CardDescription>
          </CardHeader>
          <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">เลือกแผนก</label>
                  <select
                    value={selectedDepartment}
                    onChange={(e) => setSelectedDepartment(e.target.value)}
                    className="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 h-10"
                  >
                    {departments.map((d) => (
                      <option key={d.id} value={d.id}>{d.name}</option>
                    ))}
                  </select>
                </div>
              </div>
            <div className="space-y-4">
              {sortedPriorities.map((priority, index) => (
                <div 
                  key={priority.id} 
                  className={`p-4 border rounded-lg transition-all ${
                    priority.isActive 
                      ? 'border-blue-200 bg-blue-50' 
                      : 'border-gray-200 bg-gray-50'
                  }`}
                >
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <div className="flex items-center space-x-3 mb-2">
                        <div className={`w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold ${
                          priority.isActive 
                            ? 'bg-blue-600 text-white' 
                            : 'bg-gray-400 text-white'
                        }`}>
                          {priority.order}
                        </div>
                        <h3 className={`font-medium ${
                          priority.isActive ? 'text-blue-900' : 'text-gray-700'
                        }`}>
                          {priority.name}
                        </h3>
                        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                          priority.isActive 
                            ? 'bg-green-100 text-green-800' 
                            : 'bg-red-100 text-red-800'
                        }`}>
                          {priority.isActive ? 'ใช้งาน' : 'ไม่ใช้งาน'}
                        </span>
                      </div>
                      <p className={`text-sm ml-11 ${
                        priority.isActive ? 'text-blue-700' : 'text-gray-600'
                      }`}>
                        {priority.description}
                      </p>
                    </div>
                    
                    <div className="flex items-center space-x-4 ml-4">
                      {/* Move buttons */}
                      <ButtonGroup direction="vertical" spacing="tight">
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleMoveUp(priority.id)}
                          disabled={index === 0}
                          className="h-6 w-6 p-0"
                        >
                          <ArrowUpIcon className="w-3 h-3" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleMoveDown(priority.id)}
                          disabled={index === sortedPriorities.length - 1}
                          className="h-6 w-6 p-0"
                        >
                          <ArrowDownIcon className="w-3 h-3" />
                        </Button>
                      </ButtonGroup>
                      
                      {/* Toggle switch */}
                      <div className="flex items-center space-x-2">
                        <span className="text-sm text-gray-600">
                          {priority.isActive ? 'เปิด' : 'ปิด'}
                        </span>
                        <Switch
                          checked={priority.isActive}
                          onChange={() => handleTogglePriority(priority.id, priority.isActive)}
                        />
                      </div>
                      
                       {/* Action buttons */}
                       {priority.hasSettings && (
                         <Button 
                           variant="ghost" 
                           size="sm"
                           onClick={async () => {
                             const { value: num } = await Swal.fire({
                               title: 'ตั้งค่าพารามิเตอร์',
                               input: 'number',
                               inputLabel: 'ค่าใหม่',
                               inputValue: priority.settingValue ?? 0,
                               showCancelButton: true
                             })
                             if (num !== undefined) {
                               try {
                                 await fetch(`${process.env.NEXT_PUBLIC_PRIORITY_API_URL || 'http://localhost:8086'}/api/v1/priorities/${priority.id}/setting`, {
                                   method: 'PUT',
                                   headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${localStorage.getItem('token')}` },
                                   body: JSON.stringify({ settingValue: Number(num) })
                                 })
                                 setPriorities(prev => prev.map(p => p.id === priority.id ? { ...p, settingValue: Number(num) } : p))
                                 await Swal.fire({ icon: 'success', title: 'บันทึกค่าแล้ว', timer: 1000, showConfirmButton: false })
                               } catch (e: any) {
                                 await Swal.fire({ icon: 'error', title: 'บันทึกไม่สำเร็จ', text: e.message || 'เกิดข้อผิดพลาด' })
                               }
                             }
                           }}
                         >
                           <PencilIcon className="w-4 h-4" />
                         </Button>
                       )}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Priority Settings Card */}
        <Card className="border-0 shadow-md">
          <CardHeader>
            <CardTitle className="flex items-center">
              <StarIcon className="w-5 h-5 mr-2" />
              ตั้งค่าค่าพารามิเตอร์
            </CardTitle>
            <CardDescription>
              ปรับค่าตัวเลขสำหรับความสำคัญที่เปิดใช้งานและมีการตั้งค่าได้
            </CardDescription>
          </CardHeader>
          <CardContent>
            {sortedPriorities
              .filter(priority => priority.isActive && priority.hasSettings)
              .length > 0 ? (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {sortedPriorities
                  .filter(priority => priority.isActive && priority.hasSettings)
                  .map((priority) => {
                    const config = prioritySettings[priority.settingType as keyof typeof prioritySettings] || { min: 0, max: 100, step: 1 }
                    return (
                      <div 
                        key={priority.id} 
                        className="p-4 border border-blue-200 rounded-lg bg-blue-50"
                      >
                        <div className="flex items-center space-x-2 mb-3">
                          <div className="w-6 h-6 rounded-full bg-blue-600 text-white flex items-center justify-center text-xs font-bold">
                            {priority.order}
                          </div>
                          <h4 className="font-medium text-blue-900">{priority.name}</h4>
                        </div>
                        
                        <div className="space-y-3">
                          <div>
                            <label className="block text-sm font-medium text-blue-800 mb-2">
                              {priority.settingLabel}
                            </label>
                            <div className="flex items-center space-x-3">
                              <input
                                type="range"
                                min={config.min}
                                max={config.max}
                                step={config.step}
                                value={priority.settingValue}
                                onChange={(e) => handleSettingChange(priority.id, Number(e.target.value))}
                                className="flex-1 h-2 bg-blue-200 rounded-lg appearance-none cursor-pointer"
                                style={{
                                  background: `linear-gradient(to right, #2563eb 0%, #2563eb ${((priority.settingValue! - config.min) / (config.max - config.min)) * 100}%, #dbeafe ${((priority.settingValue! - config.min) / (config.max - config.min)) * 100}%, #dbeafe 100%)`
                                }}
                              />
                              <div className="flex items-center space-x-2">
                                <Input
                                  type="number"
                                  min={config.min}
                                  max={config.max}
                                  step={config.step}
                                  value={priority.settingValue}
                                  onChange={(e) => handleSettingChange(priority.id, Number(e.target.value))}
                                  className="w-20 text-center"
                                />
                                <span className="text-sm text-blue-700 font-medium min-w-[3rem]">
                                  {priority.settingUnit}
                                </span>
                              </div>
                            </div>
                          </div>
                          
                          <div className="text-xs text-blue-600 bg-blue-100 rounded p-2">
                            <strong>ช่วงค่าที่อนุญาต:</strong> {config.min} - {config.max} {priority.settingUnit}
                          </div>
                        </div>
                      </div>
                    )
                  })}
              </div>
            ) : (
              <div className="text-center py-8 text-gray-500">
                <div className="mb-4">
                  <StarIcon className="w-12 h-12 text-gray-300 mx-auto" />
                </div>
                <p className="text-lg font-medium">ไม่มีความสำคัญที่ต้องตั้งค่า</p>
                <p className="text-sm mt-1">
                  เปิดใช้งานความสำคัญที่มีการตั้งค่าได้เพื่อปรับค่าพารามิเตอร์
                </p>
              </div>
            )}
            
            {/* Settings Summary */}
            {sortedPriorities.filter(p => p.isActive && p.hasSettings).length > 0 && (
              <div className="mt-6 p-4 bg-green-50 border border-green-200 rounded-lg">
                <h5 className="font-medium text-green-900 mb-2">สรุปการตั้งค่าปัจจุบัน</h5>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-2 text-sm">
                  {sortedPriorities
                    .filter(p => p.isActive && p.hasSettings)
                    .map((priority) => (
                      <div key={priority.id} className="flex justify-between">
                        <span className="text-green-700">{priority.settingLabel}:</span>
                        <span className="font-medium text-green-900">
                          {priority.settingValue} {priority.settingUnit}
                        </span>
                      </div>
                    ))}
                </div>
              </div>
            )}
          </CardContent>
        </Card>


      </div>
    </DashboardLayout>
  )
}
