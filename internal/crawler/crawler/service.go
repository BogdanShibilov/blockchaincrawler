package crawler

import (
	"context"
	"log"

	"github.com/bogdanshibilov/blockchaincrawler/internal/crawler/transport"
	"github.com/bogdanshibilov/blockchaincrawler/pkg/logger"
)

type Service struct {
	crawler   Crawler
	l         *logger.SugaredLogger
	blockinfo *transport.BlockInfo
}

func NewService(c Crawler, l *logger.SugaredLogger, bt *transport.BlockInfo) *Service {
	return &Service{
		crawler:   c,
		l:         l,
		blockinfo: bt,
	}
}

func (s *Service) Run(ctx context.Context) {
	var c = s.crawler
	sub, _ := c.Sub()

	for {
		select {
		case header := <-c.Headers():
			s.l.Infof("Received header with number: %v", header.Number)
			go s.blockinfo.CreateHeader(context.Background(), header)
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
