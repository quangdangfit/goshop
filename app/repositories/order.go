package repositories

import (
	"errors"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"goshop/app/models"
	"goshop/app/serializers"
	"goshop/dbs"
)

type IOrderRepository interface {
	GetOrders(query *serializers.OrderQueryParam) (*[]models.Order, error)
	GetOrderByID(id string) (*models.Order, error)
	CreateOrder(item *serializers.OrderBodyParam) (*models.Order, error)
	UpdateOrder(id string, item *serializers.OrderBodyParam) (*models.Order, error)
	AssignOrder(id string) error
}

type OrderRepo struct {
	db           *gorm.DB
	lineRepo     IOrderLineRepository
	quantityRepo IQuantityRepository
}

func NewOrderRepository() *OrderRepo {
	return &OrderRepo{db: dbs.Database, lineRepo: NewOrderLineRepository()}
}

func (r *OrderRepo) GetOrders(query *serializers.OrderQueryParam) (*[]models.Order, error) {
	var orders []models.Order
	if err := r.db.Find(&orders, query).Error; err != nil {
		return nil, err
	}

	return &orders, nil
}

func (r *OrderRepo) GetOrderByID(id string) (*models.Order, error) {
	var order models.Order
	var lines []models.OrderLine
	if err := r.db.Where("id = ?", id).First(&order).Error; err != nil {
		return nil, errors.New("not found order")
	}
	r.db.Where("order_id = ?", id).Find(&lines)
	order.Lines = lines

	return &order, nil
}

func (r *OrderRepo) CreateOrder(item *serializers.OrderBodyParam) (*models.Order, error) {
	var order models.Order
	copier.Copy(&order, &item)

	if order.Lines == nil {
		return nil, errors.New("order lines must be not empty")
	}

	if err := r.db.Create(&order).Error; err != nil {
		return nil, err
	}

	var lines []models.OrderLine
	var totalPrice float64
	for _, line := range order.Lines {
		line.OrderID = order.ID
		if err := r.db.Create(&line).Error; err != nil {
			return nil, err
		}
		lines = append(lines, line)
		totalPrice += line.Price
	}
	order.TotalPrice = totalPrice
	order.Lines = lines

	if err := r.db.Save(&order).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepo) UpdateOrder(id string, item *serializers.OrderBodyParam) (*models.Order, error) {
	order, err := r.GetOrderByID(id)
	if err != nil {
		return nil, err
	}

	copier.Copy(order, &item)
	if err := r.db.Save(&order).Error; err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepo) AssignOrder(id string) error {
	order, err := r.GetOrderByID(id)
	if err != nil {
		return err
	}

	for _, line := range order.Lines {
		quantity, err := r.quantityRepo.GetQuantityProductID(line.ProductID)
		if err != nil || quantity.Quantity < line.Quantity {
			return errors.New("product quantity is not enough")
		}
	}

	order.Status = models.OrderStatusAssigned
	if err := r.db.Save(&order).Error; err != nil {
		return err
	}

	return nil
}
