import { useMemo } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import toast from 'react-hot-toast'
import { Bell } from 'lucide-react'

import { notificationsApi, type NotificationPreference } from '@/api/notifications'
import LoadingSpinner from '@/components/LoadingSpinner'

// Event types and channels are seeded by the BE event bus; the FE renders the matrix using
// this static catalog. Server-side, missing rows default to enabled — toggling a row creates
// or updates the explicit preference.
const EVENT_TYPES: { id: string; label: string; description: string }[] = [
  { id: 'order_placed', label: 'Order placed', description: 'Confirmation when you place an order.' },
  { id: 'order_status_changed', label: 'Order status updates', description: 'When your order ships or is fulfilled.' },
]
const CHANNELS: { id: string; label: string }[] = [
  { id: 'email', label: 'Email' },
]

export default function NotificationPreferencesPage() {
  const qc = useQueryClient()

  const { data: prefs, isLoading } = useQuery({
    queryKey: ['notification-preferences'],
    queryFn: notificationsApi.list,
  })

  // Build a lookup keyed by event|channel for O(1) row resolution; missing rows mean default
  // enabled (the BE policy).
  const lookup = useMemo(() => {
    const m = new Map<string, NotificationPreference>()
    for (const p of prefs ?? []) m.set(`${p.event_type}|${p.channel}`, p)
    return m
  }, [prefs])

  const setMutation = useMutation({
    mutationFn: ({ event, channel, enabled }: { event: string; channel: string; enabled: boolean }) =>
      notificationsApi.set(event, channel, enabled),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['notification-preferences'] })
    },
    onError: () => toast.error('Failed to update preference'),
  })

  if (isLoading) return <LoadingSpinner className="min-h-[400px]" size="lg" />

  return (
    <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="text-2xl font-bold text-gray-900 mb-2 flex items-center gap-2">
        <Bell className="h-6 w-6 text-primary-600" />
        Notification preferences
      </h1>
      <p className="text-sm text-gray-500 mb-6">
        Choose which alerts you want to receive on each channel. Changes are saved immediately.
      </p>

      <div className="card">
        <table className="w-full text-sm">
          <thead>
            <tr className="text-left text-gray-500">
              <th className="pb-3">Event</th>
              {CHANNELS.map((c) => (
                <th key={c.id} className="pb-3 text-center w-32">{c.label}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {EVENT_TYPES.map((evt) => (
              <tr key={evt.id} className="border-t border-gray-100">
                <td className="py-4">
                  <p className="font-medium text-gray-900">{evt.label}</p>
                  <p className="text-xs text-gray-500">{evt.description}</p>
                </td>
                {CHANNELS.map((c) => {
                  const row = lookup.get(`${evt.id}|${c.id}`)
                  const enabled = row?.enabled ?? true
                  return (
                    <td key={c.id} className="text-center">
                      <input
                        type="checkbox"
                        checked={enabled}
                        disabled={setMutation.isPending}
                        onChange={(e) =>
                          setMutation.mutate({ event: evt.id, channel: c.id, enabled: e.target.checked })
                        }
                        className="h-4 w-4 text-primary-600"
                      />
                    </td>
                  )
                })}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
