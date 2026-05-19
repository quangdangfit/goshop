import { useEffect, useRef } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import toast from 'react-hot-toast'
import { setTokens } from '@/api/client'
import { useAuth } from '@/context/AuthContext'

/**
 * AuthCallbackPage receives the redirect from the GoShop API after a
 * successful OIDC login at Authentik. The backend appends the freshly minted
 * GoShop JWT pair to the query string; we persist it and let AuthContext load
 * the profile.
 */
export default function AuthCallbackPage() {
  const navigate = useNavigate()
  const [params] = useSearchParams()
  const { refreshUser } = useAuth()
  // StrictMode mounts effects twice in dev — guard against double-handling.
  const handled = useRef(false)

  useEffect(() => {
    if (handled.current) return
    handled.current = true

    const accessToken = params.get('access_token')
    const refreshToken = params.get('refresh_token')

    if (!accessToken || !refreshToken) {
      toast.error('Sign-in failed: missing tokens')
      navigate('/login', { replace: true })
      return
    }

    setTokens(accessToken, refreshToken)
    refreshUser()
      .then(() => {
        toast.success('Welcome!')
        navigate('/', { replace: true })
      })
      .catch(() => {
        toast.error('Sign-in failed')
        navigate('/login', { replace: true })
      })
  }, [params, navigate, refreshUser])

  return (
    <div className="min-h-screen flex items-center justify-center">
      <p className="text-gray-500">Signing you in…</p>
    </div>
  )
}
