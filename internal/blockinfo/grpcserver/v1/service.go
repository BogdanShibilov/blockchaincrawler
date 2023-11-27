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
