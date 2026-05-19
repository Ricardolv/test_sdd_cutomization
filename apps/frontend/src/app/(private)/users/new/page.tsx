'use client'

import React from 'react'
import { getMessage } from '@/shared/i18n'
import { UserForm } from '@/modules/auth/components/user-form'

export default function NewUserPage() {
  return (
    <div>
      <h1 className="mb-6 text-2xl font-semibold text-zinc-900 dark:text-zinc-50">
        {getMessage('user.new')}
      </h1>
      <UserForm />
    </div>
  )
}
