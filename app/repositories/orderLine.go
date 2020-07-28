package repositories

import (
	"errors"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"

	"goshop/app/models"
	"goshop/app/schema"
	"goshop/dbs"
)

type OrderLineRepository interface {
	GetOrderLines(query *schema.OrderLineQueryParam) (*[]models.OrderLine, error)
	GetOrderLineByID(uuid string) (*models.OrderLine, error)
	CreateOrderLine(item *schema.OrderLineBodyParam) (*models.OrderLine, error)
	UpdateOrderLine(uuid string, item *schema.OrderLineBodyParam) (*models.OrderLine, error)
}

type lineRepo struct {
	db *gorm.DB
}

func NewOrderLineRepository() OrderLineRepository {
	return &lineRepo{db: dbs.Database}
}

func (line *lineRepo) GetOrderLines(query *schema.OrderLineQueryParam) (*[]models.OrderLine, error) {
	var orderLines []models.OrderLine
	if line.db.Find(&orderLines, query).RecordNotFound() {
		return nil, nil
	}

	return &orderLines, nil
}

func (line *lineRepo) GetOrderLineByID(uuid string) (*models.OrderLine, error) {
	var orderLine models.OrderLine
	if line.db.Where("uuid = ?", uuid).First(&orderLine).RecordNotFound() {
		return nil, errors.New("not found orderLine")
	}

	return &orderLine, nil
}

func (line *lineRepo) CreateOrderLine(item *schema.OrderLineBodyParam) (*models.OrderLine, error) {
	var orderLine models.OrderLine
	copier.Copy(&orderLine, &item)

	if err := line.db.Create(&orderLine).Error; err != nil {
		return nil, err
	}

	return &orderLine, nil
}

func (line *lineRepo) UpdateOrderLine(uuid string, item *schema.OrderLineBodyParam) (*models.OrderLine, error) {
	orderLine, err := line.GetOrderLineByID(uuid)
	if err != nil {
		return nil, err
	}

	copier.Copy(orderLine, &item)
	if err := line.db.Save(&orderLine).Error; err != nil {
		return nil, err
	}

	return orderLine, nil
}
