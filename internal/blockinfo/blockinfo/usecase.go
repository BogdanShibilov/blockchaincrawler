package blockinfo

import (
	"context"

	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/entity"
)

type UseCase interface {
	CreateBlock(ctx context.Context, header *entity.Header) error
	GetBlockByHash(ctx context.Context, hash string) (*entity.Block, error)
	GetAllBlocks(ctx context.Context) ([]*entity.Block, error)
	GetBlockHeaderByHash(ctx context.Context, hash string) (*entity.Header, error)
}
