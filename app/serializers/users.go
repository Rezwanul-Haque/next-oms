package serializers

import (
	"time"
)

type User struct {
	ID          uint       `json:"id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email"`
	Password    *string    `json:"password,omitempty"`
	Phone       string     `json:"phone"`
	ProfilePic  *string    `json:"profile_pic"`
	LastLoginAt *time.Time `json:"last_login_at"`
	FirstLogin  bool       `json:"first_login" gorm:"column:first_login;default:true"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type ResolveUserResponse struct {
	CompanyName  string             `json:"company_name"`
	CompanyID    uint               `json:"company_id"`
	Subordinates []*ResolveUserResp `json:"subordinates"`
}

type UserReq struct {
	UserName   string  `json:"user_name,omitempty"`
	FirstName  string  `json:"first_name,omitempty"`
	LastName   string  `json:"last_name,omitempty"`
	Email      string  `json:"email,omitempty"`
	Password   *string `json:"password,omitempty"`
	ProfilePic *string `json:"profile_pic,omitempty"`
	Phone      string  `json:"phone,omitempty"`
	CompanyID  uint    `json:"company_id"`
}

type LoggedInUser struct {
	ID          int    `json:"user_id"`
	AccessUuid  string `json:"access_uuid"`
	RefreshUuid string `json:"refresh_uuid"`
}

type UserResp struct {
	ID          int        `json:"id"`
	UserName    string     `json:"user_name"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email"`
	Phone       *string    `json:"phone"`
	ProfilePic  *string    `json:"profile_pic"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	LastLoginAt *time.Time `json:"last_login_at"`
	FirstLogin  bool       `json:"first_login"`
}

type ResolveUserResp struct {
	ID          int        `json:"id"`
	UserName    string     `json:"user_name"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email"`
	Phone       *string    `json:"phone"`
	ProfilePic  *string    `json:"profile_pic"`
	AppKey      string     `json:"app_key,omitempty"`
	RoleID      uint       `json:"role_id"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	LastLoginAt *time.Time `json:"last_login_at"`
	FirstLogin  bool       `json:"first_login"`
}

type UserWithParamsResp struct {
	UserResp
}

type VerifyTokenResp struct {
	ID           int      `json:"id"`
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	Email        string   `json:"email"`
	Phone        *string  `json:"phone"`
	ProfilePic   *string  `json:"profile_pic"`
	BusinessID   *int     `json:"business_id"`
	BusinessName string   `json:"business_name"`
	CompanyID    *int     `json:"company_id"`
	CompanyName  string   `json:"company_name"`
	Permissions  []string `json:"permissions"`
	Admin        bool     `json:"admin"`
}

type UserWithLocations struct {
	UserID   uint   `json:"user_id"`
	UserName string `json:"user_name"`
}
