package apigateway

import (
	"context"
	"encoding/json"

	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/controller/http/v1/dto"
	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/transport"
)

type ApiGateway struct {
	blocks *transport.BlockInfo
}

func NewApi(bt *transport.BlockInfo) *ApiGateway {
	return &ApiGateway{
		blocks: bt,
	}
}

func (a *ApiGateway) GetHeaders(ctx context.Context, page int, pageSize int) (*dto.PagedDto, error) {
	res, err := a.blocks.GetHeaders(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	var headers []*dto.HeaderDTO
	err = json.Unmarshal(res.Headers, &headers)
	if err != nil {
		return nil, err
	}

	return &dto.PagedDto{
		Page:       int(res.Page),
		TotalPages: int(res.TotalPages),
		Value:      headers,
	}, nil
}
