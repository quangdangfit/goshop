package service

import (
	"context"
	"errors"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"

	"goshop/internal/order/dto"
	"goshop/internal/order/model"
	"goshop/internal/order/repository"
	"goshop/pkg/notification"
	"goshop/pkg/paging"
	"goshop/pkg/utils"
)

//go:generate mockery --name=OrderService
type OrderService interface {
	PlaceOrder(ctx context.Context, req *dto.PlaceOrderReq) (*model.Order, error)
	GetOrderByID(ctx context.Context, id string) (*model.Order, error)
	GetMyOrders(ctx context.Context, req *dto.ListOrderReq) ([]*model.Order, *paging.Pagination, error)
	CancelOrder(ctx context.Context, orderID, userID string) (*model.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status model.OrderStatus) (*model.Order, error)
}

type orderService struct {
	validator   validation.Validation
	repo        repository.OrderRepository
	productRepo repository.ProductRepository
	userRepo    repository.UserRepository
	couponSvc   CouponService
	notifier    notification.Notifier
}

func NewOrderService(
	validator validation.Validation,
	repo repository.OrderRepository,
	productRepo repository.ProductRepository,
	userRepo repository.UserRepository,
	couponSvc CouponService,
	notifier notification.Notifier,
) OrderService {
	return &orderService{
		validator:   validator,
		repo:        repo,
		productRepo: productRepo,
		userRepo:    userRepo,
		couponSvc:   couponSvc,
		notifier:    notifier,
	}
}

func (s *orderService) PlaceOrder(ctx context.Context, req *dto.PlaceOrderReq) (*model.Order, error) {
	if err := s.validator.ValidateStruct(req); err != nil {
		return nil, err
	}

	var lines []*model.OrderLine
	utils.Copy(&lines, &req.Lines)

	productMap := make(map[string]*model.Product)
	for _, line := range lines {
		product, err := s.productRepo.GetProductByID(ctx, line.ProductID)
		if err != nil {
			return nil, err
		}
		line.Price = product.Price * float64(line.Quantity)
		productMap[line.ProductID] = product
	}

	var totalPrice float64
	for _, line := range lines {
		totalPrice += line.Price
	}

	// Apply coupon if provided
	var discountAmount float64
	var couponCode string
	if req.CouponCode != "" {
		discount, coupon, err := s.couponSvc.Apply(ctx, req.CouponCode, totalPrice)
		if err != nil {
			return nil, err
		}
		discountAmount = discount
		couponCode = coupon.Code
		// Increment usage count (best effort)
		if err := s.couponSvc.IncrUsedCount(ctx, coupon.ID); err != nil {
			logger.Error("Failed to increment coupon usage: ", err)
		}
	}

	order, err := s.repo.CreateOrder(ctx, req.UserID, lines, couponCode, discountAmount)
	if err != nil {
		return nil, err
	}

	for _, line := range order.Lines {
		line.Product = productMap[line.ProductID]
	}

	// Decrement stock (best effort — non-blocking)
	for _, line := range lines {
		if err := s.productRepo.DecrementStock(ctx, line.ProductID, int(line.Quantity)); err != nil {
			logger.Errorf("Failed to decrement stock for product %s: %s", line.ProductID, err)
		}
	}

	// Send notification (best effort)
	go func() {
		user, err := s.userRepo.GetUserByID(ctx, req.UserID)
		if err != nil {
			logger.Error("Failed to get user for notification: ", err)
			return
		}
		if err := s.notifier.SendOrderPlaced(ctx, order.ID, user.Email); err != nil {
			logger.Error("Failed to send order placed notification: ", err)
		}
	}()

	return order, nil
}

func (s *orderService) GetOrderByID(ctx context.Context, id string) (*model.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, id, true)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *orderService) GetMyOrders(ctx context.Context, req *dto.ListOrderReq) ([]*model.Order, *paging.Pagination, error) {
	orders, pagination, err := s.repo.GetMyOrders(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	return orders, pagination, err
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, orderID string, status model.OrderStatus) (*model.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, orderID, false)
	if err != nil {
		return nil, err
	}

	order.Status = status
	if err := s.repo.UpdateOrder(ctx, order); err != nil {
		return nil, err
	}

	// Send notification (best effort)
	go func() {
		user, err := s.userRepo.GetUserByID(ctx, order.UserID)
		if err != nil {
			logger.Error("Failed to get user for notification: ", err)
			return
		}
		if err := s.notifier.SendOrderStatusChanged(ctx, order.ID, user.Email, string(status)); err != nil {
			logger.Error("Failed to send order status changed notification: ", err)
		}
	}()

	return order, nil
}

func (s *orderService) CancelOrder(ctx context.Context, orderID, userID string) (*model.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, orderID, false)
	if err != nil {
		return nil, err
	}

	if userID != order.UserID {
		return nil, errors.New("permission denied")
	}

	if order.Status == model.OrderStatusDone || order.Status == model.OrderStatusCancelled {
		return nil, errors.New("invalid order status")
	}

	order.Status = model.OrderStatusCancelled
	err = s.repo.UpdateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}
