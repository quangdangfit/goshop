package service

import (
	"context"

	"goshop/internal/user/domain"
	"goshop/internal/user/model"
	"goshop/internal/user/repository"
)

//go:generate mockery --name=WishlistService
type WishlistService interface {
	GetWishlist(ctx context.Context, userID string) ([]*model.Wishlist, error)
	AddProduct(ctx context.Context, userID string, req *domain.AddToWishlistReq) error
	RemoveProduct(ctx context.Context, userID, productID string) error
}

type wishlistSvc struct {
	repo repository.WishlistRepository
}

func NewWishlistService(repo repository.WishlistRepository) WishlistService {
	return &wishlistSvc{repo: repo}
}

func (s *wishlistSvc) GetWishlist(ctx context.Context, userID string) ([]*model.Wishlist, error) {
	return s.repo.GetWishlist(ctx, userID)
}

func (s *wishlistSvc) AddProduct(ctx context.Context, userID string, req *domain.AddToWishlistReq) error {
	return s.repo.Add(ctx, userID, req.ProductID)
}

func (s *wishlistSvc) RemoveProduct(ctx context.Context, userID, productID string) error {
	return s.repo.Remove(ctx, userID, productID)
}
