package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"next-oms/app/serializers"
	"next-oms/app/svc"
	"next-oms/infra/errors"
	"next-oms/infra/logger"
)

type orders struct {
	lc   logger.LogClient
	oSvc svc.IOrders
}

// NewOrdersController will initialize the controllers
func NewOrdersController(grp interface{}, lc logger.LogClient, oSvc svc.IOrders) {
	oc := &orders{
		lc:   lc,
		oSvc: oSvc,
	}

	g := grp.(*echo.Group)

	g.POST("/v1/orders", oc.Create)
	g.GET("/v1/orders/all", oc.GetOrders)
}

// swagger:route POST /v1/orders OrderReq CreateOrder
// Create a new order
// responses:
//	201: OrderCreatedResponse
//	400: errorResponse
//	404: errorResponse
//	500: errorResponse

// Create handles POST requests and create a new order
func (ctr *orders) Create(c echo.Context) error {
	var order serializers.OrderReq

	if err := c.Bind(&order); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	if err := order.Validate(); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		return c.JSON(restErr.Status, restErr)
	}

	resp, saveErr := ctr.oSvc.CreateOrder(&order)
	if saveErr != nil {
		return c.JSON(saveErr.Status, saveErr)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Order Created Successfully",
		"type":    "success",
		"code":    200,
		"data":    resp,
	})
}

// GetOrders handles GET requests and all the orders
func (ctr *orders) GetOrders(c echo.Context) error {
	listParams := &serializers.ListFilters{}
	listParams.GenerateFilters(c.QueryParams())
	listParams.BasePath = c.Request().URL.Path

	result, saveErr := ctr.oSvc.GetOrders(listParams)
	if saveErr != nil {
		return c.JSON(saveErr.Status, saveErr)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Orders successfully fetched.",
		"type":    "success",
		"code":    200,
		"data":    result,
	})
}
