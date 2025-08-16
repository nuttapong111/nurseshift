import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatDate(date: Date): string {
  return new Intl.DateTimeFormat('th-TH', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  }).format(date)
}

export function formatTime(date: Date): string {
  return new Intl.DateTimeFormat('th-TH', {
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

export function getTimeRemaining(daysLeft: number): string {
  if (daysLeft <= 0) return "หมดอายุแล้ว"
  if (daysLeft === 1) return "เหลืออีก 1 วัน"
  return `เหลืออีก ${daysLeft} วัน`
}

// Ensure env base URL is absolute with protocol
export function normalizeBaseUrl(value: string | undefined, fallback?: string): string {
  const v = (value || '').trim()
  if (v.startsWith('http://') || v.startsWith('https://')) return v.replace(/\/$/, '')
  if (v) return `https://${v.replace(/\/$/, '')}`
  return (fallback || '').replace(/\/$/, '')
}
