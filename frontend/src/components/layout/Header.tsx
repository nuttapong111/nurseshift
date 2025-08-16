'use client'

import { useState, useEffect } from 'react'
import { useRouter, usePathname } from 'next/navigation'
import Link from 'next/link'
import { Menu, MenuButton, MenuItem, MenuItems, Transition } from '@headlessui/react'
import { 
  CalendarDaysIcon,
  BellIcon,
  UserCircleIcon,
  ChevronDownIcon,
  ClockIcon,
  UserIcon,
  ArrowRightOnRectangleIcon,
  UserGroupIcon,
  CogIcon,
  CalendarIcon,
  StarIcon,
  CreditCardIcon
} from '@heroicons/react/24/outline'
import { Button } from '@/components/ui/Button'
import { cn, getTimeRemaining } from '@/lib/utils'
import Swal from 'sweetalert2'
import { normalizeBaseUrl } from '@/lib/utils'
import type { User } from '@/types'

const navigation = [
  {
    name: 'ตารางเวร',
    href: '/dashboard/schedule',
    icon: CalendarDaysIcon
  },
  {
    name: 'จัดการแผนกและพนักงาน',
    href: '/dashboard/departments',
    icon: UserGroupIcon
  },
  {
    name: 'ตั้งค่าแผนก',
    href: '/dashboard/department-settings',
    icon: CogIcon
  },
  {
    name: 'จัดการวันหยุดพนักงาน',
    href: '/dashboard/employee-leaves',
    icon: CalendarIcon
  },
  {
    name: 'จัดการความสำคัญ',
    href: '/dashboard/priorities',
    icon: StarIcon
  },
  {
    name: 'แพ็คเกจสมาชิก',
    href: '/dashboard/packages',
    icon: CreditCardIcon
  }
]

interface HeaderProps {
  user?: User | null
}

