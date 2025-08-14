'use client'

import { useState } from 'react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import {
  CalendarDaysIcon,
  UserGroupIcon,
  CogIcon,
  CalendarIcon,
  StarIcon,
  CreditCardIcon,
  Bars3Icon,
  XMarkIcon
} from '@heroicons/react/24/outline'
import { cn } from '@/lib/utils'

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

interface SidebarProps {
  className?: string
}

export default function Sidebar({ className }: SidebarProps) {
  const pathname = usePathname()
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)

  return (
    <>
      {/* Mobile menu button */}
      <div className="lg:hidden fixed top-4 left-4 z-50">
        <button
          onClick={() => setIsMobileMenuOpen(true)}
          className="p-2 rounded-lg bg-blue-500 shadow-md text-white hover:bg-blue-600"
        >
          <Bars3Icon className="w-5 h-5" />
        </button>
      </div>

      {/* Mobile menu overlay */}
      {isMobileMenuOpen && (
        <div className="lg:hidden fixed inset-0 z-40 bg-black bg-opacity-50" onClick={() => setIsMobileMenuOpen(false)} />
      )}

      {/* Sidebar */}
      <div className={cn(
        'fixed inset-y-0 left-0 z-50 w-40 lg:w-40 bg-white shadow-lg border-r border-gray-200 transform transition-transform duration-300 ease-in-out lg:translate-x-0 lg:static lg:inset-0 h-full min-h-screen',
        isMobileMenuOpen ? 'translate-x-0 w-64' : '-translate-x-full',
        className
      )}>
        {/* Mobile close button */}
        <div className="lg:hidden absolute top-4 right-4">
          <button
            onClick={() => setIsMobileMenuOpen(false)}
            className="p-2 rounded-lg text-gray-400 hover:text-gray-600"
          >
            <XMarkIcon className="w-6 h-6" />
          </button>
        </div>

        {/* Sidebar content */}
        <div className="flex flex-col h-full pt-20 lg:pt-6 min-h-screen">
          <nav className="flex-1 px-2 py-6 space-y-2">
            {navigation.map((item, index) => {
              const isActive = pathname === item.href
              return (
                <div key={item.name} className="relative group">
                  <Link
                    href={item.href}
                    onClick={() => setIsMobileMenuOpen(false)}
                    className={cn(
                      'flex items-center lg:flex-col lg:justify-center p-3 rounded-xl text-xs font-medium transition-all duration-200',
                      isActive
                        ? 'bg-blue-500 text-white shadow-lg transform scale-105'
                        : 'text-gray-500 hover:bg-gray-100 hover:text-gray-700',
                      isMobileMenuOpen ? 'flex-row space-x-3' : 'flex-col'
                    )}
                  >
                    <item.icon
                      className={cn(
                        'w-6 h-6',
                        isActive ? 'text-white' : 'text-gray-400',
                        isMobileMenuOpen ? 'mb-0' : 'lg:mb-2'
                      )}
                    />
                    <span className={cn(
                      'text-center leading-tight',
                      isMobileMenuOpen ? 'text-sm' : 'text-xs hidden lg:block'
                    )}>
                      {item.name}
                    </span>
                  </Link>
                  
                  {/* Tooltip for desktop only */}
                  <div className="hidden lg:block absolute left-full ml-2 px-3 py-2 bg-gray-900 text-white text-sm rounded-lg opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 whitespace-nowrap z-50 top-1/2 transform -translate-y-1/2">
                    {item.name}
                    <div className="absolute right-full top-1/2 transform -translate-y-1/2 w-0 h-0 border-t-4 border-b-4 border-r-4 border-transparent border-r-gray-900"></div>
                  </div>
                </div>
              )
            })}
          </nav>

          {/* Footer */}
          <div className="p-4 border-t border-gray-200 bg-gray-50">
            <div className="text-xs text-gray-500 text-center">
              <p>NurseShift v1.0</p>
              <p className="mt-1">ระบบจัดตารางเวรพยาบาล</p>
            </div>
          </div>
        </div>
      </div>
    </>
  )
}
