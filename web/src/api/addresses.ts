import apiClient from './client'
import type { Address, CreateAddressRequest, UpdateAddressRequest } from '@/types'

export const addressesApi = {
  getAddresses: async (): Promise<Address[]> => {
    const response = await apiClient.get('/addresses')
    return response.data.result
  },

  getAddress: async (id: string): Promise<Address> => {
    const response = await apiClient.get(`/addresses/${id}`)
    return response.data.result
  },

  createAddress: async (data: CreateAddressRequest): Promise<Address> => {
    const response = await apiClient.post('/addresses', data)
    return response.data.result
  },

  updateAddress: async (
    id: string,
    data: UpdateAddressRequest
  ): Promise<Address> => {
    const response = await apiClient.put(`/addresses/${id}`, data)
    return response.data.result
  },

  deleteAddress: async (id: string): Promise<void> => {
    await apiClient.delete(`/addresses/${id}`)
  },

  setDefaultAddress: async (id: string): Promise<Address> => {
    const response = await apiClient.put(`/addresses/${id}/default`)
    return response.data.result
  },
}
