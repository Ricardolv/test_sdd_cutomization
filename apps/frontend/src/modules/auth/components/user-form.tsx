'use client'

import React, { useState } from 'react'
import { useRouter } from 'next/navigation'
import { api } from '@/shared/lib/api'
import { getMessage } from '@/shared/i18n'
import { Button } from '@/shared/components'
import { FormSectionLayout } from '@/shared/components/ui/form-section-layout'

interface UserFormProps {
  initialData?: {
    id: string
    name: string
    email: string
  }
  isEdit?: boolean
}

export function UserForm({ initialData, isEdit = false }: UserFormProps) {
  const router = useRouter()
  const [name, setName] = useState(initialData?.name ?? '')
  const [email, setEmail] = useState(initialData?.email ?? '')
  const [password, setPassword] = useState('')
  const [passwordConfirm, setPasswordConfirm] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError(null)

    if (password !== passwordConfirm) {
      setError(getMessage('user.passwords_mismatch'))
      return
    }

    setLoading(true)
    try {
      const body: Record<string, string> = { name, email }
      if (password) {
        body.password = password
      }

      if (isEdit && initialData) {
        await api.put(`/users/${initialData.id}`, body)
      } else {
        await api.post('/users', body)
      }

      router.push('/users')
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Unexpected error'
      setError(message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <FormSectionLayout
        title={getMessage('user.basic_data')}
      >
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300">
              {getMessage('user.name_label')}
            </label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              className="mt-1 w-full rounded-md border border-zinc-300 px-3 py-2 text-sm dark:border-zinc-700 dark:bg-zinc-900"
              placeholder={getMessage('user.name_placeholder')}
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300">
              {getMessage('user.email_label')}
            </label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              className="mt-1 w-full rounded-md border border-zinc-300 px-3 py-2 text-sm dark:border-zinc-700 dark:bg-zinc-900"
              placeholder={getMessage('user.email_placeholder')}
            />
          </div>
        </div>
      </FormSectionLayout>

      <FormSectionLayout
        title={getMessage('user.password_section')}
        description={isEdit ? getMessage('user.password_description') : undefined}
      >
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300">
              {getMessage('user.password_label')}
            </label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required={!isEdit}
              className="mt-1 w-full rounded-md border border-zinc-300 px-3 py-2 text-sm dark:border-zinc-700 dark:bg-zinc-900"
              placeholder={getMessage('user.password_placeholder')}
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300">
              {getMessage('user.password_confirm_label')}
            </label>
            <input
              type="password"
              value={passwordConfirm}
              onChange={(e) => setPasswordConfirm(e.target.value)}
              required={!isEdit}
              className="mt-1 w-full rounded-md border border-zinc-300 px-3 py-2 text-sm dark:border-zinc-700 dark:bg-zinc-900"
              placeholder={getMessage('user.password_confirm_placeholder')}
            />
          </div>
        </div>
      </FormSectionLayout>

      {error && (
        <p className="text-sm text-red-600">{error}</p>
      )}

      <div className="flex gap-2">
        <Button type="submit" disabled={loading}>
          {loading ? 'Saving...' : getMessage('user.save')}
        </Button>
        <Button type="button" variant="outline" onClick={() => router.push('/users')}>
          {getMessage('user.cancel')}
        </Button>
      </div>
    </form>
  )
}
