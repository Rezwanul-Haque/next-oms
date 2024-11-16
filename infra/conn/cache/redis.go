package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"next-oms/infra/config"
	"next-oms/infra/logger"
)

type CacheClient struct {
	Redis *redis.Client
}

func connectRedis(lc logger.LogClient) {
	conf := config.Cache().Redis

	lc.Info("connecting to Redis at " + conf.Host + ":" + conf.Port + "...")

	c := redis.NewClient(&redis.Options{
		Addr:     conf.Host + ":" + conf.Port,
		Password: conf.Pass,
		DB:       conf.Db,
	})

	client.Redis = c

	if _, err := client.Redis.Ping(context.Background()).Result(); err != nil {
		lc.Error("failed to connect Redis: ", err)
		panic(err)
	}

	lc.Info("Redis connection successful...")
}
