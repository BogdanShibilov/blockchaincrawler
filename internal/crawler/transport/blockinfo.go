package transport

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bogdanshibilov/blockchaincrawler/internal/crawler/config"
	pb "github.com/bogdanshibilov/blockchaincrawler/pkg/protobuf/blockinfo/gw"
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

func (b *BlockInfo) CreateBlock(ctx context.Context, req *pb.CreateBlockRequest) error {
	res, err := b.client.CreateBlock(ctx, req)
	if err != nil {
		return err
	}

	_ = res
	return nil
}
