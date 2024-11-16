package impl

import (
	"context"
	"fmt"
	"next-oms/app/repository"
	"next-oms/infra/conn/cache"
	"next-oms/infra/conn/db"
	"next-oms/infra/logger"
)

type system struct {
	ctx   context.Context
	lc    logger.LogClient
	DB    db.DatabaseClient
	Cache cache.CacheClient
}

// NewSystemRepository will create an object that represent the System.Repository implementations
func NewSystemRepository(ctx context.Context, lc logger.LogClient, dbc db.DatabaseClient, c cache.CacheClient) repository.ISystem {
	return &system{
		ctx:   ctx,
		lc:    lc,
		DB:    dbc,
		Cache: c,
	}
}

func (r *system) DBCheck() (bool, error) {
	dB, _ := r.DB.DB.DB()
	if err := dB.Ping(); err != nil {
		return false, err
	}

	return true, nil
}

func (r *system) CacheCheck() bool {
	pong, err := r.Cache.Redis.Ping(r.ctx).Result()
	if err != nil {
		return false
	}

	r.lc.Info(fmt.Sprintf("%v from cache", pong))

	return true
}
