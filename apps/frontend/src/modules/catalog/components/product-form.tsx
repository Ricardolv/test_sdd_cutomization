'use client'

import React, { useState } from 'react'
import { useRouter } from 'next/navigation'
import { api } from '@/shared/lib/api'
import { getMessage } from '@/shared/i18n'
import { Button } from '@/shared/components'
import { FormSectionLayout } from '@/shared/components/ui/form-section-layout'

interface ProductFormProps {
  initialData?: {
    id: string
    name: string
    description: string
    price: number
    status: string
    availableOnline: boolean
    featured: boolean
    allowsPreOrder: boolean
  }
  isEdit?: boolean
}

export function ProductForm({ initialData, isEdit = false }: ProductFormProps) {
  const router = useRouter()
  const [name, setName] = useState(initialData?.name ?? '')
  const [description, setDescription] = useState(initialData?.description ?? '')
  const [price, setPrice] = useState(initialData?.price?.toString() ?? '')
  const [status, setStatus] = useState(initialData?.status ?? 'draft')
  const [availableOnline, setAvailableOnline] = useState(initialData?.availableOnline ?? false)
  const [featured, setFeatured] = useState(initialData?.featured ?? false)
  const [allowsPreOrder, setAllowsPreOrder] = useState(initialData?.allowsPreOrder ?? false)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError(null)

    const priceValue = parseFloat(price)
    if (isNaN(priceValue) || priceValue < 0) {
      setError(getMessage('product.price_invalid'))
      return
    }

    setLoading(true)
    try {
      const body = {
        name,
        description,
        price: priceValue,
        status,
        availableOnline,
        featured,
        allowsPreOrder,
      }

      if (isEdit && initialData) {
        await api.put(`/products/${initialData.id}`, body)
      } else {
        await api.post('/products', body)
      }

      router.push('/catalog/products')
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
        title={getMessage('product.basic_data')}
      >
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300">
              {getMessage('product.name_label')}
            </label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              minLength={2}
              maxLength={120}
              className="mt-1 w-full rounded-md border border-zinc-300 px-3 py-2 text-sm dark:border-zinc-700 dark:bg-zinc-900"
              placeholder={getMessage('product.name_placeholder')}
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300">
              {getMessage('product.description_label')}
            </label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              maxLength={500}
              rows={3}
              className="mt-1 w-full rounded-md border border-zinc-300 px-3 py-2 text-sm dark:border-zinc-700 dark:bg-zinc-900"
              placeholder={getMessage('product.description_placeholder')}
            />
          </div>
        </div>
      </FormSectionLayout>

      <FormSectionLayout
        title={getMessage('product.price_status_section')}
      >
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300">
              {getMessage('product.price_label')}
            </label>
            <input
              type="number"
              value={price}
              onChange={(e) => setPrice(e.target.value)}
              required
              min={0}
              step="0.01"
              className="mt-1 w-full rounded-md border border-zinc-300 px-3 py-2 text-sm dark:border-zinc-700 dark:bg-zinc-900"
              placeholder="0.00"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300">
              {getMessage('product.status_label')}
            </label>
            <select
              value={status}
              onChange={(e) => setStatus(e.target.value)}
              className="mt-1 w-full rounded-md border border-zinc-300 px-3 py-2 text-sm dark:border-zinc-700 dark:bg-zinc-900"
            >
              <option value="active">{getMessage('product.status.active')}</option>
              <option value="inactive">{getMessage('product.status.inactive')}</option>
              <option value="draft">{getMessage('product.status.draft')}</option>
            </select>
          </div>
        </div>
      </FormSectionLayout>

      <FormSectionLayout
        title={getMessage('product.availability_section')}
      >
        <div className="space-y-3">
          <label className="flex items-center gap-2">
            <input
              type="checkbox"
              checked={availableOnline}
              onChange={(e) => setAvailableOnline(e.target.checked)}
              className="h-4 w-4"
            />
            <span className="text-sm text-zinc-700 dark:text-zinc-300">
              {getMessage('product.available_online_label')}
            </span>
          </label>
          <label className="flex items-center gap-2">
            <input
              type="checkbox"
              checked={featured}
              onChange={(e) => setFeatured(e.target.checked)}
              className="h-4 w-4"
            />
            <span className="text-sm text-zinc-700 dark:text-zinc-300">
              {getMessage('product.featured_label')}
            </span>
          </label>
          <label className="flex items-center gap-2">
            <input
              type="checkbox"
              checked={allowsPreOrder}
              onChange={(e) => setAllowsPreOrder(e.target.checked)}
              className="h-4 w-4"
            />
            <span className="text-sm text-zinc-700 dark:text-zinc-300">
              {getMessage('product.allows_pre_order_label')}
            </span>
          </label>
        </div>
      </FormSectionLayout>

      {error && (
        <p className="text-sm text-red-600">{error}</p>
      )}

      <div className="flex gap-2">
        <Button type="submit" disabled={loading}>
          {loading ? 'Saving...' : getMessage('product.save')}
        </Button>
        <Button type="button" variant="outline" onClick={() => router.push('/catalog/products')}>
          {getMessage('product.cancel')}
        </Button>
      </div>
    </form>
  )
}
