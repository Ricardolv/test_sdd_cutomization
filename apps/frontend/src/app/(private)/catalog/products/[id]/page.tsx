'use client'

import React, { useEffect, useState } from 'react'
import { getMessage } from '@/shared/i18n'
import { ProductForm } from '@/modules/catalog/components/product-form'
import { api } from '@/shared/lib/api'

interface Product {
  id: string
  name: string
  description: string
  price: number
  status: string
  availableOnline: boolean
  featured: boolean
  allowsPreOrder: boolean
}

export default function EditProductPage({ params }: { params: { id: string } }) {
  const [product, setProduct] = useState<Product | null>(null)
  const [loading, setLoading] = useState(true)
  const [notFound, setNotFound] = useState(false)

  useEffect(() => {
    async function fetchProduct() {
      try {
        const data = await api.get<Product>(`/products/${params.id}`)
        setProduct(data)
      } catch (e) {
        setNotFound(true)
      } finally {
        setLoading(false)
      }
    }
    fetchProduct()
  }, [params.id])

  if (loading) {
    return <p className="text-zinc-500">Loading...</p>
  }

  if (notFound || !product) {
    return <p className="text-zinc-500">{getMessage('product.not_found')}</p>
  }

  return (
    <div>
      <h1 className="mb-6 text-2xl font-semibold text-zinc-900 dark:text-zinc-50">
        {getMessage('product.edit')}
      </h1>
      <ProductForm
        initialData={{
          id: product.id,
          name: product.name,
          description: product.description,
          price: product.price,
          status: product.status,
          availableOnline: product.availableOnline,
          featured: product.featured,
          allowsPreOrder: product.allowsPreOrder,
        }}
        isEdit
      />
    </div>
  )
}
