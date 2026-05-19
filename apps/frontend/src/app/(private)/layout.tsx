import React from 'react'
import { Sidebar } from '@/shared/components'
import { AuthGuard } from '@/modules/auth'
import { LayoutDashboard, Users, Settings } from 'lucide-react'
import { PrivateShell } from '@/modules/auth/components/private-shell'

const navItems = [
  { href: '/dashboard', label: 'Dashboard', icon: <LayoutDashboard /> },
  { href: '/customers', label: 'Customers', icon: <Users /> },
  { href: '/settings', label: 'Settings', icon: <Settings /> },
]

export default function PrivateLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <AuthGuard>
      <PrivateShell>
        <Sidebar items={navItems} />
        <main className="flex-1 overflow-auto p-6">{children}</main>
      </PrivateShell>
    </AuthGuard>
  )
}
