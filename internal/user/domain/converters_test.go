package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"goshop/internal/user/model"
)

func TestUserFromModel_Nil(t *testing.T) {
	require.Nil(t, UserFromModel(nil))
}

func TestUserFromModel_PopulatesFields(t *testing.T) {
	got := UserFromModel(&model.User{ID: "u1", Email: "x@example.com"})
	require.Equal(t, "u1", got.ID)
	require.Equal(t, "x@example.com", got.Email)
}

func TestAddressFromModel_Nil(t *testing.T) {
	require.Nil(t, AddressFromModel(nil))
}

func TestAddressFromModel_PopulatesFields(t *testing.T) {
	got := AddressFromModel(&model.Address{ID: "a1", UserID: "u1", Name: "Home", Phone: "1", Street: "s", City: "c", Country: "vn", IsDefault: true})
	require.Equal(t, "a1", got.ID)
	require.True(t, got.IsDefault)
}

func TestAddressesFromModel(t *testing.T) {
	out := AddressesFromModel([]*model.Address{{ID: "a1"}, {ID: "a2"}})
	require.Len(t, out, 2)
}

func TestWishlistItemFromModel_Nil(t *testing.T) {
	require.Nil(t, WishlistItemFromModel(nil))
}

func TestWishlistItemFromModel_PopulatesFields(t *testing.T) {
	got := WishlistItemFromModel(&model.Wishlist{ID: "w1", ProductID: "p1"})
	require.Equal(t, "w1", got.ID)
	require.Equal(t, "p1", got.ProductID)
}

func TestWishlistItemsFromModel(t *testing.T) {
	out := WishlistItemsFromModel([]*model.Wishlist{{ID: "w1"}, {ID: "w2"}})
	require.Len(t, out, 2)
}
