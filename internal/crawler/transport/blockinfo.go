package transport

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
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

func (b *BlockInfo) CreateTransactions(ctx context.Context, txs []*types.Transaction, blockHash string) error {
	stream, err := b.client.CreateTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to create stream: %w", err)
	}

	var totalSent int32 = 0
	for _, tx := range txs {
		txJson, err := tx.MarshalJSON()
		if err != nil {
			return fmt.Errorf("failed to marshal transaction: %w", err)
		}

		req := &pb.CreateTransactionRequest{
			BlockHash:   blockHash,
			Transaction: txJson,
		}

		if err := stream.Send(req); err != nil {
			return fmt.Errorf("failed to send transaction into stream: %w", err)
		}
		totalSent++
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("failed to receive transaction stream response: %w", err)
	}

	if res.TotalCreated != totalSent {
		return errors.New("total sent and created are diffrent")
	}

	return nil
}

func (b *BlockInfo) CreateWithdrawals(ctx context.Context, ws []*types.Withdrawal, blockHash string) error {
	stream, err := b.client.CreateWithdrawal(ctx)
	if err != nil {
		return fmt.Errorf("failed to create stream: %w", err)
	}

	var totalSent int32 = 0
	for _, w := range ws {
		wJson, err := w.MarshalJSON()
		if err != nil {
			return fmt.Errorf("failed to marshal withdrawal: %w", err)
		}

		req := &pb.CreateWithdrawalRequest{
			BlockHash:  blockHash,
			Withdrawal: wJson,
		}

		if err := stream.Send(req); err != nil {
			return fmt.Errorf("failed to send withdrawal into stream: %w", err)
		}
		totalSent++
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("failed to receive withdrawal stream response: %w", err)
	}

	if res.TotalCreated != totalSent {
		return errors.New("total sent and created are diffrent")
	}

	return nil
}
