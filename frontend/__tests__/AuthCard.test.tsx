import { render, screen } from '@testing-library/react'
import { AuthCard } from '@/components/AuthCard'
import { useRouter } from 'next/navigation'

jest.mock('next/navigation', () => ({
  useRouter: jest.fn(),
}))

jest.mock('@/lib/api', () => ({
  authApi: {
    signup: jest.fn(),
    login: jest.fn(),
  },
}))

jest.mock('@/lib/store', () => ({
  useAuthStore: jest.fn(() => ({
    user: null,
    token: null,
    hydrated: true,
    setAuth: jest.fn(),
    clearAuth: jest.fn(),
    loadAuth: jest.fn(),
  })),
}))

describe('AuthCard', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    ;(useRouter as jest.Mock).mockReturnValue({
      push: jest.fn(),
    })
  })

  it('renders login form by default', () => {
    render(<AuthCard />)
    expect(screen.getByText('Task Manager')).toBeInTheDocument()
    expect(screen.getByText('Sign in to your account')).toBeInTheDocument()
  })

  it('has email and password input fields', () => {
    render(<AuthCard />)
    expect(screen.getByPlaceholderText('you@example.com')).toBeInTheDocument()
    expect(screen.getByPlaceholderText('••••••••')).toBeInTheDocument()
  })

  it('has submit button', () => {
    render(<AuthCard />)
    expect(screen.getByText('Sign Up')).toBeInTheDocument()
  })

  it('has toggle to sign up', () => {
    render(<AuthCard />)
    expect(screen.getByText("Don't have an account?")).toBeInTheDocument()
  })
})
