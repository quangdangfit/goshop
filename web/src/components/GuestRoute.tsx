import { Navigate } from 'react-router-dom'
import { useAuth } from '@/context/AuthContext'
import LoadingSpinner from './LoadingSpinner'

interface GuestRouteProps {
  children: React.ReactNode
}

export default function GuestRoute({ children }: GuestRouteProps) {
  const { isAuthenticated, isLoading } = useAuth()

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner size="lg" />
      </div>
    )
  }

  if (isAuthenticated) {
    return <Navigate to="/" replace />
  }

  return <>{children}</>
}
