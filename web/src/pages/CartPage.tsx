import { Minus, Plus, ShoppingCart, Trash2, ArrowRight } from 'lucide-react'
import { Link, useNavigate } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import toast from 'react-hot-toast'
import { cartApi } from '@/api/cart'
import LoadingSpinner from '@/components/LoadingSpinner'

export default function CartPage() {
  const queryClient = useQueryClient()
  const navigate = useNavigate()

  const { data: cart, isLoading } = useQuery({
    queryKey: ['cart'],
    queryFn: cartApi.getCart,
  })

  const addMutation = useMutation({
    mutationFn: ({ productId, quantity }: { productId: string; quantity: number }) =>
      cartApi.addToCart({ product_id: productId, quantity }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['cart'] }),
    onError: () => toast.error('Failed to update cart'),
  })

  const removeMutation = useMutation({
    mutationFn: (productId: string) => cartApi.removeFromCart(productId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cart'] })
      toast.success('Item removed')
    },
    onError: () => toast.error('Failed to remove item'),
  })

  if (isLoading) return <LoadingSpinner className="min-h-[400px]" size="lg" />

  const lines = cart?.lines || []
  const subtotal = lines.reduce(
    (sum, line) => sum + (line.product?.price || 0) * line.quantity,
    0
  )

  if (lines.length === 0) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16 text-center">
        <ShoppingCart className="h-16 w-16 text-gray-300 mx-auto mb-4" />
        <h2 className="text-xl font-bold text-gray-900 mb-2">Your cart is empty</h2>
        <p className="text-gray-500 mb-6">Add some items to get started</p>
        <Link to="/products" className="btn-primary">
          Browse Products
        </Link>
      </div>
    )
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="text-2xl font-bold text-gray-900 mb-6">
        Shopping Cart ({lines.length} {lines.length === 1 ? 'item' : 'items'})
      </h1>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Cart Items */}
        <div className="lg:col-span-2 space-y-3">
          {lines.map((line) => (
            <div
              key={line.id}
              className="bg-white rounded-xl border border-gray-100 p-4 flex gap-4"
            >
              {/* Image */}
              <div className="w-20 h-20 bg-gray-100 rounded-lg flex-shrink-0 overflow-hidden flex items-center justify-center">
                {line.product?.images?.[0] ? (
                  <img
                    src={line.product.images[0]}
                    alt={line.product.name}
                    className="w-full h-full object-cover"
                  />
                ) : (
                  <ShoppingCart className="h-8 w-8 text-gray-300" />
                )}
              </div>

              {/* Info */}
              <div className="flex-1 min-w-0">
                <Link
                  to={`/products/${line.product_id}`}
                  className="font-medium text-gray-900 hover:text-primary-600 line-clamp-2 text-sm"
                >
                  {line.product?.name}
                </Link>
                <p className="text-sm text-gray-500 mt-0.5">
                  ${line.product?.price?.toFixed(2)} each
                </p>

                <div className="flex items-center justify-between mt-2">
                  {/* Quantity controls */}
                  <div className="flex items-center border border-gray-200 rounded-lg overflow-hidden">
                    <button
                      onClick={() => {
                        if (line.quantity <= 1) {
                          removeMutation.mutate(line.product_id)
                        } else {
                          addMutation.mutate({
                            productId: line.product_id,
                            quantity: line.quantity - 1,
                          })
                        }
                      }}
                      disabled={addMutation.isPending || removeMutation.isPending}
                      className="p-1.5 hover:bg-gray-100 transition-colors"
                    >
                      <Minus className="h-3.5 w-3.5" />
                    </button>
                    <span className="w-8 text-center text-sm font-medium">
                      {line.quantity}
                    </span>
                    <button
                      onClick={() =>
                        addMutation.mutate({
                          productId: line.product_id,
                          quantity: line.quantity + 1,
                        })
                      }
                      disabled={addMutation.isPending}
                      className="p-1.5 hover:bg-gray-100 transition-colors"
                    >
                      <Plus className="h-3.5 w-3.5" />
                    </button>
                  </div>

                  <div className="flex items-center gap-3">
                    <span className="font-semibold text-gray-900">
                      ${((line.product?.price || 0) * line.quantity).toFixed(2)}
                    </span>
                    <button
                      onClick={() => removeMutation.mutate(line.product_id)}
                      disabled={removeMutation.isPending}
                      className="p-1 text-gray-400 hover:text-red-500 transition-colors"
                    >
                      <Trash2 className="h-4 w-4" />
                    </button>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>

        {/* Order Summary */}
        <div>
          <div className="bg-white rounded-xl border border-gray-100 p-5 sticky top-20">
            <h2 className="font-bold text-gray-900 mb-4">Order Summary</h2>

            <div className="space-y-2 text-sm mb-4">
              <div className="flex justify-between text-gray-600">
                <span>Subtotal ({lines.length} items)</span>
                <span>${subtotal.toFixed(2)}</span>
              </div>
              <div className="flex justify-between text-gray-600">
                <span>Shipping</span>
                <span className="text-green-600">
                  {subtotal >= 50 ? 'Free' : '$5.99'}
                </span>
              </div>
              <div className="border-t border-gray-100 pt-2 flex justify-between font-bold text-gray-900 text-base">
                <span>Total</span>
                <span>
                  ${(subtotal >= 50 ? subtotal : subtotal + 5.99).toFixed(2)}
                </span>
              </div>
            </div>

            {subtotal < 50 && (
              <p className="text-xs text-gray-500 mb-3 bg-gray-50 px-3 py-2 rounded-lg">
                Add ${(50 - subtotal).toFixed(2)} more for free shipping!
              </p>
            )}

            <button
              onClick={() => navigate('/checkout')}
              className="btn-primary w-full py-3"
            >
              Proceed to Checkout
              <ArrowRight className="h-4 w-4" />
            </button>

            <Link
              to="/products"
              className="block text-center text-sm text-primary-600 hover:text-primary-700 mt-3"
            >
              Continue Shopping
            </Link>
          </div>
        </div>
      </div>
    </div>
  )
}
