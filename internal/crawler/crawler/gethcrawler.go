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
	log.Println("Closed channels succusfully")
}

func (c *gethCrawler) Unsub() {
	if c.client.Client().SupportsSubscriptions() {
		c.sub.Unsubscribe()
	}
}

func (c *gethCrawler) Headers() <-chan *types.Header {
	return c.headers
}

func (c *gethCrawler) Blocks() <-chan *types.Block {
	return c.blocks
}

func (c *gethCrawler) Errors() chan error {
	return c.errors
}

func (c *gethCrawler) Sub() (ethereum.Subscription, error) {
	if !c.client.Client().SupportsSubscriptions() {
		return nil, rpc.ErrNotificationsUnsupported
	}
	return c.sub, nil
}

func (c *gethCrawler) RetryConnection() {
	client, err := c.conn.Dial()
	if err != nil {
		c.errors <- err
	}

	sub, err := client.SubscribeNewHead(context.Background(), c.headers)
	if err != nil {
		c.errors <- err
	}

	c.sub = sub
}
