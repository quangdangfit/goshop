import apiClient from './client'
import type {
  Category,
  CreateCategoryRequest,
  UpdateCategoryRequest,
} from '@/types'

export const categoriesApi = {
  getCategories: async (): Promise<Category[]> => {
    const response = await apiClient.get('/categories')
    return response.data.result
  },

  getCategory: async (id: string): Promise<Category> => {
    const response = await apiClient.get(`/categories/${id}`)
    return response.data.result
  },

  createCategory: async (data: CreateCategoryRequest): Promise<Category> => {
    const response = await apiClient.post('/categories', data)
    return response.data.result
  },

  updateCategory: async (
    id: string,
    data: UpdateCategoryRequest
  ): Promise<Category> => {
    const response = await apiClient.put(`/categories/${id}`, data)
    return response.data.result
  },

  deleteCategory: async (id: string): Promise<void> => {
    await apiClient.delete(`/categories/${id}`)
  },
}
