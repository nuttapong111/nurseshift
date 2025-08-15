import { useState } from 'react'
import { userService } from '@/services/userService'
import type { 
  SendVerificationEmailResponse, 
  VerifyEmailResponse
} from '@/types'

export const useEmailVerification = () => {
  // State for send verification email
  const [sendLoading, setSendLoading] = useState(false)
  const [sendError, setSendError] = useState<string | null>(null)

  // State for verify email
  const [verifyLoading, setVerifyLoading] = useState(false)
  const [verifyError, setVerifyError] = useState<string | null>(null)

  // State for check verification
  const [checkLoading, setCheckLoading] = useState(false)
  const [checkError, setCheckError] = useState<string | null>(null)

  const sendVerificationEmail = async (email: string): Promise<SendVerificationEmailResponse | null> => {
    setSendLoading(true)
    setSendError(null)
    
    try {
      const response = await userService.sendVerificationEmail(email)
      return response
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'เกิดข้อผิดพลาดในการส่งอีเมลยืนยัน'
      setSendError(errorMessage)
      return null
    } finally {
      setSendLoading(false)
    }
  }

  const verifyEmail = async (token: string): Promise<VerifyEmailResponse | null> => {
    setVerifyLoading(true)
    setVerifyError(null)
    
    try {
      const response = await userService.verifyEmail(token)
      return response
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'เกิดข้อผิดพลาดในการยืนยันอีเมล'
      setVerifyError(errorMessage)
      return null
    } finally {
      setVerifyLoading(false)
    }
  }

  const checkEmailVerification = async (email: string): Promise<boolean> => {
    setCheckLoading(true)
    setCheckError(null)
    
    try {
      const response = await userService.checkEmailVerification(email)
      return response.isVerified || false
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'เกิดข้อผิดพลาดในการตรวจสอบสถานะ'
      setCheckError(errorMessage)
      return false
    } finally {
      setCheckLoading(false)
    }
  }

  const clearError = () => {
    setSendError(null)
    setVerifyError(null)
    setCheckError(null)
  }

  return {
    // Send verification email
    sendLoading,
    sendError,
    sendVerificationEmail,

    // Verify email
    verifyLoading,
    verifyError,
    verifyEmail,

    // Check verification status
    checkLoading,
    checkError,
    checkEmailVerification,

    // Utility
    clearError
  }
}

export default useEmailVerification
