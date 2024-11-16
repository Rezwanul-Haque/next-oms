package impl

import (
	"context"
	"next-oms/app/domain"
	"next-oms/app/repository"
	"next-oms/infra/conn/db"
	"next-oms/infra/errors"
	"next-oms/infra/logger"
	"time"
)

type users struct {
	ctx context.Context
	lc  logger.LogClient
	DB  db.DatabaseClient
}

// NewUsersRepository will create an object that represent the User.Repository implementations
func NewUsersRepository(ctx context.Context, lc logger.LogClient, dbc db.DatabaseClient) repository.IUsers {
	return &users{
		ctx: ctx,
		lc:  lc,
		DB:  dbc,
	}
}

func (r *users) SaveUser(user *domain.User) (*domain.User, *errors.RestErr) {
	return r.DB.SaveUser(user)
}

func (r *users) GetUserByID(userID uint) (*domain.User, *errors.RestErr) {
	return r.DB.GetUserByID(userID)
}

func (r *users) UpdateUser(user *domain.User) *errors.RestErr {
	return r.DB.UpdateUser(user)
}

func (r *users) UpdatePassword(userID uint, companyID uint, updateValues map[string]interface{}) *errors.RestErr {
	return r.DB.UpdatePassword(userID, companyID, updateValues)
}

func (r *users) GetUserByEmail(email string) (*domain.User, error) {
	return r.DB.GetUserByEmail(email)
}

func (r *users) SetLastLoginAt(user *domain.User) error {
	utc := time.Now().UTC()
	user.LastLoginAt = &utc

	return r.DB.SetLastLoginAt(user)
}

func (r *users) ResetPassword(userID int, hashedPass []byte) error {
	return r.DB.ResetPassword(userID, hashedPass)
}

func (r *users) GetTokenUser(id uint) (*domain.VerifyTokenResp, *errors.RestErr) {
	return r.DB.GetTokenUser(id)
}
