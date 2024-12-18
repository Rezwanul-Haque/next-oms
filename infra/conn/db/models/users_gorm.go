package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	UserName    string     `json:"user_name"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email"`
	Password    *string    `json:"password,omitempty"`
	Phone       string     `json:"phone"`
	ProfilePic  *string    `json:"profile_pic"`
	LastLoginAt *time.Time `json:"last_login_at"`
	FirstLogin  bool       `json:"first_login" gorm:"column:first_login;default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
