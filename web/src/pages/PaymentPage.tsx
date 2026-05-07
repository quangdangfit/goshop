import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import toast from 'react-hot-toast'
import { CreditCard, Loader2 } from 'lucide-react'

import { paymentsApi } from '@/api/payments'
import { ordersApi } from '@/api/orders'
import LoadingSpinner from '@/components/LoadingSpinner'

// PaymentPage hosts the Stripe PaymentElement after an order has been placed but before
// payment clears. We deliberately delay loading @stripe/stripe-js until the publishable key
// is fetched so the bundle stays small for users who never reach checkout.

interface StripeLike {
  elements(opts: { clientSecret: string; appearance?: object }): StripeElementsLike
  confirmPayment(args: {
    elements: StripeElementsLike
    confirmParams: { return_url: string }
    redirect?: 'if_required'
  }): Promise<{ error?: { message?: string }; paymentIntent?: { status: string } }>
}

interface StripeElementsLike {
  create(type: string): { mount(selector: string): void }
  submit(): Promise<{ error?: { message?: string } }>
}

declare global {
  interface Window {
    Stripe?: (key: string) => StripeLike
  }
}

async function loadStripeJS(): Promise<void> {
  if (window.Stripe) return
  await new Promise<void>((resolve, reject) => {
    const s = document.createElement('script')
    s.src = 'https://js.stripe.com/v3/'
    s.onload = () => resolve()
    s.onerror = () => reject(new Error('failed to load Stripe.js'))
    document.head.appendChild(s)
  })
}

export default function PaymentPage() {
  const { orderId } = useParams<{ orderId: string }>()
  const navigate = useNavigate()
  const [stripe, setStripe] = useState<StripeLike | null>(null)
  const [elements, setElements] = useState<StripeElementsLike | null>(null)
  const [submitting, setSubmitting] = useState(false)

  const { data: cfg } = useQuery({
    queryKey: ['publicConfig'],
    queryFn: paymentsApi.publicConfig,
  })

  const { data: intent, isLoading: intentLoading } = useQuery({
    queryKey: ['payment-intent', orderId],
    queryFn: () => paymentsApi.createIntent(orderId!),
    enabled: !!orderId,
  })

  // Reservation TTL is 15 minutes from order creation. We approximate from order.created_at;
  // when expired, the sweeper will cancel server-side, but we surface the deadline immediately
  // so the customer knows why the payment form may stop accepting input.
  const [now, setNow] = useState(() => Date.now())
  useEffect(() => {
    const t = window.setInterval(() => setNow(Date.now()), 1000)
    return () => window.clearInterval(t)
  }, [])

  // Poll order status so the UI flips to "paid" once the webhook lands server-side.
  const { data: order } = useQuery({
    queryKey: ['order', orderId],
    queryFn: () => ordersApi.getOrder(orderId!),
    enabled: !!orderId,
    refetchInterval: (q) => (q.state.data?.status === 'paid' ? false : 3000),
  })

  useEffect(() => {
    if (!cfg?.stripe_publishable_key || !intent?.client_secret) return
    let cancelled = false
    ;(async () => {
      try {
        await loadStripeJS()
        if (cancelled || !window.Stripe) return
        const s = window.Stripe(cfg.stripe_publishable_key)
        const els = s.elements({ clientSecret: intent.client_secret })
        const card = els.create('payment')
        card.mount('#payment-element')
        setStripe(s)
        setElements(els)
      } catch (e) {
        toast.error((e as Error).message)
      }
    })()
    return () => {
      cancelled = true
    }
  }, [cfg?.stripe_publishable_key, intent?.client_secret])

  useEffect(() => {
    if (order?.status === 'paid') {
      toast.success('Payment confirmed')
      navigate(`/orders/${orderId}`, { replace: true })
    }
  }, [order?.status, orderId, navigate])

  const handlePay = async () => {
    if (!stripe || !elements) return
    setSubmitting(true)
    const submitResult = await elements.submit()
    if (submitResult.error) {
      toast.error(submitResult.error.message || 'payment input invalid')
      setSubmitting(false)
      return
    }
    const result = await stripe.confirmPayment({
      elements,
      confirmParams: { return_url: window.location.origin + `/orders/${orderId}` },
      redirect: 'if_required',
    })
    if (result.error) {
      toast.error(result.error.message || 'payment failed')
      setSubmitting(false)
      return
    }
    // PaymentIntent typically returns "succeeded" or "processing"; webhook confirms server-side.
    toast.success('Payment submitted, waiting for confirmation…')
  }

  if (intentLoading) return <LoadingSpinner className="min-h-[400px]" size="lg" />

  return (
    <div className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="text-2xl font-bold text-gray-900 mb-2 flex items-center gap-2">
        <CreditCard className="h-6 w-6 text-primary-600" />
        Complete Payment
      </h1>
      <p className="text-sm text-gray-500 mb-6">
        Order #{orderId} — {(() => {
          if (!order?.created_at) return 'reserved for 15 minutes.'
          const deadline = new Date(order.created_at).getTime() + 15 * 60 * 1000
          const remainingMs = deadline - now
          if (remainingMs <= 0) return 'reservation expired — order will be cancelled shortly.'
          const m = Math.floor(remainingMs / 60000)
          const s = Math.floor((remainingMs % 60000) / 1000)
          return `reserved for ${m}:${String(s).padStart(2, '0')}.`
        })()}
      </p>

      <div className="card">
        <div id="payment-element" className="min-h-[200px]" />
        <button
          onClick={handlePay}
          disabled={!stripe || submitting || order?.status === 'paid'}
          className="btn-primary w-full mt-6 py-3"
        >
          {submitting ? (
            <>
              <Loader2 className="h-4 w-4 animate-spin" />
              Processing…
            </>
          ) : (
            <>Pay ${((intent?.amount ?? 0) / 100).toFixed(2)}</>
          )}
        </button>
      </div>

      {order?.status === 'payment_failed' && (
        <p className="mt-4 text-sm text-red-600">
          Payment failed. The order will be auto-cancelled and stock released shortly.
        </p>
      )}
    </div>
  )
}
