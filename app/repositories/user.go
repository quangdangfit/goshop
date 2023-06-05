package repositories

import (
	"errors"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	"goshop/app/models"
	"goshop/app/schema"
	"goshop/dbs"
	"goshop/pkg/utils"
)

type IUserRepository interface {
	Login(item *schema.Login) (*models.User, error)
	Register(item *schema.Register) (*models.User, error)
	GetUserByID(uuid string) (*models.User, error)
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepo {
	return &UserRepo{db: dbs.Database}
}

func (u *UserRepo) Login(item *schema.Login) (*models.User, error) {
	user := &models.User{}
	if dbs.Database.Where("username = ? ", item.Username).First(&user).RecordNotFound() {
		return nil, errors.New("user not found")
	}

	passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(item.Password))
	if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
		return nil, errors.New("wrong password")
	}

	return user, nil
}

func (u *UserRepo) Register(item *schema.Register) (*models.User, error) {
	var user models.User
	copier.Copy(&user, &item)
	hashedPassword := utils.HashAndSalt([]byte(item.Password))
	user.Password = hashedPassword

	if err := u.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) GetUserByID(uuid string) (*models.User, error) {
	user := models.User{}
	if dbs.Database.Where("uuid = ? ", uuid).First(&user).RecordNotFound() {
		return nil, errors.New("user not found")
	}

	return &user, nil
}
