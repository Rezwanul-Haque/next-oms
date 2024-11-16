package svc

import (
	"next-oms/app/serializers"
)

type IToken interface {
	CreateToken(userID uint) (*serializers.JwtToken, error)
	StoreTokenUuid(userID uint, token *serializers.JwtToken) error
	DeleteTokenUuid(uuid ...string) error
}
