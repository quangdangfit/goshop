package model

type User struct {
	ID    string `json:"id" gorm:"unique;not null;index;primary_key"`
	Email string `json:"email" gorm:"unique;not null;index:idx_user_email"`
}
