package logger

import "next-oms/app/domain"

var client LogClient

func NewLogClient(lvl string) domain.ILogger {
	connectZap(lvl)

	return &LogClient{}
}

func Client() LogClient {
	return client
}
