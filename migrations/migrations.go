package migrations

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"goshop/dbs"
	"goshop/models"
	"goshop/repositories"
)

func createAdmin() {
	roleRepo := repositories.NewRoleRepository()
	role, _ := roleRepo.CreateRole(&models.RoleRequest{Name: "admin", Description: "Admin"})

	userRepo := repositories.NewUserRepository()
	userRepo.Register(&models.RegisterRequest{
		Username: "admin",
		Password: "admin",
		Email:    "admin@admin.com",
		RoleUUID: role.UUID,
	})
}

func Migrate() {
	Product := models.Product{}
	Pategory := models.Category{}
	Order := models.Order{}
	OrderLine := models.OrderLine{}
	User := models.User{}
	Role := models.Role{}
	Warehouse := models.Warehouse{}

	dbs.Database.AutoMigrate(&Product, &Pategory, &Order, &OrderLine, &User, &Role, &Warehouse)
	dbs.Database.Model(&Product).AddForeignKey("categ_uuid", "categories(uuid)", "RESTRICT", "RESTRICT")
	dbs.Database.Model(&User).AddForeignKey("role_uuid", "roles(uuid)", "RESTRICT", "RESTRICT")

	//createAdmin()
}
