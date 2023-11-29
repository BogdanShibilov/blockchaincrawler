package repository

import (
	"context"
	"math"

	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/database/postgres"
	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/entity"
)

type Block struct {
	main *postgres.Pg
}

func NewBlock(db *postgres.Pg) *Block {
	return &Block{db}
}

func (b *Block) CreateBlock(ctx context.Context, hash string) error {
	newBlock := &entity.Block{
		Hash: hash,
	}
	res := b.main.DB.WithContext(ctx).Create(newBlock)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (b *Block) CreateHeader(ctx context.Context, header *entity.Header) error {
	res := b.main.DB.WithContext(ctx).Create(header)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (b *Block) CreateTransaction(ctx context.Context, tx *entity.Transaction) error {
	res := b.main.DB.WithContext(ctx).Create(tx)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (b *Block) CreateWithdrawal(ctx context.Context, w *entity.Withdrawal) error {
	res := b.main.DB.WithContext(ctx).Create(w)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (b *Block) GetHeaders(ctx context.Context, page int, pageSize int) (*PagedResult, error) {
	var headers []*entity.Header
	res := b.main.DB.WithContext(ctx).Scopes(paginate(page, pageSize)).Find(&headers)
	if res.Error != nil {
		return nil, res.Error
	}

	var totalRows int64
	b.main.DB.Model(new(entity.Header)).Count(&totalRows)

	return &PagedResult{
		Data:       headers,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(pageSize))),
		Page:       page,
	}, nil
}

func (b *Block) GetTxsByBlockHash(ctx context.Context, hash string, page int, pageSize int) (*PagedResult, error) {
	var txs []*entity.Transaction
	res := b.main.DB.WithContext(ctx).Where("block_hash = ?", hash).Scopes(paginate(page, pageSize)).Find(&txs)
	if res.Error != nil {
		return nil, res.Error
	}

	var totalRows int64
	b.main.DB.Model(new(entity.Transaction)).Where("block_hash = ?", hash).Count(&totalRows)

	return &PagedResult{
		Data:       txs,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(pageSize))),
		Page:       page,
	}, nil
}

func (b *Block) GetWsByBlockHash(ctx context.Context, hash string, page int, pageSize int) (*PagedResult, error) {
	var ws []*entity.Withdrawal
	res := b.main.DB.WithContext(ctx).Where("block_hash = ?", hash).Scopes(paginate(page, pageSize)).Find(&ws)
	if res.Error != nil {
		return nil, res.Error
	}

	var totalRows int64
	b.main.DB.Model(new(entity.Withdrawal)).Where("block_hash = ?", hash).Count(&totalRows)

	return &PagedResult{
		Data:       ws,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(pageSize))),
		Page:       page,
	}, nil
}
