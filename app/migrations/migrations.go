package migrations

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"goshop/app/models"
	"goshop/app/repositories"
	"goshop/app/schema"
	"goshop/dbs"
)

func createAdmin() {
	roleRepo := repositories.NewRoleRepository()
	role, _ := roleRepo.CreateRole(&schema.RoleBodyParam{Name: "admin", Description: "Admin"})

	userRepo := repositories.NewUserRepository()
	userRepo.Register(&schema.Register{
		Username: "admin",
		Password: "admin",
		Email:    "admin@admin.com",
		RoleUUID: role.ID,
	})
}

func Migrate() {
	Product := models.Product{}
	Category := models.Category{}
	Order := models.Order{}
	OrderLine := models.OrderLine{}
	User := models.User{}
	Role := models.Role{}
	Warehouse := models.Warehouse{}

	dbs.Database.AutoMigrate(&Product, &Category, &Order, &OrderLine, &User, &Role, &Warehouse)
	dbs.Database.Model(&Product).AddForeignKey("categ_id", "categories(id)", "RESTRICT", "RESTRICT")
	dbs.Database.Model(&User).AddForeignKey("role_id", "roles(id)", "RESTRICT", "RESTRICT")
	dbs.Database.Model(&OrderLine).AddForeignKey("product_id", "products(id)", "RESTRICT", "RESTRICT")
	dbs.Database.Model(&OrderLine).AddForeignKey("order_id", "orders(id)", "RESTRICT", "RESTRICT")

	createAdmin()
}
