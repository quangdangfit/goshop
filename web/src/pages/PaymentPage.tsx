import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import toast from 'react-hot-toast'
import { Clock, CreditCard, Loader2 } from 'lucide-react'

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
  const queryClient = useQueryClient()
  const [stripe, setStripe] = useState<StripeLike | null>(null)
  const [elements, setElements] = useState<StripeElementsLike | null>(null)
  const [submitting, setSubmitting] = useState(false)
  const [awaitingWebhook, setAwaitingWebhook] = useState(false)
  const [webhookWaitStart, setWebhookWaitStart] = useState<number | null>(null)

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
      // Drop any stale list cache so the orders page reflects the new paid order on next mount.
      queryClient.invalidateQueries({ queryKey: ['orders'] })
      navigate(`/orders/${orderId}`, { replace: true })
    }
  }, [order?.status, orderId, navigate, queryClient])

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
    setSubmitting(false)
    setAwaitingWebhook(true)
    setWebhookWaitStart(Date.now())
  }

  // If we've been waiting on the webhook for too long, the most common cause locally is
  // that `stripe listen --forward-to localhost:8888/api/v1/webhooks/stripe` isn't running.
  const webhookSlow =
    awaitingWebhook && webhookWaitStart !== null && now - webhookWaitStart > 20_000

  if (intentLoading) return <LoadingSpinner className="min-h-[400px]" size="lg" />

  const deadlineMs = order?.created_at
    ? new Date(order.created_at).getTime() + 15 * 60 * 1000
    : null
  const remainingMs = deadlineMs ? deadlineMs - now : null
  const expired = remainingMs !== null && remainingMs <= 0
  const warning = remainingMs !== null && remainingMs > 0 && remainingMs < 2 * 60 * 1000

  let timerLabel = 'Reserved for 15 minutes'
  if (remainingMs !== null) {
    if (expired) {
      timerLabel = 'Reservation expired — order will be cancelled shortly'
    } else {
      const m = Math.floor(remainingMs / 60000)
      const s = Math.floor((remainingMs % 60000) / 1000)
      timerLabel = `Reserved for ${m}:${String(s).padStart(2, '0')}`
    }
  }
  const timerStyle = expired
    ? 'bg-red-50 border-red-200 text-red-700'
    : warning
      ? 'bg-amber-50 border-amber-200 text-amber-700'
      : 'bg-primary-50 border-primary-200 text-primary-700'

  return (
    <div className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="text-2xl font-bold text-gray-900 mb-2 flex items-center gap-2">
        <CreditCard className="h-6 w-6 text-primary-600" />
        Complete Payment
      </h1>
      <p className="text-sm text-gray-500 mb-4">Order #{orderId}</p>

      <div className={`flex items-center gap-2 border rounded-lg px-4 py-3 mb-6 text-sm font-medium ${timerStyle}`}>
        <Clock className="h-4 w-4" />
        <span>{timerLabel}</span>
      </div>

      <div className="card">
        <div id="payment-element" className="min-h-[200px]" />
        <button
          onClick={handlePay}
          disabled={!stripe || submitting || awaitingWebhook || expired || order?.status === 'paid'}
          className="btn-primary w-full mt-6 py-3"
        >
          {submitting ? (
            <>
              <Loader2 className="h-4 w-4 animate-spin" />
              Processing…
            </>
          ) : awaitingWebhook ? (
            <>
              <Loader2 className="h-4 w-4 animate-spin" />
              Waiting for confirmation…
            </>
          ) : (
            <>Pay ${((intent?.amount ?? 0) / 100).toFixed(2)}</>
          )}
        </button>
      </div>

      {webhookSlow && order?.status !== 'paid' && (
        <p className="mt-4 text-sm text-amber-700">
          Still waiting for Stripe to confirm the payment. In local dev, make sure
          <code className="mx-1 px-1 py-0.5 bg-amber-100 rounded">stripe listen --forward-to localhost:8888/api/v1/webhooks/stripe</code>
          is running so webhooks reach the backend.
        </p>
      )}

      {order?.status === 'payment_failed' && (
        <p className="mt-4 text-sm text-red-600">
          Payment failed. The order will be auto-cancelled and stock released shortly.
        </p>
      )}
    </div>
  )
}
