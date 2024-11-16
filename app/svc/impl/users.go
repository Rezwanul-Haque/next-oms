package impl

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"next-oms/app/domain"
	"next-oms/app/repository"
	"next-oms/app/serializers"
	"next-oms/app/svc"
	"next-oms/app/utils/methodsutil"
	"next-oms/app/utils/msgutil"
	"next-oms/infra/config"
	"next-oms/infra/conn/cache"
	"next-oms/infra/errors"
	"next-oms/infra/logger"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type users struct {
	ctx   context.Context
	lc    logger.LogClient
	urepo repository.IUsers
}

func NewUsersService(ctx context.Context, lc logger.LogClient, urepo repository.IUsers) svc.IUsers {
	return &users{
		ctx:   ctx,
		lc:    lc,
		urepo: urepo,
	}
}

func (u *users) CreateAdminUser(user domain.User) (*domain.User, *errors.RestErr) {
	resp, saveErr := u.urepo.SaveUser(&user)
	if saveErr != nil {
		return nil, saveErr
	}
	return resp, nil
}

func (u *users) CreateUser(user domain.User) (*domain.User, *errors.RestErr) {
	resp, saveErr := u.urepo.SaveUser(&user)
	if saveErr != nil {
		return nil, saveErr
	}
	return resp, nil
}

func (u *users) GetUserById(userId uint) (*domain.User, *errors.RestErr) {
	resp, getErr := u.urepo.GetUserByID(userId)
	if getErr != nil {
		return nil, getErr
	}
	return resp, nil
}

func (u *users) GetUserByEmail(userName string) (*domain.User, error) {
	resp, getErr := u.urepo.GetUserByEmail(userName)
	if getErr != nil {
		return nil, getErr
	}
	return resp, nil
}

func (u *users) UpdateUser(userID uint, req serializers.UserReq) *errors.RestErr {
	var user domain.User

	err := methodsutil.StructToStruct(req, &user)
	if err != nil {
		u.lc.Error(msgutil.EntityStructToStructFailedMsg("update user"), err)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	user.ID = userID

	if updateErr := u.urepo.UpdateUser(&user); updateErr != nil {
		return updateErr
	}

	if err := u.deleteUserCache(int(userID)); err != nil {
		restErr := errors.NewInternalServerError(errors.ErrSomethingWentWrong)
		return restErr
	}
	return nil
}

func (u *users) ChangePassword(id int, data *serializers.ChangePasswordReq) error {
	user, getErr := u.urepo.GetUserByID(uint(id))
	if getErr != nil {
		return errors.NewError(getErr.Message)
	}

	currentPass := []byte(*user.Password)
	if err := bcrypt.CompareHashAndPassword(currentPass, []byte(data.OldPassword)); err != nil {
		u.lc.Error(msgutil.EntityGenericFailedMsg("comparing hash and old password"), err)
		return errors.ErrInvalidPassword
	}

	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(data.NewPassword), 8)

	updates := map[string]interface{}{
		"password":    hashedPass,
		"first_login": false,
	}

	upErr := u.urepo.UpdatePassword(user.ID, 1, updates)
	if upErr != nil {
		return errors.NewError(upErr.Message)
	}

	return nil
}

func (u *users) ForgotPassword(email string) error {
	user, err := u.urepo.GetUserByEmail(email)
	if err != nil {
		return err
	}

	secret := passwordResetSecret(user)

	payload := jwt.MapClaims{}
	payload["email"] = user.Email

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := token.SignedString([]byte(secret))

	if err != nil {
		u.lc.Error("error occur when getting complete signed token", err)
		return err
	}

	// TODO: Send Mail
	u.lc.Info(signedToken)
	// fpassReq := &serializers.ForgetPasswordMailReq{
	// 	To:     user.Email,
	// 	UserID: user.ID,
	// 	Token:  signedToken,
	// }

	// if err := u.msvc.SendForgotPasswordEmail(*fpassReq); err != nil {
	// 	return errors.ErrSendingEmail
	// }

	return nil
}

func (u *users) VerifyResetPassword(req *serializers.VerifyResetPasswordReq) error {
	user, getErr := u.urepo.GetUserByID(uint(req.ID))
	if getErr != nil {
		return errors.NewError(getErr.Message)
	}

	secret := passwordResetSecret(user)

	parsedToken, err := methodsutil.ParseJwtToken(req.Token, secret)
	if err != nil {
		u.lc.Error("error occur when parse jwt token with secret", err)
		return errors.ErrParseJwt
	}

	if _, ok := parsedToken.Claims.(jwt.Claims); !ok && !parsedToken.Valid {
		return errors.ErrInvalidPasswordResetToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return errors.ErrInvalidPasswordResetToken
	}

	parsedEmail := claims["email"].(string)
	if user.Email != parsedEmail {
		return errors.ErrInvalidPasswordResetToken
	}

	return nil
}

func (u *users) ResetPassword(req *serializers.ResetPasswordReq) error {
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 8)

	if err := u.urepo.ResetPassword(req.ID, hashedPass); err != nil {
		return err
	}

	return nil
}

func (u *users) deleteUserCache(userID int) error {
	if err := cache.Client().Del(
		u.ctx,
		config.Cache().Redis.UserPrefix+strconv.Itoa(userID),
		config.Cache().Redis.TokenPrefix+strconv.Itoa(userID),
	); err != nil {
		u.lc.Error("error occur when deleting cached user after user update", err)
		return err
	}

	return nil
}

func passwordResetSecret(user *domain.User) string {
	return *user.Password + strconv.Itoa(int(user.CreatedAt.Unix()))
}
