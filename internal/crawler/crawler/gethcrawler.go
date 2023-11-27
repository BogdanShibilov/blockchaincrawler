package crawler

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type gethCrawler struct {
	conn    gethConnection
	client  *ethclient.Client
	headers chan *types.Header
	blocks  chan *types.Block
	errors  chan error
	sub     ethereum.Subscription
}

func (c *gethCrawler) GetBlockByHash(ctx context.Context, hash common.Hash) {
	block, err := c.client.BlockByHash(ctx, hash)
	if err != nil {
		c.errors <- err
	} else {
		c.blocks <- block
	}
}

func (c *gethCrawler) Close() {
	close(c.headers)
	close(c.blocks)
	close(c.errors)
	if c.client.Client().SupportsSubscriptions() {
		c.sub.Unsubscribe()
	}
	log.Println("Closed channels succusfully")
}

func (c *gethCrawler) Headers() <-chan *types.Header {
	return c.headers
}

func (c *gethCrawler) Blocks() <-chan *types.Block {
	return c.blocks
}

func (c *gethCrawler) Errors() <-chan error {
	return c.errors
}

func (c *gethCrawler) Sub() (ethereum.Subscription, error) {
	if !c.client.Client().SupportsSubscriptions() {
		return nil, rpc.ErrNotificationsUnsupported
	}
	return c.sub, nil
}
