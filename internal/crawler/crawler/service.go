package crawler

import (
	"context"
	"log"
)

type Service struct {
	crawler Crawler
}

func NewService(crawler Crawler) *Service {
	return &Service{
		crawler: crawler,
	}
}

func (s *Service) Run(ctx context.Context) {
	var c = s.crawler
	sub, _ := c.Sub()

	for {
		select {
		case header := <-c.Headers():
			log.Printf("Got header with number %v\n", header.Number)
			go c.GetBlockByHash(ctx, header.Hash())
		case block := <-c.Blocks():
			log.Printf("Got block with %d transactions\n", len(block.Transactions()))
		case err := <-c.Errors():
			log.Printf("Error from channel: %v", err)
		case err := <-sub.Err():
			log.Printf("Error from sub: %v", err)
		}
	}
}
