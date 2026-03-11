import { CheckCircle, MapPin, Plus, Tag, X } from 'lucide-react'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import toast from 'react-hot-toast'
import { cartApi } from '@/api/cart'
import { ordersApi } from '@/api/orders'
import { addressesApi } from '@/api/addresses'
import { couponsApi } from '@/api/coupons'
import LoadingSpinner from '@/components/LoadingSpinner'
import type { Address, Coupon } from '@/types'

export default function CheckoutPage() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [selectedAddressId, setSelectedAddressId] = useState<string | null>(null)
  const [couponCode, setCouponCode] = useState('')
  const [appliedCoupon, setAppliedCoupon] = useState<Coupon | null>(null)
  const [couponError, setCouponError] = useState('')
  const [showAddAddress, setShowAddAddress] = useState(false)
  const [newAddress, setNewAddress] = useState({
    name: '', phone: '', street: '', city: '', country: '',
  })

  const { data: cart, isLoading: cartLoading } = useQuery({
    queryKey: ['cart'],
    queryFn: cartApi.getCart,
  })

  const { data: addresses, isLoading: addressLoading } = useQuery({
    queryKey: ['addresses'],
    queryFn: addressesApi.getAddresses,
  })

  const createAddressMutation = useMutation({
    mutationFn: addressesApi.createAddress,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['addresses'] })
      setShowAddAddress(false)
      setNewAddress({ name: '', phone: '', street: '', city: '', country: '' })
      toast.success('Address added!')
    },
    onError: () => toast.error('Failed to add address'),
  })

  const placeOrderMutation = useMutation({
    mutationFn: ordersApi.createOrder,
    onSuccess: (order) => {
      queryClient.invalidateQueries({ queryKey: ['cart'] })
      queryClient.invalidateQueries({ queryKey: ['orders'] })
      toast.success('Order placed successfully!')
      navigate(`/orders/${order.id}`)
    },
    onError: () => toast.error('Failed to place order'),
  })

  const applyCoupon = async () => {
    if (!couponCode.trim()) return
    setCouponError('')
    try {
      const coupon = await couponsApi.getCoupon(couponCode.trim())
      setAppliedCoupon(coupon)
      toast.success(`Coupon "${coupon.code}" applied!`)
    } catch {
      setCouponError('Invalid or expired coupon code')
      setAppliedCoupon(null)
    }
  }

  const lines = cart?.lines || []
  const subtotal = lines.reduce(
    (sum, line) => sum + (line.product?.price || 0) * line.quantity,
    0
  )

  const discount = appliedCoupon
    ? appliedCoupon.discount_type === 'percentage'
      ? subtotal * (appliedCoupon.discount_value / 100)
      : appliedCoupon.discount_value
    : 0

  const shipping = subtotal >= 50 ? 0 : 5.99
  const total = Math.max(0, subtotal - discount) + shipping

  const handlePlaceOrder = () => {
    if (lines.length === 0) {
      toast.error('Your cart is empty')
      return
    }
    placeOrderMutation.mutate({
      coupon_code: appliedCoupon?.code,
      lines: lines.map((l) => ({
        product_id: l.product_id,
        quantity: l.quantity,
      })),
    })
  }

  if (cartLoading || addressLoading) {
    return <LoadingSpinner className="min-h-[400px]" size="lg" />
  }

  const defaultAddress = addresses?.find((a: Address) => a.is_default)
  const currentAddress = addresses?.find((a: Address) => a.id === selectedAddressId) || defaultAddress

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Checkout</h1>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-5">
          {/* Address Selection */}
          <div className="card">
            <h2 className="font-bold text-gray-900 mb-4 flex items-center gap-2">
              <MapPin className="h-5 w-5 text-primary-600" />
              Delivery Address
            </h2>

            {addresses && addresses.length > 0 ? (
              <div className="space-y-2 mb-4">
                {addresses.map((addr: Address) => (
                  <label
                    key={addr.id}
                    className={`flex items-start gap-3 p-3 rounded-xl border cursor-pointer transition-colors ${
                      (selectedAddressId === addr.id ||
                        (!selectedAddressId && addr.is_default))
                        ? 'border-primary-500 bg-primary-50'
                        : 'border-gray-200 hover:border-gray-300'
                    }`}
                  >
                    <input
                      type="radio"
                      name="address"
                      value={addr.id}
                      checked={
                        selectedAddressId === addr.id ||
                        (!selectedAddressId && addr.is_default)
                      }
                      onChange={() => setSelectedAddressId(addr.id)}
                      className="mt-0.5 text-primary-600"
                    />
                    <div>
                      <p className="font-medium text-sm text-gray-900">
                        {addr.name}
                        {addr.is_default && (
                          <span className="ml-2 badge badge-info">Default</span>
                        )}
                      </p>
                      <p className="text-sm text-gray-500">{addr.phone}</p>
                      <p className="text-sm text-gray-500">
                        {addr.street}, {addr.city}, {addr.country}
                      </p>
                    </div>
                  </label>
                ))}
              </div>
            ) : (
              <p className="text-sm text-gray-500 mb-3">No addresses saved. Add one below.</p>
            )}

            <button
              onClick={() => setShowAddAddress(!showAddAddress)}
              className="flex items-center gap-1 text-sm text-primary-600 hover:text-primary-700 font-medium"
            >
              <Plus className="h-4 w-4" />
              Add new address
            </button>

            {showAddAddress && (
              <div className="mt-4 p-4 bg-gray-50 rounded-xl space-y-3">
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <label className="label">Full Name</label>
                    <input
                      className="input"
                      value={newAddress.name}
                      onChange={(e) => setNewAddress((p) => ({ ...p, name: e.target.value }))}
                    />
                  </div>
                  <div>
                    <label className="label">Phone</label>
                    <input
                      className="input"
                      value={newAddress.phone}
                      onChange={(e) => setNewAddress((p) => ({ ...p, phone: e.target.value }))}
                    />
                  </div>
                </div>
                <div>
                  <label className="label">Street Address</label>
                  <input
                    className="input"
                    value={newAddress.street}
                    onChange={(e) => setNewAddress((p) => ({ ...p, street: e.target.value }))}
                  />
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <label className="label">City</label>
                    <input
                      className="input"
                      value={newAddress.city}
                      onChange={(e) => setNewAddress((p) => ({ ...p, city: e.target.value }))}
                    />
                  </div>
                  <div>
                    <label className="label">Country</label>
                    <input
                      className="input"
                      value={newAddress.country}
                      onChange={(e) => setNewAddress((p) => ({ ...p, country: e.target.value }))}
                    />
                  </div>
                </div>
                <div className="flex gap-2">
                  <button
                    onClick={() => createAddressMutation.mutate(newAddress)}
                    disabled={createAddressMutation.isPending}
                    className="btn-primary text-sm"
                  >
                    Save Address
                  </button>
                  <button
                    onClick={() => setShowAddAddress(false)}
                    className="btn-secondary text-sm"
                  >
                    Cancel
                  </button>
                </div>
              </div>
            )}
          </div>

          {/* Coupon */}
          <div className="card">
            <h2 className="font-bold text-gray-900 mb-4 flex items-center gap-2">
              <Tag className="h-5 w-5 text-primary-600" />
              Coupon Code
            </h2>

            {appliedCoupon ? (
              <div className="flex items-center justify-between bg-green-50 border border-green-200 rounded-xl px-4 py-3">
                <div>
                  <p className="font-medium text-green-800 text-sm">{appliedCoupon.code}</p>
                  <p className="text-xs text-green-600">
                    {appliedCoupon.discount_type === 'percentage'
                      ? `${appliedCoupon.discount_value}% off`
                      : `$${appliedCoupon.discount_value} off`}
                  </p>
                </div>
                <button
                  onClick={() => {
                    setAppliedCoupon(null)
                    setCouponCode('')
                  }}
                  className="text-green-600 hover:text-green-800"
                >
                  <X className="h-4 w-4" />
                </button>
              </div>
            ) : (
              <div className="flex gap-2">
                <input
                  type="text"
                  placeholder="Enter coupon code"
                  value={couponCode}
                  onChange={(e) => setCouponCode(e.target.value.toUpperCase())}
                  onKeyDown={(e) => e.key === 'Enter' && applyCoupon()}
                  className="input flex-1"
                />
                <button onClick={applyCoupon} className="btn-secondary">
                  Apply
                </button>
              </div>
            )}
            {couponError && (
              <p className="mt-1.5 text-xs text-red-500">{couponError}</p>
            )}
          </div>

          {/* Order Items */}
          <div className="card">
            <h2 className="font-bold text-gray-900 mb-4">Order Items</h2>
            <div className="space-y-3">
              {lines.map((line) => (
                <div key={line.id} className="flex justify-between text-sm py-2 border-b border-gray-50 last:border-0">
                  <div>
                    <p className="font-medium text-gray-900">{line.product?.name}</p>
                    <p className="text-gray-500">Qty: {line.quantity}</p>
                  </div>
                  <span className="font-medium text-gray-900">
                    ${((line.product?.price || 0) * line.quantity).toFixed(2)}
                  </span>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Summary */}
        <div>
          <div className="card sticky top-20">
            <h2 className="font-bold text-gray-900 mb-4">Order Summary</h2>

            {currentAddress && (
              <div className="bg-gray-50 rounded-lg p-3 mb-4 text-sm">
                <p className="font-medium text-gray-700 flex items-center gap-1 mb-1">
                  <MapPin className="h-3.5 w-3.5" />
                  Delivery to:
                </p>
                <p className="text-gray-600">{currentAddress.name}</p>
                <p className="text-gray-500">
                  {currentAddress.street}, {currentAddress.city}
                </p>
              </div>
            )}

            <div className="space-y-2 text-sm mb-5">
              <div className="flex justify-between text-gray-600">
                <span>Subtotal</span>
                <span>${subtotal.toFixed(2)}</span>
              </div>
              {discount > 0 && (
                <div className="flex justify-between text-green-600">
                  <span>Discount</span>
                  <span>-${discount.toFixed(2)}</span>
                </div>
              )}
              <div className="flex justify-between text-gray-600">
                <span>Shipping</span>
                <span className={shipping === 0 ? 'text-green-600' : ''}>
                  {shipping === 0 ? 'Free' : `$${shipping.toFixed(2)}`}
                </span>
              </div>
              <div className="border-t border-gray-100 pt-2 flex justify-between font-bold text-gray-900 text-base">
                <span>Total</span>
                <span>${total.toFixed(2)}</span>
              </div>
            </div>

            <button
              onClick={handlePlaceOrder}
              disabled={placeOrderMutation.isPending || lines.length === 0}
              className="btn-primary w-full py-3"
            >
              <CheckCircle className="h-5 w-5" />
              {placeOrderMutation.isPending ? 'Placing Order...' : 'Place Order'}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
