package migrations

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"goshop/dbs"
	"goshop/objects/category"
	"goshop/objects/order"
	"goshop/objects/orderLine"
	"goshop/objects/product"
	"goshop/objects/role"
	"goshop/objects/user"
)

func createAdmin() {
	roleRepo := role.NewRepository()
	role, _ := roleRepo.CreateRole(&role.RoleRequest{Name: "admin", Description: "Admin"})

	userRepo := user.NewRepository()
	userRepo.Register(&user.RegisterRequest{
		Username: "admin",
		Password: "admin",
		Email:    "admin@admin.com",
		RoleUUID: role.UUID,
	})
}

func Migrate() {
	Product := product.Product{}
	Pategory := category.Category{}
	Order := order.Order{}
	OrderLine := orderLine.OrderLine{}
	User := user.User{}
	Role := role.Role{}

	dbs.Database.AutoMigrate(&Product, &Pategory, &Order, &OrderLine, &User, &Role)
	dbs.Database.Model(&Product).AddForeignKey("categ_uuid", "categories(uuid)", "RESTRICT", "RESTRICT")
	dbs.Database.Model(&User).AddForeignKey("role_uuid", "roles(uuid)", "RESTRICT", "RESTRICT")

	//createAdmin()
}
