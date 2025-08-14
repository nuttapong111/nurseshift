'use client'

import { useState } from 'react'
import EmailVerification from '@/components/EmailVerification'
import useEmailVerification from '@/hooks/useEmailVerification'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card'

export default function EmailVerificationPage() {
  const [showVerificationStatus, setShowVerificationStatus] = useState(false)
  const [verificationEmail, setVerificationEmail] = useState('')
  const { checkEmailVerification, loading: checkLoading, error: checkError } = useEmailVerification()

  const handleCheckVerification = async () => {
    if (!verificationEmail.trim()) return
    
    const response = await checkEmailVerification(verificationEmail.trim())
    if (response) {
      setShowVerificationStatus(true)
    }
  }

  const handleVerificationComplete = () => {
    console.log('Email verification completed successfully!')
  }

  return (
    <DashboardLayout>
      <div className="space-y-8">
        <div className="text-center">
          <h1 className="text-3xl font-bold text-gray-900 mb-4">
            ระบบยืนยันอีเมล
          </h1>
          <p className="text-lg text-gray-600">
            ทดสอบการทำงานของ API การยืนยันอีเมล
          </p>
        </div>

        <div className="grid md:grid-cols-2 gap-8">
          {/* Email Verification Component */}
          <Card>
            <CardHeader>
              <CardTitle>ยืนยันอีเมล</CardTitle>
              <CardDescription>
                ส่งอีเมลยืนยันและยืนยันด้วยรหัส 6 หลัก
              </CardDescription>
            </CardHeader>
            <CardContent>
              <EmailVerification onVerificationComplete={handleVerificationComplete} />
            </CardContent>
          </Card>

          {/* Check Verification Status */}
          <Card>
            <CardHeader>
              <CardTitle>ตรวจสอบสถานะการยืนยัน</CardTitle>
              <CardDescription>
                ตรวจสอบว่าอีเมลได้รับการยืนยันแล้วหรือไม่
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div>
                  <label htmlFor="checkEmail" className="block text-sm font-medium text-gray-700 mb-2">
                    อีเมลที่ต้องการตรวจสอบ
                  </label>
                  <input
                    id="checkEmail"
                    type="email"
                    value={verificationEmail}
                    onChange={(e) => setVerificationEmail(e.target.value)}
                    placeholder="example@email.com"
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </div>

                <button
                  onClick={handleCheckVerification}
                  disabled={checkLoading || !verificationEmail.trim()}
                  className="w-full bg-green-600 text-white px-4 py-2 rounded-md hover:bg-green-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
                >
                  {checkLoading ? 'กำลังตรวจสอบ...' : 'ตรวจสอบสถานะ'}
                </button>

                {checkError && (
                  <div className="p-3 bg-red-100 border border-red-400 text-red-700 rounded">
                    {checkError}
                  </div>
                )}

                {showVerificationStatus && (
                  <div className="p-4 bg-blue-50 border border-blue-200 rounded-md">
                    <h3 className="font-medium text-blue-900 mb-2">
                      สถานะการยืนยันอีเมล
                    </h3>
                    <p className="text-blue-700">
                      อีเมล: <span className="font-medium">{verificationEmail}</span>
                    </p>
                    <p className="text-blue-700">
                      สถานะ: <span className="font-medium">ได้รับการตรวจสอบแล้ว</span>
                    </p>
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        </div>

        {/* API Documentation */}
        <Card>
          <CardHeader>
            <CardTitle>API Endpoints ที่ใช้</CardTitle>
            <CardDescription>
              รายการ API endpoints สำหรับการยืนยันอีเมล
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-3 text-sm">
              <div className="flex items-center space-x-2">
                <span className="bg-blue-100 text-blue-800 px-2 py-1 rounded text-xs font-medium">
                  POST
                </span>
                <code className="bg-gray-100 px-2 py-1 rounded">
                  /api/v1/users/send-verification-email
                </code>
                <span className="text-gray-600">ส่งอีเมลยืนยัน</span>
              </div>
              <div className="flex items-center space-x-2">
                <span className="bg-green-100 text-green-800 px-2 py-1 rounded text-xs font-medium">
                  POST
                </span>
                <code className="bg-gray-100 px-2 py-1 rounded">
                  /api/v1/users/verify-email
                </code>
                <span className="text-gray-600">ยืนยันอีเมลด้วย token</span>
              </div>
              <div className="flex items-center space-x-2">
                <span className="bg-yellow-100 text-yellow-800 px-2 py-1 rounded text-xs font-medium">
                  GET
                </span>
                <code className="bg-gray-100 px-2 py-1 rounded">
                  /api/v1/users/check-email-verification/:email
                </code>
                <span className="text-gray-600">ตรวจสอบสถานะการยืนยัน</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </DashboardLayout>
  )
}
