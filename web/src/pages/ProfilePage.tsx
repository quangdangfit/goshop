import {
  Eye,
  EyeOff,
  Heart,
  MapPin,
  Pencil,
  Plus,
  Trash2,
  User,
  X,
} from 'lucide-react'
import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import toast from 'react-hot-toast'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useAuth } from '@/context/AuthContext'
import { authApi } from '@/api/auth'
import { addressesApi } from '@/api/addresses'
import { wishlistApi } from '@/api/wishlist'
import LoadingSpinner from '@/components/LoadingSpinner'
import type { Address, CreateAddressRequest } from '@/types'

type Tab = 'profile' | 'addresses' | 'wishlist' | 'password'

const passwordSchema = z
  .object({
    password: z.string().min(1, 'Current password is required'),
    new_password: z.string().min(6, 'New password must be at least 6 characters'),
    confirm_password: z.string(),
  })
  .refine((d) => d.new_password === d.confirm_password, {
    message: "Passwords don't match",
    path: ['confirm_password'],
  })

type PasswordForm = z.infer<typeof passwordSchema>

const addressSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  phone: z.string().min(1, 'Phone is required'),
  street: z.string().min(1, 'Street is required'),
  city: z.string().min(1, 'City is required'),
  country: z.string().min(1, 'Country is required'),
})

type AddressForm = z.infer<typeof addressSchema>

