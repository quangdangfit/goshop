import apiClient from './client'

export interface PaymentIntent {
  intent_id: string
  client_secret: string
  amount: number
  currency: string
  status: string
}

export interface PublicConfig {
  stripe_publishable_key: string
}

export const paymentsApi = {
  createIntent: async (orderID: string): Promise<PaymentIntent> => {
    const response = await apiClient.post(`/orders/${orderID}/payment-intent`)
    return response.data.result
  },

  publicConfig: async (): Promise<PublicConfig> => {
    const response = await apiClient.get('/config/public')
    return response.data.result
  },
}
