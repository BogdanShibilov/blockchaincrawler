package v1

import (
	"context"
	"fmt"
	"io"

	"go.uber.org/zap"

	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/blockinfo"
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

func (s *Service) CreateHeader(ctx context.Context, req *pb.CreateHeaderRequest) (*pb.Empty, error) {
	err := s.usecase.CreateHeader(ctx, req.HeaderJson)
	if err != nil {
		s.logger.Errorf("failed to create header: %v", err)
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *Service) CreateTransaction(stream pb.BlockInfoService_CreateTransactionServer) error {
	var total int32 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.CreateTransactionResponse{
				TotalCreated: total,
			})
		}
		if err != nil {
			return fmt.Errorf("error while streaming tx: %w", err)
		}

		err = s.usecase.CreateTransaction(context.Background(), req.Transaction, req.BlockHash)
		if err != nil {
			return fmt.Errorf("failed to create tx: %w", err)
		}

		total++
	}
}

func (s *Service) GetHeaders(ctx context.Context, req *pb.GetHeadersRequest) (*pb.GetHeadersResponse, error) {
	pagedHeaders, err := s.usecase.GetHeaders(ctx, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	res := &pb.GetHeadersResponse{
		Headers:    pagedHeaders.Data,
		TotalPages: pagedHeaders.TotalPages,
		Page:       pagedHeaders.Page,
	}

	return res, nil
}

func (s *Service) GetTxsByBlockHash(ctx context.Context, req *pb.TxsByBlockHashRequest) (*pb.TxsByBlockHashResponse, error) {
	pagedTxs, err := s.usecase.GetTxsByBlockHash(ctx, req.BlockHash, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	res := &pb.TxsByBlockHashResponse{
		Txs:        pagedTxs.Data,
		TotalPages: pagedTxs.TotalPages,
		Page:       pagedTxs.Page,
	}

	return res, nil
}

func (s *Service) GetWsByBlockHash(ctx context.Context, req *pb.WsByBlockHashRequest) (*pb.WsByBlockHashResponse, error) {
	pagedWs, err := s.usecase.GetWsByBlockHash(ctx, req.BlockHash, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	res := &pb.WsByBlockHashResponse{
		Ws:         pagedWs.Data,
		TotalPages: pagedWs.TotalPages,
		Page:       pagedWs.Page,
	}

	return res, nil
}
