import React from 'react'
import { cn } from '@/lib/utils'
import { NavItem } from './nav-item'

export interface NavItemConfig {
  href: string
  label: string
  icon?: React.ReactNode
}

interface SidebarProps {
  items: NavItemConfig[]
  className?: string
}

export function Sidebar({ items, className }: SidebarProps) {
  return (
    <aside
      className={cn(
        'flex flex-col w-64 min-h-screen bg-zinc-50 dark:bg-zinc-900 border-r border-zinc-200 dark:border-zinc-800 p-4',
        className
      )}
    >
      <nav className="flex flex-col gap-1">
        {items.map((item) => (
          <NavItem key={item.href} href={item.href} icon={item.icon}>
            {item.label}
          </NavItem>
        ))}
      </nav>
    </aside>
  )
}
