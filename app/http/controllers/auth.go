package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"next-oms/app/serializers"
	"next-oms/app/svc"
	"next-oms/infra/errors"
	"next-oms/infra/logger"
)

type auth struct {
	lc      logger.LogClient
	authSvc svc.IAuth
	userSvc svc.IUsers
}

// NewAuthController will initialize the controllers
func NewAuthController(grp interface{}, lc logger.LogClient, authSvc svc.IAuth, userSvc svc.IUsers) {
	ac := &auth{
		lc:      lc,
		authSvc: authSvc,
		userSvc: userSvc,
	}

	g := grp.(*echo.Group)

	g.POST("/v1/login", ac.Login)
	g.POST("/v1/logout", ac.Logout)
	g.POST("/v1/token/refresh", ac.RefreshToken)
	g.GET("/v1/token/verify", ac.VerifyToken)
}

// swagger:route POST /v1/login Auth Login
// logged in a user
// responses:
//	201: LoginResp
//	400: errorResponse
//	404: errorResponse
//	500: errorResponse

// Login handles POST requests and logged in a user
func (ctr *auth) Login(c echo.Context) error {
	var cred *serializers.LoginReq
	var resp *serializers.LoginResp
	var err error

	if err = c.Bind(&cred); err != nil {
		bodyErr := errors.NewBadRequestError("failed to parse request body")
		ctr.lc.Error("failed to parse request body", err)
		return c.JSON(bodyErr.Status, bodyErr)
	}

	if resp, err = ctr.authSvc.Login(cred); err != nil {
		switch err {
		case errors.ErrInvalidEmail, errors.ErrInvalidPassword, errors.ErrNotAdmin:
			unAuthErr := errors.NewUnauthorizedError("The user credentials were incorrect.")
			return c.JSON(unAuthErr.Status, unAuthErr)
		case errors.ErrCreateJwt:
			serverErr := errors.NewInternalServerError("failed to create jwt token")
			return c.JSON(serverErr.Status, serverErr)
		case errors.ErrStoreTokenUuid:
			serverErr := errors.NewInternalServerError("failed to store jwt token uuid")
			return c.JSON(serverErr.Status, serverErr)
		default:
			serverErr := errors.NewInternalServerError(errors.ErrSomethingWentWrong)
			return c.JSON(serverErr.Status, serverErr)
		}
	}

	return c.JSON(http.StatusOK, resp)
}

// swagger:route POST /v1/logout Auth Logout
// logged in a user
// responses:
//	200: LoginResp
//	400: errorResponse
//	404: errorResponse
//	500: errorResponse

// Logout handles POST requests and logged out a user
func (ctr *auth) Logout(c echo.Context) error {
	var user *serializers.LoggedInUser
	var err error

	if user, err = GetUserFromContext(c); err != nil {
		ctr.lc.Error(err.Error(), err)
		serverErr := errors.NewInternalServerError("no logged-in user found")
		return c.JSON(serverErr.Status, serverErr)
	}

	if err := ctr.authSvc.Logout(user); err != nil {
		ctr.lc.Error(err.Error(), err)
		serverErr := errors.NewInternalServerError("failed to logout")
		return c.JSON(serverErr.Status, serverErr)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Successfully logged out"})
}

func (ctr *auth) RefreshToken(c echo.Context) error {
	var token *serializers.TokenRefreshReq
	var res *serializers.LoginResp
	var err error

	if err = c.Bind(&token); err != nil {
		ctr.lc.Error("failed to parse request body", err)
		bodyErr := errors.NewBadRequestError("failed to parse request body")
		return c.JSON(bodyErr.Status, bodyErr)
	}

	if res, err = ctr.authSvc.RefreshToken(token.RefreshToken); err != nil {
		switch err {
		case errors.ErrParseJwt,
			errors.ErrInvalidRefreshToken,
			errors.ErrInvalidRefreshUuid:
			unAuthErr := errors.NewUnauthorizedError("invalid refresh_token")
			return c.JSON(unAuthErr.Status, unAuthErr)
		case errors.ErrCreateJwt:
			serverErr := errors.NewInternalServerError("failed to create new jwt token")
			return c.JSON(serverErr.Status, serverErr)
		default:
			serverErr := errors.NewInternalServerError(errors.ErrSomethingWentWrong)
			return c.JSON(serverErr.Status, serverErr)
		}
	}

	return c.JSON(http.StatusOK, res)
}

func (ctr *auth) VerifyToken(c echo.Context) error {
	accessToken, err := AccessTokenFromHeader(c)

	if err != nil {
		unAuthErr := errors.NewUnauthorizedError("invalid access_token")
		return c.JSON(unAuthErr.Status, unAuthErr)
	}

	res, err := ctr.authSvc.VerifyToken(accessToken)
	if err != nil {
		switch err {
		case errors.ErrParseJwt,
			errors.ErrInvalidAccessToken,
			errors.ErrInvalidAccessUuid:
			unAuthErr := errors.NewUnauthorizedError("invalid access_token")
			return c.JSON(unAuthErr.Status, unAuthErr)
		default:
			serverErr := errors.NewInternalServerError(errors.ErrSomethingWentWrong)
			return c.JSON(serverErr.Status, serverErr)
		}
	}

	return c.JSON(http.StatusOK, res)
}

func AccessTokenFromHeader(c echo.Context) (string, error) {
	header := "Authorization"
	authScheme := "Bearer"

	auth := c.Request().Header.Get(header)
	l := len(authScheme)

	if len(auth) > l+1 && auth[:l] == authScheme {
		return auth[l+1:], nil
	}

	return "", errors.ErrInvalidAccessToken
}
