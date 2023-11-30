package apigateway

import (
	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/transport"
)

type ApiGateway struct {
	blocks *transport.BlockInfo
	auth   *transport.Auth
}

func NewApi(bt *transport.BlockInfo, a *transport.Auth) *ApiGateway {
	return &ApiGateway{
		blocks: bt,
		auth:   a,
	}
}
