package repositories

import (
	"errors"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"

	"goshop/app/models"
	"goshop/dbs"
)

type OrderLineRepository interface {
	GetOrderLines(map[string]interface{}) (*[]models.OrderLine, error)
	GetOrderLineByID(uuid string) (*models.OrderLine, error)
	CreateOrderLine(req *models.OrderLineBodyRequest) (*models.OrderLine, error)
	UpdateOrderLine(uuid string, req *models.OrderLineBodyRequest) (*models.OrderLine, error)
}

type orderLineRepo struct {
	db *gorm.DB
}

func NewOrderLineRepository() OrderLineRepository {
	return &orderLineRepo{db: dbs.Database}
}

func (line *orderLineRepo) GetOrderLines(query map[string]interface{}) (*[]models.OrderLine, error) {
	var orderLines []models.OrderLine
	if line.db.Find(&orderLines, query).RecordNotFound() {
		return nil, nil
	}

	return &orderLines, nil
}

func (line *orderLineRepo) GetOrderLineByID(uuid string) (*models.OrderLine, error) {
	var orderLine models.OrderLine
	if line.db.Where("uuid = ?", uuid).First(&orderLine).RecordNotFound() {
		return nil, errors.New("not found orderLine")
	}

	return &orderLine, nil
}

func (line *orderLineRepo) CreateOrderLine(req *models.OrderLineBodyRequest) (*models.OrderLine, error) {
	var orderLine models.OrderLine
	copier.Copy(&orderLine, &req)

	if err := line.db.Create(&orderLine).Error; err != nil {
		return nil, err
	}

	return &orderLine, nil
}

func (line *orderLineRepo) UpdateOrderLine(uuid string, req *models.OrderLineBodyRequest) (*models.OrderLine, error) {
	orderLine, err := line.GetOrderLineByID(uuid)
	if err != nil {
		return nil, err
	}

	copier.Copy(orderLine, &req)
	if err := line.db.Save(&orderLine).Error; err != nil {
		return nil, err
	}

	return orderLine, nil
}
