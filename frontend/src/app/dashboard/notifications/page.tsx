'use client'

import { useState } from 'react'
import { 
  BellIcon,
  CheckCircleIcon,
  ExclamationTriangleIcon,
  InformationCircleIcon,
  XMarkIcon,
  EyeIcon,
  TrashIcon,
  FunnelIcon
} from '@heroicons/react/24/outline'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card'
import { Button } from '@/components/ui/Button'
import { ButtonGroup } from '@/components/ui/ButtonGroup'
import DashboardLayout from '@/components/layout/DashboardLayout'
import Swal from 'sweetalert2'

// Mock data for notifications
const mockNotifications = [
  {
    id: 1,
    type: 'schedule' as const,
    title: 'ตารางเวรใหม่ประจำเดือนมีนาคม',
    message: 'ตารางเวรประจำเดือนมีนาคม 2024 ได้รับการอนุมัติและเผยแพร่แล้ว กรุณาตรวจสอบตารางเวรของคุณ',
    timestamp: '2024-03-01T09:00:00Z',
    isRead: false,
    priority: 'high' as const,
    actionUrl: '/dashboard/schedule'
  },
  {
    id: 2,
    type: 'leave' as const,
    title: 'คำขอลาป่วยได้รับการอนุมัติ',
    message: 'คำขอลาป่วยวันที่ 15 มีนาคม 2024 ได้รับการอนุมัติจากหัวหน้าแผนกแล้ว',
    timestamp: '2024-03-01T14:30:00Z',
    isRead: true,
    priority: 'medium' as const,
    actionUrl: '/dashboard/employee-leaves'
  },
  {
    id: 3,
    type: 'system' as const,
    title: 'การอัปเดตระบบ',
    message: 'ระบบจะมีการปรับปรุงในวันที่ 20 มีนาคม 2024 เวลา 02:00-04:00 น. อาจมีการหยุดให้บริการชั่วคราว',
    timestamp: '2024-02-28T16:00:00Z',
    isRead: false,
    priority: 'low' as const,
    actionUrl: null
  },
  {
    id: 4,
    type: 'payment' as const,
    title: 'การชำระเงินแพ็คเกจ',
    message: 'การชำระเงินแพ็คเกจมาตรฐานได้รับการอนุมัติแล้ว บัญชีของคุณได้รับการต่ออายุเป็น 30 วัน',
    timestamp: '2024-02-25T11:15:00Z',
    isRead: true,
    priority: 'high' as const,
    actionUrl: '/dashboard/packages'
  },
  {
    id: 5,
    type: 'reminder' as const,
    title: 'แจ้งเตือนเปลี่ยนเวร',
    message: 'คุณมีเวรเช้าในวันพรุ่งนี้ (16 มีนาคม 2024) เวลา 07:00-15:00 น. ที่แผนกผู้ป่วยใน',
    timestamp: '2024-03-15T20:00:00Z',
    isRead: false,
    priority: 'medium' as const,
    actionUrl: '/dashboard/schedule'
  },
  {
    id: 6,
    type: 'holiday' as const,
    title: 'วันหยุดประจำปี',
    message: 'วันสงกรานต์ (13-15 เมษายน 2024) ได้รับการอนุมัติเป็นวันหยุดประจำปีแล้ว',
    timestamp: '2024-02-20T10:00:00Z',
    isRead: true,
    priority: 'low' as const,
    actionUrl: '/dashboard/department-settings'
  }
]

