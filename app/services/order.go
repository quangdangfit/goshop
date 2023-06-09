package services

import (
	"context"

	"github.com/jinzhu/copier"

	"goshop/app/models"
	"goshop/app/repositories"
	"goshop/app/serializers"
)

type IOrderService interface {
	PlaceOrder(ctx context.Context, req *serializers.PlaceOrderReq) (*models.Order, error)
	GetOrderByID(ctx context.Context, id string) (*models.Order, error)

	GetOrders(ctx context.Context, query *serializers.OrderQueryParam) (*[]models.Order, error)
	UpdateOrder(ctx context.Context, id string, req *serializers.PlaceOrderReq) (*models.Order, error)
}

type OrderService struct {
	repo        repositories.IOrderRepository
	productRepo repositories.IProductRepository
}

func NewOrderService(
	repo repositories.IOrderRepository,
	productRepo repositories.IProductRepository,
) IOrderService {
	return &OrderService{
		repo:        repo,
		productRepo: productRepo,
	}
}

func (s *OrderService) PlaceOrder(ctx context.Context, req *serializers.PlaceOrderReq) (*models.Order, error) {
	var lines []*models.OrderLine
	err := copier.Copy(&lines, &req.Lines)
	if err != nil {
		return nil, err
	}

	productMap := make(map[string]*models.Product)
	for _, line := range lines {
		product, err := s.productRepo.GetProductByID(ctx, line.ProductID)
		if err != nil {
			return nil, err
		}
		line.Price = product.Price * float64(line.Quantity)
		productMap[line.ProductID] = product
	}

	order, err := s.repo.CreateOrder(ctx, lines)
	if err != nil {
		return nil, err
	}

	for _, line := range order.Lines {
		line.Product = *productMap[line.ProductID]
	}

	return order, nil
}

func (s *OrderService) GetOrderByID(ctx context.Context, id string) (*models.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetOrders(ctx context.Context, query *serializers.OrderQueryParam) (*[]models.Order, error) {
	orders, err := s.repo.GetOrders(query)
	if err != nil {
		return nil, err
	}

	return orders, err
}

func (s *OrderService) UpdateOrder(ctx context.Context, id string, req *serializers.PlaceOrderReq) (*models.Order, error) {
	order, err := s.repo.UpdateOrder(id, req)
	if err != nil {
		return nil, err
	}

	return order, nil
}
