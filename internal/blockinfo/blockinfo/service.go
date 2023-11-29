package blockinfo

import (
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/entity"
	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/repository"
)

type PagedResult struct {
	Data       []byte
	Page       int32
	TotalPages int32
}

type Service struct {
	blocks repository.BlockRepo
}

func New(blocks repository.BlockRepo) UseCase {
	return &Service{blocks: blocks}
}

func (s *Service) CreateHeader(ctx context.Context, headerJson []byte) error {
	gethHeader := new(types.Header)
	err := gethHeader.UnmarshalJSON(headerJson)
	if err != nil {
		return err
	}

	header := new(entity.Header)
	header.From(gethHeader)

	err = s.blocks.CreateBlock(ctx, header.BlockHash)
	if err != nil {
		return err
	}

	err = s.blocks.CreateHeader(ctx, header)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CreateTransaction(ctx context.Context, txJson []byte, blockHash string) error {
	gethTx := new(types.Transaction)
	err := gethTx.UnmarshalJSON(txJson)
	if err != nil {
		return err
	}

	tx := new(entity.Transaction)
	tx.SetBlockHash(blockHash).From(gethTx)

	err = s.blocks.CreateTransaction(ctx, tx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CreateWithdrawal(ctx context.Context, withdrawalJson []byte, blockHash string) error {
	gethW := new(types.Withdrawal)
	err := gethW.UnmarshalJSON(withdrawalJson)
	if err != nil {
		return err
	}

	w := new(entity.Withdrawal)
	w.SetBlockHash(blockHash).From(gethW)

	err = s.blocks.CreateWithdrawal(ctx, w)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetHeaders(ctx context.Context, page int, pageSize int) (*PagedResult, error) {
	if page <= 0 {
		page = 1
	}

	switch {
	case pageSize > 50:
		pageSize = 50
	case pageSize <= 0:
		pageSize = 10
	}

	headers, err := s.blocks.GetHeaders(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	headersJson, err := json.Marshal(headers)
	if err != nil {
		return nil, err
	}

	totalPages := s.blocks.GetTotalPagesFor(new(entity.Header), pageSize)

	return &PagedResult{
		Data:       headersJson,
		Page:       int32(page),
		TotalPages: totalPages,
	}, nil
}
