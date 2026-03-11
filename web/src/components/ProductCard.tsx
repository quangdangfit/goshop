import { Heart, ShoppingCart } from 'lucide-react'
import { Link } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import toast from 'react-hot-toast'
import { cartApi } from '@/api/cart'
import { wishlistApi } from '@/api/wishlist'
import { useAuth } from '@/context/AuthContext'
import type { Product } from '@/types'
import StarRating from './StarRating'

interface ProductCardProps {
  product: Product
}

export default function ProductCard({ product }: ProductCardProps) {
  const { isAuthenticated } = useAuth()
  const queryClient = useQueryClient()

  const { data: wishlist } = useQuery({
    queryKey: ['wishlist'],
    queryFn: wishlistApi.getWishlist,
    enabled: isAuthenticated,
  })

  const isWishlisted = wishlist?.some((w) => w.product_id === product.id)

  const addToCartMutation = useMutation({
    mutationFn: () =>
      cartApi.addToCart({ product_id: product.id, quantity: 1 }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cart'] })
      toast.success('Added to cart!')
    },
    onError: () => toast.error('Failed to add to cart'),
  })

  const wishlistMutation = useMutation({
    mutationFn: async (): Promise<void> => {
      if (isWishlisted) {
        await wishlistApi.removeFromWishlist(product.id)
      } else {
        await wishlistApi.addToWishlist(product.id)
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['wishlist'] })
      toast.success(isWishlisted ? 'Removed from wishlist' : 'Added to wishlist!')
    },
    onError: () => toast.error('Failed to update wishlist'),
  })

  const handleAddToCart = (e: React.MouseEvent) => {
    e.preventDefault()
    if (!isAuthenticated) {
      toast.error('Please login to add items to cart')
      return
    }
    addToCartMutation.mutate()
  }

  const handleWishlist = (e: React.MouseEvent) => {
    e.preventDefault()
    if (!isAuthenticated) {
      toast.error('Please login to use wishlist')
      return
    }
    wishlistMutation.mutate()
  }

  return (
    <Link
      to={`/products/${product.id}`}
      className="group bg-white rounded-xl shadow-sm border border-gray-100 hover:shadow-md transition-shadow overflow-hidden"
    >
      {/* Image */}
      <div className="relative h-48 bg-gradient-to-br from-gray-100 to-gray-200 overflow-hidden">
        {product.images?.[0] ? (
          <img
            src={product.images[0]}
            alt={product.name}
            className="absolute inset-0 w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
          />
        ) : (
          <div className="absolute inset-0 flex items-center justify-center">
            <ShoppingCart className="h-16 w-16 text-gray-300" />
          </div>
        )}
        <button
          onClick={handleWishlist}
          className={`absolute top-2 right-2 p-1.5 rounded-full shadow transition-all ${
            isWishlisted
              ? 'bg-red-500 text-white'
              : 'bg-white text-gray-400 hover:text-red-500'
          }`}
        >
          <Heart className={`h-4 w-4 ${isWishlisted ? 'fill-current' : ''}`} />
        </button>
        {product.stock_quantity === 0 && (
          <div className="absolute bottom-0 left-0 right-0 bg-black/50 text-white text-xs text-center py-1">
            Out of Stock
          </div>
        )}
      </div>

      <div className="p-4">
        {product.category && (
          <span className="text-xs text-primary-600 font-medium uppercase tracking-wide">
            {product.category.name}
          </span>
        )}
        <h3 className="mt-1 text-sm font-semibold text-gray-900 line-clamp-2 group-hover:text-primary-600 transition-colors">
          {product.name}
        </h3>

        <div className="mt-1 flex items-center gap-1">
          <StarRating rating={Math.round(product.average_rating || 0)} size="sm" />
          <span className="text-xs text-gray-500">
            ({product.review_count || 0})
          </span>
        </div>

        <div className="mt-3 flex items-center justify-between">
          <span className="text-lg font-bold text-gray-900">
            ${product.price?.toFixed(2)}
          </span>
          <button
            onClick={handleAddToCart}
            disabled={
              product.stock_quantity === 0 || addToCartMutation.isPending
            }
            className="flex items-center gap-1 px-3 py-1.5 bg-primary-600 text-white text-xs font-medium rounded-lg hover:bg-primary-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <ShoppingCart className="h-3.5 w-3.5" />
            Add
          </button>
        </div>
      </div>
    </Link>
  )
}