export default function ProfilePage() {
  const { user, refreshUser } = useAuth()
  const queryClient = useQueryClient()
  const [activeTab, setActiveTab] = useState<Tab>('profile')
  const [editingAddressId, setEditingAddressId] = useState<string | null>(null)
  const [showAddAddress, setShowAddAddress] = useState(false)
  const [showPasswords, setShowPasswords] = useState(false)

  const { data: addresses, isLoading: addressLoading } = useQuery({
    queryKey: ['addresses'],
    queryFn: addressesApi.getAddresses,
    enabled: activeTab === 'addresses',
  })

  const { data: wishlist, isLoading: wishlistLoading } = useQuery({
    queryKey: ['wishlist'],
    queryFn: wishlistApi.getWishlist,
    enabled: activeTab === 'wishlist',
  })

  const {
    register: registerAddr,
    handleSubmit: handleAddrSubmit,
    reset: resetAddr,
    formState: { errors: addrErrors },
  } = useForm<AddressForm>({ resolver: zodResolver(addressSchema) })

  const {
    register: registerPwd,
    handleSubmit: handlePwdSubmit,
    reset: resetPwd,
    formState: { errors: pwdErrors, isSubmitting: pwdSubmitting },
  } = useForm<PasswordForm>({ resolver: zodResolver(passwordSchema) })

  const createAddressMutation = useMutation({
    mutationFn: addressesApi.createAddress,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['addresses'] })
      setShowAddAddress(false)
      resetAddr()
      toast.success('Address added!')
    },
    onError: () => toast.error('Failed to add address'),
  })

  const updateAddressMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<CreateAddressRequest> }) =>
      addressesApi.updateAddress(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['addresses'] })
      setEditingAddressId(null)
      toast.success('Address updated!')
    },
    onError: () => toast.error('Failed to update address'),
  })

  const deleteAddressMutation = useMutation({
    mutationFn: addressesApi.deleteAddress,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['addresses'] })
      toast.success('Address deleted')
    },
    onError: () => toast.error('Failed to delete address'),
  })

  const setDefaultMutation = useMutation({
    mutationFn: addressesApi.setDefaultAddress,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['addresses'] })
      toast.success('Default address updated!')
    },
    onError: () => toast.error('Failed to update default address'),
  })

  const removeWishlistMutation = useMutation({
    mutationFn: wishlistApi.removeFromWishlist,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['wishlist'] })
      toast.success('Removed from wishlist')
    },
    onError: () => toast.error('Failed to remove from wishlist'),
  })

  const onChangePassword = async (data: PasswordForm) => {
    try {
      await authApi.changePassword({
        password: data.password,
        new_password: data.new_password,
      })
      await refreshUser()
      resetPwd()
      toast.success('Password changed successfully!')
    } catch {
      toast.error('Failed to change password. Check your current password.')
    }
  }

  const tabs: { id: Tab; label: string; icon: typeof User }[] = [
    { id: 'profile', label: 'Profile', icon: User },
    { id: 'addresses', label: 'Addresses', icon: MapPin },
    { id: 'wishlist', label: 'Wishlist', icon: Heart },
    { id: 'password', label: 'Password', icon: Eye },
  ]

  return (
    <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="text-2xl font-bold text-gray-900 mb-6">My Profile</h1>

      <div className="flex flex-col md:flex-row gap-6">
        {/* Sidebar */}
        <aside className="md:w-56 flex-shrink-0">
          <div className="card p-2">
            {/* Avatar */}
            <div className="px-3 py-4 text-center border-b border-gray-100 mb-2">
              <div className="h-16 w-16 bg-primary-100 rounded-full flex items-center justify-center mx-auto mb-2">
                <User className="h-8 w-8 text-primary-600" />
              </div>
              <p className="font-semibold text-gray-900 truncate">
                {user?.username || 'User'}
              </p>
              <p className="text-xs text-gray-500 truncate">{user?.email}</p>
              {user?.role === 'admin' && (
                <span className="mt-1 badge badge-info">Admin</span>
              )}
            </div>

            <nav className="space-y-0.5">
              {tabs.map(({ id, label, icon: Icon }) => (
                <button
                  key={id}
                  onClick={() => setActiveTab(id)}
                  className={`w-full flex items-center gap-2.5 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors ${
                    activeTab === id
                      ? 'bg-primary-50 text-primary-700'
                      : 'text-gray-600 hover:bg-gray-50'
                  }`}
                >
                  <Icon className="h-4 w-4" />
                  {label}
                </button>
              ))}
            </nav>
          </div>
        </aside>

        {/* Content */}
        <div className="flex-1">
          {/* Profile Tab */}
          {activeTab === 'profile' && (
            <div className="card">
              <h2 className="font-bold text-gray-900 mb-5">Profile Information</h2>
              <div className="space-y-4 text-sm">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <p className="label">Username</p>
                    <p className="text-gray-900 font-medium">{user?.username || '—'}</p>
                  </div>
                  <div>
                    <p className="label">Role</p>
                    <p className="text-gray-900 font-medium capitalize">{user?.role}</p>
                  </div>
                </div>
                <div>
                  <p className="label">Email</p>
                  <p className="text-gray-900 font-medium">{user?.email}</p>
                </div>
                <div>
                  <p className="label">Member Since</p>
                  <p className="text-gray-900 font-medium">
                    {user?.created_at
                      ? new Date(user.created_at).toLocaleDateString('en-US', {
                          year: 'numeric',
                          month: 'long',
                          day: 'numeric',
                        })
                      : '—'}
                  </p>
                </div>
              </div>
            </div>
          )}

          {/* Addresses Tab */}
          {activeTab === 'addresses' && (
            <div className="card">
              <div className="flex items-center justify-between mb-5">
                <h2 className="font-bold text-gray-900">My Addresses</h2>
                <button
                  onClick={() => {
                    setShowAddAddress(!showAddAddress)
                    setEditingAddressId(null)
                  }}
                  className="btn-primary text-sm py-1.5"
                >
                  <Plus className="h-4 w-4" />
                  Add Address
                </button>
              </div>

              {showAddAddress && (
                <AddressFormComponent
                  register={registerAddr}
                  errors={addrErrors}
                  onSubmit={handleAddrSubmit((data) => createAddressMutation.mutate(data))}
                  onCancel={() => {
                    setShowAddAddress(false)
                    resetAddr()
                  }}
                  isPending={createAddressMutation.isPending}
                  submitLabel="Save Address"
                />
              )}

              {addressLoading ? (
                <LoadingSpinner className="py-8" />
              ) : !addresses?.length ? (
                <div className="text-center py-8 text-gray-500">
                  <MapPin className="h-10 w-10 text-gray-300 mx-auto mb-2" />
                  <p>No addresses yet</p>
                </div>
              ) : (
                <div className="space-y-3">
                  {addresses.map((addr: Address) => (
                    <div
                      key={addr.id}
                      className={`rounded-xl border p-4 ${
                        addr.is_default ? 'border-primary-300 bg-primary-50' : 'border-gray-200'
                      }`}
                    >
                      {editingAddressId === addr.id ? (
                        <EditAddressForm
                          address={addr}
                          onSave={(data) =>
                            updateAddressMutation.mutate({ id: addr.id, data })
                          }
                          onCancel={() => setEditingAddressId(null)}
                          isPending={updateAddressMutation.isPending}
                        />
                      ) : (
                        <div className="flex items-start justify-between">
                          <div>
                            <p className="font-medium text-gray-900 text-sm">
                              {addr.name}
                              {addr.is_default && (
                                <span className="ml-2 badge badge-info">Default</span>
                              )}
                            </p>
                            <p className="text-sm text-gray-500">{addr.phone}</p>
                            <p className="text-sm text-gray-500">
                              {addr.street}, {addr.city}, {addr.country}
                            </p>
                          </div>
                          <div className="flex items-center gap-1 ml-3">
                            {!addr.is_default && (
                              <button
                                onClick={() => setDefaultMutation.mutate(addr.id)}
                                className="text-xs px-2 py-1 text-primary-600 border border-primary-200 rounded hover:bg-primary-50"
                              >
                                Set Default
                              </button>
                            )}
                            <button
                              onClick={() => setEditingAddressId(addr.id)}
                              className="p-1.5 text-gray-400 hover:text-primary-600 rounded"
                            >
                              <Pencil className="h-3.5 w-3.5" />
                            </button>
                            <button
                              onClick={() => deleteAddressMutation.mutate(addr.id)}
                              className="p-1.5 text-gray-400 hover:text-red-500 rounded"
                            >
                              <Trash2 className="h-3.5 w-3.5" />
                            </button>
                          </div>
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}

          {/* Wishlist Tab */}
          {activeTab === 'wishlist' && (
            <div className="card">
              <h2 className="font-bold text-gray-900 mb-5">My Wishlist</h2>
              {wishlistLoading ? (
                <LoadingSpinner className="py-8" />
              ) : !wishlist?.length ? (
                <div className="text-center py-8 text-gray-500">
                  <Heart className="h-10 w-10 text-gray-300 mx-auto mb-2" />
                  <p>Your wishlist is empty</p>
                  <Link to="/products" className="btn-primary mt-3 inline-flex">
                    Browse Products
                  </Link>
                </div>
              ) : (
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                  {wishlist.map((item) => (
                    <div
                      key={item.id}
                      className="flex items-center gap-3 p-3 border border-gray-100 rounded-xl hover:border-gray-200 transition-colors"
                    >
                      <div className="h-14 w-14 bg-gray-100 rounded-lg flex-shrink-0 flex items-center justify-center">
                        <Heart className="h-5 w-5 text-gray-300" />
                      </div>
                      <div className="flex-1 min-w-0">
                        <Link
                          to={`/products/${item.product_id}`}
                          className="text-sm font-medium text-gray-900 hover:text-primary-600 line-clamp-1"
                        >
                          {item.product?.name || 'Product'}
                        </Link>
                        <p className="text-sm font-bold text-primary-600">
                          ${item.product?.price?.toFixed(2)}
                        </p>
                      </div>
                      <button
                        onClick={() => removeWishlistMutation.mutate(item.product_id)}
                        className="p-1.5 text-gray-400 hover:text-red-500 transition-colors"
                      >
                        <X className="h-4 w-4" />
                      </button>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}

          {/* Password Tab */}
          {activeTab === 'password' && (
            <div className="card">
              <h2 className="font-bold text-gray-900 mb-5">Change Password</h2>
              <form
                onSubmit={handlePwdSubmit(onChangePassword)}
                className="space-y-4 max-w-sm"
              >
                <div>
                  <label className="label">Current Password</label>
                  <div className="relative">
                    <input
                      type={showPasswords ? 'text' : 'password'}
                      {...registerPwd('password')}
                      className="input pr-10"
                      placeholder="••••••••"
                    />
                    <button
                      type="button"
                      onClick={() => setShowPasswords(!showPasswords)}
                      className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400"
                    >
                      {showPasswords ? (
                        <EyeOff className="h-4 w-4" />
                      ) : (
                        <Eye className="h-4 w-4" />
                      )}
                    </button>
                  </div>
                  {pwdErrors.password && (
                    <p className="mt-1 text-xs text-red-500">{pwdErrors.password.message}</p>
                  )}
                </div>

                <div>
                  <label className="label">New Password</label>
                  <input
                    type={showPasswords ? 'text' : 'password'}
                    {...registerPwd('new_password')}
                    className="input"
                    placeholder="••••••••"
                  />
                  {pwdErrors.new_password && (
                    <p className="mt-1 text-xs text-red-500">{pwdErrors.new_password.message}</p>
                  )}
                </div>

                <div>
                  <label className="label">Confirm New Password</label>
                  <input
                    type={showPasswords ? 'text' : 'password'}
                    {...registerPwd('confirm_password')}
                    className="input"
                    placeholder="••••••••"
                  />
                  {pwdErrors.confirm_password && (
                    <p className="mt-1 text-xs text-red-500">
                      {pwdErrors.confirm_password.message}
                    </p>
                  )}
                </div>

                <button
                  type="submit"
                  disabled={pwdSubmitting}
                  className="btn-primary"
                >
                  {pwdSubmitting ? 'Updating...' : 'Update Password'}
                </button>
              </form>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

// Helper sub-components

interface AddressFormComponentProps {
  register: ReturnType<typeof useForm<AddressForm>>['register']
  errors: ReturnType<typeof useForm<AddressForm>>['formState']['errors']
  onSubmit: (e?: React.BaseSyntheticEvent) => Promise<void>
  onCancel: () => void
  isPending: boolean
  submitLabel: string
}

function AddressFormComponent({
  register,
  errors,
  onSubmit,
  onCancel,
  isPending,
  submitLabel,
}: AddressFormComponentProps) {
  return (
    <form onSubmit={onSubmit} className="bg-gray-50 rounded-xl p-4 mb-4 space-y-3">
      <div className="grid grid-cols-2 gap-3">
        <div>
          <label className="label">Full Name</label>
          <input {...register('name')} className="input" placeholder="John Doe" />
          {errors.name && <p className="mt-1 text-xs text-red-500">{errors.name.message}</p>}
        </div>
        <div>
          <label className="label">Phone</label>
          <input {...register('phone')} className="input" placeholder="+1234567890" />
          {errors.phone && <p className="mt-1 text-xs text-red-500">{errors.phone.message}</p>}
        </div>
      </div>
      <div>
        <label className="label">Street</label>
        <input {...register('street')} className="input" placeholder="123 Main St" />
        {errors.street && <p className="mt-1 text-xs text-red-500">{errors.street.message}</p>}
      </div>
      <div className="grid grid-cols-2 gap-3">
        <div>
          <label className="label">City</label>
          <input {...register('city')} className="input" placeholder="New York" />
          {errors.city && <p className="mt-1 text-xs text-red-500">{errors.city.message}</p>}
        </div>
        <div>
          <label className="label">Country</label>
          <input {...register('country')} className="input" placeholder="USA" />
          {errors.country && <p className="mt-1 text-xs text-red-500">{errors.country.message}</p>}
        </div>
      </div>
      <div className="flex gap-2">
        <button type="submit" disabled={isPending} className="btn-primary text-sm">
          {isPending ? 'Saving...' : submitLabel}
        </button>
        <button type="button" onClick={onCancel} className="btn-secondary text-sm">
          Cancel
        </button>
      </div>
    </form>
  )
}

interface EditAddressFormProps {
  address: Address
  onSave: (data: Partial<Address>) => void
  onCancel: () => void
  isPending: boolean
}

function EditAddressForm({ address, onSave, onCancel, isPending }: EditAddressFormProps) {
  const [form, setForm] = useState({
    name: address.name,
    phone: address.phone,
    street: address.street,
    city: address.city,
    country: address.country,
  })

  return (
    <div className="space-y-3">
      <div className="grid grid-cols-2 gap-3">
        <div>
          <label className="label">Full Name</label>
          <input
            value={form.name}
            onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
            className="input"
          />
        </div>
        <div>
          <label className="label">Phone</label>
          <input
            value={form.phone}
            onChange={(e) => setForm((f) => ({ ...f, phone: e.target.value }))}
            className="input"
          />
        </div>
      </div>
      <div>
        <label className="label">Street</label>
        <input
          value={form.street}
          onChange={(e) => setForm((f) => ({ ...f, street: e.target.value }))}
          className="input"
        />
      </div>
      <div className="grid grid-cols-2 gap-3">
        <div>
          <label className="label">City</label>
          <input
            value={form.city}
            onChange={(e) => setForm((f) => ({ ...f, city: e.target.value }))}
            className="input"
          />
        </div>
        <div>
          <label className="label">Country</label>
          <input
            value={form.country}
            onChange={(e) => setForm((f) => ({ ...f, country: e.target.value }))}
            className="input"
          />
        </div>
      </div>
      <div className="flex gap-2">
        <button
          onClick={() => onSave(form)}
          disabled={isPending}
          className="btn-primary text-sm"
        >
          {isPending ? 'Saving...' : 'Save Changes'}
        </button>
        <button onClick={onCancel} className="btn-secondary text-sm">
          Cancel
        </button>
      </div>
    </div>
  )
}
