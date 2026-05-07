import apiClient from './client'

export interface NotificationPreference {
  user_id: string
  event_type: string
  channel: string
  enabled: boolean
}

export const notificationsApi = {
  list: async (): Promise<NotificationPreference[]> => {
    const response = await apiClient.get('/me/notification-preferences')
    return response.data.result ?? []
  },

  set: async (
    eventType: string,
    channel: string,
    enabled: boolean,
  ): Promise<NotificationPreference> => {
    const response = await apiClient.put('/me/notification-preferences', {
      event_type: eventType,
      channel,
      enabled,
    })
    return response.data.result
  },
}
