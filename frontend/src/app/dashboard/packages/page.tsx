'use client'

import { useState } from 'react'
import { 
  CreditCardIcon,
  CheckCircleIcon,
  ClockIcon,
  DocumentIcon,
  EyeIcon,
  StarIcon,
  PencilIcon
} from '@heroicons/react/24/outline'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card'
import { Button } from '@/components/ui/Button'
import { ButtonGroup } from '@/components/ui/ButtonGroup'
import { Switch } from '@/components/ui/Switch'
import DashboardLayout from '@/components/layout/DashboardLayout'
import Swal from 'sweetalert2'

// Mock data
const mockPackages = [
  {
    id: 1,
    name: 'แพ็คเกจมาตรฐาน',
    price: 990,
    duration: 30,
    description: 'เหมาะสำหรับแผนกขนาดกลาง',
    features: [
      'จัดการหลายแผนก',
      'พนักงานไม่จำกัด',
      'ตารางเวรอัตโนมัติ',
      'การแจ้งเตือนแบบเรียลไทม์',
      'รายงานและสถิติ'
    ],
    isPopular: true,
    isActive: true
  },
  {
    id: 2,
    name: 'แพ็คเกจระดับองค์กร',
    price: 2990,
    duration: 90,
    description: 'เหมาะสำหรับองค์กรขนาดใหญ่',
    features: [
      'จัดการหลายแผนกไม่จำกัด',
      'พนักงานไม่จำกัด',
      'ตารางเวรอัตโนมัติด้วย AI',
      'การแจ้งเตือนแบบเรียลไทม์',
      'รายงานและสถิติขั้นสูง',
      'การสำรองข้อมูล',
      'การสนับสนุนลูกค้าแบบพิเศษ'
    ],
    isPopular: false,
    isActive: true
  }
]

const mockPayments = [
  {
    id: 1,
    packageName: 'แพ็คเกจมาตรฐาน',
    amount: 990,
    paymentDate: '2024-03-01',
    status: 'approved' as const,
    evidence: 'payment_evidence_1.jpg',
    approvedDate: '2024-03-02',
    extendedDays: 30
  },
  {
    id: 2,
    packageName: 'แพ็คเกจระดับองค์กร',
    amount: 2990,
    paymentDate: '2024-03-15',
    status: 'pending' as const,
    evidence: 'payment_evidence_2.jpg',
    approvedDate: null,
    extendedDays: null
  },
  {
    id: 3,
    packageName: 'แพ็คเกจมาตรฐาน',
    amount: 990,
    paymentDate: '2024-02-15',
    status: 'rejected' as const,
    evidence: 'payment_evidence_3.jpg',
    approvedDate: null,
    extendedDays: null,
    rejectReason: 'หลักฐานการโอนเงินไม่ชัดเจน'
  }
]

