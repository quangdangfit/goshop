import {
  ArrowLeft,
  Heart,
  Minus,
  Plus,
  ShoppingCart,
  Star,
  Trash2,
} from 'lucide-react'
import { useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import {
  useMutation,
  useQuery,
  useQueryClient,
} from '@tanstack/react-query'
import toast from 'react-hot-toast'
import { productsApi } from '@/api/products'
import { cartApi } from '@/api/cart'
import { wishlistApi } from '@/api/wishlist'
import { useAuth } from '@/context/AuthContext'
import StarRating from '@/components/StarRating'
import LoadingSpinner from '@/components/LoadingSpinner'
import Pagination from '@/components/Pagination'
import type { CreateReviewRequest } from '@/types'

export default function ProductDetailPage() {
  const { id } = useParams<{ id: string }>()
  const { isAuthenticated, user } = useAuth()
  const queryClient = useQueryClient()
  const [quantity, setQuantity] = useState(1)
  const [selectedImage, setSelectedImage] = useState(0)
  const [reviewPage, setReviewPage] = useState(1)
  const [reviewForm, setReviewForm] = useState<CreateReviewRequest>({
    rating: 5,
    comment: '',
  })
  const [editingReviewId, setEditingReviewId] = useState<string | null>(null)

  const { data: product, isLoading } = useQuery({
    queryKey: ['product', id],
    queryFn: () => productsApi.getProduct(id!),
    enabled: !!id,
  })

  const { data: reviewsData } = useQuery({
    queryKey: ['reviews', id, reviewPage],
    queryFn: () =>
      productsApi.getReviews(id!, { page: reviewPage, limit: 10 }),
    enabled: !!id,
  })

  const { data: wishlist } = useQuery({
    queryKey: ['wishlist'],
    queryFn: wishlistApi.getWishlist,
    enabled: isAuthenticated,
  })

  const isWishlisted = wishlist?.some((w) => w.product_id === id)

  const addToCartMutation = useMutation({
    mutationFn: () =>
      cartApi.addToCart({ product_id: id!, quantity }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cart'] })
      toast.success('Added to cart!')
    },
    onError: () => toast.error('Failed to add to cart'),
  })

  const wishlistMutation = useMutation({
    mutationFn: async (): Promise<void> => {
      if (isWishlisted) {
        await wishlistApi.removeFromWishlist(id!)
      } else {
        await wishlistApi.addToWishlist(id!)
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['wishlist'] })
      toast.success(isWishlisted ? 'Removed from wishlist' : 'Added to wishlist!')
    },
    onError: () => toast.error('Failed to update wishlist'),
  })

  const createReviewMutation = useMutation({
    mutationFn: (data: CreateReviewRequest) =>
      productsApi.createReview(id!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['reviews', id] })
      queryClient.invalidateQueries({ queryKey: ['product', id] })
      setReviewForm({ rating: 5, comment: '' })
      toast.success('Review submitted!')
    },
    onError: () => toast.error('Failed to submit review'),
  })

  const updateReviewMutation = useMutation({
    mutationFn: ({ reviewId, data }: { reviewId: string; data: CreateReviewRequest }) =>
      productsApi.updateReview(id!, reviewId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['reviews', id] })
      setEditingReviewId(null)
      toast.success('Review updated!')
    },
    onError: () => toast.error('Failed to update review'),
  })

  const deleteReviewMutation = useMutation({
    mutationFn: (reviewId: string) =>
      productsApi.deleteReview(id!, reviewId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['reviews', id] })
      queryClient.invalidateQueries({ queryKey: ['product', id] })
      toast.success('Review deleted!')
    },
    onError: () => toast.error('Failed to delete review'),
  })

  if (isLoading) {
    return <LoadingSpinner className="min-h-[400px]" size="lg" />
  }

  if (!product) {
    return (
      <div className="max-w-7xl mx-auto px-4 py-12 text-center">
        <p className="text-gray-500">Product not found</p>
        <Link to="/products" className="btn-primary mt-4 inline-flex">
          Back to Products
        </Link>
      </div>
    )
  }

  const userReview = reviewsData?.items?.find((r) => r.user_id === user?.id)

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Breadcrumb */}
      <Link
        to="/products"
        className="inline-flex items-center gap-1.5 text-sm text-gray-500 hover:text-primary-600 mb-6 transition-colors"
      >
        <ArrowLeft className="h-4 w-4" />
        Back to Products
      </Link>

      {/* Product Info */}
      <div className="bg-white rounded-2xl border border-gray-100 p-6 md:p-8 mb-8">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
          {/* Image gallery */}
          <div>
            <div className="bg-gradient-to-br from-gray-100 to-gray-200 rounded-xl h-80 md:h-96 overflow-hidden flex items-center justify-center">
              {product.images?.[selectedImage] ? (
                <img
                  src={product.images[selectedImage]}
                  alt={product.name}
                  className="w-full h-full object-cover"
                />
              ) : (
                <ShoppingCart className="h-24 w-24 text-gray-300" />
              )}
            </div>
            {product.images?.length > 1 && (
              <div className="flex gap-2 mt-3">
                {product.images.map((img, i) => (
                  <button
                    key={i}
                    onClick={() => setSelectedImage(i)}
                    className={`w-16 h-16 rounded-lg overflow-hidden border-2 transition-colors flex-shrink-0 ${
                      selectedImage === i ? 'border-primary-500' : 'border-gray-200 hover:border-gray-300'
                    }`}
                  >
                    <img src={img} alt={`${product.name} ${i + 1}`} className="w-full h-full object-cover" />
                  </button>
                ))}
              </div>
            )}
          </div>

          {/* Details */}
          <div>
            {product.category && (
              <Link
                to={`/products?category_id=${product.category.id}`}
                className="text-sm text-primary-600 font-medium hover:text-primary-700"
              >
                {product.category.name}
              </Link>
            )}

            <h1 className="text-2xl md:text-3xl font-bold text-gray-900 mt-1 mb-3">
              {product.name}
            </h1>

            <div className="flex items-center gap-2 mb-4">
              <StarRating rating={Math.round(product.average_rating || 0)} />
              <span className="text-sm text-gray-500">
                {product.average_rating?.toFixed(1) || '0.0'} ({product.review_count || 0} reviews)
              </span>
            </div>

            <p className="text-3xl font-extrabold text-gray-900 mb-4">
              ${product.price?.toFixed(2)}
            </p>

            <p className="text-gray-600 leading-relaxed mb-6">
              {product.description || 'No description available.'}
            </p>

            <div className="flex items-center gap-2 mb-4">
              <span
                className={`badge ${
                  product.stock_quantity > 0 ? 'badge-success' : 'badge-danger'
                }`}
              >
                {product.stock_quantity > 0
                  ? `${product.stock_quantity} in stock`
                  : 'Out of stock'}
              </span>
              <span className="text-sm text-gray-500">SKU: {product.code}</span>
            </div>

            {/* Quantity + Add to Cart */}
            {product.stock_quantity > 0 && (
              <div className="flex items-center gap-3 mb-4">
                <div className="flex items-center border border-gray-300 rounded-lg overflow-hidden">
                  <button
                    onClick={() => setQuantity((q) => Math.max(1, q - 1))}
                    className="p-2 hover:bg-gray-100 transition-colors"
                  >
                    <Minus className="h-4 w-4" />
                  </button>
                  <span className="w-10 text-center font-medium">{quantity}</span>
                  <button
                    onClick={() =>
                      setQuantity((q) => Math.min(product.stock_quantity, q + 1))
                    }
                    className="p-2 hover:bg-gray-100 transition-colors"
                  >
                    <Plus className="h-4 w-4" />
                  </button>
                </div>

                <button
                  onClick={() => {
                    if (!isAuthenticated) {
                      toast.error('Please login to add items to cart')
                      return
                    }
                    addToCartMutation.mutate()
                  }}
                  disabled={addToCartMutation.isPending}
                  className="flex-1 btn-primary"
                >
                  <ShoppingCart className="h-4 w-4" />
                  {addToCartMutation.isPending ? 'Adding...' : 'Add to Cart'}
                </button>

                <button
                  onClick={() => {
                    if (!isAuthenticated) {
                      toast.error('Please login to use wishlist')
                      return
                    }
                    wishlistMutation.mutate()
                  }}
                  className={`p-2.5 border rounded-lg transition-colors ${
                    isWishlisted
                      ? 'bg-red-50 border-red-300 text-red-500'
                      : 'border-gray-300 text-gray-500 hover:bg-gray-50'
                  }`}
                >
                  <Heart
                    className={`h-5 w-5 ${isWishlisted ? 'fill-current' : ''}`}
                  />
                </button>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Reviews Section */}
      <div className="bg-white rounded-2xl border border-gray-100 p-6 md:p-8">
        <h2 className="text-xl font-bold text-gray-900 mb-6">
          Customer Reviews
        </h2>

        {/* Write a review */}
        {isAuthenticated && !userReview && (
          <div className="bg-gray-50 rounded-xl p-5 mb-6">
            <h3 className="font-semibold text-gray-900 mb-3">Write a Review</h3>
            <div className="mb-3">
              <label className="label">Your Rating</label>
              <StarRating
                rating={reviewForm.rating}
                interactive
                onRatingChange={(r) =>
                  setReviewForm((f) => ({ ...f, rating: r }))
                }
                size="lg"
              />
            </div>
            <div className="mb-3">
              <label className="label">Comment</label>
              <textarea
                rows={3}
                value={reviewForm.comment}
                onChange={(e) =>
                  setReviewForm((f) => ({ ...f, comment: e.target.value }))
                }
                placeholder="Share your experience..."
                className="input resize-none"
              />
            </div>
            <button
              onClick={() => createReviewMutation.mutate(reviewForm)}
              disabled={createReviewMutation.isPending || !reviewForm.comment}
              className="btn-primary"
            >
              {createReviewMutation.isPending ? 'Submitting...' : 'Submit Review'}
            </button>
          </div>
        )}

        {/* Reviews list */}
        {reviewsData?.items?.length === 0 ? (
          <div className="text-center py-8 text-gray-500">
            <Star className="h-10 w-10 text-gray-300 mx-auto mb-2" />
            <p>No reviews yet. Be the first to review!</p>
          </div>
        ) : (
          <div className="space-y-4">
            {reviewsData?.items?.map((review) => (
              <div key={review.id} className="border-b border-gray-100 pb-4 last:border-0">
                {editingReviewId === review.id ? (
                  <div className="bg-gray-50 rounded-xl p-4">
                    <div className="mb-3">
                      <label className="label">Rating</label>
                      <StarRating
                        rating={reviewForm.rating}
                        interactive
                        onRatingChange={(r) =>
                          setReviewForm((f) => ({ ...f, rating: r }))
                        }
                        size="lg"
                      />
                    </div>
                    <textarea
                      rows={3}
                      value={reviewForm.comment}
                      onChange={(e) =>
                        setReviewForm((f) => ({ ...f, comment: e.target.value }))
                      }
                      className="input resize-none mb-3"
                    />
                    <div className="flex gap-2">
                      <button
                        onClick={() =>
                          updateReviewMutation.mutate({
                            reviewId: review.id,
                            data: reviewForm,
                          })
                        }
                        disabled={updateReviewMutation.isPending}
                        className="btn-primary text-sm"
                      >
                        Save
                      </button>
                      <button
                        onClick={() => setEditingReviewId(null)}
                        className="btn-secondary text-sm"
                      >
                        Cancel
                      </button>
                    </div>
                  </div>
                ) : (
                  <div>
                    <div className="flex items-start justify-between">
                      <div>
                        <p className="font-medium text-gray-900 text-sm">
                          {review.user?.username || review.user?.email || 'Anonymous'}
                        </p>
                        <StarRating rating={review.rating} size="sm" />
                      </div>
                      <div className="flex items-center gap-2">
                        <span className="text-xs text-gray-400">
                          {new Date(review.created_at).toLocaleDateString()}
                        </span>
                        {user?.id === review.user_id && (
                          <>
                            <button
                              onClick={() => {
                                setEditingReviewId(review.id)
                                setReviewForm({
                                  rating: review.rating,
                                  comment: review.comment,
                                })
                              }}
                              className="text-xs text-primary-600 hover:text-primary-700"
                            >
                              Edit
                            </button>
                            <button
                              onClick={() => deleteReviewMutation.mutate(review.id)}
                              className="text-xs text-red-500 hover:text-red-600"
                            >
                              <Trash2 className="h-3.5 w-3.5" />
                            </button>
                          </>
                        )}
                      </div>
                    </div>
                    <p className="mt-1.5 text-sm text-gray-600">{review.comment}</p>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}

        {reviewsData?.pagination && (
          <Pagination
            pagination={reviewsData.pagination}
            onPageChange={setReviewPage}
          />
        )}
      </div>
    </div>
  )
}
