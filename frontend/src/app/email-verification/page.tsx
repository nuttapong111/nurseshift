'use client'

import { useState, useEffect, Suspense } from 'react'
import { useSearchParams, useRouter } from 'next/navigation'
import EmailVerification from '@/components/EmailVerification'
import useEmailVerification from '@/hooks/useEmailVerification'

export default function EmailVerificationPage() {
	return (
		<Suspense fallback={<div className="min-h-screen flex items-center justify-center">กำลังโหลด...</div>}>
			<EmailVerificationContent />
		</Suspense>
	)
}

function EmailVerificationContent() {
	const searchParams = useSearchParams()
	const router = useRouter()
	const token = searchParams.get('token')
	
	const [showSuccessPopup, setShowSuccessPopup] = useState(false)
	const [verificationToken, setVerificationToken] = useState('')
	const { verifyEmail, verifyLoading, verifyError } = useEmailVerification()

	useEffect(() => {
		if (token) {
			setVerificationToken(token)
			// Auto-verify if token is present in URL
			handleAutoVerify(token)
		}
	}, [token])

	const handleAutoVerify = async (token: string) => {
		const response = await verifyEmail(token)
		if (response) {
			setShowSuccessPopup(true)
			// Redirect to login page after 3 seconds
			setTimeout(() => {
				router.push('/auth/login')
			}, 3000)
		}
	}

	const handleVerificationComplete = () => {
		setShowSuccessPopup(true)
		// Redirect to login page after 3 seconds
		setTimeout(() => {
			router.push('/auth/login')
		}, 3000)
	}

	const handleVerifyWithToken = async () => {
		if (!verificationToken.trim()) return
		
		const response = await verifyEmail(verificationToken.trim())
		if (response) {
			setShowSuccessPopup(true)
			// Redirect to login page after 3 seconds
			setTimeout(() => {
				router.push('/auth/login')
			}, 3000)
		}
	}

	const handleGoToLogin = () => {
		router.push('/auth/login')
	}

	return (
		<div className="min-h-screen bg-gray-50 py-12 px-4">
			<div className="max-w-md mx-auto">
				{/* Success Popup */}
				{showSuccessPopup && (
					<div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
						<div className="bg-white rounded-lg p-8 max-w-sm mx-4 text-center">
							<div className="mx-auto flex items-center justify-center h-16 w-16 rounded-full bg-green-100 mb-4">
								<svg className="h-8 w-8 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
									<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
								</svg>
							</div>
							<h3 className="text-lg font-medium text-gray-900 mb-2">ยืนยันอีเมลสำเร็จ!</h3>
							<p className="text-sm text-gray-500 mb-6">
								อีเมลของคุณได้รับการยืนยันเรียบร้อยแล้ว ระบบจะนำคุณไปยังหน้าเข้าสู่ระบบในอีกสักครู่
							</p>
							<button
								onClick={handleGoToLogin}
								className="w-full bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 transition-colors"
							>
								ไปหน้าเข้าสู่ระบบ
							</button>
						</div>
					</div>
				)}

				<div className="text-center mb-8">
					<h1 className="text-3xl font-bold text-gray-900 mb-4">
						ยืนยันอีเมล
					</h1>
					<p className="text-lg text-gray-600">
						กรุณายืนยันอีเมลของคุณเพื่อเข้าสู่ระบบ
					</p>
				</div>

				{/* Email Verification Component */}
				<div className="mb-8">
					<EmailVerification onVerificationComplete={handleVerificationComplete} />
				</div>

				{/* Manual Token Verification (if no auto-token) */}
				{!token && (
					<div className="bg-white rounded-lg shadow-md p-6">
						<h2 className="text-xl font-semibold text-gray-900 mb-4 text-center">
							หรือใส่รหัสยืนยันด้วยตนเอง
						</h2>
						<div className="space-y-4">
							<div>
								<label htmlFor="verificationToken" className="block text-sm font-medium text-gray-700 mb-2">
									รหัสยืนยันจากอีเมล
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
								onClick={handleVerifyWithToken}
								disabled={verifyLoading || !verificationToken.trim()}
								className="w-full bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
							>
								{verifyLoading ? 'กำลังยืนยัน...' : 'ยืนยันอีเมล'}
							</button>

							{verifyError && (
								<div className="p-3 bg-red-100 border border-red-400 text-red-700 rounded">
									{verifyError}
								</div>
							)}
						</div>
					</div>
				)}

				{/* Back to Login Link */}
				<div className="text-center mt-8">
					<button
						onClick={handleGoToLogin}
						className="text-blue-600 hover:text-blue-800 underline"
					>
						← กลับไปหน้าเข้าสู่ระบบ
					</button>
				</div>
			</div>
		</div>
	)
}
