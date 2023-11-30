package apigateway

import (
	"context"
	"encoding/json"

	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/controller/http/v1/dto"
)

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

func (a *ApiGateway) GetTxsByBlockHash(ctx context.Context, hash string, page int, pageSize int) (*dto.PagedDto, error) {
	res, err := a.blocks.GetTxsByBlockHash(ctx, hash, page, pageSize)
	if err != nil {
		return nil, err
	}

	var txs []*dto.TxDto
	err = json.Unmarshal(res.Txs, &txs)
	if err != nil {
		return nil, err
	}

	return &dto.PagedDto{
		Page:       int(res.Page),
		TotalPages: int(res.TotalPages),
		Value:      txs,
	}, nil
}

func (a *ApiGateway) GetWsByBlockHash(ctx context.Context, hash string, page int, pageSize int) (*dto.PagedDto, error) {
	res, err := a.blocks.GetWsByBlockHash(ctx, hash, page, pageSize)
	if err != nil {
		return nil, err
	}

	var ws []*dto.WithdrawalDto
	err = json.Unmarshal(res.Ws, &ws)
	if err != nil {
		return nil, err
	}

	return &dto.PagedDto{
		Page:       int(res.Page),
		TotalPages: int(res.TotalPages),
		Value:      ws,
	}, nil
}
