package repositories

import (
	"errors"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	"goshop/app/models"
	"goshop/dbs"
	"goshop/pkg/utils"
)

type UserRepository interface {
	Login(req *models.LoginRequest) (*models.User, error)
	Register(req *models.RegisterRequest) (*models.User, error)
	GetUserByID(uuid string) (*models.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepo{db: dbs.Database}
}

func (u *userRepo) Login(req *models.LoginRequest) (*models.User, error) {
	user := &models.User{}
	if dbs.Database.Where("username = ? ", req.Username).First(&user).RecordNotFound() {
		return nil, errors.New("user not found")
	}

	passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
		return nil, errors.New("wrong password")
	}

	return user, nil
}

func (u *userRepo) Register(req *models.RegisterRequest) (*models.User, error) {
	var user models.User
	copier.Copy(&user, &req)
	hashedPassword := utils.HashAndSalt([]byte(req.Password))
	user.Password = hashedPassword

	if err := u.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepo) GetUserByID(uuid string) (*models.User, error) {
	user := models.User{}
	if dbs.Database.Where("uuid = ? ", uuid).First(&user).RecordNotFound() {
		return nil, errors.New("user not found")
	}

	return &user, nil
}
