package user

import (
	"errors"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"goshop/dbs"
	"goshop/utils"
)

type Repository interface {
	Login(req *LoginRequest) (*User, error)
	Register(req *RegisterRequest) (*User, error)
	GetUserByID(uuid string) (*User, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepository() Repository {
	return &repo{db: dbs.Database}
}

func (r *repo) Login(req *LoginRequest) (*User, error) {
	user := &User{}
	if dbs.Database.Where("username = ? ", req.Username).First(&user).RecordNotFound() {
		return nil, errors.New("user not found")
	}

	passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
		return nil, errors.New("wrong password")
	}

	return user, nil
}

func (r *repo) Register(req *RegisterRequest) (*User, error) {
	var user User
	copier.Copy(&user, &req)
	hashedPassword := utils.HashAndSalt([]byte(req.Password))
	user.Password = hashedPassword

	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repo) GetUserByID(uuid string) (*User, error) {
	user := User{}
	if dbs.Database.Where("uuid = ? ", uuid).First(&user).RecordNotFound() {
		return nil, errors.New("user not found")
	}

	return &user, nil
}
