package apigateway

import (
	"context"

	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/controller/http/v1/dto"
)

type UseCase interface {
	GetHeaders(ctx context.Context, page int, pageSize int) (*dto.PagedDto, error)
	GetTxsByBlockHash(ctx context.Context, hash string, page int, pageSize int) (*dto.PagedDto, error)
	GetWsByBlockHash(ctx context.Context, hash string, page int, pageSize int) (*dto.PagedDto, error)
}
