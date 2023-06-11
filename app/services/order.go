package services

import (
	"context"
	"errors"

	"goshop/app/models"
	"goshop/app/repositories"
	"goshop/app/serializers"
	"goshop/pkg/paging"
	"goshop/pkg/utils"
)

type IOrderService interface {
	PlaceOrder(ctx context.Context, req *serializers.PlaceOrderReq) (*models.Order, error)
	GetOrderByID(ctx context.Context, id string) (*models.Order, error)
	GetMyOrders(ctx context.Context, req *serializers.ListOrderReq) ([]*models.Order, *paging.Pagination, error)
	CancelOrder(ctx context.Context, orderID, userID string) (*models.Order, error)
}

type OrderService struct {
	repo        repositories.IOrderRepository
	productRepo repositories.IProductRepository
}

func NewOrderService(
	repo repositories.IOrderRepository,
	productRepo repositories.IProductRepository,
) *OrderService {
	return &OrderService{
		repo:        repo,
		productRepo: productRepo,
	}
}

func (s *OrderService) PlaceOrder(ctx context.Context, req *serializers.PlaceOrderReq) (*models.Order, error) {
	var lines []*models.OrderLine
	utils.Copy(&lines, &req.Lines)

	productMap := make(map[string]*models.Product)
	for _, line := range lines {
		product, err := s.productRepo.GetProductByID(ctx, line.ProductID)
		if err != nil {
			return nil, err
		}
		line.Price = product.Price * float64(line.Quantity)
		productMap[line.ProductID] = product
	}

	order, err := s.repo.CreateOrder(ctx, req.UserID, lines)
	if err != nil {
		return nil, err
	}

	for _, line := range order.Lines {
		line.Product = productMap[line.ProductID]
	}

	return order, nil
}

func (s *OrderService) GetOrderByID(ctx context.Context, id string) (*models.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, id, true)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetMyOrders(ctx context.Context, req *serializers.ListOrderReq) ([]*models.Order, *paging.Pagination, error) {
	orders, pagination, err := s.repo.GetMyOrders(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	return orders, pagination, err
}

func (s *OrderService) CancelOrder(ctx context.Context, orderID, userID string) (*models.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, orderID, false)
	if err != nil {
		return nil, err
	}

	if userID != order.UserID {
		return nil, errors.New("permission denied")
	}

	if order.Status == models.OrderStatusDone || order.Status == models.OrderStatusCancelled {
		return nil, errors.New("invalid order status")
	}

	order.Status = models.OrderStatusCancelled
	err = s.repo.UpdateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}
