package impl

import (
	"next-oms/app/repository"
	"next-oms/app/serializers"
	"next-oms/app/svc"
)

type system struct {
	repo repository.ISystem
}

func NewSystemService(sysrepo repository.ISystem) svc.ISystem {
	return &system{
		repo: sysrepo,
	}
}

func (sys *system) GetHealth() (*serializers.HealthResp, error) {
	resp := serializers.HealthResp{}

	// check cache
	cacheOnline := sys.repo.CacheCheck()
	resp.CacheOnline = cacheOnline
	// check db
	dbOnline, err := sys.repo.DBCheck()
	resp.DBOnline = dbOnline

	if err != nil {
		return &resp, err
	}

	return &resp, nil
}
