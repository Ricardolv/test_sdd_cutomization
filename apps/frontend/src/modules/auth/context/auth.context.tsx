'use client'

import { createContext, useContext, useState, useEffect, type ReactNode } from 'react'
import Cookies from 'js-cookie'
import { decodeJwtPayload, type JwtPayload } from '../util/jwt.util'

interface User {
  id: string
  name: string
  email: string
}

type AuthStatus = 'loading' | 'authenticated' | 'unauthenticated'

interface AuthContextType {
  user: User | null
  token: string | null
  status: AuthStatus
  login: (token: string) => void
  logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [status, setStatus] = useState<AuthStatus>('loading')

  useEffect(() => {
    const cookieToken = Cookies.get('auth_token')
    if (cookieToken) {
      const payload = decodeJwtPayload(cookieToken)
      if (payload) {
        setUser({ id: payload.sub, name: payload.name, email: payload.email })
        setToken(cookieToken)
        setStatus('authenticated')
      } else {
        Cookies.remove('auth_token')
        setStatus('unauthenticated')
      }
    } else {
      setStatus('unauthenticated')
    }
  }, [])

  function login(newToken: string) {
    const expires = new Date()
    expires.setDate(expires.getDate() + 7)
    Cookies.set('auth_token', newToken, {
      expires,
      sameSite: 'lax',
      secure: process.env.NODE_ENV === 'production',
    })
    const payload = decodeJwtPayload(newToken)
    if (payload) {
      setUser({ id: payload.sub, name: payload.name, email: payload.email })
      setToken(newToken)
      setStatus('authenticated')
    }
  }

  function logout() {
    Cookies.remove('auth_token')
    setUser(null)
    setToken(null)
    setStatus('unauthenticated')
  }

  return (
    <AuthContext.Provider value={{ user, token, status, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth(): AuthContextType {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
