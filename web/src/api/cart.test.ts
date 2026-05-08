import { afterEach, beforeEach, describe, expect, it } from 'vitest'

import { cartStore } from './cart'
import type { Product } from '@/types'

const product = (id: string, price = 10): Product => ({
  id,
  name: `Product ${id}`,
  code: `P-${id}`,
  description: '',
  price,
  stock_quantity: 100,
  images: [],
  category_id: null,
  active: true,
  avg_rating: 0,
  review_count: 0,
  created_at: '',
  updated_at: '',
} as unknown as Product)

describe('cartStore', () => {
  beforeEach(() => {
    localStorage.clear()
  })
  afterEach(() => {
    localStorage.clear()
  })

  it('starts empty', () => {
    expect(cartStore.get().items).toEqual([])
    expect(cartStore.itemCount()).toBe(0)
  })

  it('add merges quantity for an existing product', () => {
    cartStore.add(product('a'), 2)
    cartStore.add(product('a'), 3)
    const cart = cartStore.get()
    expect(cart.items).toHaveLength(1)
    expect(cart.items[0].quantity).toBe(5)
    expect(cartStore.itemCount()).toBe(5)
  })

  it('setQuantity updates and removes when <= 0', () => {
    cartStore.add(product('b'), 4)
    cartStore.setQuantity('b', 7)
    expect(cartStore.get().items[0].quantity).toBe(7)
    cartStore.setQuantity('b', 0)
    expect(cartStore.get().items).toHaveLength(0)
  })

  it('remove drops the matching line', () => {
    cartStore.add(product('a'))
    cartStore.add(product('b'))
    cartStore.remove('a')
    const ids = cartStore.get().items.map((i) => i.product_id)
    expect(ids).toEqual(['b'])
  })

  it('clear empties the cart', () => {
    cartStore.add(product('a'))
    cartStore.clear()
    expect(cartStore.itemCount()).toBe(0)
  })

  it('subscribe fires on writes and unsubscribe stops them', () => {
    let calls = 0
    const off = cartStore.subscribe(() => {
      calls++
    })
    cartStore.add(product('a'))
    cartStore.add(product('a'))
    expect(calls).toBe(2)
    off()
    cartStore.add(product('a'))
    expect(calls).toBe(2)
  })

  it('discards malformed localStorage payloads', () => {
    localStorage.setItem('goshop:cart:v1', 'not json')
    expect(cartStore.get().items).toEqual([])
    localStorage.setItem('goshop:cart:v1', JSON.stringify({ version: 999, items: 'nope' }))
    expect(cartStore.get().items).toEqual([])
  })
})
