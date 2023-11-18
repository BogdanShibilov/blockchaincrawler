package v1

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/blockinfo"
	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/entity"
	pb "github.com/bogdanshibilov/blockchaincrawler/pkg/protobuf/blockinfo/gw"
)

type Service struct {
	pb.UnimplementedBlockInfoServiceServer
	usecase blockinfo.UseCase
	logger  *zap.SugaredLogger
}

func NewService(u blockinfo.UseCase, l *zap.SugaredLogger) *Service {
	return &Service{
		usecase: u,
		logger:  l,
	}
}

func (s *Service) CreateBlock(ctx context.Context, req *pb.CreateBlockRequest) (*pb.CreateBlockResponse, error) {
	header := &entity.Header{
		ParentHash:  req.Header.ParentHash,
		UncleHash:   req.Header.UncleHash,
		Miner:       req.Header.Miner,
		Root:        req.Header.Root,
		TxHash:      req.Header.TxHash,
		ReceiptHash: req.Header.ReceiptHash,
		Bloom:       req.Header.Bloom,
		Difficulty:  req.Header.Difficulty,
		GasLimit:    req.Header.GasLimit,
		GasUsed:     req.Header.GasUsed,
		Extra:       req.Header.Extra,
		MixDigest:   req.Header.MixDigest,
		Nonce:       req.Header.Nonce,

		BlockHash: req.Header.BlockHash,
		Block: entity.Block{
			Hash:      req.Header.Block.Hash,
			Number:    req.Header.Block.Number,
			Timestamp: req.Header.Block.Timestamp,
		},
	}

	err := s.usecase.CreateBlock(ctx, header)
	if err != nil {
		s.logger.Errorf("failed to create header err: %v", err)
		return nil, fmt.Errorf("failed to create header err: %w", err)
	}

	return &pb.CreateBlockResponse{}, nil
}

func (s *Service) GetBlockByHash(ctx context.Context, req *pb.GetBlockByHashRequest) (*pb.GetBlockByHashResponse, error) {
	block, err := s.usecase.GetBlockByHash(ctx, req.Hash)
	if err != nil {
		s.logger.Errorf("failed to get block by hash err: %v", err)
		return nil, fmt.Errorf("failed to get block by hash err: %w", err)
	}

	return &pb.GetBlockByHashResponse{
		Block: &pb.Block{
			Hash:      block.Hash,
			Number:    block.Number,
			Timestamp: block.Timestamp,
		},
	}, nil
}

func (s *Service) GetAllBlocks(ctx context.Context, req *pb.GetAllBlocksRequest) (*pb.GetAllBlocksResponse, error) {
	blocks, err := s.usecase.GetAllBlocks(ctx)
	if err != nil {
		s.logger.Errorf("failed to get all blocks err: %v", err)
		return nil, fmt.Errorf("failed to get all blocks err: %w", err)
	}

	var pbBlocks []*pb.Block
	for _, block := range blocks {
		pbBlocks = append(pbBlocks, &pb.Block{
			Hash:      block.Hash,
			Number:    block.Number,
			Timestamp: block.Timestamp,
		})
	}

	return &pb.GetAllBlocksResponse{
		Blocks: pbBlocks,
	}, nil
}

func (s *Service) GetBlockHeaderByHash(ctx context.Context, req *pb.GetBlockHeaderByHashRequest) (*pb.GetBlockHeaderByHashResponse, error) {
	header, err := s.usecase.GetBlockHeaderByHash(ctx, req.Hash)
	if err != nil {
		s.logger.Errorf("failed to get header by hash err: %v", err)
		return nil, fmt.Errorf("failed to get header by hash err: %w", err)
	}

	return &pb.GetBlockHeaderByHashResponse{
		Header: &pb.Header{
			ParentHash:  header.ParentHash,
			UncleHash:   header.UncleHash,
			Miner:       header.Miner,
			Root:        header.Root,
			TxHash:      header.TxHash,
			ReceiptHash: header.ReceiptHash,
			Bloom:       header.Bloom,
			Difficulty:  header.Difficulty,
			GasLimit:    header.GasLimit,
			GasUsed:     header.GasUsed,
			Extra:       header.Extra,
			MixDigest:   header.MixDigest,
			Nonce:       header.Nonce,
			BlockHash:   header.BlockHash,
			Block: &pb.Block{
				Hash:      header.Block.Hash,
				Number:    header.Block.Number,
				Timestamp: header.Block.Timestamp,
			},
		},
	}, nil
}
