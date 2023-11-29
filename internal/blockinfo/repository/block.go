package repository

import (
	"context"

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

func (b *Block) GetHeaders(ctx context.Context, page int, pageSize int) (headers []*entity.Header, err error) {
	res := b.main.DB.WithContext(ctx).Scopes(paginate(page, pageSize)).Find(&headers)
	if res.Error != nil {
		return nil, res.Error
	}

	return headers, nil
}
