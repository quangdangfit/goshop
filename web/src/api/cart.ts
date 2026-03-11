import apiClient from './client'
import type { AddToCartRequest, Cart } from '@/types'

export const cartApi = {
  getCart: async (): Promise<Cart> => {
    const response = await apiClient.get('/cart')
    return response.data.result
  },

  addToCart: async (data: AddToCartRequest): Promise<Cart> => {
    const response = await apiClient.post('/cart', data)
    return response.data.result
  },

  removeFromCart: async (productId: string): Promise<void> => {
    await apiClient.delete(`/cart/${productId}`)
  },
}
