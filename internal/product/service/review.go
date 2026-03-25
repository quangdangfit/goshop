package service

import (
	"context"

	"github.com/quangdangfit/gocommon/validation"

	"goshop/internal/product/dto"
	"goshop/internal/product/model"
	"goshop/internal/product/repository"
	"goshop/pkg/apperror"
	"goshop/pkg/paging"
	"goshop/pkg/utils"
)

//go:generate mockery --name=ReviewService
type ReviewService interface {
	ListReviews(ctx context.Context, productID string, page, limit int64) ([]*model.Review, *paging.Pagination, error)
	CreateReview(ctx context.Context, productID, userID string, req *dto.CreateReviewReq) (*model.Review, error)
	UpdateReview(ctx context.Context, id, userID string, req *dto.UpdateReviewReq) (*model.Review, error)
	DeleteReview(ctx context.Context, id, userID string) error
}

type reviewSvc struct {
	validator   validation.Validation
	repo        repository.ReviewRepository
	productRepo repository.ProductRepository
}

func NewReviewService(
	validator validation.Validation,
	repo repository.ReviewRepository,
	productRepo repository.ProductRepository,
) ReviewService {
	return &reviewSvc{validator: validator, repo: repo, productRepo: productRepo}
}

func (s *reviewSvc) ListReviews(ctx context.Context, productID string, page, limit int64) ([]*model.Review, *paging.Pagination, error) {
	return s.repo.ListByProduct(ctx, productID, page, limit)
}

func (s *reviewSvc) CreateReview(ctx context.Context, productID, userID string, req *dto.CreateReviewReq) (*model.Review, error) {
	if err := s.validator.ValidateStruct(req); err != nil {
		return nil, err
	}
	var review model.Review
	utils.Copy(&review, req)
	review.ProductID = productID
	review.UserID = userID
	if err := s.repo.Create(ctx, &review); err != nil {
		return nil, err
	}
	s.updateProductRating(ctx, productID)
	return &review, nil
}

func (s *reviewSvc) UpdateReview(ctx context.Context, id, userID string, req *dto.UpdateReviewReq) (*model.Review, error) {
	review, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if review.UserID != userID {
		return nil, apperror.ErrForbidden
	}
	utils.Copy(review, req)
	if err := s.repo.Update(ctx, review); err != nil {
		return nil, err
	}
	s.updateProductRating(ctx, review.ProductID)
	return review, nil
}

func (s *reviewSvc) DeleteReview(ctx context.Context, id, userID string) error {
	review, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if review.UserID != userID {
		return apperror.ErrForbidden
	}
	if err := s.repo.Delete(ctx, id, userID); err != nil {
		return err
	}
	s.updateProductRating(ctx, review.ProductID)
	return nil
}

func (s *reviewSvc) updateProductRating(ctx context.Context, productID string) {
	avg, count, err := s.repo.GetAggregates(ctx, productID)
	if err != nil {
		return
	}
	_ = s.productRepo.UpdateRating(ctx, productID, avg, count)
}
