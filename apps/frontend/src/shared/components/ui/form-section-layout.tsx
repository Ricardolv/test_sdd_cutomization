import React from 'react'

interface FormSectionLayoutProps {
  title: string
  description?: string
  children: React.ReactNode
}

export function FormSectionLayout({ title, description, children }: FormSectionLayoutProps) {
  return (
    <div className="border-b border-zinc-200 pb-6 dark:border-zinc-800">
      <h3 className="text-lg font-medium text-zinc-900 dark:text-zinc-50">{title}</h3>
      {description && (
        <p className="mt-1 text-sm text-zinc-500 dark:text-zinc-400">{description}</p>
      )}
      <div className="mt-4">{children}</div>
    </div>
  )
}
