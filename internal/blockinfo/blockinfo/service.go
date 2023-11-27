package blockinfo

import (
	"context"

	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/entity"
	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/repository"
	"github.com/ethereum/go-ethereum/core/types"
)

type Service struct {
	blocks repository.BlockRepo
}

func New(blocks repository.BlockRepo) UseCase {
	return &Service{blocks: blocks}
}

func (s *Service) CreateHeader(ctx context.Context, headerJson []byte) error {
	gethHeader := &types.Header{}
	gethHeader.UnmarshalJSON(headerJson)

	header := &entity.Header{}
	header.From(gethHeader)

	err := s.blocks.CreateBlock(ctx, header.BlockHash)
	if err != nil {
		return err
	}

	err = s.blocks.CreateHeader(ctx, header)
	if err != nil {
		return err
	}

	return nil
}
