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
	page, pageSize = validatePagingOptions(page, pageSize)

	pagedHeaders, err := s.blocks.GetHeaders(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	headersJson, err := json.Marshal(pagedHeaders.Data)
	if err != nil {
		return nil, err
	}

	return &PagedResult{
		Data:       headersJson,
		Page:       int32(page),
		TotalPages: int32(pagedHeaders.TotalPages),
	}, nil
}

func (s *Service) GetTxsByBlockHash(ctx context.Context, hash string, page int, pageSize int) (*PagedResult, error) {
	page, pageSize = validatePagingOptions(page, pageSize)

	pagedTxs, err := s.blocks.GetTxsByBlockHash(ctx, hash, page, pageSize)
	if err != nil {
		return nil, err
	}

	txsJson, err := json.Marshal(pagedTxs.Data)
	if err != nil {
		return nil, err
	}

	return &PagedResult{
		Data:       txsJson,
		Page:       int32(page),
		TotalPages: int32(pagedTxs.TotalPages),
	}, nil
}

func (s *Service) GetWsByBlockHash(ctx context.Context, hash string, page int, pageSize int) (*PagedResult, error) {
	page, pageSize = validatePagingOptions(page, pageSize)

	pagedWs, err := s.blocks.GetWsByBlockHash(ctx, hash, page, pageSize)
	if err != nil {
		return nil, err
	}

	wsJson, err := json.Marshal(pagedWs.Data)
	if err != nil {
		return nil, err
	}

	return &PagedResult{
		Data:       wsJson,
		Page:       int32(page),
		TotalPages: int32(pagedWs.TotalPages),
	}, nil
}
func (s *Service) GetLastNBlocks(ctx context.Context, count int) ([]byte, error) {
	if count <= 0 || count > 20 {
		count = 20
	}

	blocks, err := s.blocks.GetLastNBlocks(ctx, count)
	if err != nil {
		return nil, err
	}
	jsonBlocks, err := json.Marshal(blocks)
	if err != nil {
		return nil, err
	}

	return jsonBlocks, nil
}

func validatePagingOptions(page int, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}

	switch {
	case pageSize > 50:
		pageSize = 50
	case pageSize <= 0:
		pageSize = 10
	}

	return page, pageSize
}
