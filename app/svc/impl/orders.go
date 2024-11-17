package impl

import (
	"context"
	"fmt"
	"math/rand/v2"
	"next-oms/app/domain"
	"next-oms/app/repository"
	"next-oms/app/serializers"
	"next-oms/app/svc"
	"next-oms/app/utils/consts"
	"next-oms/infra/errors"
	"next-oms/infra/logger"
)

type orders struct {
	ctx   context.Context
	lc    logger.LogClient
	orepo repository.IOrders
}

func NewOrdersService(ctx context.Context, lc logger.LogClient, orepo repository.IOrders) svc.IOrders {
	return &orders{
		ctx:   ctx,
		lc:    lc,
		orepo: orepo,
	}
}

func (o *orders) CreateOrder(order *serializers.OrderReq) (*serializers.OrderResp, *errors.RestErr) {
	orderType := 1
	ord := domain.Order{
		ConsignmentID:    fmt.Sprintf("CONS-%d-%s-%d", order.StoreID, order.RecipientName, rand.Int()),
		Description:      order.ItemDescription,
		MerchantOrderID:  order.MerchantOrderID,
		RecipientName:    order.RecipientName,
		RecipientAddress: order.RecipientAddress,
		RecipientPhone:   order.RecipientPhone,
		Amount:           order.AmountToCollect,
		TotalFee:         consts.CalculateTotalFee(*order),
		Instruction:      order.SpecialInstruction,
		OrderTypeID:      orderType,
		CodFee:           0,
		PromoDiscount:    0,
		Discount:         0,
		DeliveryFee:      consts.CalculateDeliveryFee(orderType, order.AmountToCollect),
		Status:           consts.OrderPending,
		OrderType:        consts.GetOrderTypeDescription(orderType),
		ItemType:         consts.GetItemTypeDescription(order.ItemType),
	}

	result, saveErr := o.orepo.SaveOrder(&ord)
	if saveErr != nil {
		return nil, saveErr
	}

	//TODO: Create Shipment info : for now adding it to order

	//TODO: Store order History

	resp := &serializers.OrderResp{
		ConsignmentID:   result.ConsignmentID,
		MerchantOrderID: result.MerchantOrderID,
		OrderStatus:     result.Status,
		DeliveryFee:     result.DeliveryFee,
	}

	return resp, nil
}

func (o *orders) GetOrders(filters *serializers.ListFilters) (*serializers.ListFilters, *errors.RestErr) {
	orders, err := o.orepo.GetOrders(filters)
	if err != nil {
		return nil, err
	}

	filters.Results = orders
	return filters, nil
}
