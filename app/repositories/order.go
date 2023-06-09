package repositories

import (
	"context"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"goshop/app/models"
	"goshop/app/serializers"
	"goshop/config"
	"goshop/dbs"
)

type IOrderRepository interface {
	CreateOrder(ctx context.Context, lines []*models.OrderLine) (*models.Order, error)
	GetOrderByID(ctx context.Context, id string) (*models.Order, error)

	GetOrders(query *serializers.OrderQueryParam) (*[]models.Order, error)
	UpdateOrder(id string, req *serializers.PlaceOrderReq) (*models.Order, error)
	AssignOrder(id string) error
}

type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepository() *OrderRepo {
	return &OrderRepo{db: dbs.Database}
}

func (r *OrderRepo) CreateOrder(ctx context.Context, lines []*models.OrderLine) (*models.Order, error) {
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
		order.Lines = lines

		return nil
	})
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepo) GetOrderByID(ctx context.Context, id string) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	var order models.Order
	if err := r.db.Where("id = ?", id).
		Preload("Lines").
		Preload("Lines.Product").
		First(&order).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepo) GetOrders(query *serializers.OrderQueryParam) (*[]models.Order, error) {
	var orders []models.Order
	if err := r.db.Find(&orders, query).Error; err != nil {
		return nil, err
	}

	return &orders, nil
}

func (r *OrderRepo) UpdateOrder(id string, req *serializers.PlaceOrderReq) (*models.Order, error) {
	order, err := r.GetOrderByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	copier.Copy(order, &req)
	if err := r.db.Save(&order).Error; err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepo) AssignOrder(id string) error {
	order, err := r.GetOrderByID(context.Background(), id)
	if err != nil {
		return err
	}

	order.Status = models.OrderStatusInProgress
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
