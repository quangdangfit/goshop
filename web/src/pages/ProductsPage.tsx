import { Filter, Search, SlidersHorizontal, X } from 'lucide-react'
import { useEffect, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useSearchParams } from 'react-router-dom'
import { productsApi } from '@/api/products'
import { categoriesApi } from '@/api/categories'
import ProductCard from '@/components/ProductCard'
import Pagination from '@/components/Pagination'
import LoadingSpinner from '@/components/LoadingSpinner'
import type { ProductsQueryParams } from '@/types'

export default function ProductsPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const [sidebarOpen, setSidebarOpen] = useState(false)

  const [filters, setFilters] = useState<ProductsQueryParams>({
    name: searchParams.get('name') || '',
    category_id: searchParams.get('category_id') || '',
    page: Number(searchParams.get('page')) || 1,
    limit: 12,
    order_by: searchParams.get('order_by') || '',
    order_desc: searchParams.get('order_desc') === 'true',
  })

  const { data, isLoading } = useQuery({
    queryKey: ['products', filters],
    queryFn: () =>
      productsApi.getProducts({
        ...filters,
        name: filters.name || undefined,
        category_id: filters.category_id || undefined,
        order_by: filters.order_by || undefined,
      }),
  })

  const { data: categories } = useQuery({
    queryKey: ['categories'],
    queryFn: categoriesApi.getCategories,
  })

  useEffect(() => {
    const params: Record<string, string> = {}
    if (filters.name) params.name = filters.name
    if (filters.category_id) params.category_id = filters.category_id
    if (filters.page && filters.page > 1) params.page = String(filters.page)
    if (filters.order_by) params.order_by = filters.order_by
    if (filters.order_desc) params.order_desc = 'true'
    setSearchParams(params)
  }, [filters, setSearchParams])

  const updateFilter = <K extends keyof ProductsQueryParams>(
    key: K,
    value: ProductsQueryParams[K]
  ) => {
    setFilters((prev) => ({ ...prev, [key]: value, page: key !== 'page' ? 1 : (value as number) }))
  }

  const clearFilters = () => {
    setFilters({ page: 1, limit: 12 })
  }

  const hasActiveFilters = !!(filters.name || filters.category_id || filters.order_by)

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Products</h1>
          {data && (
            <p className="text-sm text-gray-500 mt-0.5">
              {data.pagination.total} products found
            </p>
          )}
        </div>
        <button
          onClick={() => setSidebarOpen(!sidebarOpen)}
          className="flex items-center gap-2 px-4 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 lg:hidden"
        >
          <Filter className="h-4 w-4" />
          Filters
          {hasActiveFilters && (
            <span className="h-2 w-2 bg-primary-600 rounded-full" />
          )}
        </button>
      </div>

      <div className="flex gap-6">
        {/* Sidebar Filters */}
        <aside
          className={`${
            sidebarOpen ? 'block' : 'hidden'
          } lg:block w-full lg:w-64 flex-shrink-0`}
        >
          <div className="bg-white rounded-xl border border-gray-100 p-5 sticky top-20">
            <div className="flex items-center justify-between mb-4">
              <h2 className="font-semibold text-gray-900 flex items-center gap-2">
                <SlidersHorizontal className="h-4 w-4" />
                Filters
              </h2>
              {hasActiveFilters && (
                <button
                  onClick={clearFilters}
                  className="text-xs text-red-500 hover:text-red-600 flex items-center gap-1"
                >
                  <X className="h-3 w-3" />
                  Clear
                </button>
              )}
            </div>

            {/* Search */}
            <div className="mb-5">
              <label className="label">Search</label>
              <div className="relative">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
                <input
                  type="text"
                  placeholder="Search products..."
                  value={filters.name || ''}
                  onChange={(e) => updateFilter('name', e.target.value)}
                  className="input pl-9"
                />
              </div>
            </div>

            {/* Categories */}
            <div className="mb-5">
              <label className="label">Category</label>
              <div className="space-y-1">
                <button
                  onClick={() => updateFilter('category_id', '')}
                  className={`w-full text-left text-sm px-3 py-2 rounded-lg transition-colors ${
                    !filters.category_id
                      ? 'bg-primary-50 text-primary-700 font-medium'
                      : 'text-gray-600 hover:bg-gray-50'
                  }`}
                >
                  All Categories
                </button>
                {categories?.map((cat) => (
                  <button
                    key={cat.id}
                    onClick={() => updateFilter('category_id', cat.id)}
                    className={`w-full text-left text-sm px-3 py-2 rounded-lg transition-colors ${
                      filters.category_id === cat.id
                        ? 'bg-primary-50 text-primary-700 font-medium'
                        : 'text-gray-600 hover:bg-gray-50'
                    }`}
                  >
                    {cat.name}
                  </button>
                ))}
              </div>
            </div>

            {/* Sort */}
            <div>
              <label className="label">Sort By</label>
              <select
                value={`${filters.order_by || ''}:${filters.order_desc ? 'desc' : 'asc'}`}
                onChange={(e) => {
                  const [field, dir] = e.target.value.split(':')
                  setFilters((prev) => ({
                    ...prev,
                    order_by: field || undefined,
                    order_desc: dir === 'desc',
                    page: 1,
                  }))
                }}
                className="input"
              >
                <option value=":asc">Default</option>
                <option value="price:asc">Price: Low to High</option>
                <option value="price:desc">Price: High to Low</option>
                <option value="created_at:desc">Newest First</option>
                <option value="name:asc">Name: A-Z</option>
                <option value="name:desc">Name: Z-A</option>
              </select>
            </div>
          </div>
        </aside>

        {/* Product Grid */}
        <div className="flex-1">
          {isLoading ? (
            <LoadingSpinner className="py-20" />
          ) : !data?.items?.length ? (
            <div className="text-center py-20 text-gray-500">
              <ShoppingBagEmpty />
              <p className="mt-3 font-medium">No products found</p>
              <p className="text-sm">Try adjusting your filters</p>
              {hasActiveFilters && (
                <button
                  onClick={clearFilters}
                  className="mt-4 btn-primary text-sm"
                >
                  Clear Filters
                </button>
              )}
            </div>
          ) : (
            <>
              <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-5">
                {data.items.map((product) => (
                  <ProductCard key={product.id} product={product} />
                ))}
              </div>
              {data.pagination && (
                <Pagination
                  pagination={data.pagination}
                  onPageChange={(page) => updateFilter('page', page)}
                />
              )}
            </>
          )}
        </div>
      </div>
    </div>
  )
}

function ShoppingBagEmpty() {
  return (
    <svg
      className="mx-auto h-16 w-16 text-gray-300"
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={1}
        d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"
      />
    </svg>
  )
}
