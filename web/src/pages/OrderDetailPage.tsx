import { ArrowLeft, Package, ShoppingCart } from 'lucide-react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import toast from 'react-hot-toast'
import { ordersApi } from '@/api/orders'
import LoadingSpinner from '@/components/LoadingSpinner'

const STATUS_COLORS: Record<string, string> = {
  pending: 'badge-warning',
  processing: 'badge-info',
  shipped: 'badge-info',
  delivered: 'badge-success',
  cancelled: 'badge-danger',
}

const STATUS_STEPS = ['pending', 'processing', 'shipped', 'delivered']

export default function OrderDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const { data: order, isLoading } = useQuery({
    queryKey: ['order', id],
    queryFn: () => ordersApi.getOrder(id!),
    enabled: !!id,
  })

  const cancelMutation = useMutation({
    mutationFn: () => ordersApi.cancelOrder(id!),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['order', id] })
      queryClient.invalidateQueries({ queryKey: ['orders'] })
      toast.success('Order cancelled')
    },
    onError: () => toast.error('Failed to cancel order'),
  })

  if (isLoading) return <LoadingSpinner className="min-h-[400px]" size="lg" />

  if (!order) {
    return (
      <div className="max-w-4xl mx-auto px-4 py-12 text-center">
        <p className="text-gray-500">Order not found</p>
        <button onClick={() => navigate('/orders')} className="btn-primary mt-4">
          Back to Orders
        </button>
      </div>
    )
  }

  const currentStep = STATUS_STEPS.indexOf(order.status)

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <Link
        to="/orders"
        className="inline-flex items-center gap-1.5 text-sm text-gray-500 hover:text-primary-600 mb-6 transition-colors"
      >
        <ArrowLeft className="h-4 w-4" />
        Back to Orders
      </Link>

      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Order #{order.code}</h1>
          <p className="text-sm text-gray-500 mt-0.5">
            Placed on{' '}
            {new Date(order.created_at).toLocaleDateString('en-US', {
              year: 'numeric',
              month: 'long',
              day: 'numeric',
            })}
          </p>
        </div>
        <span className={`badge text-sm ${STATUS_COLORS[order.status] || 'badge-default'}`}>
          {order.status}
        </span>
      </div>

      {/* Progress Steps */}
      {order.status !== 'cancelled' && (
        <div className="card mb-5">
          <div className="flex items-center justify-between">
            {STATUS_STEPS.map((step, idx) => (
              <div key={step} className="flex items-center flex-1">
                <div className="flex flex-col items-center">
                  <div
                    className={`h-8 w-8 rounded-full flex items-center justify-center text-xs font-bold transition-colors ${
                      idx <= currentStep
                        ? 'bg-primary-600 text-white'
                        : 'bg-gray-100 text-gray-400'
                    }`}
                  >
                    {idx < currentStep ? '✓' : idx + 1}
                  </div>
                  <span className="text-xs text-gray-500 mt-1 capitalize">{step}</span>
                </div>
                {idx < STATUS_STEPS.length - 1 && (
                  <div
                    className={`flex-1 h-0.5 mx-2 ${
                      idx < currentStep ? 'bg-primary-600' : 'bg-gray-200'
                    }`}
                  />
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-5">
        <div className="lg:col-span-2 space-y-5">
          {/* Order Items */}
          <div className="card">
            <h2 className="font-bold text-gray-900 mb-4">Order Items</h2>
            <div className="space-y-3">
              {order.lines?.map((line) => (
                <div
                  key={line.id}
                  className="flex items-center gap-3 py-2 border-b border-gray-50 last:border-0"
                >
                  <div className="h-14 w-14 bg-gray-100 rounded-lg flex-shrink-0 flex items-center justify-center">
                    <ShoppingCart className="h-6 w-6 text-gray-300" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <Link
                      to={`/products/${line.product_id}`}
                      className="text-sm font-medium text-gray-900 hover:text-primary-600 line-clamp-1"
                    >
                      {line.product?.name}
                    </Link>
                    <p className="text-xs text-gray-500">Qty: {line.quantity}</p>
                  </div>
                  <p className="font-semibold text-sm text-gray-900">
                    ${((line.price || line.product?.price || 0) * line.quantity).toFixed(2)}
                  </p>
                </div>
              ))}
            </div>
          </div>
        </div>

        <div className="space-y-5">
          {/* Summary */}
          <div className="card">
            <h2 className="font-bold text-gray-900 mb-4">Summary</h2>
            <div className="space-y-2 text-sm">
              <div className="flex justify-between text-gray-600">
                <span>Subtotal</span>
                <span>${order.total_price?.toFixed(2)}</span>
              </div>
              {order.coupon_code && (
                <div className="flex justify-between text-green-600">
                  <span>Coupon: {order.coupon_code}</span>
                  <span>Applied</span>
                </div>
              )}
              <div className="border-t border-gray-100 pt-2 flex justify-between font-bold text-gray-900">
                <span>Total</span>
                <span>${order.total_price?.toFixed(2)}</span>
              </div>
            </div>
          </div>

          {/* Actions */}
          {order.status === 'pending' && (
            <div className="card">
              <h2 className="font-bold text-gray-900 mb-3">Actions</h2>
              <button
                onClick={() => cancelMutation.mutate()}
                disabled={cancelMutation.isPending}
                className="btn-danger w-full"
              >
                <Package className="h-4 w-4" />
                {cancelMutation.isPending ? 'Cancelling...' : 'Cancel Order'}
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
