package impl

import (
	"context"
	"next-oms/app/domain"
	"next-oms/app/repository"
	"next-oms/app/serializers"
	"next-oms/infra/conn/db"
	"next-oms/infra/errors"
	"next-oms/infra/logger"
)

type orders struct {
	ctx context.Context
	lc  logger.LogClient
	DB  db.DatabaseClient
}

// NewOrdersRepository will create an object that represent the Orders.Repository implementations
func NewOrdersRepository(ctx context.Context, lc logger.LogClient, dbc db.DatabaseClient) repository.IOrders {
	return &orders{
		ctx: ctx,
		lc:  lc,
		DB:  dbc,
	}
}

func (r *orders) SaveOrder(user *domain.Order) (*domain.Order, *errors.RestErr) {
	return r.DB.SaveOrder(user)
}

func (r *orders) GetOrders(filters *serializers.ListFilters) (domain.Orders, *errors.RestErr) {
	return r.DB.GetOrders(filters)
}
