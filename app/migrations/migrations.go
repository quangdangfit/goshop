package migrations

import (
	"goshop/app/models"
	"goshop/dbs"
)

func Migrate() {
	User := models.User{}
	Product := models.Product{}
	Order := models.Order{}
	OrderLine := models.OrderLine{}
	//Warehouse := models.Warehouse{}

	dbs.Database.AutoMigrate(&Product, &User, Order, OrderLine)
	//dbs.Database.Model(&User).AddForeignKey("role_id", "roles(id)", "RESTRICT", "RESTRICT")
	//dbs.Database.Model(&OrderLine).AddForeignKey("product_id", "products(id)", "RESTRICT", "RESTRICT")
	//dbs.Database.Model(&OrderLine).AddForeignKey("order_id", "orders(id)", "RESTRICT", "RESTRICT")
}
