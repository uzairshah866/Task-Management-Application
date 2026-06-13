'use client'

import { useEffect, useState } from 'react'
import { useAuthStore } from '@/lib/store'

export function Providers({ children }: { children: React.ReactNode }) {
  const [mounted, setMounted] = useState(false)
  const loadAuth = useAuthStore((state) => state.loadAuth)

  useEffect(() => {
    loadAuth()
    setMounted(true)
  }, [loadAuth])

  useEffect(() => {
    const theme = localStorage.getItem('theme') || 'light'
    const root = document.documentElement
    if (theme === 'dark') {
      root.classList.add('dark')
    } else {
      root.classList.remove('dark')
    }
  }, [])

  if (!mounted) return <>{children}</>

  return <>{children}</>
}
