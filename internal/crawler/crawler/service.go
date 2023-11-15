package crawler

import (
	"context"
	"fmt"

	"github.com/bogdanshibilov/blockchaincrawler/pkg/logger"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Crawler struct {
	client *ethclient.Client
	l      *logger.ZapLogger
}

func NewCrawler(rawurl string, l *logger.ZapLogger) (*Crawler, error) {
	c, err := ethclient.DialContext(context.Background(), rawurl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %v err: %w", rawurl, err)
	}

	return &Crawler{
		client: c,
		l:      l,
	}, nil
}

func (c *Crawler) CrawlNewBlocks(blocks chan<- *Result) error {
	headers := make(chan *types.Header)
	sub, err := c.client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		return fmt.Errorf("could not subcribe new head: %w", err)
	}

	go func() {
		defer sub.Unsubscribe()
		defer close(headers)

		for {
			select {
			case err := <-sub.Err():
				blocks <- &Result{nil, err}
			case header := <-headers:
				block, err := c.client.BlockByHash(context.TODO(), header.Hash())
				if err != nil {
					blocks <- &Result{nil, err}
				}

				blocks <- &Result{block, nil}
			}
		}
	}()

	return nil
}
