package migrations

import (
	"goshop/app/models"
	"goshop/dbs"
)

func Migrate() {
	User := models.User{}
	Product := models.Product{}
	//Category := models.Category{}
	//Order := models.Order{}
	//OrderLine := models.OrderLine{}
	//Role := models.Role{}
	//Warehouse := models.Warehouse{}

	dbs.Database.AutoMigrate(&Product, &User)
	//dbs.Database.Model(&User).AddForeignKey("role_id", "roles(id)", "RESTRICT", "RESTRICT")
	//dbs.Database.Model(&OrderLine).AddForeignKey("product_id", "products(id)", "RESTRICT", "RESTRICT")
	//dbs.Database.Model(&OrderLine).AddForeignKey("order_id", "orders(id)", "RESTRICT", "RESTRICT")
}
