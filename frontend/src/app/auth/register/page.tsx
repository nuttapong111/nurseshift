'use client'

import { useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { CalendarDaysIcon, EyeIcon, EyeSlashIcon } from '@heroicons/react/24/outline'
import { Button } from '@/components/ui/Button'
import { normalizeBaseUrl } from '@/lib/utils'
import { ButtonGroup } from '@/components/ui/ButtonGroup'
import { Input } from '@/components/ui/Input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card'
import Swal from 'sweetalert2'

export default function RegisterPage() {
  const router = useRouter()
  const [formData, setFormData] = useState({
    firstName: '',
    lastName: '',
    position: '',
    email: '',
    password: '',
    confirmPassword: ''
  })
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const [isLoading, setIsLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)

    try {
      // Validation
      if (formData.password !== formData.confirmPassword) {
        throw new Error('รหัสผ่านไม่ตรงกัน')
      }

      if (formData.password.length < 6) {
        throw new Error('รหัสผ่านต้องมีอย่างน้อย 6 ตัวอักษร')
      }

      // Call real API
      const response = await fetch(`${normalizeBaseUrl(process.env.NEXT_PUBLIC_AUTH_SERVICE_URL)}/api/v1/auth/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          firstName: formData.firstName,
          lastName: formData.lastName,
          position: formData.position,
          email: formData.email,
          password: formData.password,
        }),
      })

      const data = await response.json()

      if (response.ok) {
        await Swal.fire({
          icon: 'success',
          title: 'สมัครสมาชิกสำเร็จ!',
          html: `
            <p>ยินดีต้อนรับสู่ระบบ NurseShift</p>
            <p class="text-green-600 font-medium">คุณได้รับเวลาทดลองใช้ฟรี 30 วัน</p>
          `,
          confirmButtonText: 'เริ่มใช้งาน',
          confirmButtonColor: '#2563eb'
        })

        router.push('/auth/login')
      } else {
        // Handle specific error cases
        let errorMessage = 'สมัครสมาชิกไม่สำเร็จ'
        
        if (data.error === 'email already exists') {
          errorMessage = 'อีเมลนี้มีในระบบอยู่แล้ว กรุณาใช้อีเมลอื่นหรือเข้าสู่ระบบ'
        } else if (data.error === 'invalid email') {
          errorMessage = 'รูปแบบอีเมลไม่ถูกต้อง กรุณาตรวจสอบอีกครั้ง'
        } else if (data.error === 'password too short') {
          errorMessage = 'รหัสผ่านต้องมีอย่างน้อย 6 ตัวอักษร'
        } else if (data.message) {
          errorMessage = data.message
        }
        
        throw new Error(errorMessage)
      }
    } catch (error) {
      await Swal.fire({
        icon: 'error',
        title: 'สมัครสมาชิกไม่สำเร็จ',
        text: error instanceof Error ? error.message : 'เกิดข้อผิดพลาดไม่ทราบสาเหตุ',
        confirmButtonText: 'ลองใหม่',
        confirmButtonColor: '#dc2626'
      })
    } finally {
      setIsLoading(false)
    }
  }

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    })
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
          <h2 className="text-2xl font-bold text-gray-900 mb-2">สมัครสมาชิก</h2>
          <p className="text-gray-600">เริ่มต้นใช้งานระบบจัดตารางเวรพยาบาล</p>
        </div>

        {/* Register Form */}
        <Card className="shadow-lg border-0">
          <CardHeader>
            <CardTitle className="text-center">สมัครสมาชิก</CardTitle>
            <CardDescription className="text-center">
              ทดลองใช้ฟรี 90 วัน ไม่ต้องใช้บัตรเครดิต
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label htmlFor="firstName" className="block text-sm font-medium text-gray-700 mb-2">
                    ชื่อ
                  </label>
                  <Input
                    id="firstName"
                    name="firstName"
                    type="text"
                    required
                    value={formData.firstName}
                    onChange={handleChange}
                    placeholder="ชื่อ"
                  />
                </div>
                
                <div>
                  <label htmlFor="lastName" className="block text-sm font-medium text-gray-700 mb-2">
                    นามสกุล
                  </label>
                  <Input
                    id="lastName"
                    name="lastName"
                    type="text"
                    required
                    value={formData.lastName}
                    onChange={handleChange}
                    placeholder="นามสกุล"
                  />
                </div>
              </div>

              <div>
                <label htmlFor="position" className="block text-sm font-medium text-gray-700 mb-2">
                  ตำแหน่ง
                </label>
                <Input
                  id="position"
                  name="position"
                  type="text"
                  required
                  value={formData.position}
                  onChange={handleChange}
                  placeholder="เช่น หัวหน้าพยาบาล, พยาบาลวิชาชีพ"
                />
              </div>

              <div>
                <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">
                  อีเมล
                </label>
                <Input
                  id="email"
                  name="email"
                  type="email"
                  required
                  value={formData.email}
                  onChange={handleChange}
                  placeholder="กรอกอีเมลของคุณ"
                />
              </div>
              
              <div>
                <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-2">
                  รหัสผ่าน
                </label>
                <div className="relative">
                  <Input
                    id="password"
                    name="password"
                    type={showPassword ? 'text' : 'password'}
                    required
                    value={formData.password}
                    onChange={handleChange}
                    placeholder="สร้างรหัสผ่าน (อย่างน้อย 6 ตัวอักษร)"
                    className="pr-10"
                  />
                  <button
                    type="button"
                    className="absolute inset-y-0 right-0 pr-3 flex items-center"
                    onClick={() => setShowPassword(!showPassword)}
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
                  ยืนยันรหัสผ่าน
                </label>
                <div className="relative">
                  <Input
                    id="confirmPassword"
                    name="confirmPassword"
                    type={showConfirmPassword ? 'text' : 'password'}
                    required
                    value={formData.confirmPassword}
                    onChange={handleChange}
                    placeholder="ยืนยันรหัสผ่านอีกครั้ง"
                    className="pr-10"
                  />
                  <button
                    type="button"
                    className="absolute inset-y-0 right-0 pr-3 flex items-center"
                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                  >
                    {showConfirmPassword ? (
                      <EyeSlashIcon className="h-5 w-5 text-gray-400" />
                    ) : (
                      <EyeIcon className="h-5 w-5 text-gray-400" />
                    )}
                  </button>
                </div>
              </div>

              <div className="flex items-center">
                <input
                  id="agree-terms"
                  name="agree-terms"
                  type="checkbox"
                  required
                  className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                />
                <label htmlFor="agree-terms" className="ml-2 block text-sm text-gray-900">
                  ยอมรับ{' '}
                  <Link href="#" className="text-blue-600 hover:text-blue-500">
                    เงื่อนไขการใช้งาน
                  </Link>{' '}
                  และ{' '}
                  <Link href="#" className="text-blue-600 hover:text-blue-500">
                    นโยบายความเป็นส่วนตัว
                  </Link>
                </label>
              </div>

              <ButtonGroup direction="horizontal" spacing="normal" className="w-full">
                <Button
                  type="submit"
                  className="w-full bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700"
                  disabled={isLoading}
                >
                  {isLoading ? 'กำลังสมัครสมาชิก...' : 'สมัครสมาชิก'}
                </Button>
              </ButtonGroup>
            </form>

            <div className="mt-6">
              <div className="relative">
                <div className="absolute inset-0 flex items-center">
                  <div className="w-full border-t border-gray-300" />
                </div>
                <div className="relative flex justify-center text-sm">
                  <span className="px-2 bg-white text-gray-500">หรือ</span>
                </div>
              </div>

              <div className="mt-6 text-center">
                <p className="text-sm text-gray-600">
                  มีบัญชีอยู่แล้ว?{' '}
                  <Link href="/auth/login" className="font-medium text-blue-600 hover:text-blue-500">
                    เข้าสู่ระบบ
                  </Link>
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Free Trial Info */}
        <div className="mt-4 text-center">
          <div className="bg-green-50 border border-green-200 rounded-lg p-4">
            <p className="text-sm text-green-800">
              🎉 <strong>ฟรี 90 วัน!</strong> ใช้งานครบทุกฟีเจอร์ ไม่ต้องใช้บัตรเครดิต
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
