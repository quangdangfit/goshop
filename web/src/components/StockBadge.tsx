import type { Product } from '@/types'

export const LOW_STOCK_THRESHOLD = 5

export function availableStock(product: Pick<Product, 'stock_quantity' | 'reserved_quantity'>): number {
  return Math.max(0, product.stock_quantity - (product.reserved_quantity ?? 0))
}

interface StockBadgeProps {
  product: Pick<Product, 'stock_quantity' | 'reserved_quantity'>
  showCount?: boolean
  className?: string
}

export default function StockBadge({ product, showCount, className = '' }: StockBadgeProps) {
  const available = availableStock(product)

  if (available === 0) {
    return <span className={`badge badge-danger ${className}`}>Out of stock</span>
  }
  if (available <= LOW_STOCK_THRESHOLD) {
    return (
      <span className={`badge badge-warning ${className}`}>
        Low stock{showCount ? ` — only ${available} left` : ''}
      </span>
    )
  }
  return (
    <span className={`badge badge-success ${className}`}>
      {showCount ? `${available} in stock` : 'In stock'}
    </span>
  )
}
