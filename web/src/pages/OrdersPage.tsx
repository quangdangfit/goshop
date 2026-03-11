import { Package, Search } from 'lucide-react'
import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import toast from 'react-hot-toast'
import { ordersApi } from '@/api/orders'
import LoadingSpinner from '@/components/LoadingSpinner'
import Pagination from '@/components/Pagination'
import type { OrdersQueryParams } from '@/types'

const STATUS_COLORS: Record<string, string> = {
  pending: 'badge-warning',
  processing: 'badge-info',
  shipped: 'badge-info',
  delivered: 'badge-success',
  cancelled: 'badge-danger',
}

export default function OrdersPage() {
  const queryClient = useQueryClient()
  const [params, setParams] = useState<OrdersQueryParams>({ page: 1, limit: 10 })
  const [searchCode, setSearchCode] = useState('')

  const { data, isLoading } = useQuery({
    queryKey: ['orders', params],
    queryFn: () => ordersApi.getOrders(params),
  })

  const cancelMutation = useMutation({
    mutationFn: ordersApi.cancelOrder,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['orders'] })
      toast.success('Order cancelled')
    },
    onError: () => toast.error('Failed to cancel order'),
  })

  const handleSearch = () => {
    setParams((p) => ({ ...p, code: searchCode || undefined, page: 1 }))
  }

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="text-2xl font-bold text-gray-900 mb-6">My Orders</h1>

      {/* Search */}
      <div className="flex gap-2 mb-5">
        <div className="relative flex-1">
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
          className="input w-40"
        >
          <option value="">All Status</option>
          <option value="pending">Pending</option>
          <option value="processing">Processing</option>
          <option value="shipped">Shipped</option>
          <option value="delivered">Delivered</option>
          <option value="cancelled">Cancelled</option>
        </select>
        <button onClick={handleSearch} className="btn-primary">
          Search
        </button>
      </div>

      {isLoading ? (
        <LoadingSpinner className="py-12" />
      ) : !data?.items?.length ? (
        <div className="text-center py-16 text-gray-500">
          <Package className="h-14 w-14 text-gray-300 mx-auto mb-3" />
          <p className="font-medium">No orders found</p>
          <Link to="/products" className="btn-primary mt-4 inline-flex">
            Start Shopping
          </Link>
        </div>
      ) : (
        <>
          <div className="space-y-3">
            {data.items.map((order) => (
              <div
                key={order.id}
                className="bg-white rounded-xl border border-gray-100 p-5"
              >
                <div className="flex items-start justify-between mb-3">
                  <div>
                    <Link
                      to={`/orders/${order.id}`}
                      className="font-semibold text-gray-900 hover:text-primary-600 transition-colors"
                    >
                      #{order.code}
                    </Link>
                    <p className="text-xs text-gray-500 mt-0.5">
                      {new Date(order.created_at).toLocaleDateString('en-US', {
                        year: 'numeric',
                        month: 'long',
                        day: 'numeric',
                      })}
                    </p>
                  </div>
                  <div className="text-right">
                    <span
                      className={`badge ${STATUS_COLORS[order.status] || 'badge-default'}`}
                    >
                      {order.status}
                    </span>
                    <p className="font-bold text-gray-900 mt-1">
                      ${order.total_price?.toFixed(2)}
                    </p>
                  </div>
                </div>

                {/* Items preview */}
                <div className="text-sm text-gray-600 mb-3">
                  {order.lines?.slice(0, 2).map((line) => (
                    <span key={line.id} className="mr-2">
                      {line.product?.name} x{line.quantity}
                    </span>
                  ))}
                  {(order.lines?.length || 0) > 2 && (
                    <span className="text-gray-400">
                      +{order.lines.length - 2} more
                    </span>
                  )}
                </div>

                <div className="flex items-center gap-2">
                  <Link
                    to={`/orders/${order.id}`}
                    className="btn-secondary text-xs py-1.5 px-3"
                  >
                    View Details
                  </Link>
                  {order.status === 'pending' && (
                    <button
                      onClick={() => cancelMutation.mutate(order.id)}
                      disabled={cancelMutation.isPending}
                      className="text-xs px-3 py-1.5 border border-red-200 text-red-600 rounded-lg hover:bg-red-50 transition-colors disabled:opacity-50"
                    >
                      Cancel Order
                    </button>
                  )}
                </div>
              </div>
            ))}
          </div>
          {data.pagination && (
            <Pagination
              pagination={data.pagination}
              onPageChange={(page) => setParams((p) => ({ ...p, page }))}
            />
          )}
        </>
      )}
    </div>
  )
}
