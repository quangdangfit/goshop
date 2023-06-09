package repositories

import (
	"context"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"goshop/app/dbs"
	"goshop/app/models"
	"goshop/app/serializers"
	"goshop/config"
	"goshop/pkg/paging"
)

type IOrderRepository interface {
	CreateOrder(ctx context.Context, userID string, lines []*models.OrderLine) (*models.Order, error)
	GetOrderByID(ctx context.Context, id string, preload bool) (*models.Order, error)
	GetMyOrders(ctx context.Context, req *serializers.ListOrderReq) ([]*models.Order, *paging.Pagination, error)
	UpdateOrder(ctx context.Context, order *models.Order) error
}

type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepository() *OrderRepo {
	return &OrderRepo{db: dbs.Database}
}

func (r *OrderRepo) CreateOrder(ctx context.Context, userID string, lines []*models.OrderLine) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	order := new(models.Order)
	err := r.WithTransaction(func(*gorm.DB) error {
		// Create Order
		var totalPrice float64
		for _, line := range lines {
			totalPrice += line.Price
		}
		order.TotalPrice = totalPrice
		order.UserID = userID

		if err := r.db.Create(order).Error; err != nil {
			return err
		}

		// Create order lines
		for _, line := range lines {
			line.OrderID = order.ID
		}
		if err := r.db.CreateInBatches(&lines, len(lines)).Error; err != nil {
			return err
		}

		err := copier.Copy(&order.Lines, &lines)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepo) GetOrderByID(ctx context.Context, id string, preload bool) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	var order models.Order
	query := r.db.Where("id = ?", id)
	if preload {
		query = query.Preload("Lines").Preload("Lines.Product")
	}
	if err := query.First(&order).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepo) GetMyOrders(ctx context.Context, req *serializers.ListOrderReq) ([]*models.Order, *paging.Pagination, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	query := r.db.Where("user_id = ?", req.UserID)
	order := "id"
	if req.Code != "" {
		query = query.Where("code = ?", req.Code)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.OrderBy != "" {
		order = req.OrderBy
		if req.OrderDesc {
			order += " DESC"
		}
	}

	var total int64
	if err := query.Model(&models.Order{}).Count(&total).Error; err != nil {
		return nil, nil, err
	}

	pagination := paging.New(req.Page, req.Limit, total)

	var orders []*models.Order
	if err := query.Preload("Lines").
		Preload("Lines.Product").
		Find(&orders).Error; err != nil {
		return nil, nil, err
	}

	return orders, pagination, nil
}

func (r *OrderRepo) UpdateOrder(ctx context.Context, order *models.Order) error {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	if err := r.db.Save(&order).Error; err != nil {
		return err
	}

	return nil
}

func (r *OrderRepo) WithTransaction(callback func(*gorm.DB) error) error {
	tx := r.db.Begin()

	if err := callback(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
