package repository

import (
	"context"

	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/entity"
)

type BlockRepository interface {
	CreateBlock(ctx context.Context, header *entity.Header) error
	GetBlockByHash(ctx context.Context, hash string) (block *entity.Block, err error)
	GetAllBlocks(ctx context.Context) (blocks []*entity.Block, err error)
	GetBlockHeaderByHash(ctx context.Context, hash string) (header *entity.Header, err error)
}
