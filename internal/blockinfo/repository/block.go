package repository

import (
	"context"
	"fmt"

	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/database/postgres"
	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/entity"
)

type Block struct {
	main *postgres.Pg
}

func NewBlock(db *postgres.Pg) *Block {
	return &Block{db}
}

func (b *Block) CreateBlock(ctx context.Context, header *entity.Header) error {
	res := b.main.DB.WithContext(ctx).Create(header)
	if res.Error != nil {
		return fmt.Errorf("failed to create block err: %w", res.Error)
	}

	return nil
}

func (b *Block) GetBlockByHash(ctx context.Context, hash string) (block *entity.Block, err error) {
	res := b.main.DB.WithContext(ctx).Where("hash = ?", hash).First(&block)
	if res.Error != nil {
		return nil, fmt.Errorf("failed to get block err: %w", res.Error)
	}

	return block, nil
}

func (b *Block) GetAllBlocks(ctx context.Context) (blocks []*entity.Block, err error) {
	res := b.main.DB.WithContext(ctx).Find(&blocks)
	if res.Error != nil {
		return nil, fmt.Errorf("failed to get all blocks err: %w", res.Error)
	}

	return blocks, nil
}

func (b *Block) GetBlockHeaderByHash(ctx context.Context, hash string) (header *entity.Header, err error) {
	res := b.main.DB.WithContext(ctx).Where("block_hash = ?", hash).Preload("Block").First(&header)
	if res.Error != nil {
		return nil, fmt.Errorf("failed to get block header err: %w", err)
	}

	return header, nil
}
