package repository

import (
	"context"
	"math"
	"time"

	"gorm.io/gorm"

	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/database/postgres"
	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/entity"
)

type Block struct {
	main *postgres.Pg
}

func NewBlock(db *postgres.Pg) *Block {
	return &Block{db}
}

func (b *Block) CreateHeader(ctx context.Context, header *entity.Header) error {
	err := b.main.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		newBlock := &entity.Block{
			Hash: header.BlockHash,
		}
		if err := tx.Create(&newBlock).Error; err != nil {
			return err
		}
		if err := tx.Create(&header).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (b *Block) CreateTransaction(ctx context.Context, tx *entity.Transaction) error {
	if err := b.main.DB.WithContext(ctx).First(&entity.Block{Hash: tx.BlockHash}).Error; err != nil {
		time.Sleep(time.Millisecond * 500)
	}

	res := b.main.DB.WithContext(ctx).Create(tx)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (b *Block) CreateWithdrawal(ctx context.Context, w *entity.Withdrawal) error {
	if err := b.main.DB.WithContext(ctx).First(&entity.Block{Hash: w.BlockHash}).Error; err != nil {
		time.Sleep(time.Millisecond * 500)
	}

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

func (b *Block) GetLastNBlocks(ctx context.Context, count int) (blocks []*entity.Block, err error) {
	err = b.main.DB.WithContext(ctx).
		Preload("Header").
		Preload("Transactions").
		Preload("Withdrawals").
		Limit(count).
		Find(&blocks).Error
	if err != nil {
		return nil, err
	}

	return blocks, nil
}
