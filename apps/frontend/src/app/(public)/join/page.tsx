'use client'

import { useState, type FormEvent } from 'react'
import { toast } from 'sonner'
import { getMessage } from '@/shared/i18n'
import type { ApiErrorResponse } from '@/shared/types/api-error.type'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:9090'

type Mode = 'register' | 'login'

export default function JoinPage() {
  const [mode, setMode] = useState<Mode>('register')
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleRegister(e: FormEvent) {
    e.preventDefault()
    setLoading(true)

    try {
      const res = await fetch(`${API_URL}/auth/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, email, password }),
      })

      if (res.status === 201) {
        toast.success(getMessage('auth.register_success'))
        setName('')
        setEmail('')
        setPassword('')
        return
      }

      const data: ApiErrorResponse = await res.json()
      if (data.errors && data.errors.length > 0) {
        for (const code of data.errors) {
          toast.error(getMessage(code))
        }
      } else if (data.error) {
        toast.error(getMessage(data.error))
      }
    } catch {
      toast.error('Unexpected error')
    } finally {
      setLoading(false)
    }
  }

  function handleLogin(e: FormEvent) {
    e.preventDefault()
    toast.info(getMessage('auth.login_coming_soon'))
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-50">
      <div className="w-full max-w-md rounded-xl border border-zinc-200 bg-white p-8 shadow-sm">
        <h1 className="mb-6 text-center text-2xl font-semibold text-zinc-900">
          {mode === 'register' ? getMessage('join.title_register') : getMessage('join.title_login')}
        </h1>

        {mode === 'register' ? (
          <form onSubmit={handleRegister} className="space-y-4">
            <div>
              <label htmlFor="name" className="mb-1 block text-sm font-medium text-zinc-700">
                {getMessage('join.name_label')}
              </label>
              <input
                id="name"
                type="text"
                required
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder={getMessage('join.name_placeholder')}
                className="w-full rounded-md border border-zinc-300 px-3 py-2 text-sm outline-none focus:border-zinc-500 focus:ring-1 focus:ring-zinc-500"
              />
            </div>

            <div>
              <label htmlFor="email" className="mb-1 block text-sm font-medium text-zinc-700">
                {getMessage('join.email_label')}
              </label>
              <input
                id="email"
                type="email"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder={getMessage('join.email_placeholder')}
                className="w-full rounded-md border border-zinc-300 px-3 py-2 text-sm outline-none focus:border-zinc-500 focus:ring-1 focus:ring-zinc-500"
              />
            </div>

            <div>
              <label htmlFor="password" className="mb-1 block text-sm font-medium text-zinc-700">
                {getMessage('join.password_label')}
              </label>
              <input
                id="password"
                type="password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder={getMessage('join.password_placeholder')}
                className="w-full rounded-md border border-zinc-300 px-3 py-2 text-sm outline-none focus:border-zinc-500 focus:ring-1 focus:ring-zinc-500"
              />
            </div>

            <button
              type="submit"
              disabled={loading}
              className="w-full rounded-md bg-zinc-900 px-4 py-2 text-sm font-medium text-white transition hover:bg-zinc-800 disabled:opacity-50"
            >
              {loading ? '...' : getMessage('join.register_button')}
            </button>
          </form>
        ) : (
          <form onSubmit={handleLogin} className="space-y-4">
            <div>
              <label htmlFor="login-email" className="mb-1 block text-sm font-medium text-zinc-700">
                {getMessage('join.email_label')}
              </label>
              <input
                id="login-email"
                type="email"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder={getMessage('join.email_placeholder')}
                className="w-full rounded-md border border-zinc-300 px-3 py-2 text-sm outline-none focus:border-zinc-500 focus:ring-1 focus:ring-zinc-500"
              />
            </div>

            <div>
              <label htmlFor="login-password" className="mb-1 block text-sm font-medium text-zinc-700">
                {getMessage('join.password_label')}
              </label>
              <input
                id="login-password"
                type="password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder={getMessage('join.password_placeholder')}
                className="w-full rounded-md border border-zinc-300 px-3 py-2 text-sm outline-none focus:border-zinc-500 focus:ring-1 focus:ring-zinc-500"
              />
            </div>

            <button
              type="submit"
              className="w-full rounded-md bg-zinc-900 px-4 py-2 text-sm font-medium text-white transition hover:bg-zinc-800"
            >
              {getMessage('join.login_button')}
            </button>
          </form>
        )}

        <div className="mt-6 text-center">
          <button
            type="button"
            onClick={() => setMode(mode === 'register' ? 'login' : 'register')}
            className="text-sm text-zinc-500 underline hover:text-zinc-700"
          >
            {mode === 'register' ? getMessage('join.switch_to_login') : getMessage('join.switch_to_register')}
          </button>
        </div>
      </div>
    </div>
  )
}
