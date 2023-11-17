package blockinfo

import (
	"context"
	"fmt"

	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/entity"
	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/repository"
)

type Service struct {
	blocks repository.BlockRepository
}

func New(blocks repository.BlockRepository) UseCase {
	return &Service{blocks: blocks}
}

func (s *Service) CreateBlock(ctx context.Context, header *entity.Header) error {
	err := s.blocks.CreateBlock(ctx, header)
	if err != nil {
		return fmt.Errorf("failed to create block err: %w", err)
	}

	return nil
}

func (s *Service) GetBlockByHash(ctx context.Context, hash string) (*entity.Block, error) {
	block, err := s.blocks.GetBlockByHash(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get block err: %w", err)
	}

	return block, nil
}

func (s *Service) GetAllBlocks(ctx context.Context) ([]*entity.Block, error) {
	blocks, err := s.blocks.GetAllBlocks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all blocks err: %w", err)
	}

	return blocks, nil
}

func (s *Service) GetBlockHeaderByHash(ctx context.Context, hash string) (*entity.Header, error) {
	header, err := s.blocks.GetBlockHeaderByHash(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get header err: %w", err)
	}

	return header, nil
}