export default function NotificationsPage() {
  const [notifications, setNotifications] = useState(mockNotifications)
  const [filter, setFilter] = useState<'all' | 'unread' | 'read'>('all')
  const [typeFilter, setTypeFilter] = useState<'all' | 'schedule' | 'leave' | 'system' | 'payment' | 'reminder' | 'holiday'>('all')

  const handleMarkAsRead = (id: number) => {
    setNotifications(notifications.map(notif => 
      notif.id === id ? { ...notif, isRead: true } : notif
    ))
  }

  const handleMarkAllAsRead = async () => {
    const result = await Swal.fire({
      title: 'ทำเครื่องหมายอ่านแล้วทั้งหมด?',
      text: 'การแจ้งเตือนทั้งหมดจะถูกทำเครื่องหมายว่าอ่านแล้ว',
      icon: 'question',
      showCancelButton: true,
      confirmButtonColor: '#2563eb',
      cancelButtonColor: '#6b7280',
      confirmButtonText: 'ดำเนินการ',
      cancelButtonText: 'ยกเลิก'
    })

    if (result.isConfirmed) {
      setNotifications(notifications.map(notif => ({ ...notif, isRead: true })))
      await Swal.fire({
        icon: 'success',
        title: 'ทำเครื่องหมายสำเร็จ!',
        text: 'การแจ้งเตือนทั้งหมดถูกทำเครื่องหมายว่าอ่านแล้ว',
        confirmButtonColor: '#2563eb'
      })
    }
  }

  const handleDeleteNotification = async (id: number) => {
    const result = await Swal.fire({
      title: 'ลบการแจ้งเตือน?',
      text: 'การแจ้งเตือนนี้จะถูกลบอย่างถาวร',
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#dc2626',
      cancelButtonColor: '#6b7280',
      confirmButtonText: 'ลบ',
      cancelButtonText: 'ยกเลิก'
    })

    if (result.isConfirmed) {
      setNotifications(notifications.filter(notif => notif.id !== id))
      await Swal.fire({
        icon: 'success',
        title: 'ลบสำเร็จ!',
        confirmButtonColor: '#2563eb'
      })
    }
  }

  const handleDeleteAllRead = async () => {
    const readNotifications = notifications.filter(notif => notif.isRead)
    if (readNotifications.length === 0) {
      await Swal.fire({
        icon: 'info',
        title: 'ไม่พบการแจ้งเตือนที่อ่านแล้ว',
        confirmButtonColor: '#2563eb'
      })
      return
    }

    const result = await Swal.fire({
      title: 'ลบการแจ้งเตือนที่อ่านแล้วทั้งหมด?',
      text: `จะลบการแจ้งเตือนที่อ่านแล้ว ${readNotifications.length} รายการ`,
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#dc2626',
      cancelButtonColor: '#6b7280',
      confirmButtonText: 'ลบทั้งหมด',
      cancelButtonText: 'ยกเลิก'
    })

    if (result.isConfirmed) {
      setNotifications(notifications.filter(notif => !notif.isRead))
      await Swal.fire({
        icon: 'success',
        title: 'ลบสำเร็จ!',
        text: `ลบการแจ้งเตือนแล้ว ${readNotifications.length} รายการ`,
        confirmButtonColor: '#2563eb'
      })
    }
  }

  const getNotificationIcon = (type: string) => {
    switch (type) {
      case 'schedule':
        return <BellIcon className="w-5 h-5 text-blue-500" />
      case 'leave':
        return <ExclamationTriangleIcon className="w-5 h-5 text-yellow-500" />
      case 'system':
        return <InformationCircleIcon className="w-5 h-5 text-gray-500" />
      case 'payment':
        return <CheckCircleIcon className="w-5 h-5 text-green-500" />
      case 'reminder':
        return <BellIcon className="w-5 h-5 text-orange-500" />
      case 'holiday':
        return <InformationCircleIcon className="w-5 h-5 text-purple-500" />
      default:
        return <BellIcon className="w-5 h-5 text-gray-500" />
    }
  }

  const getPriorityBadge = (priority: string) => {
    switch (priority) {
      case 'high':
        return 'bg-red-100 text-red-800 border-red-200'
      case 'medium':
        return 'bg-yellow-100 text-yellow-800 border-yellow-200'
      case 'low':
        return 'bg-gray-100 text-gray-800 border-gray-200'
      default:
        return 'bg-gray-100 text-gray-800 border-gray-200'
    }
  }

  const getPriorityText = (priority: string) => {
    switch (priority) {
      case 'high':
        return 'สูง'
      case 'medium':
        return 'กลาง'
      case 'low':
        return 'ต่ำ'
      default:
        return 'ไม่ระบุ'
    }
  }

  const getTypeText = (type: string) => {
    switch (type) {
      case 'schedule':
        return 'ตารางเวร'
      case 'leave':
        return 'การลา'
      case 'system':
        return 'ระบบ'
      case 'payment':
        return 'การชำระเงิน'
      case 'reminder':
        return 'แจ้งเตือน'
      case 'holiday':
        return 'วันหยุด'
      default:
        return 'ทั่วไป'
    }
  }

  const filteredNotifications = notifications.filter(notif => {
    const matchesReadFilter = filter === 'all' || 
      (filter === 'read' && notif.isRead) || 
      (filter === 'unread' && !notif.isRead)
    
    const matchesTypeFilter = typeFilter === 'all' || notif.type === typeFilter

    return matchesReadFilter && matchesTypeFilter
  })

  const unreadCount = notifications.filter(notif => !notif.isRead).length

  return (
    <DashboardLayout>
      <div className="space-y-8">
        {/* Header */}
        <div className="space-y-4">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">การแจ้งเตือนทั้งหมด</h1>
            <p className="text-gray-600 mt-2">
              จัดการและติดตามการแจ้งเตือนต่างๆ ในระบบ
              {unreadCount > 0 && (
                <span className="ml-2 inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
                  ยังไม่อ่าน {unreadCount} รายการ
                </span>
              )}
            </p>
          </div>
          <div className="flex flex-col sm:flex-row space-y-2 sm:space-y-0 sm:space-x-3">
            <ButtonGroup direction="horizontal" spacing="normal" className="w-full sm:w-auto">
              <Button variant="outline" onClick={handleMarkAllAsRead} className="w-full sm:w-auto">
                <CheckCircleIcon className="w-4 h-4 mr-2" />
                <span className="hidden sm:inline">ทำเครื่องหมายอ่านแล้วทั้งหมด</span>
                <span className="sm:hidden">อ่านแล้วทั้งหมด</span>
              </Button>
              <Button variant="outline" onClick={handleDeleteAllRead} className="w-full sm:w-auto">
                <TrashIcon className="w-4 h-4 mr-2" />
                <span className="hidden sm:inline">ลบที่อ่านแล้ว</span>
                <span className="sm:hidden">ลบที่อ่านแล้ว</span>
              </Button>
            </ButtonGroup>
          </div>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
          <Card className="border-0 shadow-md">
            <CardContent className="pt-6">
              <div className="flex items-center">
                <div className="p-2 bg-blue-100 rounded-lg">
                  <BellIcon className="w-6 h-6 text-blue-600" />
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">ทั้งหมด</p>
                  <p className="text-2xl font-bold text-gray-900">{notifications.length}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card className="border-0 shadow-md">
            <CardContent className="pt-6">
              <div className="flex items-center">
                <div className="p-2 bg-red-100 rounded-lg">
                  <ExclamationTriangleIcon className="w-6 h-6 text-red-600" />
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">ยังไม่อ่าน</p>
                  <p className="text-2xl font-bold text-gray-900">{unreadCount}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card className="border-0 shadow-md">
            <CardContent className="pt-6">
              <div className="flex items-center">
                <div className="p-2 bg-green-100 rounded-lg">
                  <CheckCircleIcon className="w-6 h-6 text-green-600" />
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">อ่านแล้ว</p>
                  <p className="text-2xl font-bold text-gray-900">{notifications.filter(n => n.isRead).length}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card className="border-0 shadow-md">
            <CardContent className="pt-6">
              <div className="flex items-center">
                <div className="p-2 bg-yellow-100 rounded-lg">
                  <ExclamationTriangleIcon className="w-6 h-6 text-yellow-600" />
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">ความสำคัญสูง</p>
                  <p className="text-2xl font-bold text-gray-900">{notifications.filter(n => n.priority === 'high').length}</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Filters */}
        <Card className="border-0 shadow-md">
          <CardHeader>
            <CardTitle className="flex items-center">
              <FunnelIcon className="w-5 h-5 mr-2" />
              ตัวกรอง
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">สถานะการอ่าน</label>
                <select
                  value={filter}
                  onChange={(e) => setFilter(e.target.value as 'all' | 'unread' | 'read')}
                  className="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                >
                  <option value="all">ทั้งหมด</option>
                  <option value="unread">ยังไม่อ่าน</option>
                  <option value="read">อ่านแล้ว</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">ประเภทการแจ้งเตือน</label>
                <select
                  value={typeFilter}
                  onChange={(e) => setTypeFilter(e.target.value as any)}
                  className="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                >
                  <option value="all">ทุกประเภท</option>
                  <option value="schedule">ตารางเวร</option>
                  <option value="leave">การลา</option>
                  <option value="payment">การชำระเงิน</option>
                  <option value="reminder">แจ้งเตือน</option>
                  <option value="holiday">วันหยุด</option>
                  <option value="system">ระบบ</option>
                </select>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Notifications List */}
        <Card className="border-0 shadow-md">
          <CardHeader>
            <CardTitle>รายการการแจ้งเตือน</CardTitle>
            <CardDescription>
              แสดงการแจ้งเตือน {filteredNotifications.length} รายการ
            </CardDescription>
          </CardHeader>
          <CardContent>
            {filteredNotifications.length > 0 ? (
              <div className="space-y-4">
                {filteredNotifications.map((notification) => (
                  <div 
                    key={notification.id} 
                    className={`p-4 border rounded-lg transition-all ${
                      notification.isRead 
                        ? 'border-gray-200 bg-white' 
                        : 'border-blue-200 bg-blue-50'
                    }`}
                  >
                    <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between space-y-3 sm:space-y-0">
                      <div className="flex-1">
                        <div className="flex items-start space-x-3 mb-2">
                          {getNotificationIcon(notification.type)}
                          <div className="flex-1 min-w-0">
                            <h3 className={`font-medium ${notification.isRead ? 'text-gray-900' : 'text-blue-900'}`}>
                              {notification.title}
                            </h3>
                            <div className="flex flex-wrap gap-2 mt-1">
                              <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${getPriorityBadge(notification.priority)}`}>
                                {getPriorityText(notification.priority)}
                              </span>
                              <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                                {getTypeText(notification.type)}
                              </span>
                              {!notification.isRead && (
                                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
                                  ใหม่
                                </span>
                              )}
                            </div>
                          </div>
                        </div>
                        <p className={`text-sm mb-2 ${notification.isRead ? 'text-gray-600' : 'text-blue-800'}`}>
                          {notification.message}
                        </p>
                        <p className="text-xs text-gray-500">
                          {new Date(notification.timestamp).toLocaleDateString('th-TH', {
                            year: 'numeric',
                            month: 'long',
                            day: 'numeric',
                            hour: '2-digit',
                            minute: '2-digit'
                          })}
                        </p>
                      </div>
                      <div className="flex flex-col sm:flex-row sm:items-center space-y-2 sm:space-y-0 sm:ml-4">
                        <ButtonGroup direction="horizontal" spacing="tight" className="w-full sm:w-auto">
                          {!notification.isRead && (
                            <Button 
                              variant="outline" 
                              size="sm"
                              onClick={() => handleMarkAsRead(notification.id)}
                              className="w-full sm:w-auto"
                            >
                              <EyeIcon className="w-4 h-4 mr-1" />
                              <span className="hidden sm:inline">ทำเครื่องหมายอ่านแล้ว</span>
                              <span className="sm:hidden">อ่านแล้ว</span>
                            </Button>
                          )}
                          {notification.actionUrl && (
                            <Button 
                              variant="outline" 
                              size="sm"
                              onClick={() => window.location.href = notification.actionUrl!}
                              className="w-full sm:w-auto"
                            >
                              ดูรายละเอียด
                            </Button>
                          )}
                          <Button 
                            variant="outline" 
                            size="sm"
                            className="text-red-600 border-red-300 hover:bg-red-50 w-full sm:w-auto"
                            onClick={() => handleDeleteNotification(notification.id)}
                          >
                            <TrashIcon className="w-4 h-4" />
                            <span className="sm:hidden ml-2">ลบ</span>
                          </Button>
                        </ButtonGroup>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-8 text-gray-500">
                <BellIcon className="w-12 h-12 mx-auto text-gray-300 mb-4" />
                <p>ไม่พบการแจ้งเตือนตามเงื่อนไขที่เลือก</p>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </DashboardLayout>
  )
}
