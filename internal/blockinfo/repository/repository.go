package repository

import (
	"context"

	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/entity"
)

type BlockRepo interface {
	// CreateBlock(ctx context.Context, hash string) error
	CreateHeader(ctx context.Context, header *entity.Header) error
	CreateTransaction(ctx context.Context, tx *entity.Transaction) error
	CreateWithdrawal(ctx context.Context, w *entity.Withdrawal) error
	GetHeaders(ctx context.Context, page int, pageSize int) (*PagedResult, error)
	GetTxsByBlockHash(ctx context.Context, hash string, page int, pageSize int) (*PagedResult, error)
	GetWsByBlockHash(ctx context.Context, hash string, page int, pageSize int) (*PagedResult, error)
	GetLastNBlocks(ctx context.Context, count int) (blocks []*entity.Block, err error)
}

type PagedResult struct {
	Data       any `json:"data"`
	Page       int `json:"page"`
	TotalPages int `json:"totalPages"`
}
