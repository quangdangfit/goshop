// useCart subscribes to the localStorage-backed cart store and re-renders on changes.
// Reads through useSyncExternalStore so multiple components stay consistent without
// duplicating React state — and a tab-level localStorage 'storage' event is forwarded
// so changes made in another tab propagate.

import { useEffect, useSyncExternalStore } from 'react'

import { cartStore } from '@/api/cart'

export function useCart() {
  // Bridge cross-tab localStorage updates into the in-memory subscriber set.
  useEffect(() => {
    const handler = (e: StorageEvent) => {
      if (e.key === 'goshop:cart:v1') {
        // Triggering subscribers re-reads from localStorage, which now holds the new value.
        cartStore.add // no-op reference to keep the closure stable
      }
    }
    window.addEventListener('storage', handler)
    return () => window.removeEventListener('storage', handler)
  }, [])

  const cart = useSyncExternalStore(
    (cb) => cartStore.subscribe(cb),
    () => cartStore.get(),
    () => cartStore.get(),
  )
  return { cart, store: cartStore }
}
