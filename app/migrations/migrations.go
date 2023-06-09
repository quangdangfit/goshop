package migrations

import (
	"goshop/app/dbs"
	"goshop/app/models"
)

func Migrate() {
	User := models.User{}
	Product := models.Product{}
	Order := models.Order{}
	OrderLine := models.OrderLine{}

	dbs.Database.AutoMigrate(&Product, &User, Order, OrderLine)
}
