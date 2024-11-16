package container

import (
	"context"
	"next-oms/app/http/controllers"
	repoImpl "next-oms/app/repository/impl"
	svcImpl "next-oms/app/svc/impl"
	"next-oms/infra/conn/cache"
	"next-oms/infra/conn/db"
	"next-oms/infra/logger"
)

func Init(g interface{}, lc logger.LogClient) {
	basectx := context.Background()
	dbc := db.Client()
	cachec := cache.Client()

	// register all repos impl, services impl, controllers
	sysRepo := repoImpl.NewSystemRepository(basectx, lc, dbc, cachec)
	userRepo := repoImpl.NewUsersRepository(basectx, lc, dbc)

	sysSvc := svcImpl.NewSystemService(sysRepo)
	userSvc := svcImpl.NewUsersService(basectx, lc, userRepo)
	tokenSvc := svcImpl.NewTokenService(basectx, lc, userRepo)
	authSvc := svcImpl.NewAuthService(basectx, lc, userRepo, tokenSvc)

	controllers.NewSystemController(g, lc, sysSvc)
	controllers.NewAuthController(g, lc, authSvc, userSvc)
	controllers.NewUsersController(g, lc, userSvc)
}
