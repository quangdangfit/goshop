package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
	"gorm.io/gorm"

	"goshop/internal/order/domain"
	"goshop/internal/order/model"
	orderRepo "goshop/internal/order/repository"
	"goshop/pkg/apperror"
	"goshop/pkg/dbs"
	"goshop/pkg/eventbus"
	"goshop/pkg/notification"
	"goshop/pkg/paging"
)

// ReservationTTL is how long a placed order holds reserved stock before the sweeper releases
// it and cancels the order. Tuned for an interactive checkout — long enough to clear payment,
// short enough that abandoned carts don't starve other shoppers.
const ReservationTTL = 15 * time.Minute

// LowStockThreshold is the available-stock floor at or below which a LowStock event is
// published after a commit. Matches the FE badge threshold so admin alerts and customer
// "Low stock" badges fire on the same boundary.
const LowStockThreshold = 5

//go:generate mockery --name=OrderService
type OrderService interface {
	PlaceOrder(ctx context.Context, req *domain.PlaceOrderReq) (*model.Order, error)
	GetOrderByID(ctx context.Context, id string) (*model.Order, error)
	GetMyOrders(ctx context.Context, req *domain.ListOrderReq) ([]*model.Order, *paging.Pagination, error)
	CancelOrder(ctx context.Context, orderID, userID string) (*model.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status model.OrderStatus) (*model.Order, error)
	// MarkOrderPaid commits the order's stock reservations and flips its status to paid. Idempotent
	// per order: a second call on an already-paid order is a no-op.
	MarkOrderPaid(ctx context.Context, orderID string) (*model.Order, error)
	// SweepExpiredReservations releases reservations past their TTL whose parent order is still
	// pending payment, and cancels those orders. Returns the number of reservations released.
	SweepExpiredReservations(ctx context.Context, batchSize int) (int, error)
}

type orderService struct {
	validator       validation.Validation
	db              dbs.Database
	repo            orderRepo.OrderRepository
	productRepo     orderRepo.ProductRepository
	userRepo        orderRepo.UserRepository
	reservationRepo orderRepo.ReservationRepository
	couponSvc       CouponService
	notifier        notification.Notifier
	bus             eventbus.Bus
}

// SetEventBus wires an optional event bus. When set, the service publishes domain
// events (currently LowStock after commit) without coupling notification transports
// into the order service itself. Safe to leave unset for tests.
func (s *orderService) SetEventBus(bus eventbus.Bus) { s.bus = bus }

func NewOrderService(
	validator validation.Validation,
	db dbs.Database,
	repo orderRepo.OrderRepository,
	productRepo orderRepo.ProductRepository,
	userRepo orderRepo.UserRepository,
	reservationRepo orderRepo.ReservationRepository,
	couponSvc CouponService,
	notifier notification.Notifier,
) OrderService {
	return &orderService{
		validator:       validator,
		db:              db,
		repo:            repo,
		productRepo:     productRepo,
		userRepo:        userRepo,
		reservationRepo: reservationRepo,
		couponSvc:       couponSvc,
		notifier:        notifier,
	}
}

