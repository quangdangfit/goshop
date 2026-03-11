import apiClient from './client'
import type { WishlistItem } from '@/types'

export const wishlistApi = {
  getWishlist: async (): Promise<WishlistItem[]> => {
    const response = await apiClient.get('/wishlist')
    return response.data.result
  },

  addToWishlist: async (productId: string): Promise<WishlistItem> => {
    const response = await apiClient.post('/wishlist', { product_id: productId })
    return response.data.result
  },

  removeFromWishlist: async (productId: string): Promise<void> => {
    await apiClient.delete(`/wishlist/${productId}`)
  },
}
