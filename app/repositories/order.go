package repositories

import (
	"errors"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"

	"goshop/app/models"
	"goshop/app/schema"
	"goshop/dbs"
)

type IOrderRepository interface {
	GetOrders(query *schema.OrderQueryParam) (*[]models.Order, error)
	GetOrderByID(uuid string) (*models.Order, error)
	CreateOrder(item *schema.OrderBodyParam) (*models.Order, error)
	UpdateOrder(uuid string, item *schema.OrderBodyParam) (*models.Order, error)
	AssignOrder(uuid string) error
}

type OrderRepo struct {
	db           *gorm.DB
	lineRepo     IOrderLineRepository
	quantityRepo IQuantityRepository
}

func NewOrderRepository() *OrderRepo {
	return &OrderRepo{db: dbs.Database, lineRepo: NewOrderLineRepository()}
}

func (r *OrderRepo) GetOrders(query *schema.OrderQueryParam) (*[]models.Order, error) {
	var orders []models.Order
	if r.db.Find(&orders, query).RecordNotFound() {
		return nil, nil
	}

	return &orders, nil
}

func (r *OrderRepo) GetOrderByID(uuid string) (*models.Order, error) {
	var order models.Order
	var lines []models.OrderLine
	if r.db.Where("uuid = ?", uuid).First(&order).RecordNotFound() {
		return nil, errors.New("not found order")
	}
	r.db.Where("order_uuid = ?", uuid).Find(&lines)
	order.Lines = lines

	return &order, nil
}

func (r *OrderRepo) CreateOrder(item *schema.OrderBodyParam) (*models.Order, error) {
	var order models.Order
	copier.Copy(&order, &item)

	if order.Lines == nil {
		return nil, errors.New("order lines must be not empty")
	}

	if err := r.db.Create(&order).Error; err != nil {
		return nil, err
	}

	var lines []models.OrderLine
	var totalPrice uint
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

func (r *OrderRepo) UpdateOrder(uuid string, item *schema.OrderBodyParam) (*models.Order, error) {
	order, err := r.GetOrderByID(uuid)
	if err != nil {
		return nil, err
	}

	copier.Copy(order, &item)
	if err := r.db.Save(&order).Error; err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepo) AssignOrder(uuid string) error {
	order, err := r.GetOrderByID(uuid)
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
