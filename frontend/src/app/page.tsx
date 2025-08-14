'use client'

import { useState } from 'react'
import Link from 'next/link'
import { 
  CalendarDaysIcon, 
  UserGroupIcon, 
  ClockIcon, 
  ShieldCheckIcon,
  CheckCircleIcon,
  SparklesIcon
} from '@heroicons/react/24/outline'
import { Button } from '@/components/ui/Button'
import { ButtonGroup } from '@/components/ui/ButtonGroup'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card'

const features = [
  {
    icon: CalendarDaysIcon,
    title: 'จัดตารางเวรอัตโนมัติ',
    description: 'ระบบ AI ช่วยจัดตารางเวรที่เหมาะสมตามความต้องการและกฎระเบียบของแผนก'
  },
  {
    icon: UserGroupIcon,
    title: 'จัดการแผนกและพนักงาน',
    description: 'จัดการข้อมูลแผนก พนักงาน และการกำหนดสิทธิ์การเข้าถึงอย่างมีระบบ'
  },
  {
    icon: ClockIcon,
    title: 'ตั้งค่าเวรและเวลาทำงาน',
    description: 'กำหนดเวลาทำงาน วันหยุด และความสำคัญในการจัดเวรได้อย่างยืดหยุ่น'
  },
  {
    icon: ShieldCheckIcon,
    title: 'ระบบความปลอดภัยสูง',
    description: 'ระบบรักษาความปลอดภัยข้อมูลด้วย JWT และการเข้ารหัสขั้นสูง'
  }
]

const plans = [
  {
    name: 'แพ็คเกจทดลองใช้',
    price: 'ฟรี',
    duration: '90 วัน',
    features: [
      'จัดการ 1 แผนก',
      'พนักงานไม่จำกัด',
      'ตารางเวรพื้นฐาน',
      'การแจ้งเตือนพื้นฐาน'
    ]
  },
  {
    name: 'แพ็คเกจมาตรฐาน',
    price: '990',
    duration: '1 เดือน',
    features: [
      'จัดการหลายแผนก',
      'พนักงานไม่จำกัด',
      'ตารางเวรอัตโนมัติ',
      'การแจ้งเตือนแบบเรียลไทม์',
      'รายงานและสถิติ'
    ]
  },
  {
    name: 'แพ็คเกจระดับองค์กร',
    price: '2,990',
    duration: '3 เดือน',
    features: [
      'จัดการหลายแผนกไม่จำกัด',
      'พนักงานไม่จำกัด',
      'ตารางเวรอัตโนมัติด้วย AI',
      'การแจ้งเตือนแบบเรียลไทม์',
      'รายงานและสถิติขั้นสูง',
      'การสำรองข้อมูล',
      'การสนับสนุนลูกค้าแบบพิเศษ'
    ]
  }
]

