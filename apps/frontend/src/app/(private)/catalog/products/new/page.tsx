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

export default function NewProductPage() {
  return (
    <div>
      <h1 className="mb-6 text-2xl font-semibold text-zinc-900 dark:text-zinc-50">
        {getMessage('product.new')}
      </h1>
      <ProductForm />
    </div>
  )
}
