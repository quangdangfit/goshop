package service

import (
	"context"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"

	"goshop/internal/cart/dto"
	"goshop/internal/cart/model"
	"goshop/internal/cart/repository"
)

//go:generate mockery --name=ICartService
type ICartService interface {
	AddProduct(ctx context.Context, req *dto.AddProductReq) (*model.Cart, error)
	GetCartByUserID(ctx context.Context, userID string) (*model.Cart, error)
	RemoveProduct(ctx context.Context, req *dto.RemoveProductReq) (*model.Cart, error)
}

type CartService struct {
	validator validation.Validation
	repo      repository.ICartRepository
}

func NewCartService(
	validator validation.Validation,
	repo repository.ICartRepository,
) *CartService {
	return &CartService{
		validator: validator,
		repo:      repo,
	}
}

func (p *CartService) GetCartByUserID(ctx context.Context, userID string) (*model.Cart, error) {
	cart, err := p.repo.GetCartByUserID(ctx, userID)
	if err != nil {
		cart = &model.Cart{
			UserID: userID,
		}
		err = p.repo.Create(ctx, cart)
		if err != nil {
			return nil, err
		}
		return cart, err
	}

	return cart, nil
}

func (p *CartService) AddProduct(ctx context.Context, req *dto.AddProductReq) (*model.Cart, error) {
	if err := p.validator.ValidateStruct(req); err != nil {
		return nil, err
	}

	cart, err := p.repo.GetCartByUserID(ctx, req.UserID)
	if err != nil {
		cart = &model.Cart{
			UserID: req.UserID,
			Lines: []*model.CartLine{{
				ProductID: req.Line.ProductID,
				Quantity:  req.Line.Quantity,
			}},
		}
		err = p.repo.Create(ctx, cart)
		if err != nil {
			return nil, err
		}
		return cart, err
	}

	for _, line := range cart.Lines {
		if line.ProductID == req.Line.ProductID {
			return cart, nil
		}
	}

	cart.Lines = append(cart.Lines, &model.CartLine{
		ProductID: req.Line.ProductID,
		Quantity:  req.Line.Quantity,
	})

	err = p.repo.Update(ctx, cart)
	if err != nil {
		logger.Errorf("AddProductReq.Update fail, userID: %s, error: %s", req.UserID, err)
		return nil, err
	}

	return cart, nil
}

func (p *CartService) RemoveProduct(ctx context.Context, req *dto.RemoveProductReq) (*model.Cart, error) {
	if err := p.validator.ValidateStruct(req); err != nil {
		return nil, err
	}

	cart, err := p.repo.GetCartByUserID(ctx, req.UserID)
	if err != nil {
		cart = &model.Cart{
			UserID: req.UserID,
		}
		err = p.repo.Create(ctx, cart)
		if err != nil {
			return nil, err
		}
		return cart, err
	}

	for i, line := range cart.Lines {
		if line.ProductID == req.ProductID {
			cart.Lines = append(cart.Lines[:i], cart.Lines[i+1:]...)
			break
		}
	}

	err = p.repo.Update(ctx, cart)
	if err != nil {
		logger.Errorf("RemoveProductReq.Update fail, userID: %s, error: %s", req.UserID, err)
		return nil, err
	}

	return cart, nil
}
