'use client'

import { useState } from 'react'
import useEmailVerification from '@/hooks/useEmailVerification'

interface EmailVerificationProps {
  email?: string
  onVerificationComplete?: () => void
}

export const EmailVerification: React.FC<EmailVerificationProps> = ({ 
  email = '', 
  onVerificationComplete 
}) => {
  const [inputEmail, setInputEmail] = useState(email)
  const [verificationToken, setVerificationToken] = useState('')
  const [step, setStep] = useState<'input' | 'verification'>('input')
  
  const { 
    sendLoading, 
    sendError, 
    verifyLoading,
    verifyError,
    sendVerificationEmail, 
    verifyEmail, 
    clearError 
  } = useEmailVerification()

  const handleSendVerificationEmail = async () => {
    if (!inputEmail.trim()) {
      return
    }

    clearError()
    const response = await sendVerificationEmail(inputEmail.trim())
    
    if (response) {
      setStep('verification')
    }
  }

  const handleVerifyEmail = async () => {
    if (!verificationToken.trim()) {
      return
    }

    clearError()
    const response = await verifyEmail(verificationToken.trim())
    
    if (response) {
      onVerificationComplete?.()
    }
  }

  const handleBackToInput = () => {
    setStep('input')
    setVerificationToken('')
    clearError()
  }

  if (step === 'verification') {
    return (
      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="text-center mb-6">
          <h3 className="text-lg font-medium text-gray-900">ยืนยันอีเมล</h3>
          <p className="mt-2 text-sm text-gray-500">
            กรุณาใส่รหัสยืนยันที่ส่งไปยัง <span className="font-medium">{inputEmail}</span>
          </p>
        </div>

        {(sendError || verifyError) && (
          <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
            {sendError || verifyError}
          </div>
        )}

        <div className="space-y-4">
          <div>
            <label htmlFor="verificationToken" className="block text-sm font-medium text-gray-700 mb-2">
              รหัสยืนยัน
            </label>
            <input
              id="verificationToken"
              type="text"
              value={verificationToken}
              onChange={(e) => setVerificationToken(e.target.value)}
              placeholder="ใส่รหัสยืนยัน"
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>

          <button
            onClick={handleVerifyEmail}
            disabled={verifyLoading || !verificationToken.trim()}
            className="w-full bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
          >
            {verifyLoading ? 'กำลังยืนยัน...' : 'ยืนยันอีเมล'}
          </button>

          <button
            onClick={handleBackToInput}
            className="w-full bg-gray-200 text-gray-700 px-4 py-2 rounded-md hover:bg-gray-300 transition-colors"
          >
            เปลี่ยนอีเมล
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <div className="text-center mb-6">
        <h3 className="text-lg font-medium text-gray-900">ส่งอีเมลยืนยัน</h3>
        <p className="mt-2 text-sm text-gray-500">
          กรุณาใส่อีเมลที่ต้องการยืนยัน
        </p>
      </div>

      {(sendError || verifyError) && (
        <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
          {sendError || verifyError}
        </div>
      )}

      <div className="space-y-4">
        <div>
          <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">
            อีเมล
          </label>
          <input
            id="email"
            type="email"
            value={inputEmail}
            onChange={(e) => setInputEmail(e.target.value)}
            placeholder="example@email.com"
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        <button
          onClick={handleSendVerificationEmail}
          disabled={sendLoading || !inputEmail.trim()}
          className="w-full bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
        >
          {sendLoading ? 'กำลังส่ง...' : 'ส่งอีเมลยืนยัน'}
        </button>
      </div>
    </div>
  )
}

export default EmailVerification
