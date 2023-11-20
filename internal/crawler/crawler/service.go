package crawler

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/bogdanshibilov/blockchaincrawler/internal/crawler/transport"
	"github.com/bogdanshibilov/blockchaincrawler/pkg/logger"
	pb "github.com/bogdanshibilov/blockchaincrawler/pkg/protobuf/blockinfo/gw"
)

type Crawler struct {
	client    *ethclient.Client
	l         *logger.ZapLogger
	blockInfo *transport.BlockInfo
}

func NewCrawler(rawurl string, l *logger.ZapLogger, biTransport *transport.BlockInfo) (*Crawler, error) {
	c, err := ethclient.DialContext(context.Background(), rawurl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %v err: %w", rawurl, err)
	}

	return &Crawler{
		client:    c,
		l:         l,
		blockInfo: biTransport,
	}, nil
}

func (c *Crawler) Crawl(
	ctx context.Context,
	blocks chan *types.Block,
	headers chan *types.Header,
	errCh chan error,
) (ethereum.Subscription, error) {
	c.l.Infoln("starting to crawl")
	sub, err := c.client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		return nil, fmt.Errorf("failed to subcribe for new headers: %w", err)
	}

	go func() {
		for {
			select {
			case err := <-sub.Err():
				errCh <- err
			case header := <-headers:
				c.l.Infof("got new header with number: %v\n", header.Number)
				block, err := c.client.BlockByHash(ctx, header.Hash())
				if err != nil {
					errCh <- err
				}
				blocks <- block
			}
		}
	}()

	return sub, nil
}

func (c *Crawler) WriteBlocks(ctx context.Context, blocks chan *types.Block, errCh chan error) {
	for block := range blocks {
		c.l.Infof("writing new block with number: %v\n", block.Number())
		req := &pb.CreateBlockRequest{
			Header: &pb.Header{
				ParentHash:  block.ParentHash().Hex(),
				UncleHash:   block.UncleHash().Hex(),
				Miner:       block.Coinbase().Hex(),
				Root:        block.Root().Hex(),
				TxHash:      block.TxHash().Hex(),
				ReceiptHash: block.ReceiptHash().Hex(),
				Bloom:       block.Bloom().Big().String(),
				Difficulty:  block.Difficulty().Uint64(),
				GasLimit:    block.GasLimit(),
				GasUsed:     block.GasUsed(),
				MixDigest:   block.MixDigest().Hex(),
				Nonce:       block.Nonce(),
				BlockHash:   block.Hash().Hex(),
				Block: &pb.Block{
					Hash:      block.Hash().Hex(),
					Number:    block.NumberU64(),
					Timestamp: block.Time(),
				},
			},
		}
		err := c.blockInfo.CreateBlock(ctx, req)
		if err != nil {
			errCh <- fmt.Errorf("failed to create a new block: %w", err)
		}
	}
}
