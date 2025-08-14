'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Header from './Header'
import type { User } from '@/types'
import userService from '@/services/userService'
import Swal from 'sweetalert2'

interface DashboardLayoutProps {
  children: React.ReactNode
}

export default function DashboardLayout({ children }: DashboardLayoutProps) {
  const router = useRouter()
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const initializeUser = async () => {
      try {
        // Check authentication
        const token = localStorage.getItem('token')
        const userData = localStorage.getItem('user')

        if (!token || !userData) {
          router.push('/auth/login')
          return
        }

        try {
          const parsedUser = JSON.parse(userData) as User
          setUser(parsedUser)
          
          // Fetch fresh user data from API
          const freshUserData = await userService.getProfile()
          setUser(freshUserData)
          
          // Update localStorage with fresh data
          localStorage.setItem('user', JSON.stringify(freshUserData))
        } catch (error) {
          console.error('Error fetching user data:', error)
          
          // Show error message
          await Swal.fire({
            icon: 'error',
            title: 'เกิดข้อผิดพลาด',
            text: 'ไม่สามารถดึงข้อมูลผู้ใช้ได้ กรุณาเข้าสู่ระบบใหม่',
            confirmButtonColor: '#dc2626'
          })
          
          // Clear localStorage and redirect to login
          localStorage.removeItem('token')
          localStorage.removeItem('user')
          router.push('/auth/login')
          return
        }
      } catch (error) {
        console.error('Error initializing user:', error)
        
        // Show error message
        await Swal.fire({
          icon: 'error',
          title: 'เกิดข้อผิดพลาด',
          text: 'ไม่สามารถเชื่อมต่อกับระบบได้ กรุณาลองใหม่อีกครั้ง',
          confirmButtonColor: '#dc2626'
        })
        
        // Clear localStorage and redirect to login
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        router.push('/auth/login')
        return
      } finally {
        setIsLoading(false)
      }
    }

    initializeUser()
  }, [router])

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">กำลังโหลด...</p>
        </div>
      </div>
    )
  }

  if (!user) {
    return null
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header with navigation */}
      <Header user={user} />
      
      {/* Main content */}
      <main className="flex-1">
        <div className="px-4 sm:px-6 lg:px-8 py-8">
          {children}
        </div>
      </main>
    </div>
  )
}
