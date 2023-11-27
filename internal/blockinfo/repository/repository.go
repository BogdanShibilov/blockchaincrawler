package repository

import (
	"context"

	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/entity"
)

type BlockRepo interface {
	CreateHeader(ctx context.Context, header *entity.Header) error
	CreateBlock(ctx context.Context, hash string) error
}
