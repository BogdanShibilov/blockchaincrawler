package crawler

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type CrawlerBuilder interface {
	setClient() error
	setHeadersCh()
	setBlocksCh()
	setErrorsCh()
	setSubcription(ctx context.Context) error
	buildCrawler() Crawler
}

type gethCrawlerBuilder struct {
	conn    gethConnection
	client  *ethclient.Client
	headers chan *types.Header
	blocks  chan *types.Block
	errors  chan error
	sub     ethereum.Subscription
}

func NewGethCrawlerBuilder(conn gethConnection) CrawlerBuilder {
	return &gethCrawlerBuilder{
		conn: conn,
	}
}

func (b *gethCrawlerBuilder) setClient() error {
	client, err := b.conn.Dial()
	if err != nil {
		return err
	}
	b.client = client
	return nil
}

func (b *gethCrawlerBuilder) setHeadersCh() {
	b.headers = make(chan *types.Header)
}

func (b *gethCrawlerBuilder) setBlocksCh() {
	b.blocks = make(chan *types.Block)
}

func (b *gethCrawlerBuilder) setErrorsCh() {
	b.errors = make(chan error)
}

func (b *gethCrawlerBuilder) setSubcription(ctx context.Context) error {
	if b.client.Client().SupportsSubscriptions() {
		sub, err := b.client.SubscribeNewHead(ctx, b.headers)
		if err != nil {
			return err
		}
		b.sub = sub
	}
	return nil
}

func (b *gethCrawlerBuilder) buildCrawler() Crawler {
	return &gethCrawler{
		conn:    b.conn,
		client:  b.client,
		headers: b.headers,
		blocks:  b.blocks,
		errors:  b.errors,
		sub:     b.sub,
	}
}
