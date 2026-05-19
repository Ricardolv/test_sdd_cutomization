'use client'

import { useState } from 'react'
import { useAuth } from '@/modules/auth'

export function PrivateShell({ children }: { children: React.ReactNode }) {
  const { user, logout } = useAuth()
  const [dropdownOpen, setDropdownOpen] = useState(false)

  function handleLogout() {
    logout()
    window.location.href = '/join'
  }

  return (
    <div className="flex h-screen flex-col">
      <header className="flex h-14 items-center justify-end border-b border-zinc-200 bg-white px-4">
        <div className="relative">
          <button
            type="button"
            onClick={() => setDropdownOpen(!dropdownOpen)}
            className="flex items-center gap-2 rounded-md px-3 py-1.5 text-sm text-zinc-700 hover:bg-zinc-100"
          >
            <span className="font-medium">{user?.name ?? 'User'}</span>
            <span className="text-zinc-400">{user?.email ?? ''}</span>
          </button>

          {dropdownOpen && (
            <div className="absolute right-0 z-50 mt-1 w-48 rounded-md border border-zinc-200 bg-white py-1 shadow-sm">
              <button
                type="button"
                onClick={handleLogout}
                className="w-full px-4 py-2 text-left text-sm text-red-600 hover:bg-zinc-50"
              >
                Logout
              </button>
            </div>
          )}
        </div>
      </header>

      <div className="flex flex-1 overflow-hidden">
        {children}
      </div>
    </div>
  )
}
