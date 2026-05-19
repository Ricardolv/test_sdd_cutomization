'use client'

import React, { useEffect, useState } from 'react'
import Link from 'next/link'
import { Pencil, Trash2, Plus } from 'lucide-react'
import { api } from '@/shared/lib/api'
import { getMessage } from '@/shared/i18n'
import { Button } from '@/shared/components'
import { DeleteConfirmationDialog } from '@/shared/components/ui/delete-confirmation-dialog'

interface User {
  id: string
  name: string
  email: string
}

interface ListResponse {
  items: User[]
  total: number
  page: number
  limit: number
}

export default function UsersListPage() {
  const [users, setUsers] = useState<User[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(true)
  const [deleteTarget, setDeleteTarget] = useState<User | null>(null)

  const limit = 10

  useEffect(() => {
    fetchUsers()
  }, [page])

  async function fetchUsers() {
    setLoading(true)
    try {
      const data = await api.get<ListResponse>(`/users?page=${page}&limit=${limit}`)
      setUsers(data.items)
      setTotal(data.total)
    } catch (e) {
      console.error('Failed to fetch users', e)
    } finally {
      setLoading(false)
    }
  }

  async function handleDelete() {
    if (!deleteTarget) return
    try {
      await api.delete(`/users/${deleteTarget.id}`)
      setDeleteTarget(null)
      fetchUsers()
    } catch (e) {
      console.error('Failed to delete user', e)
    }
  }

  const totalPages = Math.ceil(total / limit)

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <h1 className="text-2xl font-semibold text-zinc-900 dark:text-zinc-50">
          {getMessage('user.list_title')}
        </h1>
        <Link href="/users/new">
          <Button>
            <Plus className="mr-2 h-4 w-4" />
            {getMessage('user.new')}
          </Button>
        </Link>
      </div>

      {loading ? (
        <p className="text-zinc-500">Loading...</p>
      ) : (
        <>
          <div className="rounded-lg border border-zinc-200 dark:border-zinc-800">
            <table className="w-full">
              <thead className="border-b border-zinc-200 bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900">
                <tr>
                  <th className="px-4 py-3 text-left text-sm font-medium text-zinc-600 dark:text-zinc-400">
                    {getMessage('user.name_label')}
                  </th>
                  <th className="px-4 py-3 text-left text-sm font-medium text-zinc-600 dark:text-zinc-400">
                    {getMessage('user.email_label')}
                  </th>
                  <th className="px-4 py-3 text-left text-sm font-medium text-zinc-600 dark:text-zinc-400">
                    {getMessage('user.actions')}
                  </th>
                </tr>
              </thead>
              <tbody>
                {users.map((user) => (
                  <tr key={user.id} className="border-b border-zinc-100 dark:border-zinc-800">
                    <td className="px-4 py-3 text-sm text-zinc-900 dark:text-zinc-50">{user.name}</td>
                    <td className="px-4 py-3 text-sm text-zinc-600 dark:text-zinc-400">{user.email}</td>
                    <td className="px-4 py-3">
                      <div className="flex gap-2">
                        <Link href={`/users/${user.id}`}>
                          <Button variant="ghost" size="sm">
                            <Pencil className="h-4 w-4" />
                          </Button>
                        </Link>
                        <Button variant="ghost" size="sm" onClick={() => setDeleteTarget(user)}>
                          <Trash2 className="h-4 w-4 text-red-600" />
                        </Button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {totalPages > 1 && (
            <div className="mt-4 flex items-center justify-between">
              <p className="text-sm text-zinc-600 dark:text-zinc-400">
                Page {page} of {totalPages} ({total} total)
              </p>
              <div className="flex gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  disabled={page <= 1}
                  onClick={() => setPage((p) => p - 1)}
                >
                  Previous
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  disabled={page >= totalPages}
                  onClick={() => setPage((p) => p + 1)}
                >
                  Next
                </Button>
              </div>
            </div>
          )}
        </>
      )}

      <DeleteConfirmationDialog
        open={!!deleteTarget}
        title={getMessage('user.delete_title')}
        description={getMessage('user.delete_description')}
        onConfirm={handleDelete}
        onCancel={() => setDeleteTarget(null)}
      />
    </div>
  )
}