export default function Header({ user }: HeaderProps) {
  const router = useRouter()
  const pathname = usePathname()
  const [notifications] = useState([
    { id: 1, message: 'ตารางเวรประจำเดือนมีนาคมได้รับการสร้างแล้ว', time: '5 นาทีที่แล้ว', read: false },
    { id: 2, message: 'พนักงานใหม่ถูกเพิ่มในแผนกอายุรกรรม', time: '1 ชั่วโมงที่แล้ว', read: true },
    { id: 3, message: 'การตั้งค่าเวรได้รับการอัปเดต', time: '3 ชั่วโมงที่แล้ว', read: true }
  ])
  const [unreadCount, setUnreadCount] = useState(0)

  useEffect(() => {
    const count = notifications.filter(n => !n.read).length
    setUnreadCount(count)
  }, [notifications])

  const handleLogout = async () => {
    const result = await Swal.fire({
      title: 'ออกจากระบบ?',
      text: 'คุณต้องการออกจากระบบหรือไม่?',
      icon: 'question',
      showCancelButton: true,
      confirmButtonColor: '#dc2626',
      cancelButtonColor: '#6b7280',
      confirmButtonText: 'ออกจากระบบ',
      cancelButtonText: 'ยกเลิก'
    })

    if (result.isConfirmed) {
      try {
        // Call logout API
        const token = localStorage.getItem('token')
        if (token) {
          await fetch(`${normalizeBaseUrl(process.env.NEXT_PUBLIC_AUTH_SERVICE_URL)}/api/v1/auth/logout`, {
            method: 'POST',
            headers: {
              'Authorization': `Bearer ${token}`,
              'Content-Type': 'application/json',
            },
          })
        }
      } catch (error) {
        console.error('Logout API error:', error)
      }

      // Clear local storage
      localStorage.removeItem('token')
      localStorage.removeItem('refreshToken')
      localStorage.removeItem('user')
      
      await Swal.fire({
        icon: 'success',
        title: 'ออกจากระบบสำเร็จ',
        text: 'ขอบคุณที่ใช้งานระบบ NurseShift',
        confirmButtonText: 'ตกลง',
        confirmButtonColor: '#2563eb'
      })

      router.push('/auth/login')
    }
  }

  const remainingTime = user?.remainingDays ? getTimeRemaining(user.remainingDays) : null

  return (
    <header className="bg-white shadow-sm border-b border-gray-200">
      <div className="px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <div className="flex items-center space-x-8">
            <div className="flex items-center space-x-2">
              <div className="w-8 h-8 bg-gradient-to-r from-blue-600 to-purple-600 rounded-lg flex items-center justify-center">
                <CalendarDaysIcon className="w-5 h-5 text-white" />
              </div>
              <h1 className="text-xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
                NurseShift
              </h1>
            </div>

            {/* Navigation Menu */}
            <nav className="hidden md:flex items-center space-x-1">
              {navigation.map((item) => {
                const isActive = pathname === item.href
                return (
                  <Link
                    key={item.name}
                    href={item.href}
                    className={cn(
                      'flex items-center px-3 py-2 rounded-md text-sm font-medium transition-colors',
                      isActive
                        ? 'bg-blue-100 text-blue-700'
                        : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
                    )}
                  >
                    <item.icon className={cn(
                      'w-4 h-4 mr-2',
                      isActive ? 'text-blue-600' : 'text-gray-400'
                    )} />
                    {item.name}
                  </Link>
                )
              })}
            </nav>
          </div>

          {/* Right side */}
          <div className="flex items-center space-x-4">
            {/* Mobile menu button */}
            <Menu as="div" className="relative md:hidden">
              <MenuButton className="p-2 text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 rounded-md">
                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                </svg>
              </MenuButton>
              <Transition
                enter="transition ease-out duration-100"
                enterFrom="transform opacity-0 scale-95"
                enterTo="transform opacity-100 scale-100"
                leave="transition ease-in duration-75"
                leaveFrom="transform opacity-100 scale-100"
                leaveTo="transform opacity-0 scale-95"
              >
                <MenuItems className="absolute right-0 z-10 mt-2 w-56 origin-top-right rounded-md bg-white py-1 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                  {navigation.map((item) => (
                    <MenuItem key={item.name}>
                      {({ focus }) => (
                        <Link
                          href={item.href}
                          className={cn(
                            focus ? 'bg-gray-50' : '',
                            'flex items-center px-4 py-2 text-sm text-gray-700'
                          )}
                        >
                          <item.icon className="w-4 h-4 mr-3 text-gray-400" />
                          {item.name}
                        </Link>
                      )}
                    </MenuItem>
                  ))}
                </MenuItems>
              </Transition>
            </Menu>
            {/* Notifications */}
            <Menu as="div" className="relative">
              <MenuButton className="relative p-2 text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 rounded-full">
                <BellIcon className="w-6 h-6" />
                {unreadCount > 0 && (
                  <span className="absolute -top-1 -right-1 h-5 w-5 bg-red-500 text-white text-xs rounded-full flex items-center justify-center">
                    {unreadCount}
                  </span>
                )}
              </MenuButton>
              <Transition
                enter="transition ease-out duration-100"
                enterFrom="transform opacity-0 scale-95"
                enterTo="transform opacity-100 scale-100"
                leave="transition ease-in duration-75"
                leaveFrom="transform opacity-100 scale-100"
                leaveTo="transform opacity-0 scale-95"
              >
                <MenuItems className="absolute right-0 z-10 mt-2 w-80 origin-top-right rounded-md bg-white py-1 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                  <div className="px-4 py-3 border-b border-gray-200">
                    <h3 className="text-sm font-medium text-gray-900">การแจ้งเตือน</h3>
                  </div>
                  {notifications.map((notification) => (
                    <MenuItem key={notification.id}>
                      {({ focus }) => (
                        <div
                          className={cn(
                            focus ? 'bg-gray-50' : '',
                            'px-4 py-3 border-b border-gray-100 last:border-b-0'
                          )}
                        >
                          <div className="flex items-start space-x-3">
                            <div className={cn(
                              'w-2 h-2 rounded-full mt-2 flex-shrink-0',
                              notification.read ? 'bg-gray-300' : 'bg-blue-500'
                            )} />
                            <div className="flex-1 min-w-0">
                              <p className="text-sm text-gray-900">{notification.message}</p>
                              <p className="text-xs text-gray-500 mt-1">{notification.time}</p>
                            </div>
                          </div>
                        </div>
                      )}
                    </MenuItem>
                  ))}
                  <div className="px-4 py-2 border-t border-gray-200">
                    <Button variant="ghost" size="sm" className="w-full text-blue-600" onClick={() => router.push('/dashboard/notifications')}>
                      ดูการแจ้งเตือนทั้งหมด
                    </Button>
                  </div>
                </MenuItems>
              </Transition>
            </Menu>

            {/* User menu */}
            <Menu as="div" className="relative">
              <MenuButton className="flex items-center space-x-3 p-2 rounded-lg hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2">
                <UserCircleIcon className="w-8 h-8 text-gray-400" />
                <div className="hidden sm:block text-left">
                  <p className="text-sm font-medium text-gray-900">
                    {user ? `${user.firstName} ${user.lastName}` : 'ผู้ใช้งาน'}
                  </p>
                  {user?.position && (
                    <p className="text-xs text-gray-500">{user.position}</p>
                  )}
                </div>
                <ChevronDownIcon className="w-4 h-4 text-gray-400" />
              </MenuButton>
              <Transition
                enter="transition ease-out duration-100"
                enterFrom="transform opacity-0 scale-95"
                enterTo="transform opacity-100 scale-100"
                leave="transition ease-in duration-75"
                leaveFrom="transform opacity-100 scale-100"
                leaveTo="transform opacity-0 scale-95"
              >
                <MenuItems className="absolute right-0 z-10 mt-2 w-56 origin-top-right rounded-md bg-white py-1 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                  <MenuItem>
                    {({ focus }) => (
                      <button
                        className={cn(
                          focus ? 'bg-gray-50' : '',
                          'flex items-center px-4 py-2 text-sm text-gray-700 w-full text-left'
                        )}
                        onClick={() => router.push('/dashboard/profile')}
                      >
                        <UserIcon className="w-4 h-4 mr-3" />
                        โปรไฟล์ส่วนตัว
                      </button>
                    )}
                  </MenuItem>
                  
                  {remainingTime && (
                    <MenuItem>
                      {({ focus }) => (
                        <div
                          className={cn(
                            focus ? 'bg-gray-50' : '',
                            'flex items-center px-4 py-2 text-sm text-gray-700'
                          )}
                        >
                          <ClockIcon className="w-4 h-4 mr-3" />
                          <div>
                            <div className="text-sm">เวลาการใช้งาน</div>
                            <div className={cn(
                              'text-xs',
                              user?.remainingDays && user.remainingDays <= 7 
                                ? 'text-red-600' 
                                : 'text-green-600'
                            )}>
                              {remainingTime}
                            </div>
                          </div>
                        </div>
                      )}
                    </MenuItem>
                  )}

                  <div className="border-t border-gray-100 my-1" />
                  
                  <MenuItem>
                    {({ focus }) => (
                      <button
                        onClick={handleLogout}
                        className={cn(
                          focus ? 'bg-gray-50' : '',
                          'flex items-center px-4 py-2 text-sm text-red-700 w-full text-left'
                        )}
                      >
                        <ArrowRightOnRectangleIcon className="w-4 h-4 mr-3" />
                        ออกจากระบบ
                      </button>
                    )}
                  </MenuItem>
                </MenuItems>
              </Transition>
            </Menu>
          </div>
        </div>
      </div>
    </header>
  )
}
