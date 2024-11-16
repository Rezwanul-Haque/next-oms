package svc

import "next-oms/app/serializers"

type ISystem interface {
	GetHealth() (*serializers.HealthResp, error)
}
