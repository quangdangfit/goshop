import apiClient from './client'
import type {
  CreateProductRequest,
  PaginatedResponse,
  Product,
  ProductsQueryParams,
  Review,
  CreateReviewRequest,
  UpdateReviewRequest,
  UpdateProductRequest,
} from '@/types'

export const productsApi = {
  getProducts: async (
    params?: ProductsQueryParams
  ): Promise<PaginatedResponse<Product>> => {
    const response = await apiClient.get('/products', { params })
    return response.data.result
  },

  getProduct: async (id: string): Promise<Product> => {
    const response = await apiClient.get(`/products/${id}`)
    return response.data.result
  },

  createProduct: async (data: CreateProductRequest): Promise<Product> => {
    const response = await apiClient.post('/products', data)
    return response.data.result
  },

  updateProduct: async (
    id: string,
    data: UpdateProductRequest
  ): Promise<Product> => {
    const response = await apiClient.put(`/products/${id}`, data)
    return response.data.result
  },

  getReviews: async (
    productId: string,
    params?: { page?: number; limit?: number }
  ): Promise<PaginatedResponse<Review>> => {
    const response = await apiClient.get(`/products/${productId}/reviews`, {
      params,
    })
    return response.data.result
  },

  createReview: async (
    productId: string,
    data: CreateReviewRequest
  ): Promise<Review> => {
    const response = await apiClient.post(`/products/${productId}/reviews`, data)
    return response.data.result
  },

  updateReview: async (
    productId: string,
    reviewId: string,
    data: UpdateReviewRequest
  ): Promise<Review> => {
    const response = await apiClient.put(
      `/products/${productId}/reviews/${reviewId}`,
      data
    )
    return response.data.result
  },

  deleteReview: async (productId: string, reviewId: string): Promise<void> => {
    await apiClient.delete(`/products/${productId}/reviews/${reviewId}`)
  },
}
