import { useState } from 'react'
import { userService } from '@/services/userService'
import type { 
  SendVerificationEmailResponse, 
  VerifyEmailResponse
} from '@/types'

export const useEmailVerification = () => {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const sendVerificationEmail = async (email: string): Promise<SendVerificationEmailResponse | null> => {
    setLoading(true)
    setError(null)
    
    try {
      const response = await userService.sendVerificationEmail(email)
      return response
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'เกิดข้อผิดพลาดในการส่งอีเมลยืนยัน'
      setError(errorMessage)
      return null
    } finally {
      setLoading(false)
    }
  }

  const verifyEmail = async (token: string): Promise<VerifyEmailResponse | null> => {
    setLoading(true)
    setError(null)
    
    try {
      const response = await userService.verifyEmail(token)
      return response
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'เกิดข้อผิดพลาดในการยืนยันอีเมล'
      setError(errorMessage)
      return null
    } finally {
      setLoading(false)
    }
  }

  const clearError = () => {
    setError(null)
  }

  return {
    loading,
    error,
    sendVerificationEmail,
    verifyEmail,
    clearError
  }
}

export default useEmailVerification
