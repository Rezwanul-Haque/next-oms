package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"next-oms/app/domain"
	"next-oms/app/serializers"
	"next-oms/app/svc"
	"next-oms/app/utils/methodsutil"
	"next-oms/app/utils/msgutil"
	"next-oms/infra/errors"
	"next-oms/infra/logger"

	"golang.org/x/crypto/bcrypt"
)

type users struct {
	lc   logger.LogClient
	uSvc svc.IUsers
}

// NewUsersController will initialize the controllers
func NewUsersController(grp interface{}, lc logger.LogClient, uSvc svc.IUsers) {
	uc := &users{
		lc:   lc,
		uSvc: uSvc,
	}

	g := grp.(*echo.Group)

	g.POST("/v1/users/signup", uc.Create)
	g.PATCH("/v1/user", uc.Update)
	g.POST("/v1/password/change", uc.ChangePassword)
	g.POST("/v1/password/forgot", uc.ForgotPassword)
	g.POST("/v1/password/verifyreset", uc.VerifyResetPassword)
	g.POST("/v1/password/reset", uc.ResetPassword)
}

// swagger:route POST /v1/users/signup User CreateUser
// Create a new user
// responses:
//	201: UserCreatedResponse
//	400: errorResponse
//	404: errorResponse
//	500: errorResponse

// Create handles POST requests and create a new sales user
func (ctr *users) Create(c echo.Context) error {
	var user domain.User

	if err := c.Bind(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(*user.Password), 8)
	*user.Password = string(hashedPass)

	result, saveErr := ctr.uSvc.CreateUser(user)
	if saveErr != nil {
		return c.JSON(saveErr.Status, saveErr)
	}
	var resp serializers.UserResp
	respErr := methodsutil.StructToStruct(result, &resp)
	if respErr != nil {
		return respErr
	}

	return c.JSON(http.StatusCreated, resp)
}

func (ctr *users) Update(c echo.Context) error {
	loggedInUser, err := GetUserFromContext(c)
	if err != nil {
		ctr.lc.Error(err.Error(), err)
		restErr := errors.NewUnauthorizedError("no logged-in user found")
		return c.JSON(restErr.Status, restErr)
	}

	var user serializers.UserReq
	if err := c.Bind(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	updateErr := ctr.uSvc.UpdateUser(uint(loggedInUser.ID), user)
	if updateErr != nil {
		return c.JSON(updateErr.Status, updateErr)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": msgutil.EntityUpdateSuccessMsg("user")})
}

func (ctr *users) ChangePassword(c echo.Context) error {
	loggedInUser, err := GetUserFromContext(c)
	if err != nil {
		ctr.lc.Error(err.Error(), err)
		restErr := errors.NewUnauthorizedError("no logged-in user found")
		return c.JSON(restErr.Status, restErr)
	}
	body := &serializers.ChangePasswordReq{}
	if err := c.Bind(&body); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}
	if err = body.Validate(); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		return c.JSON(restErr.Status, restErr)
	}
	if body.OldPassword == body.NewPassword {
		restErr := errors.NewBadRequestError("password can't be same as old one")
		return c.JSON(restErr.Status, restErr)
	}
	if err := ctr.uSvc.ChangePassword(loggedInUser.ID, body); err != nil {
		switch err {
		case errors.ErrInvalidPassword:
			restErr := errors.NewBadRequestError("old password didn't match")
			return c.JSON(restErr.Status, restErr)
		default:
			restErr := errors.NewInternalServerError(errors.ErrSomethingWentWrong)
			return c.JSON(restErr.Status, restErr)
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": msgutil.EntityChangedSuccessMsg("password")})
}

func (ctr *users) ForgotPassword(c echo.Context) error {
	body := &serializers.ForgotPasswordReq{}

	if err := c.Bind(&body); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	if err := body.Validate(); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		return c.JSON(restErr.Status, restErr)
	}

	if err := ctr.uSvc.ForgotPassword(body.Email); err != nil && err == errors.ErrSendingEmail {
		restErr := errors.NewInternalServerError("failed to send password reset email")
		return c.JSON(restErr.Status, restErr)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Password reset link sent to email"})
}

func (ctr *users) VerifyResetPassword(c echo.Context) error {
	req := &serializers.VerifyResetPasswordReq{}

	if err := c.Bind(&req); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	if err := req.Validate(); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		return c.JSON(restErr.Status, restErr)
	}

	if err := ctr.uSvc.VerifyResetPassword(req); err != nil {
		switch err {
		case errors.ErrParseJwt,
			errors.ErrInvalidPasswordResetToken:
			restErr := errors.NewUnauthorizedError("failed to send reset_token email")
			return c.JSON(restErr.Status, restErr)
		default:
			restErr := errors.NewInternalServerError(errors.ErrSomethingWentWrong)
			return c.JSON(restErr.Status, restErr)
		}
	}

	return c.JSON(http.StatusOK, "reset token verified")
}

func (ctr *users) ResetPassword(c echo.Context) error {
	req := &serializers.ResetPasswordReq{}

	if err := c.Bind(&req); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	if err := req.Validate(); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		return c.JSON(restErr.Status, restErr)
	}

	verifyReq := &serializers.VerifyResetPasswordReq{
		Token: req.Token,
		ID:    req.ID,
	}

	if err := ctr.uSvc.VerifyResetPassword(verifyReq); err != nil {
		switch err {
		case errors.ErrParseJwt,
			errors.ErrInvalidPasswordResetToken:
			restErr := errors.NewUnauthorizedError("failed to send reset_token email")
			return c.JSON(restErr.Status, restErr)
		default:
			restErr := errors.NewInternalServerError(errors.ErrSomethingWentWrong)
			return c.JSON(restErr.Status, restErr)
		}
	}

	if err := ctr.uSvc.ResetPassword(req); err != nil {
		restErr := errors.NewInternalServerError(errors.ErrSomethingWentWrong)
		return c.JSON(restErr.Status, restErr)
	}

	return c.JSON(http.StatusOK, "password reset successful")
}
