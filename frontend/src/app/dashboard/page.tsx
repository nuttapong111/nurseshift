'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'

export default function DashboardPage() {
  const router = useRouter()

  useEffect(() => {
    // Redirect to schedule page
    router.replace('/dashboard/schedule')
  }, [router])

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-50">
      <p className="text-gray-600">กำลังเปลี่ยนเส้นทาง...</p>
    </div>
  )
}
