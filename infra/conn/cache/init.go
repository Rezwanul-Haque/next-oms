package cache

import (
	"next-oms/app/domain"
	"next-oms/infra/logger"
)

var client CacheClient

func NewCacheClient(lc logger.LogClient) domain.ICache {
	connectRedis(lc)

	return &CacheClient{}
}

func Client() CacheClient {
	return client
}
