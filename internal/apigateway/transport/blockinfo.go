package transport

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/config"
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

func (b *BlockInfo) GetHeaders(ctx context.Context, page int, pageSize int) (*pb.GetHeadersResponse, error) {
	req := &pb.GetHeadersRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	res, err := b.client.GetHeaders(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (b *BlockInfo) GetTxsByBlockHash(ctx context.Context, hash string, page int, pageSize int) (*pb.TxsByBlockHashResponse, error) {
	req := &pb.TxsByBlockHashRequest{
		BlockHash: hash,
		Page:      int32(page),
		PageSize:  int32(pageSize),
	}

	res, err := b.client.GetTxsByBlockHash(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (b *BlockInfo) GetWsByBlockHash(ctx context.Context, hash string, page int, pageSize int) (*pb.WsByBlockHashResponse, error) {
	req := &pb.WsByBlockHashRequest{
		BlockHash: hash,
		Page:      int32(page),
		PageSize:  int32(pageSize),
	}

	res, err := b.client.GetWsByBlockHash(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
