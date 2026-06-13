'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '@/lib/store'
import { TaskDashboard } from '@/components/TaskDashboard'

export default function Home() {
  const router = useRouter()
  const { token, user, hydrated } = useAuthStore()

  useEffect(() => {
    if (hydrated && !token) {
      router.push('/login')
    }
  }, [hydrated, token, router])

  if (!hydrated) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-gradient-to-br from-blue-50 to-blue-100 dark:from-gray-950 dark:to-blue-950">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-blue-500 border-t-transparent" />
      </main>
    )
  }

  if (!token || !user) {
    return null
  }

  return (
    <main className="min-h-screen bg-gradient-to-br from-blue-50 to-blue-100 dark:from-gray-950 dark:to-blue-950">
      <TaskDashboard user={user} />
    </main>
  )
}
