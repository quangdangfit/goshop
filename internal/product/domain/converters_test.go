package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"goshop/internal/product/model"
)

func TestCategoryFromModel_Nil(t *testing.T) {
	require.Nil(t, CategoryFromModel(nil))
}

func TestCategoryFromModel_PopulatesFields(t *testing.T) {
	out := CategoryFromModel(&model.Category{ID: "c1", Name: "Books", Slug: "books"})
	require.Equal(t, "c1", out.ID)
	require.Equal(t, "Books", out.Name)
}

func TestCategoriesFromModel(t *testing.T) {
	out := CategoriesFromModel([]*model.Category{{ID: "c1"}, {ID: "c2"}})
	require.Len(t, out, 2)
	require.Equal(t, "c1", out[0].ID)
}

func TestReviewFromModel_Nil(t *testing.T) {
	require.Nil(t, ReviewFromModel(nil))
}

func TestReviewFromModel_PopulatesFields(t *testing.T) {
	out := ReviewFromModel(&model.Review{ID: "r1", Rating: 4, Comment: "ok"})
	require.Equal(t, "r1", out.ID)
	require.Equal(t, 4, out.Rating)
}

func TestReviewsFromModel(t *testing.T) {
	out := ReviewsFromModel([]*model.Review{{ID: "r1"}, {ID: "r2"}})
	require.Len(t, out, 2)
}
