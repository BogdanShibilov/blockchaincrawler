package transport

import (
	"context"
	"fmt"

	"github.com/bogdanshibilov/blockchaincrawler/internal/crawler/config"
	pb "github.com/bogdanshibilov/blockchaincrawler/pkg/protobuf/blockinfo/gw"
	"github.com/ethereum/go-ethereum/core/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BlockInfo struct {
	client pb.BlockInfoServiceClient
	cfg    *config.BlockInfoTransport
}

func NewBlockInfo(cfg config.BlockInfoTransport) (*BlockInfo, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.Dial(cfg.Host+":"+cfg.Port, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial user grpc server on %v:%v error: %w", cfg.Host, cfg.Port, err)
	}

	client := pb.NewBlockInfoServiceClient(conn)

	return &BlockInfo{
		client: client,
		cfg:    &cfg,
	}, nil
}

func (b *BlockInfo) CreateHeader(ctx context.Context, header *types.Header) error {
	headerJson, err := header.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal header: %w", err)
	}
	req := &pb.CreateHeaderRequest{
		HeaderJson: headerJson,
	}
	res, err := b.client.CreateHeader(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to make request to blockInfo: %w", err)
	}

	_ = res
	return nil
}
