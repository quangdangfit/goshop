export interface User {
  id: string
  email: string
  username: string
  role: string
  created_at: string
  updated_at: string
}

export interface AuthTokens {
  access_token: string
  refresh_token: string
}

export interface AuthResponse {
  user: User
  access_token: string
  refresh_token: string
}

export interface Pagination {
  page: number
  limit: number
  total: number
  total_pages: number
}

export interface PaginatedResponse<T> {
  items: T[]
  pagination: Pagination
}

export interface Category {
  id: string
  name: string
  slug: string
  description: string
  created_at: string
  updated_at: string
}

export interface Product {
  id: string
  name: string
  code: string
  description: string
  price: number
  stock_quantity: number
  category_id: string
  category: Category
  average_rating: number
  review_count: number
  created_at: string
  updated_at: string
}

export interface Review {
  id: string
  product_id: string
  user_id: string
  user: User
  rating: number
  comment: string
  created_at: string
  updated_at: string
}

export interface CartLine {
  id: string
  cart_id: string
  product_id: string
  product: Product
  quantity: number
  created_at: string
  updated_at: string
}

export interface Cart {
  id: string
  user_id: string
  user: User
  lines: CartLine[]
  created_at: string
  updated_at: string
}

export interface OrderLine {
  id: string
  order_id: string
  product_id: string
  product: Product
  quantity: number
  price: number
  created_at: string
  updated_at: string
}

export interface Order {
  id: string
  code: string
  user_id: string
  user: User
  coupon_code: string
  status: string
  total_price: number
  lines: OrderLine[]
  created_at: string
  updated_at: string
}

export interface Address {
  id: string
  user_id: string
  name: string
  phone: string
  street: string
  city: string
  country: string
  is_default: boolean
  created_at: string
  updated_at: string
}

export interface WishlistItem {
  id: string
  user_id: string
  product_id: string
  product: Product
  created_at: string
}

export interface Coupon {
  id: string
  code: string
  discount_type: 'fixed' | 'percentage'
  discount_value: number
  min_order_amount: number
  max_usage: number
  used_count: number
  expires_at: string
  created_at: string
  updated_at: string
}

export interface ApiError {
  error: string
  message?: string
}

// Request types
export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
  username?: string
}

export interface ChangePasswordRequest {
  password: string
  new_password: string
}

export interface CreateProductRequest {
  name: string
  description: string
  price: number
  stock_quantity: number
  category_id: string
}

export interface UpdateProductRequest extends Partial<CreateProductRequest> {}

export interface CreateCategoryRequest {
  name: string
  slug: string
  description: string
}

export interface UpdateCategoryRequest extends Partial<CreateCategoryRequest> {}

export interface CreateReviewRequest {
  rating: number
  comment: string
}

export interface UpdateReviewRequest extends Partial<CreateReviewRequest> {}

export interface AddToCartRequest {
  product_id: string
  quantity: number
}

export interface CreateOrderRequest {
  coupon_code?: string
  lines: Array<{
    product_id: string
    quantity: number
  }>
}

export interface CreateAddressRequest {
  name: string
  phone: string
  street: string
  city: string
  country: string
}

export interface UpdateAddressRequest extends Partial<CreateAddressRequest> {}

export interface CreateCouponRequest {
  code: string
  discount_type: 'fixed' | 'percentage'
  discount_value: number
  min_order_amount: number
  max_usage: number
  expires_at: string
}

export interface ProductsQueryParams {
  name?: string
  code?: string
  category_id?: string
  page?: number
  limit?: number
  order_by?: string
  order_desc?: boolean
}

export interface OrdersQueryParams {
  code?: string
  status?: string
  page?: number
  limit?: number
}