export default function PackagesPage() {
  const [packages] = useState(mockPackages)
  type PaymentStatus = 'approved' | 'pending' | 'rejected'

  interface Payment {
    id: number
    packageName: string
    amount: number
    paymentDate: string
    status: PaymentStatus
    evidence: string
    approvedDate: string | null
    extendedDays: number | null
    rejectReason?: string
  }

  const [payments, setPayments] = useState<Payment[]>(mockPayments)
  const [showPaymentModal, setShowPaymentModal] = useState(false)
  const [showResubmitModal, setShowResubmitModal] = useState(false)
  const [selectedPackage, setSelectedPackage] = useState<typeof mockPackages[0] | null>(null)
  const [selectedPayment, setSelectedPayment] = useState<Payment | null>(null)
  const [evidence, setEvidence] = useState<File | null>(null)

  const handleSelectPackage = (pkg: typeof mockPackages[0]) => {
    setSelectedPackage(pkg)
    setShowPaymentModal(true)
  }

  const handlePaymentSubmit = async () => {
    if (!selectedPackage || !evidence) {
      await Swal.fire({
        icon: 'warning',
        title: 'กรุณาแนบหลักฐานการโอนเงิน',
        confirmButtonColor: '#2563eb'
      })
      return
    }

    const newPayment = {
      id: payments.length + 1,
      packageName: selectedPackage.name,
      amount: selectedPackage.price,
      paymentDate: new Date().toISOString().split('T')[0],
      status: 'pending' as const,
      evidence: evidence.name,
      approvedDate: null,
      extendedDays: null
    }

    setPayments([newPayment, ...payments])
    setShowPaymentModal(false)
    setSelectedPackage(null)
    setEvidence(null)

    await Swal.fire({
      icon: 'success',
      title: 'ส่งข้อมูลการชำระเงินสำเร็จ!',
      text: 'รอการอนุมัติจากผู้ดูแลระบบ',
      confirmButtonColor: '#2563eb'
    })
  }

  const handleResubmitPayment = (payment: Payment) => {
    setSelectedPayment(payment)
    setShowResubmitModal(true)
  }

  const handleResubmitSubmit = async () => {
    if (!selectedPayment || !evidence) {
      await Swal.fire({
        icon: 'warning',
        title: 'กรุณาแนบหลักฐานการโอนเงินใหม่',
        confirmButtonColor: '#2563eb'
      })
      return
    }

    // Update payment status to pending and new evidence
    const updatedPayments = payments.map(payment => 
      payment.id === selectedPayment.id 
        ? {
            ...payment,
            status: 'pending' as const,
            evidence: evidence.name,
            paymentDate: new Date().toISOString().split('T')[0], // Update submission date
            rejectReason: undefined // Clear reject reason
          }
        : payment
    )

    setPayments(updatedPayments)
    setShowResubmitModal(false)
    setSelectedPayment(null)
    setEvidence(null)

    await Swal.fire({
      icon: 'success',
      title: 'ส่งหลักฐานใหม่เรียบร้อย!',
      text: 'สถานะเปลี่ยนเป็น "รอการอนุมัติ" แล้ว',
      confirmButtonColor: '#2563eb'
    })
  }

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'approved':
        return 'bg-green-100 text-green-800'
      case 'pending':
        return 'bg-yellow-100 text-yellow-800'
      case 'rejected':
        return 'bg-red-100 text-red-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const getStatusText = (status: string) => {
    switch (status) {
      case 'approved':
        return 'อนุมัติแล้ว'
      case 'pending':
        return 'รอการอนุมัติ'
      case 'rejected':
        return 'ไม่อนุมัติ'
      default:
        return 'ไม่ทราบสถานะ'
    }
  }

  return (
    <DashboardLayout>
      <div className="space-y-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">แพ็คเกจสมาชิก</h1>
          <p className="text-gray-600 mt-2">เลือกแพ็คเกจที่เหมาะสมและจัดการการชำระเงิน</p>
        </div>

        {/* Current Subscription Status */}
        <Card className="border-0 shadow-md bg-green-50 border-green-200">
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <CheckCircleIcon className="w-8 h-8 text-green-600" />
                <div>
                  <h3 className="font-medium text-green-900">สถานะการใช้งานปัจจุบัน</h3>
                  <p className="text-green-700">แพ็คเกจทดลองใช้ฟรี - เหลืออีก 85 วัน</p>
                </div>
              </div>
              <Button className="bg-green-600 hover:bg-green-700">
                ต่ออายุเลย
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* Available Packages */}
        <div>
          <h2 className="text-2xl font-bold text-gray-900 mb-6">แพ็คเกจที่ใช้ได้</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {packages.map((pkg) => (
              <Card 
                key={pkg.id} 
                className={`border-0 shadow-md hover:shadow-lg transition-shadow relative ${
                  pkg.isPopular ? 'ring-2 ring-blue-500' : ''
                }`}
              >
                {pkg.isPopular && (
                  <div className="absolute -top-3 left-1/2 transform -translate-x-1/2">
                    <span className="bg-blue-500 text-white text-sm px-4 py-1 rounded-full flex items-center">
                      <StarIcon className="w-4 h-4 mr-1" />
                      แนะนำ
                    </span>
                  </div>
                )}
                <CardHeader className="text-center pt-8">
                  <CardTitle className="text-xl">{pkg.name}</CardTitle>
                  <div className="text-3xl font-bold text-blue-600">
                    ฿{pkg.price.toLocaleString()}
                  </div>
                  <CardDescription>/ {pkg.duration} วัน</CardDescription>
                  <p className="text-gray-600 mt-2">{pkg.description}</p>
                </CardHeader>
                <CardContent>
                  <ul className="space-y-3 mb-6">
                    {pkg.features.map((feature, index) => (
                      <li key={index} className="flex items-center">
                        <CheckCircleIcon className="w-5 h-5 text-green-500 mr-3 flex-shrink-0" />
                        <span className="text-gray-600">{feature}</span>
                      </li>
                    ))}
                  </ul>
                  <ButtonGroup direction="horizontal" spacing="normal" className="w-full">
                    <Button 
                      className={`w-full ${
                        pkg.isPopular 
                          ? 'bg-blue-600 hover:bg-blue-700' 
                          : ''
                      }`}
                      variant={pkg.isPopular ? 'default' : 'outline'}
                      onClick={() => handleSelectPackage(pkg)}
                    >
                      เลือกแพ็คเกจนี้
                    </Button>
                  </ButtonGroup>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>

        {/* Payment History */}
        <Card className="border-0 shadow-md">
          <CardHeader>
            <CardTitle className="flex items-center">
              <DocumentIcon className="w-5 h-5 mr-2" />
              ประวัติการชำระเงิน
            </CardTitle>
            <CardDescription>รายการการชำระเงินและสถานะการอนุมัติ</CardDescription>
          </CardHeader>
          <CardContent>
            {payments.length > 0 ? (
              <div className="space-y-4">
                {payments.map((payment) => (
                  <div key={payment.id} className="p-4 border border-gray-200 rounded-lg">
                    <div className="flex items-center justify-between">
                      <div className="flex-1">
                        <div className="flex items-center space-x-3 mb-2">
                          <h3 className="font-medium text-gray-900">{payment.packageName}</h3>
                          <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusBadge(payment.status)}`}>
                            {getStatusText(payment.status)}
                          </span>
                        </div>
                        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm text-gray-600">
                          <div>
                            <span className="font-medium">จำนวนเงิน:</span> ฿{payment.amount.toLocaleString()}
                          </div>
                          <div>
                            <span className="font-medium">วันที่ชำระ:</span> {new Date(payment.paymentDate).toLocaleDateString('th-TH')}
                          </div>
                          {payment.approvedDate && (
                            <div>
                              <span className="font-medium">วันที่อนุมัติ:</span> {new Date(payment.approvedDate).toLocaleDateString('th-TH')}
                            </div>
                          )}
                        </div>
                        {payment.status === 'approved' && payment.extendedDays && (
                          <div className="mt-2 p-2 bg-green-50 rounded text-sm text-green-700">
                            ✅ ต่ออายุการใช้งานแล้ว {payment.extendedDays} วัน
                          </div>
                        )}
                        {payment.status === 'rejected' && payment.rejectReason && (
                          <div className="mt-2 p-2 bg-red-50 rounded text-sm text-red-700">
                            ❌ เหตุผลที่ไม่อนุมัติ: {payment.rejectReason}
                          </div>
                        )}
                      </div>
                      <div className="flex items-center space-x-2">
                        <Button variant="outline" size="sm">
                          <EyeIcon className="w-4 h-4 mr-1" />
                          ดูหลักฐาน
                        </Button>
                        {payment.status === 'rejected' && (
                          <Button 
                            variant="outline" 
                            size="sm"
                            className="text-blue-600 border-blue-300 hover:bg-blue-50"
                            onClick={() => handleResubmitPayment(payment)}
                          >
                            <PencilIcon className="w-4 h-4 mr-1" />
                            แก้ไข
                          </Button>
                        )}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-8 text-gray-500">
                ยังไม่มีประวัติการชำระเงิน
              </div>
            )}
          </CardContent>
        </Card>

        {/* Payment Modal */}
        {showPaymentModal && selectedPackage && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg p-6 w-full max-w-md">
              <h3 className="text-lg font-medium mb-4">ชำระเงิน - {selectedPackage.name}</h3>
              
              <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-4">
                <div className="text-center">
                  <div className="text-2xl font-bold text-blue-600">฿{selectedPackage.price.toLocaleString()}</div>
                  <div className="text-blue-700">ระยะเวลา {selectedPackage.duration} วัน</div>
                </div>
              </div>

              <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mb-4">
                <h4 className="font-medium text-yellow-900 mb-2">ข้อมูลการโอนเงิน</h4>
                <div className="text-sm text-yellow-800 space-y-1">
                  <div><strong>ธนาคาร:</strong> กสิกรไทย</div>
                  <div><strong>เลขที่บัญชี:</strong> 123-4-56789-0</div>
                  <div><strong>ชื่อบัญชี:</strong> บริษัท NurseShift จำกัด</div>
                </div>
              </div>

              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    แนบหลักฐานการโอนเงิน
                  </label>
                  <input
                    type="file"
                    accept="image/*"
                    onChange={(e) => setEvidence(e.target.files?.[0] || null)}
                    className="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                  <p className="text-xs text-gray-500 mt-1">
                    รองรับไฟล์ภาพ (JPG, PNG) เท่านั้น
                  </p>
                </div>

                <div className="flex justify-end space-x-3">
                  <Button
                    variant="outline"
                    onClick={() => {
                      setShowPaymentModal(false)
                      setSelectedPackage(null)
                      setEvidence(null)
                    }}
                  >
                    ยกเลิก
                  </Button>
                  <Button onClick={handlePaymentSubmit}>
                    ส่งข้อมูลการชำระ
                  </Button>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Resubmit Payment Modal */}
        {showResubmitModal && selectedPayment && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg p-6 w-full max-w-md">
              <h3 className="text-lg font-medium mb-4">ส่งหลักฐานใหม่ - {selectedPayment.packageName}</h3>
              
              <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-4">
                <h4 className="font-medium text-red-900 mb-2">เหตุผลที่ไม่อนุมัติ</h4>
                <p className="text-sm text-red-800">{selectedPayment.rejectReason}</p>
              </div>

              <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-4">
                <div className="text-center">
                  <div className="text-2xl font-bold text-blue-600">฿{selectedPayment.amount.toLocaleString()}</div>
                  <div className="text-blue-700">{selectedPayment.packageName}</div>
                </div>
              </div>

              <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mb-4">
                <h4 className="font-medium text-yellow-900 mb-2">ข้อมูลการโอนเงิน</h4>
                <div className="text-sm text-yellow-800 space-y-1">
                  <div><strong>ธนาคาร:</strong> กสิกรไทย</div>
                  <div><strong>เลขที่บัญชี:</strong> 123-4-56789-0</div>
                  <div><strong>ชื่อบัญชี:</strong> บริษัท NurseShift จำกัด</div>
                </div>
              </div>

              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    แนบหลักฐานการโอนเงินใหม่
                  </label>
                  <input
                    type="file"
                    accept="image/*"
                    onChange={(e) => setEvidence(e.target.files?.[0] || null)}
                    className="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                  <p className="text-xs text-gray-500 mt-1">
                    รองรับไฟล์ภาพ (JPG, PNG) เท่านั้น
                  </p>
                </div>

                <div className="bg-green-50 border border-green-200 rounded-md p-3">
                  <p className="text-sm text-green-800">
                    <strong>หมายเหตุ:</strong> หลังจากส่งหลักฐานใหม่ สถานะจะเปลี่ยนเป็น "รอการอนุมัติ" และรอการตรวจสอบจากผู้ดูแลระบบอีกครั้ง
                  </p>
                </div>

                <div className="flex justify-end space-x-3">
                  <Button
                    variant="outline"
                    onClick={() => {
                      setShowResubmitModal(false)
                      setSelectedPayment(null)
                      setEvidence(null)
                    }}
                  >
                    ยกเลิก
                  </Button>
                  <Button onClick={handleResubmitSubmit}>
                    ส่งหลักฐานใหม่
                  </Button>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </DashboardLayout>
  )
}
