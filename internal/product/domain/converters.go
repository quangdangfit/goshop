package domain

import "goshop/internal/product/model"

// CategoryFromModel maps the storage model to the API DTO. Defined here to keep
// the handler layer free of conversion boilerplate and to remove the unreachable
// utils.Copy error branches that JSON-roundtrip conversion produced.
func CategoryFromModel(m *model.Category) *Category {
	if m == nil {
		return nil
	}
	return &Category{
		ID:          m.ID,
		Name:        m.Name,
		Slug:        m.Slug,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func CategoriesFromModel(rows []*model.Category) []*Category {
	out := make([]*Category, len(rows))
	for i, r := range rows {
		out[i] = CategoryFromModel(r)
	}
	return out
}

func ReviewFromModel(m *model.Review) *Review {
	if m == nil {
		return nil
	}
	return &Review{
		ID:        m.ID,
		UserID:    m.UserID,
		ProductID: m.ProductID,
		Rating:    m.Rating,
		Comment:   m.Comment,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ReviewsFromModel(rows []*model.Review) []*Review {
	out := make([]*Review, len(rows))
	for i, r := range rows {
		out[i] = ReviewFromModel(r)
	}
	return out
}
