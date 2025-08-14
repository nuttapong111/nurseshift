'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { useRouter, useSearchParams } from 'next/navigation'
import { CalendarDaysIcon, ArrowLeftIcon, EyeIcon, EyeSlashIcon } from '@heroicons/react/24/outline'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card'
import Swal from 'sweetalert2'

export default function ResetPasswordPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const token = searchParams.get('token')
  
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [isValidToken, setIsValidToken] = useState(false)

  useEffect(() => {
    if (!token) {
      Swal.fire({
        icon: 'error',
        title: 'Token ไม่ถูกต้อง',
        text: 'ลิงก์รีเซ็ตรหัสผ่านไม่ถูกต้อง กรุณาขอรหัสยืนยันใหม่',
        confirmButtonText: 'กลับไปหน้าลืมรหัสผ่าน',
        confirmButtonColor: '#dc2626'
      }).then(() => {
        router.push('/auth/forgot-password')
      })
    } else {
      setIsValidToken(true)
    }
  }, [token, router])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (newPassword !== confirmPassword) {
      await Swal.fire({
        icon: 'error',
        title: 'รหัสผ่านไม่ตรงกัน',
        text: 'กรุณากรอกรหัสผ่านให้ตรงกัน',
        confirmButtonText: 'ตกลง',
        confirmButtonColor: '#dc2626'
      })
      return
    }

    if (newPassword.length < 6) {
      await Swal.fire({
        icon: 'error',
        title: 'รหัสผ่านสั้นเกินไป',
        text: 'รหัสผ่านต้องมีความยาวอย่างน้อย 6 ตัวอักษร',
        confirmButtonText: 'ตกลง',
        confirmButtonColor: '#dc2626'
      })
      return
    }

    setIsLoading(true)

    try {
      const response = await fetch('http://localhost:8081/api/v1/auth/reset-password', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          token: token,
          newPassword: newPassword
        }),
      })

      const data = await response.json()

      if (response.ok) {
        await Swal.fire({
          icon: 'success',
          title: 'รีเซ็ตรหัสผ่านสำเร็จ!',
          text: 'รหัสผ่านของคุณได้รับการอัปเดตแล้ว กรุณาเข้าสู่ระบบด้วยรหัสผ่านใหม่',
          confirmButtonText: 'ไปหน้าเข้าสู่ระบบ',
          confirmButtonColor: '#2563eb'
        })
        router.push('/auth/login')
      } else {
        // Handle specific error cases
        let errorMessage = 'เกิดข้อผิดพลาดในการรีเซ็ตรหัสผ่าน'
        
        if (data.error === 'invalid token') {
          errorMessage = 'ลิงก์รีเซ็ตรหัสผ่านไม่ถูกต้องหรือหมดอายุแล้ว กรุณาขอรีเซ็ตรหัสผ่านใหม่'
        } else if (data.error === 'token expired') {
          errorMessage = 'ลิงก์รีเซ็ตรหัสผ่านหมดอายุแล้ว กรุณาขอรีเซ็ตรหัสผ่านใหม่'
        } else if (data.error === 'password too short') {
          errorMessage = 'รหัสผ่านต้องมีอย่างน้อย 6 ตัวอักษร'
        } else if (data.error === 'weak password') {
          errorMessage = 'รหัสผ่านไม่ปลอดภัย กรุณาใช้รหัสผ่านที่ซับซ้อนกว่านี้'
        } else if (data.message) {
          errorMessage = data.message
        }
        
        throw new Error(errorMessage)
      }
    } catch (error) {
      await Swal.fire({
        icon: 'error',
        title: 'เกิดข้อผิดพลาด',
        text: error instanceof Error ? error.message : 'ไม่สามารถรีเซ็ตรหัสผ่านได้ กรุณาลองใหม่อีกครั้ง',
        confirmButtonText: 'ลองใหม่',
        confirmButtonColor: '#dc2626'
      })
    } finally {
      setIsLoading(false)
    }
  }

  if (!isValidToken) {
    return null
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50 flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        {/* Header */}
        <div className="text-center mb-8">
          <Link href="/" className="inline-flex items-center space-x-2 mb-6">
            <div className="w-12 h-12 bg-gradient-to-r from-blue-600 to-purple-600 rounded-lg flex items-center justify-center">
              <CalendarDaysIcon className="w-7 h-7 text-white" />
            </div>
            <h1 className="text-3xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
              NurseShift
            </h1>
          </Link>
          <h2 className="text-2xl font-bold text-gray-900 mb-2">ตั้งรหัสผ่านใหม่</h2>
          <p className="text-gray-600">กรุณากรอกรหัสผ่านใหม่ที่ต้องการใช้</p>
        </div>

        {/* Back to Login Link */}
        <div className="mb-6">
          <Link 
            href="/auth/login" 
            className="inline-flex items-center space-x-2 text-sm text-blue-600 hover:text-blue-500"
          >
            <ArrowLeftIcon className="w-4 h-4" />
            <span>กลับไปหน้าเข้าสู่ระบบ</span>
          </Link>
        </div>

        {/* Reset Password Form */}
        <Card className="shadow-lg border-0">
          <CardHeader>
            <CardTitle className="text-center">ตั้งรหัสผ่านใหม่</CardTitle>
            <CardDescription className="text-center">
              กรุณากรอกรหัสผ่านใหม่ที่ต้องการใช้
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label htmlFor="newPassword" className="block text-sm font-medium text-gray-700 mb-2">
                  รหัสผ่านใหม่
                </label>
                <div className="relative">
                  <Input
                    id="newPassword"
                    name="newPassword"
                    type={showPassword ? 'text' : 'password'}
                    required
                    value={newPassword}
                    onChange={(e) => setNewPassword(e.target.value)}
                    placeholder="รหัสผ่านใหม่ (อย่างน้อย 6 ตัวอักษร)"
                    className="w-full pr-10"
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute inset-y-0 right-0 pr-3 flex items-center"
                  >
                    {showPassword ? (
                      <EyeSlashIcon className="h-5 w-5 text-gray-400" />
                    ) : (
                      <EyeIcon className="h-5 w-5 text-gray-400" />
                    )}
                  </button>
                </div>
              </div>

              <div>
                <label htmlFor="confirmPassword" className="block text-sm font-medium text-gray-700 mb-2">
                  ยืนยันรหัสผ่านใหม่
                </label>
                <div className="relative">
                  <Input
                    id="confirmPassword"
                    name="confirmPassword"
                    type={showConfirmPassword ? 'text' : 'password'}
                    required
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    placeholder="ยืนยันรหัสผ่านใหม่"
                    className="w-full pr-10"
                  />
                  <button
                    type="button"
                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                    className="absolute inset-y-0 right-0 pr-3 flex items-center"
                  >
                    {showConfirmPassword ? (
                      <EyeSlashIcon className="h-5 w-5 text-gray-400" />
                    ) : (
                      <EyeIcon className="h-5 w-5 text-gray-400" />
                    )}
                  </button>
                </div>
              </div>

              <Button
                type="submit"
                className="w-full bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700"
                disabled={isLoading}
              >
                {isLoading ? 'กำลังตั้งรหัสผ่านใหม่...' : 'ตั้งรหัสผ่านใหม่'}
              </Button>
            </form>
          </CardContent>
        </Card>

        {/* Help Section */}
        <Card className="mt-4 bg-gray-50 border-gray-200">
          <CardContent className="pt-6">
            <h3 className="font-medium text-gray-900 mb-3">ต้องการความช่วยเหลือ?</h3>
            <p className="text-sm text-gray-600 mb-3">
              หากคุณยังคงมีปัญหาในการตั้งรหัสผ่านใหม่ กรุณาติดต่อทีมสนับสนุน
            </p>
            <Button variant="outline" size="sm" className="w-full">
              ติดต่อฝ่ายสนับสนุน
            </Button>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