export default function LandingPage() {
  const [activeTab, setActiveTab] = useState('features')

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-6">
            <div className="flex items-center space-x-2">
              <div className="w-10 h-10 bg-gradient-to-r from-blue-600 to-purple-600 rounded-lg flex items-center justify-center">
                <CalendarDaysIcon className="w-6 h-6 text-white" />
              </div>
              <h1 className="text-2xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
                NurseShift
              </h1>
            </div>
            <div className="flex items-center space-x-4">
              <ButtonGroup direction="horizontal" spacing="normal">
                <Link href="/auth/login">
                  <Button variant="outline">เข้าสู่ระบบ</Button>
                </Link>
                <Link href="/auth/register">
                  <Button>สมัครสมาชิก</Button>
                </Link>
              </ButtonGroup>
            </div>
          </div>
        </div>
      </header>

      {/* Hero Section */}
      <section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
        <div className="text-center">
          <h1 className="text-5xl md:text-6xl font-bold text-gray-900 mb-6">
            ระบบจัดตารางเวร
            <span className="block bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
              สำหรับพยาบาล
            </span>
          </h1>
          <p className="text-xl text-gray-600 mb-8 max-w-3xl mx-auto leading-relaxed">
            จัดการตารางเวรพยาบาลอย่างมีประสิทธิภาพด้วยระบบอัตโนมัติที่ทันสมัย 
            ลดภาระงานและเพิ่มความแม่นยำในการจัดเวร
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <ButtonGroup direction="horizontal" spacing="normal">
              <Link href="/auth/register">
                <Button size="lg" className="bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700">
                  <SparklesIcon className="w-5 h-5 mr-2" />
                  เริ่มทดลองใช้ฟรี 90 วัน
                </Button>
              </Link>
              <Button variant="outline" size="lg">
                ดูคุณสมบัติ
              </Button>
            </ButtonGroup>
          </div>
        </div>
      </section>

      {/* Tabs Navigation */}
      <section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 mb-16">
        <div className="flex justify-center mb-12">
          <div className="bg-gray-100 p-1 rounded-lg">
            <button
              onClick={() => setActiveTab('features')}
              className={`px-6 py-3 rounded-md font-medium transition-colors ${
                activeTab === 'features'
                  ? 'bg-white text-blue-600 shadow-sm'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              คุณสมบัติ
            </button>
            <button
              onClick={() => setActiveTab('pricing')}
              className={`px-6 py-3 rounded-md font-medium transition-colors ${
                activeTab === 'pricing'
                  ? 'bg-white text-blue-600 shadow-sm'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              แพ็คเกจ
            </button>
          </div>
        </div>

        {/* Features Tab */}
        {activeTab === 'features' && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
            {features.map((feature, index) => (
              <Card key={index} className="hover:shadow-lg transition-shadow border-0 shadow-md">
                <CardHeader className="text-center">
                  <div className="w-12 h-12 bg-gradient-to-r from-blue-100 to-purple-100 rounded-lg flex items-center justify-center mx-auto mb-4">
                    <feature.icon className="w-6 h-6 text-blue-600" />
                  </div>
                  <CardTitle className="text-lg">{feature.title}</CardTitle>
                </CardHeader>
                <CardContent>
                  <CardDescription className="text-center text-gray-600">
                    {feature.description}
                  </CardDescription>
                </CardContent>
              </Card>
            ))}
          </div>
        )}

        {/* Pricing Tab */}
        {activeTab === 'pricing' && (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {plans.map((plan, index) => (
              <Card key={index} className={`hover:shadow-lg transition-shadow border-0 shadow-md ${
                index === 1 ? 'ring-2 ring-blue-500 scale-105' : ''
              }`}>
                <CardHeader className="text-center">
                  {index === 1 && (
                    <div className="bg-blue-500 text-white text-sm font-medium px-3 py-1 rounded-full mx-auto mb-4 w-fit">
                      แนะนำ
                    </div>
                  )}
                  <CardTitle className="text-xl">{plan.name}</CardTitle>
                  <div className="text-3xl font-bold text-blue-600">
                    {plan.price === 'ฟรี' ? 'ฟรี' : `฿${plan.price}`}
                  </div>
                  <CardDescription>/ {plan.duration}</CardDescription>
                </CardHeader>
                <CardContent>
                  <ul className="space-y-3">
                    {plan.features.map((feature, featureIndex) => (
                      <li key={featureIndex} className="flex items-center">
                        <CheckCircleIcon className="w-5 h-5 text-green-500 mr-3 flex-shrink-0" />
                        <span className="text-gray-600">{feature}</span>
                      </li>
                    ))}
                  </ul>
                  <div className="mt-6">
                    <Link href="/auth/register">
                      <Button 
                        className={`w-full ${
                          index === 1 
                            ? 'bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700' 
                            : ''
                        }`}
                        variant={index === 1 ? 'default' : 'outline'}
                      >
                        {plan.price === 'ฟรี' ? 'เริ่มทดลองใช้' : 'เลือกแพ็คเกจ'}
                      </Button>
                    </Link>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </section>

      {/* CTA Section */}
      <section className="bg-gradient-to-r from-blue-600 to-purple-600 py-16">
        <div className="max-w-4xl mx-auto text-center px-4 sm:px-6 lg:px-8">
          <h2 className="text-3xl md:text-4xl font-bold text-white mb-6">
            พร้อมที่จะเริ่มต้นใช้งานแล้วหรือยัง?
          </h2>
          <p className="text-xl text-blue-100 mb-8">
            เริ่มต้นด้วยการทดลองใช้ฟรี 90 วัน ไม่ต้องใช้บัตรเครดิต
          </p>
          <Link href="/auth/register">
            <Button size="lg" variant="outline" className="bg-white text-blue-600 hover:bg-gray-50">
              สมัครสมาชิกฟรี
            </Button>
          </Link>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-gray-900 text-white py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
            <div className="col-span-2">
              <div className="flex items-center space-x-2 mb-4">
                <div className="w-8 h-8 bg-gradient-to-r from-blue-600 to-purple-600 rounded-lg flex items-center justify-center">
                  <CalendarDaysIcon className="w-5 h-5 text-white" />
                </div>
                <h3 className="text-xl font-bold">NurseShift</h3>
              </div>
              <p className="text-gray-400 mb-4">
                ระบบจัดการตารางเวรพยาบาลที่ทันสมัยและมีประสิทธิภาพ 
                ช่วยให้การจัดการเวรเป็นเรื่องง่าย
              </p>
            </div>
            <div>
              <h4 className="font-semibold mb-4">ผลิตภัณฑ์</h4>
              <ul className="space-y-2 text-gray-400">
                <li>จัดตารางเวร</li>
                <li>จัดการแผนก</li>
                <li>รายงานสถิติ</li>
                <li>การแจ้งเตือน</li>
              </ul>
            </div>
            <div>
              <h4 className="font-semibold mb-4">การสนับสนุน</h4>
              <ul className="space-y-2 text-gray-400">
                <li>คู่มือการใช้งาน</li>
                <li>ติดต่อสนับสนุน</li>
                <li>FAQ</li>
                <li>ความปลอดภัย</li>
              </ul>
            </div>
          </div>
          <div className="border-t border-gray-800 mt-8 pt-8 text-center text-gray-400">
            <p>&copy; 2024 NurseShift. สงวนลิขสิทธิ์ทุกประการ</p>
          </div>
        </div>
      </footer>
    </div>
  )
}