'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '@/lib/store'
import { authApi } from '@/lib/api'
import { Buttons } from '@/app/ui/Buttons'

export function AuthCard() {
  const router = useRouter()
  const setAuth = useAuthStore((state) => state.setAuth)
  const [isSignup, setIsSignup] = useState(false)
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const response = isSignup
        ? await authApi.signup(email, password)
        : await authApi.login(email, password)

      const { token, user } = response.data
      setAuth(user, token)
      router.push('/')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Authentication failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="w-full max-w-md rounded-lg bg-white p-8 shadow-lg dark:bg-gray-900">
      <h1 className="mb-2 text-center text-2xl font-bold text-gray-900 dark:text-white">
        Task Manager
      </h1>
      <p className="mb-6 text-center text-gray-600 dark:text-gray-400">
        {isSignup ? 'Create an account' : 'Sign in to your account'}
      </p>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
            Email
          </label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            className="mt-1 w-full rounded-lg border border-gray-300 px-4 py-2 dark:border-gray-600 dark:bg-gray-800 dark:text-white"
            placeholder="you@example.com"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
            Password
          </label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            className="mt-1 w-full rounded-lg border border-gray-300 px-4 py-2 dark:border-gray-600 dark:bg-gray-800 dark:text-white"
            placeholder="••••••••"
          />
        </div>

        {error && <div className="rounded bg-red-100 p-3 text-sm text-red-700">{error}</div>}

        <Buttons
          size="lg"
          type="submit"
          variant="primary"
          className="w-full flex-1"
          disabled={loading}
          onClick={() => {
            setIsSignup(!isSignup)
            setError('')
          }}
        >
          {loading ? 'Loading...' : isSignup ? 'Sign Up' : 'Sign In'}
        </Buttons>
      </form>

      <div className="mt-6 text-center">
        <p className="text-sm text-gray-600 dark:text-gray-400">
          {isSignup ? 'Already have an account?' : "Don't have an account?"}
          <button
            onClick={() => {
              setIsSignup(!isSignup)
              setError('')
            }}
            className="ml-2 cursor-pointer text-blue-600 hover:underline dark:text-blue-400"
          >
            {isSignup ? 'Sign In' : 'Sign Up'}
          </button>
        </p>
      </div>
    </div>
  )
}
