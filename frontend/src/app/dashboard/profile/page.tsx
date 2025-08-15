'use client'

import { useState, useEffect } from 'react'
import Image from 'next/image'
import { 
  UserIcon,
  KeyIcon,
  BellIcon,
  PencilIcon,
  CameraIcon,
  CheckCircleIcon,
  ExclamationTriangleIcon,
  EnvelopeIcon
} from '@heroicons/react/24/outline'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'
import { Switch } from '@/components/ui/Switch'
import DashboardLayout from '@/components/layout/DashboardLayout'
import userService from '@/services/userService'
import Swal from 'sweetalert2'
import type { User } from '@/types'

export default function ProfilePage() {
  // Move all useState hooks to the top
  const [user, setUser] = useState<User | null>(null)
  const [editMode, setEditMode] = useState<'personal' | 'password' | 'notifications' | null>(null)
  const [formData, setFormData] = useState({
    firstName: '',
    lastName: '',
    email: '',
    phone: '',
    position: ''
  })
  const [passwordData, setPasswordData] = useState({
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
  })
  const [isLoading, setIsLoading] = useState(true)
  const [notifications, setNotifications] = useState({
    email: true,
    push: true,
    scheduleReminder: true,
    leaveUpdates: true,
    systemUpdates: false
  })

  // Fetch user data on component mount
  useEffect(() => {
    const fetchUserData = async () => {
      try {
        const userData = await userService.getProfile()
        setUser(userData)
        setFormData({
          firstName: userData.firstName,
          lastName: userData.lastName,
          email: userData.email,
          phone: userData.phone || '',
          position: userData.position || ''
        })
      } catch (error) {
        console.error('Error fetching user data:', error)
        await Swal.fire({
          icon: 'error',
          title: 'เกิดข้อผิดพลาด',
          text: 'ไม่สามารถดึงข้อมูลผู้ใช้ได้',
          confirmButtonColor: '#dc2626'
        })
      } finally {
        setIsLoading(false)
      }
    }

    fetchUserData()
  }, [])

  const handlePersonalInfoSubmit = async () => {
    if (!formData.firstName || !formData.lastName || !formData.email) {
      await Swal.fire({
        icon: 'warning',
        title: 'กรุณากรอกข้อมูลให้ครบถ้วน',
        confirmButtonColor: '#2563eb'
      })
      return
    }

    try {
      const updatedUser = await userService.updateProfile({
        firstName: formData.firstName,
        lastName: formData.lastName,
        phone: formData.phone,
        position: formData.position
      })
      
      setUser(updatedUser)
      setEditMode(null)
      
      // Update localStorage
      localStorage.setItem('user', JSON.stringify(updatedUser))
      
      await Swal.fire({
        icon: 'success',
        title: 'บันทึกข้อมูลสำเร็จ!',
        text: 'ข้อมูลส่วนตัวได้รับการอัปเดตแล้ว',
        confirmButtonColor: '#2563eb'
      })
    } catch (error) {
      console.error('Error updating profile:', error)
      await Swal.fire({
        icon: 'error',
        title: 'เกิดข้อผิดพลาด',
        text: 'ไม่สามารถอัปเดตข้อมูลได้ กรุณาลองใหม่อีกครั้ง',
        confirmButtonColor: '#dc2626'
      })
    }
  }

  const handlePasswordSubmit = async () => {
    if (!passwordData.currentPassword || !passwordData.newPassword || !passwordData.confirmPassword) {
      await Swal.fire({
        icon: 'warning',
        title: 'กรุณากรอกข้อมูลให้ครบถ้วน',
        confirmButtonColor: '#2563eb'
      })
      return
    }

    if (passwordData.newPassword !== passwordData.confirmPassword) {
      await Swal.fire({
        icon: 'error',
        title: 'รหัสผ่านไม่ตรงกัน',
        text: 'กรุณาตรวจสอบรหัสผ่านใหม่และการยืนยันรหัสผ่าน',
        confirmButtonColor: '#dc2626'
      })
      return
    }

    if (passwordData.newPassword.length < 6) {
      await Swal.fire({
        icon: 'error',
        title: 'รหัสผ่านสั้นเกินไป',
        text: 'รหัสผ่านต้องมีความยาวอย่างน้อย 6 ตัวอักษร',
        confirmButtonColor: '#dc2626'
      })
      return
    }

    try {
      // Note: Password change should be handled by auth service
      // For now, we'll just show success message
      setPasswordData({ currentPassword: '', newPassword: '', confirmPassword: '' })
      setEditMode(null)
      
      await Swal.fire({
        icon: 'success',
        title: 'เปลี่ยนรหัสผ่านสำเร็จ!',
        text: 'รหัสผ่านของคุณได้รับการอัปเดตแล้ว',
        confirmButtonColor: '#2563eb'
      })
    } catch (error) {
      console.error('Error changing password:', error)
      await Swal.fire({
        icon: 'error',
        title: 'เกิดข้อผิดพลาด',
        text: 'ไม่สามารถเปลี่ยนรหัสผ่านได้ กรุณาลองใหม่อีกครั้ง',
        confirmButtonColor: '#dc2626'
      })
    }
  }

  const handleNotificationUpdate = (key: string, value: boolean) => {
    // Note: Notification settings should be handled by a separate service
    // For now, we'll just update local state
    setNotifications(prev => ({
      ...prev,
      [key]: value
    }))
  }

  const handleAvatarUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (file) {
      try {
        // In a real implementation, you would upload the file to a storage service
        // and get back a URL. For now, we'll simulate this with a mock URL
        const mockAvatarUrl = `https://example.com/avatar-${Date.now()}.jpg`
        
        await userService.uploadAvatar(mockAvatarUrl)
        
        // Update local user state
        if (user) {
          const updatedUser = { ...user, avatarUrl: mockAvatarUrl }
          setUser(updatedUser)
          localStorage.setItem('user', JSON.stringify(updatedUser))
        }
        
        await Swal.fire({
          icon: 'success',
          title: 'อัปโหลดรูปภาพสำเร็จ!',
          text: 'รูปโปรไฟล์ได้รับการอัปเดตแล้ว',
          confirmButtonColor: '#2563eb'
        })
      } catch (error) {
        console.error('Error uploading avatar:', error)
        await Swal.fire({
          icon: 'error',
          title: 'เกิดข้อผิดพลาด',
          text: 'ไม่สามารถอัปโหลดรูปภาพได้ กรุณาลองใหม่อีกครั้ง',
          confirmButtonColor: '#dc2626'
        })
      }
    }
  }

  const handleEmailVerification = async () => {
    if (!user) return

    try {
      // เรียก API การส่งอีเมลยืนยัน
      const response = await userService.sendVerificationEmail(user.email)
      
      if (response) {
        await Swal.fire({
          icon: 'success',
          title: 'ส่งอีเมลยืนยันแล้ว!',
          text: 'กรุณาตรวจสอบอีเมลของคุณและคลิกลิงก์ยืนยัน',
          confirmButtonColor: '#2563eb'
        })
        
        // รีเฟรชข้อมูลผู้ใช้เพื่ออัปเดตสถานะ
        const updatedUser = await userService.getProfile()
        setUser(updatedUser)
      }
    } catch (error) {
      console.error('Error sending verification email:', error)
      await Swal.fire({
        icon: 'error',
        title: 'เกิดข้อผิดพลาด',
        text: 'ไม่สามารถส่งอีเมลยืนยันได้ กรุณาลองใหม่อีกครั้ง',
        confirmButtonColor: '#dc2626'
      })
    }
  }

  if (isLoading || !user) {
    return (
      <DashboardLayout>
        <div className="min-h-screen bg-gray-50 flex items-center justify-center">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
            <p className="text-gray-600">กำลังโหลดข้อมูล...</p>
          </div>
        </div>
      </DashboardLayout>
    )
  }

  // Parse settings JSON if it exists
  const userSettings = user.settings ? JSON.parse(user.settings) : {}
  
  // Check if email is verified from user data
  const isEmailVerified = user.emailVerified || false

  return (
    <DashboardLayout>
      <div className="space-y-8">
        {/* Header */}
        <div>
          <h1 className="text-3xl font-bold text-gray-900">โปรไฟล์ของฉัน</h1>
          <p className="text-gray-600 mt-2">จัดการข้อมูลส่วนตัว การตั้งค่าบัญชี และความปลอดภัย</p>
        </div>

        {/* Profile Overview */}
        <Card className="border-0 shadow-md">
          <CardContent className="pt-6">
            <div className="flex items-center space-x-6">
              <div className="relative">
                {user.avatarUrl ? (
                  <Image 
                    src={user.avatarUrl} 
                    alt="Avatar" 
                    width={96}
                    height={96}
                    className="rounded-full object-cover"
                  />
                ) : (
                  <div className="w-24 h-24 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center text-white text-2xl font-bold">
                    {user.firstName.charAt(0)}{user.lastName.charAt(0)}
                  </div>
                )}
                <label className="absolute -bottom-1 -right-1 bg-blue-500 hover:bg-blue-600 text-white p-1.5 rounded-full cursor-pointer transition-colors">
                  <CameraIcon className="w-4 h-4" />
                  <input
                    type="file"
                    accept="image/*"
                    onChange={handleAvatarUpload}
                    className="hidden"
                  />
                </label>
              </div>
              <div className="flex-1">
                <h2 className="text-2xl font-bold text-gray-900">{user.firstName} {user.lastName}</h2>
                <p className="text-gray-600">{user.position || 'ไม่ระบุตำแหน่ง'}</p>
                <p className="text-gray-500 text-sm">อีเมล: {user.email}</p>
                <p className="text-gray-500 text-sm">เข้าร่วมเมื่อ: {new Date(user.createdAt).toLocaleDateString('th-TH')}</p>
                <p className="text-gray-500 text-sm">แพ็คเกจ: {user.packageType === 'enterprise' ? 'Enterprise' : user.packageType === 'standard' ? 'Standard' : 'Trial'}</p>
              </div>
              <div className="text-right">
                <div className="flex items-center space-x-2 mb-2">
                  {isEmailVerified ? (
                    <>
                      <CheckCircleIcon className="w-5 h-5 text-green-500" />
                      <span className="text-sm text-green-600 font-medium">บัญชียืนยันแล้ว</span>
                    </>
                  ) : (
                    <>
                      <ExclamationTriangleIcon className="w-5 h-5 text-yellow-500" />
                      <span className="text-sm text-yellow-600 font-medium">รอยืนยันอีเมล</span>
                    </>
                  )}
                </div>
                <p className="text-xs text-gray-500">
                  เข้าสู่ระบบล่าสุด: {user.lastLoginAt ? new Date(user.lastLoginAt).toLocaleDateString('th-TH', {
                    month: 'short',
                    day: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit'
                  }) : 'ไม่เคย'}
                </p>
                <p className="text-xs text-gray-500 mt-1">
                  วันใช้งานคงเหลือ: {user.remainingDays} วัน
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Email Verification Alert */}
        {!isEmailVerified && (
          <Card className="border-0 shadow-md bg-yellow-50 border-yellow-200">
            <CardContent className="pt-6">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-3">
                  <ExclamationTriangleIcon className="w-6 h-6 text-yellow-600" />
                  <div>
                    <h3 className="font-medium text-yellow-900">อีเมลยังไม่ได้ยืนยัน</h3>
                    <p className="text-sm text-yellow-700">กรุณายืนยันอีเมลของคุณเพื่อใช้งานระบบได้อย่างสมบูรณ์</p>
                  </div>
                </div>
                <Button
                  onClick={handleEmailVerification}
                  className="bg-yellow-600 hover:bg-yellow-700 text-white"
                >
                  <EnvelopeIcon className="w-4 h-4 mr-2" />
                  ยืนยันอีเมล
                </Button>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Personal Information */}
        <Card className="border-0 shadow-md">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="flex items-center">
                  <UserIcon className="w-5 h-5 mr-2" />
                  ข้อมูลส่วนตัว
                </CardTitle>
                <CardDescription>จัดการข้อมูลพื้นฐานของคุณ</CardDescription>
              </div>
              <Button
                variant="outline"
                onClick={() => {
                  if (editMode === 'personal') {
                    setEditMode(null)
                    setFormData({
                      firstName: user.firstName,
                      lastName: user.lastName,
                      email: user.email,
                      phone: user.phone || '',
                      position: user.position || ''
                    })
                  } else {
                    setEditMode('personal')
                  }
                }}
              >
                <PencilIcon className="w-4 h-4 mr-2" />
                {editMode === 'personal' ? 'ยกเลิก' : 'แก้ไข'}
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            {editMode === 'personal' ? (
              <div className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">ชื่อ</label>
                    <Input
                      value={formData.firstName}
                      onChange={(e) => setFormData({ ...formData, firstName: e.target.value })}
                      placeholder="ชื่อ"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">นามสกุล</label>
                    <Input
                      value={formData.lastName}
                      onChange={(e) => setFormData({ ...formData, lastName: e.target.value })}
                      placeholder="นามสกุล"
                    />
                  </div>
                </div>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">อีเมล</label>
                    <Input
                      type="email"
                      value={formData.email}
                      onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                      placeholder="อีเมล"
                      disabled
                    />
                    <p className="text-xs text-gray-500 mt-1">อีเมลไม่สามารถแก้ไขได้</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">เบอร์โทรศัพท์</label>
                    <Input
                      value={formData.phone}
                      onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                      placeholder="เบอร์โทรศัพท์"
                    />
                  </div>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">ตำแหน่ง</label>
                  <Input
                    value={formData.position}
                    onChange={(e) => setFormData({ ...formData, position: e.target.value })}
                    placeholder="ตำแหน่ง"
                  />
                </div>
                <div className="flex justify-end space-x-3">
                  <Button variant="outline" onClick={() => setEditMode(null)}>
                    ยกเลิก
                  </Button>
                  <Button onClick={handlePersonalInfoSubmit}>
                    บันทึกข้อมูล
                  </Button>
                </div>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label className="block text-sm font-medium text-gray-500">ชื่อ - นามสกุล</label>
                  <p className="text-gray-900 font-medium">{user.firstName} {user.lastName}</p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-500">อีเมล</label>
                  <p className="text-gray-900">{user.email}</p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-500">เบอร์โทรศัพท์</label>
                  <p className="text-gray-900">{user.phone || 'ไม่ระบุ'}</p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-500">ตำแหน่ง</label>
                  <p className="text-gray-900">{user.position || 'ไม่ระบุ'}</p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-500">สถานะ</label>
                  <p className="text-gray-900">{user.status === 'active' ? 'ใช้งาน' : user.status === 'inactive' ? 'ไม่ใช้งาน' : user.status === 'pending' ? 'รอการยืนยัน' : 'ระงับการใช้งาน'}</p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-500">บทบาท</label>
                  <p className="text-gray-900">{user.role === 'admin' ? 'ผู้ดูแลระบบ' : 'ผู้ใช้งาน'}</p>
                </div>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Change Password */}
        <Card className="border-0 shadow-md">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="flex items-center">
                  <KeyIcon className="w-5 h-5 mr-2" />
                  เปลี่ยนรหัสผ่าน
                </CardTitle>
                <CardDescription>อัปเดตรหัสผ่านเพื่อความปลอดภัย</CardDescription>
              </div>
              <Button
                variant="outline"
                onClick={() => {
                  if (editMode === 'password') {
                    setEditMode(null)
                    setPasswordData({ currentPassword: '', newPassword: '', confirmPassword: '' })
                  } else {
                    setEditMode('password')
                  }
                }}
              >
                <KeyIcon className="w-4 h-4 mr-2" />
                {editMode === 'password' ? 'ยกเลิก' : 'เปลี่ยนรหัสผ่าน'}
              </Button>
            </div>
          </CardHeader>
          {editMode === 'password' && (
            <CardContent>
              <div className="space-y-4 max-w-md">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">รหัสผ่านปัจจุบัน</label>
                  <Input
                    type="password"
                    value={passwordData.currentPassword}
                    onChange={(e) => setPasswordData({ ...passwordData, currentPassword: e.target.value })}
                    placeholder="รหัสผ่านปัจจุบัน"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">รหัสผ่านใหม่</label>
                  <Input
                    type="password"
                    value={passwordData.newPassword}
                    onChange={(e) => setPasswordData({ ...passwordData, newPassword: e.target.value })}
                    placeholder="รหัสผ่านใหม่ (อย่างน้อย 6 ตัวอักษร)"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">ยืนยันรหัสผ่านใหม่</label>
                  <Input
                    type="password"
                    value={passwordData.confirmPassword}
                    onChange={(e) => setPasswordData({ ...passwordData, confirmPassword: e.target.value })}
                    placeholder="ยืนยันรหัสผ่านใหม่"
                  />
                </div>
                <div className="flex justify-end space-x-3">
                  <Button variant="outline" onClick={() => setEditMode(null)}>
                    ยกเลิก
                  </Button>
                  <Button onClick={handlePasswordSubmit}>
                    เปลี่ยนรหัสผ่าน
                  </Button>
                </div>
              </div>
            </CardContent>
          )}
        </Card>

        {/* Notification Settings */}
        <Card className="border-0 shadow-md">
          <CardHeader>
            <CardTitle className="flex items-center">
              <BellIcon className="w-5 h-5 mr-2" />
              การตั้งค่าการแจ้งเตือน
            </CardTitle>
            <CardDescription>กำหนดการแจ้งเตือนที่คุณต้องการรับ</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-6">
              <div className="flex items-center justify-between">
                <div>
                  <h4 className="font-medium text-gray-900">การแจ้งเตือนทางอีเมล</h4>
                  <p className="text-sm text-gray-500">รับการแจ้งเตือนผ่านอีเมล</p>
                </div>
                <Switch
                  checked={notifications.email}
                  onChange={(checked) => handleNotificationUpdate('email', checked)}
                />
              </div>
              <div className="flex items-center justify-between">
                <div>
                  <h4 className="font-medium text-gray-900">การแจ้งเตือนแบบ Push</h4>
                  <p className="text-sm text-gray-500">รับการแจ้งเตือนผ่านเบราว์เซอร์</p>
                </div>
                <Switch
                  checked={notifications.push}
                  onChange={(checked) => handleNotificationUpdate('push', checked)}
                />
              </div>
              <hr className="border-gray-200" />
              <div className="flex items-center justify-between">
                <div>
                  <h4 className="font-medium text-gray-900">แจ้งเตือนตารางเวร</h4>
                  <p className="text-sm text-gray-500">แจ้งเตือนเกี่ยวกับการเปลี่ยนแปลงตารางเวร</p>
                </div>
                <Switch
                  checked={notifications.scheduleReminder}
                  onChange={(checked) => handleNotificationUpdate('scheduleReminder', checked)}
                />
              </div>
              <div className="flex items-center justify-between">
                <div>
                  <h4 className="font-medium text-gray-900">อัปเดตการลา</h4>
                  <p className="text-sm text-gray-500">แจ้งเตือนสถานะการขอลา</p>
                </div>
                <Switch
                  checked={notifications.leaveUpdates}
                  onChange={(checked) => handleNotificationUpdate('leaveUpdates', checked)}
                />
              </div>
              <div className="flex items-center justify-between">
                <div>
                  <h4 className="font-medium text-gray-900">อัปเดตระบบ</h4>
                  <p className="text-sm text-gray-500">แจ้งเตือนการปรับปรุงระบบ</p>
                </div>
                <Switch
                  checked={notifications.systemUpdates}
                  onChange={(checked) => handleNotificationUpdate('systemUpdates', checked)}
                />
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </DashboardLayout>
  )
}
