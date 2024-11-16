package db

import (
	"next-oms/app/domain"
	"next-oms/infra/logger"
)

var client DatabaseClient

func NewDbClient(lc logger.LogClient) domain.IDb {
	connectMySQL(lc)

	return &DatabaseClient{}
}

func Client() DatabaseClient {
	return client
}
