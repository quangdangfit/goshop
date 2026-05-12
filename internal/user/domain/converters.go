package domain

import "goshop/internal/user/model"

// UserFromModel maps the storage user to the API DTO, hiding internal fields such as
// password and deleted_at. Defined here so handlers don't need a JSON-roundtrip
// utils.Copy and the unreachable error branches it produces.
func UserFromModel(m *model.User) *User {
	if m == nil {
		return nil
	}
	return &User{
		ID:        m.ID,
		Email:     m.Email,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func AddressFromModel(m *model.Address) *Address {
	if m == nil {
		return nil
	}
	return &Address{
		ID:        m.ID,
		UserID:    m.UserID,
		Name:      m.Name,
		Phone:     m.Phone,
		Street:    m.Street,
		City:      m.City,
		Country:   m.Country,
		IsDefault: m.IsDefault,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func AddressesFromModel(rows []*model.Address) []*Address {
	out := make([]*Address, len(rows))
	for i, r := range rows {
		out[i] = AddressFromModel(r)
	}
	return out
}

func WishlistItemFromModel(m *model.Wishlist) *WishlistItem {
	if m == nil {
		return nil
	}
	return &WishlistItem{
		ID:        m.ID,
		ProductID: m.ProductID,
		CreatedAt: m.CreatedAt,
	}
}

func WishlistItemsFromModel(rows []*model.Wishlist) []*WishlistItem {
	out := make([]*WishlistItem, len(rows))
	for i, r := range rows {
		out[i] = WishlistItemFromModel(r)
	}
	return out
}
