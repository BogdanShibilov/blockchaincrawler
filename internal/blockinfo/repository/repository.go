package repository

import (
	"context"

	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/entity"
)

type BlockRepo interface {
	CreateBlock(ctx context.Context, hash string) error
	CreateHeader(ctx context.Context, header *entity.Header) error
	CreateTransaction(ctx context.Context, tx *entity.Transaction) error
	CreateWithdrawal(ctx context.Context, w *entity.Withdrawal) error
	GetHeaders(ctx context.Context, page int, pageSize int) ([]*entity.Header, error)
	GetTotalPagesFor(entity any, pageSize int) int32
}
