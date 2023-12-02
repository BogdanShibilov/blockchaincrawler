package apigateway

import (
	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/transport"
)

type ApiGateway struct {
	blocks *transport.BlockInfo
	auth   *transport.Auth
	users  *transport.User
}

func NewApi(bt *transport.BlockInfo, a *transport.Auth, u *transport.User) *ApiGateway {
	return &ApiGateway{
		blocks: bt,
		auth:   a,
		users:  u,
	}
}
