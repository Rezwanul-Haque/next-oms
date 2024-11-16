package controllers

import (
	"github.com/labstack/echo/v4"
	"next-oms/app/serializers"
	"next-oms/infra/errors"
)

func GetUserFromContext(c echo.Context) (*serializers.LoggedInUser, error) {
	user, ok := c.Get("user").(*serializers.LoggedInUser)
	if !ok {
		return nil, errors.ErrNoContextUser
	}

	return user, nil
}
