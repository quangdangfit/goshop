import apiClient from './client'
import type { Coupon, CreateCouponRequest } from '@/types'

export const couponsApi = {
  createCoupon: async (data: CreateCouponRequest): Promise<Coupon> => {
    const response = await apiClient.post('/coupons', data)
    return response.data.result
  },

  getCoupon: async (code: string): Promise<Coupon> => {
    const response = await apiClient.get(`/coupons/${code}`)
    return response.data.result
  },
}
