import { Minus, Plus, ShoppingCart, Trash2, ArrowRight } from 'lucide-react'
import { Link, useNavigate } from 'react-router-dom'
import toast from 'react-hot-toast'

import { useCart } from '@/hooks/useCart'

export default function CartPage() {
  const navigate = useNavigate()
  const { cart, store } = useCart()
  const lines = cart.items

  const subtotal = lines.reduce(
    (sum, line) => sum + (line.snapshot?.price || 0) * line.quantity,
    0,
  )

  const handleRemove = (productId: string) => {
    store.remove(productId)
    toast.success('Item removed')
  }

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
        <div className="lg:col-span-2 space-y-3">
          {lines.map((line) => (
            <div
              key={line.product_id}
              className="bg-white rounded-xl border border-gray-100 p-4 flex gap-4"
            >
              <div className="w-20 h-20 bg-gray-100 rounded-lg flex-shrink-0 overflow-hidden flex items-center justify-center">
                {line.snapshot?.images?.[0] ? (
                  <img
                    src={line.snapshot.images[0]}
                    alt={line.snapshot.name}
                    className="w-full h-full object-cover"
                  />
                ) : (
                  <ShoppingCart className="h-8 w-8 text-gray-300" />
                )}
              </div>

              <div className="flex-1 min-w-0">
                <Link
                  to={`/products/${line.product_id}`}
                  className="font-medium text-gray-900 hover:text-primary-600 line-clamp-2 text-sm"
                >
                  {line.snapshot?.name}
                </Link>
                <p className="text-sm text-gray-500 mt-0.5">
                  ${line.snapshot?.price?.toFixed(2)} each
                </p>

                <div className="flex items-center justify-between mt-2">
                  <div className="flex items-center border border-gray-200 rounded-lg overflow-hidden">
                    <button
                      onClick={() => store.setQuantity(line.product_id, line.quantity - 1)}
                      className="p-1.5 hover:bg-gray-100 transition-colors"
                    >
                      <Minus className="h-3.5 w-3.5" />
                    </button>
                    <span className="w-8 text-center text-sm font-medium">{line.quantity}</span>
                    <button
                      onClick={() => store.setQuantity(line.product_id, line.quantity + 1)}
                      className="p-1.5 hover:bg-gray-100 transition-colors"
                    >
                      <Plus className="h-3.5 w-3.5" />
                    </button>
                  </div>

                  <div className="flex items-center gap-3">
                    <span className="font-semibold text-gray-900">
                      ${((line.snapshot?.price || 0) * line.quantity).toFixed(2)}
                    </span>
                    <button
                      onClick={() => handleRemove(line.product_id)}
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
                <span className="text-green-600">{subtotal >= 50 ? 'Free' : '$5.99'}</span>
              </div>
              <div className="border-t border-gray-100 pt-2 flex justify-between font-bold text-gray-900 text-base">
                <span>Total</span>
                <span>${(subtotal >= 50 ? subtotal : subtotal + 5.99).toFixed(2)}</span>
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