func (s *orderService) PlaceOrder(ctx context.Context, req *domain.PlaceOrderReq) (*model.Order, error) {
	if err := s.validator.ValidateStruct(req); err != nil {
		return nil, err
	}

	lines := make([]*model.OrderLine, len(req.Lines))
	for i, l := range req.Lines {
		lines[i] = &model.OrderLine{ProductID: l.ProductID, Quantity: l.Quantity}
	}

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

	var discountAmount float64
	var couponCode string
	var couponID string
	if req.CouponCode != "" {
		discount, coupon, err := s.couponSvc.Apply(ctx, req.CouponCode, totalPrice)
		if err != nil {
			return nil, err
		}
		discountAmount = discount
		couponCode = coupon.Code
		couponID = coupon.ID
	}

	// Reserve stock + create order + persist reservations + bump coupon usage atomically.
	// Reservations hold inventory until payment clears or the sweeper releases them.
	var order *model.Order
	expiresAt := time.Now().Add(ReservationTTL)
	txErr := s.db.WithTransaction(func() error {
		o, err := s.repo.CreateOrder(ctx, req.UserID, lines, couponCode, discountAmount)
		if err != nil {
			return err
		}
		o.Status = model.OrderStatusPendingPayment
		if err := s.repo.UpdateOrder(ctx, o); err != nil {
			return err
		}

		reservations := make([]*model.StockReservation, 0, len(lines))
		for _, line := range lines {
			qty := int(line.Quantity) //nolint:gosec // bounded by validation (lte=5 lines, uint qty)
			if err := s.productRepo.ReserveStock(ctx, line.ProductID, qty); err != nil {
				if errors.Is(err, orderRepo.ErrInsufficientStock) {
					return &InsufficientStockError{ProductID: line.ProductID, Requested: qty}
				}
				return fmt.Errorf("reserve stock for product %s: %w", line.ProductID, err)
			}
			reservations = append(reservations, &model.StockReservation{
				OrderID:   o.ID,
				ProductID: line.ProductID,
				Quantity:  qty,
				Status:    model.ReservationStatusActive,
				ExpiresAt: expiresAt,
			})
		}
		if err := s.reservationRepo.CreateMany(ctx, reservations); err != nil {
			return fmt.Errorf("persist reservations: %w", err)
		}

		if couponID != "" {
			if err := s.couponSvc.IncrUsedCount(ctx, couponID); err != nil {
				return fmt.Errorf("increment coupon usage: %w", err)
			}
		}
		order = o
		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	for _, line := range order.Lines {
		line.Product = productMap[line.ProductID]
	}

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

func (s *orderService) GetMyOrders(ctx context.Context, req *domain.ListOrderReq) ([]*model.Order, *paging.Pagination, error) {
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

	if !status.IsValid() {
		return nil, apperror.ErrInvalidStatus
	}
	if !order.Status.CanTransitionTo(status) {
		return nil, apperror.ErrInvalidStatus
	}

	order.Status = status
	if err := s.repo.UpdateOrder(ctx, order); err != nil {
		return nil, err
	}

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

func (s *orderService) MarkOrderPaid(ctx context.Context, orderID string) (*model.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, orderID, true)
	if err != nil {
		return nil, err
	}
	if order.Status == model.OrderStatusPaid {
		return order, nil // idempotent: already committed
	}
	if order.Status == model.OrderStatusCancelled || order.Status == model.OrderStatusPaymentFailed {
		return nil, apperror.ErrInvalidStatus
	}

	var committedProductIDs []string
	txErr := s.db.WithTransaction(func() error {
		reservations, err := s.reservationRepo.FindActiveByOrderID(ctx, order.ID)
		if err != nil {
			return err
		}
		for _, res := range reservations {
			if err := s.productRepo.CommitReservation(ctx, res.ProductID, res.Quantity); err != nil {
				return fmt.Errorf("commit reservation %s: %w", res.ID, err)
			}
			committedProductIDs = append(committedProductIDs, res.ProductID)
		}
		ids := reservationIDs(reservations)
		if err := s.reservationRepo.UpdateStatus(ctx, ids, model.ReservationStatusCommitted); err != nil {
			return err
		}
		order.Status = model.OrderStatusPaid
		return s.repo.UpdateOrder(ctx, order)
	})
	if txErr != nil {
		return nil, txErr
	}
	s.publishLowStock(ctx, committedProductIDs)
	return order, nil
}

// publishLowStock fires a LowStock event for each product whose available stock is at
// or below the threshold. Best-effort: a fetch failure is logged but does not fail the
// commit, since payment has already cleared.
func (s *orderService) publishLowStock(ctx context.Context, productIDs []string) {
	if len(productIDs) == 0 {
		return
	}
	bus := s.bus
	if bus == nil {
		bus = eventbus.Default()
	}
	for _, pid := range productIDs {
		p, err := s.productRepo.GetProductByID(ctx, pid)
		if err != nil {
			logger.Error("low_stock check: get product failed", err)
			continue
		}
		available := p.StockQuantity - p.ReservedQuantity
		if available <= LowStockThreshold {
			bus.Publish(ctx, eventbus.LowStock{
				ProductID: pid,
				Available: available,
				Threshold: LowStockThreshold,
			})
		}
	}
}

func (s *orderService) CancelOrder(ctx context.Context, orderID, userID string) (*model.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, orderID, false)
	if err != nil {
		return nil, err
	}

	if userID != order.UserID {
		return nil, apperror.ErrForbidden
	}

	if order.Status == model.OrderStatusDone ||
		order.Status == model.OrderStatusCancelled ||
		order.Status == model.OrderStatusPaid {
		return nil, apperror.ErrInvalidStatus
	}

	txErr := s.db.WithTransaction(func() error {
		if err := s.releaseActiveReservations(ctx, order.ID); err != nil {
			return err
		}
		order.Status = model.OrderStatusCancelled
		return s.repo.UpdateOrder(ctx, order)
	})
	if txErr != nil {
		return nil, txErr
	}
	return order, nil
}

func (s *orderService) SweepExpiredReservations(ctx context.Context, batchSize int) (int, error) {
	if batchSize <= 0 {
		batchSize = 100
	}
	expired, err := s.reservationRepo.FindExpired(ctx, time.Now(), batchSize)
	if err != nil {
		return 0, err
	}
	if len(expired) == 0 {
		return 0, nil
	}

	// Group by order so we can cancel each parent once.
	byOrder := make(map[string][]*model.StockReservation, len(expired))
	for _, r := range expired {
		byOrder[r.OrderID] = append(byOrder[r.OrderID], r)
	}

	released := 0
	for orderID, group := range byOrder {
		orderMissing := false
		txErr := s.db.WithTransaction(func() error {
			order, err := s.repo.GetOrderByID(ctx, orderID, false)
			if err != nil {
				// Orphaned reservation: the order row was deleted (or soft-deleted) but the
				// reservation lived on. Still release the held stock and mark the rows so
				// they don't keep showing up in the sweep — there's just no order to cancel.
				if errors.Is(err, gorm.ErrRecordNotFound) {
					orderMissing = true
				} else {
					return err
				}
			} else if order.Status != model.OrderStatusPendingPayment && order.Status != model.OrderStatusNew {
				// Paid / fulfilled order's reservations should have been committed already,
				// not released. Skip without touching stock.
				return nil
			}
			for _, r := range group {
				if err := s.productRepo.ReleaseReservation(ctx, r.ProductID, r.Quantity); err != nil {
					return fmt.Errorf("release reservation %s: %w", r.ID, err)
				}
			}
			ids := reservationIDs(group)
			if err := s.reservationRepo.UpdateStatus(ctx, ids, model.ReservationStatusReleased); err != nil {
				return err
			}
			if orderMissing {
				return nil
			}
			order.Status = model.OrderStatusCancelled
			return s.repo.UpdateOrder(ctx, order)
		})
		if txErr != nil {
			logger.Errorf("sweep order %s: %s", orderID, txErr)
			continue
		}
		if orderMissing {
			logger.Warnf("sweep order %s: order row missing, released %d orphaned reservation(s)", orderID, len(group))
			released += len(group)
			continue
		}
		released += len(group)

		// Notify the buyer their reservation was released. Best-effort: failures only log.
		go func(orderID string) {
			order, err := s.repo.GetOrderByID(ctx, orderID, false)
			if err != nil {
				return
			}
			user, err := s.userRepo.GetUserByID(ctx, order.UserID)
			if err != nil || user == nil {
				return
			}
			if err := s.notifier.SendOrderStatusChanged(ctx, order.ID, user.Email, string(model.OrderStatusCancelled)); err != nil {
				logger.Errorf("sweep notify %s: %s", order.ID, err)
			}
		}(orderID)
	}
	return released, nil
}

func (s *orderService) releaseActiveReservations(ctx context.Context, orderID string) error {
	reservations, err := s.reservationRepo.FindActiveByOrderID(ctx, orderID)
	if err != nil {
		return err
	}
	for _, res := range reservations {
		if err := s.productRepo.ReleaseReservation(ctx, res.ProductID, res.Quantity); err != nil {
			return fmt.Errorf("release reservation %s: %w", res.ID, err)
		}
	}
	return s.reservationRepo.UpdateStatus(ctx, reservationIDs(reservations), model.ReservationStatusReleased)
}

func reservationIDs(rs []*model.StockReservation) []string {
	ids := make([]string, 0, len(rs))
	for _, r := range rs {
		ids = append(ids, r.ID)
	}
	return ids
}
