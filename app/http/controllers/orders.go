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
	g.PUT("/v1/orders/:con_id/cancel", oc.CancelOrder)
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

// CancelOrder handles PUT requests and cancel an order
func (ctr *orders) CancelOrder(c echo.Context) error {
	conID := c.Param("con_id")
	if conID == "" {
		restErr := errors.NewBadRequestError("order_consignment_id is required")
		return c.JSON(restErr.Status, restErr)
	}

	cancelErr := ctr.oSvc.CancelOrder(conID)
	if cancelErr != nil {
		return c.JSON(cancelErr.Status, cancelErr)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Order Cancelled Successfully",
		"type":    "success",
		"code":    200,
	})

}
