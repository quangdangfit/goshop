import { Plus, X } from 'lucide-react'
import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import toast from 'react-hot-toast'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { couponsApi } from '@/api/coupons'
import type { Coupon } from '@/types'

const schema = z.object({
  code: z.string().min(3, 'Code must be at least 3 characters').toUpperCase(),
  discount_type: z.enum(['fixed', 'percentage']),
  discount_value: z.coerce.number().positive('Must be positive'),
  min_order_amount: z.coerce.number().min(0, 'Must be 0 or more'),
  max_usage: z.coerce.number().int().positive('Must be positive'),
  expires_at: z.string().min(1, 'Expiry date is required'),
})

type FormData = z.infer<typeof schema>

export default function AdminCouponsPage() {
  const queryClient = useQueryClient()
  const [showModal, setShowModal] = useState(false)
  const [createdCoupons, setCreatedCoupons] = useState<Coupon[]>([])
  const [lookupCode, setLookupCode] = useState('')
  const [lookedUpCoupon, setLookedUpCoupon] = useState<Coupon | null>(null)
  const [lookupError, setLookupError] = useState('')

  const {
    register,
    handleSubmit,
    reset,
    watch,
    formState: { errors, isSubmitting },
  } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: {
      discount_type: 'percentage',
      max_usage: 100,
      min_order_amount: 0,
    },
  })

  const discountType = watch('discount_type')

  const createMutation = useMutation({
    mutationFn: couponsApi.createCoupon,
    onSuccess: (coupon) => {
      queryClient.invalidateQueries({ queryKey: ['coupons'] })
      setCreatedCoupons((prev) => [coupon, ...prev])
      closeModal()
      toast.success(`Coupon "${coupon.code}" created!`)
    },
    onError: () => toast.error('Failed to create coupon'),
  })

  const closeModal = () => {
    setShowModal(false)
    reset()
  }

  const onSubmit = (data: FormData) => {
    createMutation.mutate({
      ...data,
      expires_at: new Date(data.expires_at).toISOString(),
    })
  }

  const handleLookup = async () => {
    if (!lookupCode.trim()) return
    setLookupError('')
    setLookedUpCoupon(null)
    try {
      const coupon = await couponsApi.getCoupon(lookupCode.trim().toUpperCase())
      setLookedUpCoupon(coupon)
    } catch {
      setLookupError('Coupon not found')
    }
  }

  const isExpired = (expiresAt: string) => new Date(expiresAt) < new Date()

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Coupons</h1>
          <p className="text-sm text-gray-500">Manage discount coupons</p>
        </div>
        <button onClick={() => setShowModal(true)} className="btn-primary">
          <Plus className="h-4 w-4" />
          Create Coupon
        </button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Lookup */}
        <div className="bg-white rounded-xl border border-gray-100 p-5">
          <h2 className="font-bold text-gray-900 mb-4">Lookup Coupon</h2>
          <div className="flex gap-2 mb-3">
            <input
              type="text"
              placeholder="Enter coupon code"
              value={lookupCode}
              onChange={(e) => setLookupCode(e.target.value.toUpperCase())}
              onKeyDown={(e) => e.key === 'Enter' && handleLookup()}
              className="input flex-1"
            />
            <button onClick={handleLookup} className="btn-secondary">
              Lookup
            </button>
          </div>
          {lookupError && (
            <p className="text-sm text-red-500">{lookupError}</p>
          )}
          {lookedUpCoupon && (
            <CouponCard coupon={lookedUpCoupon} isExpired={isExpired(lookedUpCoupon.expires_at)} />
          )}
        </div>

        {/* Recently created */}
        <div className="bg-white rounded-xl border border-gray-100 p-5">
          <h2 className="font-bold text-gray-900 mb-4">Recently Created</h2>
          {createdCoupons.length === 0 ? (
            <p className="text-sm text-gray-400 text-center py-6">
              No coupons created in this session
            </p>
          ) : (
            <div className="space-y-3">
              {createdCoupons.map((coupon) => (
                <CouponCard
                  key={coupon.id}
                  coupon={coupon}
                  isExpired={isExpired(coupon.expires_at)}
                />
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Create Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-2xl shadow-xl w-full max-w-md max-h-[90vh] overflow-y-auto">
            <div className="flex items-center justify-between p-5 border-b border-gray-100">
              <h2 className="font-bold text-gray-900">Create Coupon</h2>
              <button onClick={closeModal} className="p-1 text-gray-400 hover:text-gray-600">
                <X className="h-5 w-5" />
              </button>
            </div>

            <form onSubmit={handleSubmit(onSubmit)} className="p-5 space-y-4">
              <div>
                <label className="label">Coupon Code</label>
                <input
                  {...register('code')}
                  className="input font-mono uppercase"
                  placeholder="SAVE20"
                  style={{ textTransform: 'uppercase' }}
                />
                {errors.code && (
                  <p className="mt-1 text-xs text-red-500">{errors.code.message}</p>
                )}
              </div>

              <div>
                <label className="label">Discount Type</label>
                <div className="flex gap-3">
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      {...register('discount_type')}
                      type="radio"
                      value="percentage"
                      className="text-primary-600"
                    />
                    <span className="text-sm">Percentage (%)</span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      {...register('discount_type')}
                      type="radio"
                      value="fixed"
                      className="text-primary-600"
                    />
                    <span className="text-sm">Fixed ($)</span>
                  </label>
                </div>
              </div>

              <div>
                <label className="label">
                  Discount Value ({discountType === 'percentage' ? '%' : '$'})
                </label>
                <input
                  {...register('discount_value')}
                  type="number"
                  step={discountType === 'percentage' ? '1' : '0.01'}
                  min="0"
                  max={discountType === 'percentage' ? '100' : undefined}
                  className="input"
                  placeholder={discountType === 'percentage' ? '20' : '10.00'}
                />
                {errors.discount_value && (
                  <p className="mt-1 text-xs text-red-500">{errors.discount_value.message}</p>
                )}
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="label">Min Order Amount ($)</label>
                  <input
                    {...register('min_order_amount')}
                    type="number"
                    min="0"
                    step="0.01"
                    className="input"
                    placeholder="0"
                  />
                  {errors.min_order_amount && (
                    <p className="mt-1 text-xs text-red-500">{errors.min_order_amount.message}</p>
                  )}
                </div>
                <div>
                  <label className="label">Max Usage</label>
                  <input
                    {...register('max_usage')}
                    type="number"
                    min="1"
                    className="input"
                    placeholder="100"
                  />
                  {errors.max_usage && (
                    <p className="mt-1 text-xs text-red-500">{errors.max_usage.message}</p>
                  )}
                </div>
              </div>

              <div>
                <label className="label">Expires At</label>
                <input
                  {...register('expires_at')}
                  type="datetime-local"
                  className="input"
                  min={new Date().toISOString().slice(0, 16)}
                />
                {errors.expires_at && (
                  <p className="mt-1 text-xs text-red-500">{errors.expires_at.message}</p>
                )}
              </div>

              <div className="flex gap-2 pt-2">
                <button
                  type="submit"
                  disabled={isSubmitting || createMutation.isPending}
                  className="btn-primary flex-1"
                >
                  {isSubmitting || createMutation.isPending ? 'Creating...' : 'Create Coupon'}
                </button>
                <button type="button" onClick={closeModal} className="btn-secondary">
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}

function CouponCard({ coupon, isExpired }: { coupon: Coupon; isExpired: boolean }) {
  return (
    <div
      className={`rounded-xl border p-4 ${
        isExpired ? 'border-gray-200 bg-gray-50 opacity-60' : 'border-primary-200 bg-primary-50'
      }`}
    >
      <div className="flex items-center justify-between mb-2">
        <code className="font-bold text-lg font-mono text-primary-700">{coupon.code}</code>
        <span
          className={`badge ${
            isExpired ? 'badge-danger' : coupon.used_count >= coupon.max_usage ? 'badge-warning' : 'badge-success'
          }`}
        >
          {isExpired ? 'Expired' : coupon.used_count >= coupon.max_usage ? 'Exhausted' : 'Active'}
        </span>
      </div>
      <div className="grid grid-cols-2 gap-2 text-xs text-gray-600">
        <div>
          <span className="font-medium">Discount: </span>
          {coupon.discount_type === 'percentage'
            ? `${coupon.discount_value}%`
            : `$${coupon.discount_value}`}
        </div>
        <div>
          <span className="font-medium">Usage: </span>
          {coupon.used_count}/{coupon.max_usage}
        </div>
        <div>
          <span className="font-medium">Min Order: </span>
          ${coupon.min_order_amount}
        </div>
        <div>
          <span className="font-medium">Expires: </span>
          {new Date(coupon.expires_at).toLocaleDateString()}
        </div>
      </div>
    </div>
  )
}
