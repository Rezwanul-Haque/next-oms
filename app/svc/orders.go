package svc

import (
	"next-oms/app/serializers"
	"next-oms/infra/errors"
)

type IOrders interface {
	CreateOrder(order *serializers.OrderReq) (*serializers.OrderResp, *errors.RestErr)
	GetOrders(filters *serializers.ListFilters) (*serializers.ListFilters, *errors.RestErr)
}
