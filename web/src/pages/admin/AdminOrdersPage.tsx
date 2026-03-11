import { ChevronDown, Search } from 'lucide-react'
import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import toast from 'react-hot-toast'
import { ordersApi } from '@/api/orders'
import LoadingSpinner from '@/components/LoadingSpinner'
import Pagination from '@/components/Pagination'
import type { Order, OrdersQueryParams } from '@/types'

const STATUS_OPTIONS = ['pending', 'processing', 'shipped', 'delivered', 'cancelled']

const STATUS_COLORS: Record<string, string> = {
  pending: 'badge-warning',
  processing: 'badge-info',
  shipped: 'badge-info',
  delivered: 'badge-success',
  cancelled: 'badge-danger',
}

export default function AdminOrdersPage() {
  const queryClient = useQueryClient()
  const [params, setParams] = useState<OrdersQueryParams>({ page: 1, limit: 20 })
  const [searchCode, setSearchCode] = useState('')
  const [updatingOrderId, setUpdatingOrderId] = useState<string | null>(null)

  const { data, isLoading } = useQuery({
    queryKey: ['admin-orders', params],
    queryFn: () => ordersApi.getOrders(params),
  })

  const updateStatusMutation = useMutation({
    mutationFn: ({ id, status }: { id: string; status: string }) =>
      ordersApi.updateOrderStatus(id, status),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-orders'] })
      setUpdatingOrderId(null)
      toast.success('Order status updated!')
    },
    onError: () => {
      toast.error('Failed to update order status')
      setUpdatingOrderId(null)
    },
  })

  const handleSearch = () => {
    setParams((p) => ({ ...p, code: searchCode || undefined, page: 1 }))
  }

  const handleStatusChange = (order: Order, newStatus: string) => {
    if (newStatus === order.status) return
    setUpdatingOrderId(order.id)
    updateStatusMutation.mutate({ id: order.id, status: newStatus })
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Orders</h1>
          {data && (
            <p className="text-sm text-gray-500">{data.pagination.total} total orders</p>
          )}
        </div>
      </div>

      {/* Filters */}
      <div className="flex gap-2 mb-5">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
          <input
            type="text"
            placeholder="Search by order code..."
            value={searchCode}
            onChange={(e) => setSearchCode(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
            className="input pl-9"
          />
        </div>
        <select
          value={params.status || ''}
          onChange={(e) =>
            setParams((p) => ({ ...p, status: e.target.value || undefined, page: 1 }))
          }
          className="input w-44"
        >
          <option value="">All Statuses</option>
          {STATUS_OPTIONS.map((s) => (
            <option key={s} value={s}>
              {s.charAt(0).toUpperCase() + s.slice(1)}
            </option>
          ))}
        </select>
        <button onClick={handleSearch} className="btn-primary">
          Search
        </button>
      </div>

      {isLoading ? (
        <LoadingSpinner className="py-16" />
      ) : !data?.items?.length ? (
        <div className="bg-white rounded-xl border border-gray-100 p-12 text-center text-gray-500">
          <p className="font-medium">No orders found</p>
        </div>
      ) : (
        <div className="bg-white rounded-xl border border-gray-100 overflow-hidden">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-gray-100 bg-gray-50">
                <th className="text-left px-4 py-3 font-semibold text-gray-600">Order Code</th>
                <th className="text-left px-4 py-3 font-semibold text-gray-600">Customer</th>
                <th className="text-left px-4 py-3 font-semibold text-gray-600">Items</th>
                <th className="text-left px-4 py-3 font-semibold text-gray-600">Total</th>
                <th className="text-left px-4 py-3 font-semibold text-gray-600">Date</th>
                <th className="text-left px-4 py-3 font-semibold text-gray-600">Status</th>
                <th className="text-right px-4 py-3 font-semibold text-gray-600">Actions</th>
              </tr>
            </thead>
            <tbody>
              {data.items.map((order) => (
                <tr
                  key={order.id}
                  className="border-b border-gray-50 hover:bg-gray-50 transition-colors"
                >
                  <td className="px-4 py-3">
                    <Link
                      to={`/orders/${order.id}`}
                      className="font-mono font-medium text-primary-600 hover:text-primary-700"
                    >
                      #{order.code}
                    </Link>
                  </td>
                  <td className="px-4 py-3 text-gray-600">
                    {order.user?.email || '—'}
                  </td>
                  <td className="px-4 py-3 text-gray-600">
                    {order.lines?.length || 0} items
                  </td>
                  <td className="px-4 py-3 font-semibold text-gray-900">
                    ${order.total_price?.toFixed(2)}
                  </td>
                  <td className="px-4 py-3 text-gray-500">
                    {new Date(order.created_at).toLocaleDateString()}
                  </td>
                  <td className="px-4 py-3">
                    <span className={`badge ${STATUS_COLORS[order.status] || 'badge-default'}`}>
                      {order.status}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-right">
                    <div className="relative inline-block">
                      <select
                        value={order.status}
                        onChange={(e) => handleStatusChange(order, e.target.value)}
                        disabled={updatingOrderId === order.id}
                        className="appearance-none bg-white border border-gray-200 rounded-lg px-3 py-1.5 text-sm pr-8 text-gray-700 cursor-pointer hover:border-primary-300 focus:outline-none focus:ring-1 focus:ring-primary-500 disabled:opacity-50"
                      >
                        {STATUS_OPTIONS.map((s) => (
                          <option key={s} value={s}>
                            {s.charAt(0).toUpperCase() + s.slice(1)}
                          </option>
                        ))}
                      </select>
                      <ChevronDown className="absolute right-2 top-1/2 -translate-y-1/2 h-3.5 w-3.5 text-gray-400 pointer-events-none" />
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>

          {data.pagination && (
            <div className="px-4 pb-4">
              <Pagination
                pagination={data.pagination}
                onPageChange={(page) => setParams((p) => ({ ...p, page }))}
              />
            </div>
          )}
        </div>
      )}
    </div>
  )
}
