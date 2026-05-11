// Client-side cart. The server no longer stores carts; this module owns the source of truth
// and exposes a small event-based API so React Query can subscribe via cartStore.subscribe().
//
// Storage key is schema-versioned (`goshop:cart:v1`) so future shape changes can migrate
// without colliding with stale entries.

import type { Product } from '@/types'

const STORAGE_KEY = 'goshop:cart:v1'

export interface CartItem {
  product_id: string
  quantity: number
  // Snapshot of the product at the time of add. Used for offline display; the server
  // re-validates price + stock at order time, so this is presentational only.
  snapshot: Pick<Product, 'name' | 'price' | 'images' | 'stock_quantity'>
}

export interface ClientCart {
  version: 1
  items: CartItem[]
  updated_at: string
}

const emptyCart = (): ClientCart => ({
  version: 1,
  items: [],
  updated_at: new Date().toISOString(),
})

function read(): ClientCart {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return emptyCart()
    const parsed = JSON.parse(raw) as ClientCart
    if (parsed.version !== 1 || !Array.isArray(parsed.items)) return emptyCart()
    return parsed
  } catch {
    return emptyCart()
  }
}

const subscribers = new Set<() => void>()

// Cached snapshot for useSyncExternalStore: getSnapshot() must return a stable reference
// while the underlying state hasn't changed, otherwise React enters an infinite render loop.
// We re-read localStorage on every call (cheap), but only re-parse and produce a new object
// when the raw bytes have actually changed.
let cachedRaw: string | null = null
let cachedSnapshot: ClientCart = emptyCart()

function snapshot(): ClientCart {
  const raw = (typeof localStorage !== 'undefined' ? localStorage.getItem(STORAGE_KEY) : null) ?? ''
  if (raw !== cachedRaw) {
    cachedRaw = raw
    cachedSnapshot = read()
  }
  return cachedSnapshot
}

function write(cart: ClientCart) {
  cart.updated_at = new Date().toISOString()
  const raw = JSON.stringify(cart)
  localStorage.setItem(STORAGE_KEY, raw)
  cachedRaw = raw
  cachedSnapshot = cart
  subscribers.forEach((fn) => fn())
}

export const cartStore = {
  get(): ClientCart {
    return snapshot()
  },

  subscribe(fn: () => void): () => void {
    subscribers.add(fn)
    return () => subscribers.delete(fn)
  },

  add(product: Product, quantity = 1): ClientCart {
    const cart = read()
    const idx = cart.items.findIndex((i) => i.product_id === product.id)
    if (idx >= 0) {
      cart.items[idx].quantity += quantity
      cart.items[idx].snapshot = pickSnapshot(product)
    } else {
      cart.items.push({
        product_id: product.id,
        quantity,
        snapshot: pickSnapshot(product),
      })
    }
    write(cart)
    return cart
  },

  setQuantity(productId: string, quantity: number): ClientCart {
    const cart = read()
    if (quantity <= 0) {
      cart.items = cart.items.filter((i) => i.product_id !== productId)
    } else {
      const item = cart.items.find((i) => i.product_id === productId)
      if (item) item.quantity = quantity
    }
    write(cart)
    return cart
  },

  remove(productId: string): ClientCart {
    const cart = read()
    cart.items = cart.items.filter((i) => i.product_id !== productId)
    write(cart)
    return cart
  },

  clear(): ClientCart {
    const cart = emptyCart()
    write(cart)
    return cart
  },

  itemCount(): number {
    return read().items.reduce((sum, i) => sum + i.quantity, 0)
  },
}

function pickSnapshot(p: Product): CartItem['snapshot'] {
  return {
    name: p.name,
    price: p.price,
    images: p.images,
    stock_quantity: p.stock_quantity,
  }
}
