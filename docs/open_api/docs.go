// Package classification next-oms API.
//
// the purpose of this service is to provide & store orders and authenticate and authorize a user
//
//	Schemes: http
//	Host: localhost:8080
//	BasePath: /api
//	Version: v1.0.0
//	License: None
//	Contact: Rezwanul-Haque<rezwanul.cse@gmail.com>
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Security:
//	- base64
//
//	SecurityDefinitions:
//	base64:
//	     type: apiKey
//	     name: ar5go-app-key
//	     in: header
//
// swagger:meta
package openapi

import (
	"next-oms/app/serializers"
	"next-oms/infra/errors"
)

// Generic error message
// swagger:response errorResponse
type errorResponseWrapper struct {
	// Description of the error
	// in: body
	Body errors.RestErr `json:"error_response"`
}

type genericSuccessResponse struct {
	// example: resource created
	Message string `json:"message"`
}

// returns a message
// swagger:response genericSuccessResponse
type genericSuccessResponseWrapper struct {
	// in: body
	genericSuccessResponse `json:"message"`
}

// Fetch users request query params
// swagger:parameters UserQueryParameters
type usersQueryParametersWrapper struct {
	// in:query
	//example: 10
	Size int64 `json:"size"`
	//example: 2
	Page int64 `json:"page"`
	//example: created_at desc
	Sort string `json:"sort"`
	//example: rezwan
	QS string `json:"qs"`
	//example: user_id.(in, contains, equals, gt, gte, lt, lte)
	ColumnOperation string `json:"column:operation"`
}

// Payload for create a user
// swagger:parameters CreateUser
type userPayloadWrapper struct {
	// in:body
	Body serializers.UserReq
}

// response after a user created
// swagger:response UserCreatedResponse
type userCreateRespWrapper struct {
	// in:body
	Body serializers.UserResp
}

// List all the users of a company
// swagger:response UserResponse
type usersRespWrapper struct {
	// in:body
	Body serializers.ListFilters
}
