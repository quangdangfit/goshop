package user

import (
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"goshop/dbs"
	"goshop/utils"
)

type Repository interface {
	//Login(username string, pass string) (*User, error)
	Register(req RegisterRequest) (*User, error)
	//GetUser(id string, jwt string) (*User, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepository() Repository {
	return &repo{db: dbs.Database}
}

//func (r *repo) Login(username string, pass string) (*User, error) {
//	// Validation before login
//	valid := validation.Validate(
//		[]validation.Validation{
//			{Value: username, Valid: "username"},
//			{Value: pass, Valid: "password"},
//		})
//
//	if valid {
//		user := &User{}
//		if database.DB.Where("username = ? ", username).First(&user).RecordNotFound() {
//			return nil, errors.New("user not found")
//		}
//
//		// Verify password
//		passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
//		if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
//			return nil, errors.New("wrong password")
//		}
//
//		return user, nil
//	} else {
//		return nil, errors.New("not valid values")
//	}
//}

func (r *repo) Register(req RegisterRequest) (*User, error) {
	var user User
	copier.Copy(&user, &req)
	hashedPassword := utils.HashAndSalt([]byte(req.Password))
	user.Password = hashedPassword

	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

//func (r *repo) GetUser(uid string, jwt string) (*User, error) {
//	userUUID, isValid := validation.ValidateToken(jwt)
//	if isValid {
//		user := &User{}
//		if database.DB.Where("uid = ? ", uid).First(&user).RecordNotFound() {
//			return nil, errors.New("user not found")
//		}
//		if user.UID == userUUID {
//			return user, nil
//		}
//	}
//	return nil, errors.New("not valid token")
//}
