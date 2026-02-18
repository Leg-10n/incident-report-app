import { useState, useCallback } from 'react'
import { AuthUser, AuthResponse, AuthFormData } from '../types/auth'

const TOKEN_KEY = 'incident_token'
const USER_KEY  = 'incident_user'

// These run once at module level to hydrate state from localStorage
function getStoredToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

function getStoredUser(): AuthUser | null {
  const raw = localStorage.getItem(USER_KEY)
  if (!raw) return null
  try { return JSON.parse(raw) } catch { return null }
}

export function useAuth() {
  // Pass initialiser functions so localStorage is only read once, not every render
  const [token, setToken] = useState<string | null>(getStoredToken)
  const [user, setUser]   = useState<AuthUser | null>(getStoredUser)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const storeAuth = (data: AuthResponse) => {
    localStorage.setItem(TOKEN_KEY, data.token)
    localStorage.setItem(USER_KEY, JSON.stringify(data.user))
    setToken(data.token)
    setUser(data.user)
    setError(null)
  }

  const login = useCallback(async (form: AuthFormData): Promise<boolean> => {
    setLoading(true)
    setError(null)
    try {
      const res = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(form),
      })
      const data = await res.json()
      if (!res.ok) throw new Error(data.error || 'Login failed')
      storeAuth(data as AuthResponse)
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed')
      return false
    } finally {
      setLoading(false)
    }
  }, [])

  const register = useCallback(async (form: AuthFormData): Promise<boolean> => {
    setLoading(true)
    setError(null)
    try {
      const res = await fetch('/api/auth/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(form),
      })
      const data = await res.json()
      if (!res.ok) throw new Error(data.error || 'Registration failed')
      storeAuth(data as AuthResponse)
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Registration failed')
      return false
    } finally {
      setLoading(false)
    }
  }, [])

  const logout = useCallback(() => {
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
    setToken(null)
    setUser(null)
  }, [])

  return { token, user, error, loading, login, register, logout }
}