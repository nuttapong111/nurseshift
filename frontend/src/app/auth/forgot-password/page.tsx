'use client'

import { useState } from 'react'
import Link from 'next/link'
import { CalendarDaysIcon, ArrowLeftIcon } from '@heroicons/react/24/outline'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card'
import Swal from 'sweetalert2'
import { normalizeBaseUrl } from '@/lib/utils'

export default function ForgotPasswordPage() {
  const [email, setEmail] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [emailSent, setEmailSent] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)

    try {
      const response = await fetch(`${normalizeBaseUrl(process.env.NEXT_PUBLIC_AUTH_SERVICE_URL)}/api/v1/auth/forgot-password`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email }),
      })

      const data = await response.json()

      if (response.ok) {
        setEmailSent(true)
        
        await Swal.fire({
          icon: 'success',
          title: 'ส่งอีเมลสำเร็จ!',
          html: `
            <p>เราได้ส่งลิงก์สำหรับรีเซ็ตรหัสผ่านไปยัง</p>
            <p class="font-medium text-blue-600">${email}</p>
            <p class="text-sm text-gray-600 mt-2">กรุณาตรวจสอบอีเมลของคุณและคลิกลิงก์เพื่อรีเซ็ตรหัสผ่าน</p>
            <p class="text-sm text-gray-600 mt-2">ลิงก์จะหมดอายุใน 15 นาที</p>
          `,
          confirmButtonText: 'ตกลง',
          confirmButtonColor: '#2563eb'
        })
      } else {
        // Handle specific error cases
        let errorMessage = 'ไม่สามารถส่งอีเมลได้'
        
        if (data.error === 'user not found') {
          errorMessage = 'ไม่พบอีเมลนี้ในระบบ กรุณาตรวจสอบอีเมลหรือสมัครสมาชิกใหม่'
        } else if (data.error === 'invalid email') {
          errorMessage = 'รูปแบบอีเมลไม่ถูกต้อง กรุณาตรวจสอบอีกครั้ง'
        } else if (data.error === 'too many requests') {
          errorMessage = 'คุณได้ขอรีเซ็ตรหัสผ่านบ่อยเกินไป กรุณารอสักครู่แล้วลองใหม่'
        } else if (data.message) {
          errorMessage = data.message
        }
        
        throw new Error(errorMessage)
      }
    } catch (error) {
      await Swal.fire({
        icon: 'error',
        title: 'เกิดข้อผิดพลาด',
        text: error instanceof Error ? error.message : 'ไม่สามารถส่งอีเมลได้ กรุณาลองใหม่อีกครั้ง',
        confirmButtonText: 'ลองใหม่',
        confirmButtonColor: '#dc2626'
      })
    } finally {
      setIsLoading(false)
    }
  }

  const handleResendEmail = async () => {
    setIsLoading(true)
    
    try {
      const response = await fetch(`${normalizeBaseUrl(process.env.NEXT_PUBLIC_AUTH_SERVICE_URL)}/api/v1/auth/forgot-password`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email }),
      })

      const data = await response.json()

      if (response.ok) {
        await Swal.fire({
          icon: 'success',
          title: 'ส่งอีเมลใหม่สำเร็จ!',
          text: 'กรุณาตรวจสอบอีเมลของคุณอีกครั้ง ลิงก์ใหม่จะหมดอายุใน 15 นาที',
          confirmButtonText: 'ตกลง',
          confirmButtonColor: '#2563eb'
        })
      } else {
        // Handle specific error cases
        let errorMessage = 'ไม่สามารถส่งอีเมลได้'
        
        if (data.error === 'user not found') {
          errorMessage = 'ไม่พบอีเมลนี้ในระบบ กรุณาตรวจสอบอีเมลหรือสมัครสมาชิกใหม่'
        } else if (data.error === 'invalid email') {
          errorMessage = 'รูปแบบอีเมลไม่ถูกต้อง กรุณาตรวจสอบอีกครั้ง'
        } else if (data.error === 'too many requests') {
          errorMessage = 'คุณได้ขอรีเซ็ตรหัสผ่านบ่อยเกินไป กรุณารอสักครู่แล้วลองใหม่'
        } else if (data.message) {
          errorMessage = data.message
        }
        
        throw new Error(errorMessage)
      }
    } catch (error) {
      await Swal.fire({
        icon: 'error',
        title: 'เกิดข้อผิดพลาด',
        text: error instanceof Error ? error.message : 'ไม่สามารถส่งอีเมลได้ กรุณาลองใหม่อีกครั้ง',
        confirmButtonText: 'ลองใหม่',
        confirmButtonColor: '#dc2626'
      })
    } finally {
      setIsLoading(false)
    }
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
          <h2 className="text-2xl font-bold text-gray-900 mb-2">ลืมรหัสผ่าน</h2>
          <p className="text-gray-600">กรอกอีเมลเพื่อรีเซ็ตรหัสผ่าน</p>
        </div>

        {/* Back to Login Link */}
        <div className="mb-6">
          <Link 
            href="/auth/login" 
            className="inline-flex items-center text-sm text-blue-600 hover:text-blue-500"
          >
            <ArrowLeftIcon className="w-4 h-4 mr-1" />
            กลับไปหน้าเข้าสู่ระบบ
          </Link>
        </div>

        {/* Forgot Password Form */}
        <Card className="shadow-lg border-0">
          <CardHeader>
            <CardTitle className="text-center">
              {emailSent ? 'ตรวจสอบอีเมลของคุณ' : 'รีเซ็ตรหัสผ่าน'}
            </CardTitle>
            <CardDescription className="text-center">
              {emailSent 
                ? 'เราได้ส่งลิงก์สำหรับรีเซ็ตรหัสผ่านไปยังอีเมลของคุณแล้ว'
                : 'กรอกอีเมลที่ใช้สมัครสมาชิก เราจะส่งลิงก์สำหรับรีเซ็ตรหัสผ่านให้คุณ'
              }
            </CardDescription>
          </CardHeader>
          <CardContent>
            {!emailSent ? (
              <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                  <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">
                    อีเมล
                  </label>
                  <Input
                    id="email"
                    name="email"
                    type="email"
                    required
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="กรอกอีเมลของคุณ"
                    className="w-full"
                  />
                </div>

                <Button
                  type="submit"
                  className="w-full bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700"
                  disabled={isLoading}
                >
                  {isLoading ? 'กำลังส่งอีเมล...' : 'ส่งลิงก์รีเซ็ตรหัสผ่าน'}
                </Button>
              </form>
            ) : (
              <div className="space-y-4">
                <div className="text-center p-6 bg-green-50 rounded-lg border border-green-200">
                  <div className="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
                    <svg className="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                  </div>
                  <h3 className="text-lg font-medium text-green-900 mb-2">อีเมลถูกส่งแล้ว!</h3>
                  <p className="text-sm text-green-700 mb-4">
                    เราได้ส่งลิงก์สำหรับรีเซ็ตรหัสผ่านไปยัง
                  </p>
                  <p className="font-medium text-green-900">{email}</p>
                  <p className="text-xs text-green-600 mt-2">ลิงก์จะหมดอายุใน 15 นาที</p>
                </div>

                <div className="text-center space-y-3">
                  <p className="text-sm text-gray-600">
                    ไม่ได้รับอีเมล? ตรวจสอบในโฟลเดอร์สแปมหรือ
                  </p>
                  <Button
                    onClick={handleResendEmail}
                    variant="outline"
                    disabled={isLoading}
                    className="w-full"
                  >
                    {isLoading ? 'กำลังส่งใหม่...' : 'ส่งอีเมลใหม่'}
                  </Button>
                </div>

                <div className="text-center">
                  <Link 
                    href="/auth/login" 
                    className="text-sm text-blue-600 hover:text-blue-500"
                  >
                    กลับไปหน้าเข้าสู่ระบบ
                  </Link>
                </div>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Help Section */}
        <Card className="mt-4 bg-gray-50 border-gray-200">
          <CardContent className="pt-6">
            <h3 className="font-medium text-gray-900 mb-3">ต้องการความช่วยเหลือ?</h3>
            <p className="text-sm text-gray-600 mb-3">
              หากคุณยังคงมีปัญหาในการรีเซ็ตรหัสผ่าน กรุณาติดต่อทีมสนับสนุน
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
