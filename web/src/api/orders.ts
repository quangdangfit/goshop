import apiClient from './client'
import type {
  CreateOrderRequest,
  Order,
  OrdersQueryParams,
  PaginatedResponse,
} from '@/types'

export const ordersApi = {
  createOrder: async (data: CreateOrderRequest): Promise<Order> => {
    const response = await apiClient.post('/orders', data)
    return response.data.result
  },

  getOrders: async (
    params?: OrdersQueryParams
  ): Promise<PaginatedResponse<Order>> => {
    const response = await apiClient.get('/orders', { params })
    return response.data.result
  },

  getOrder: async (id: string): Promise<Order> => {
    const response = await apiClient.get(`/orders/${id}`)
    return response.data.result
  },

  cancelOrder: async (id: string): Promise<Order> => {
    const response = await apiClient.put(`/orders/${id}/cancel`)
    return response.data.result
  },

  updateOrderStatus: async (id: string, status: string): Promise<Order> => {
    const response = await apiClient.put(`/orders/${id}/status`, null, {
      params: { status },
    })
    return response.data.result
  },
}
