'use client'

import React, { useEffect, useState } from 'react'
import { useParams } from 'next/navigation'
import { getMessage } from '@/shared/i18n'
import { UserForm } from '@/modules/auth/components/user-form'
import { api } from '@/shared/lib/api'

interface User {
  id: string
  name: string
  email: string
}

export default function EditUserPage() {
  const params = useParams()
  const id = params.id as string
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function fetchUser() {
      try {
        const data = await api.get<User>(`/users/${id}`)
        setUser(data)
      } catch (e) {
        console.error('Failed to fetch user', e)
      } finally {
        setLoading(false)
      }
    }
    fetchUser()
  }, [id])

  if (loading) {
    return <p className="text-zinc-500">Loading...</p>
  }

  if (!user) {
    return <p className="text-zinc-500">{getMessage('user.not_found')}</p>
  }

  return (
    <div>
      <h1 className="mb-6 text-2xl font-semibold text-zinc-900 dark:text-zinc-50">
        {getMessage('user.edit')}
      </h1>
      <UserForm initialData={user} isEdit />
    </div>
  )
}
