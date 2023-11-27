package crawler

import (
	"context"

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
			blockHash := block.Hash().Hex()
			go s.blockinfo.CreateTransactions(ctx, block.Transactions(), blockHash)
			go s.blockinfo.CreateWithdrawals(ctx, block.Withdrawals(), blockHash)
		case err := <-c.Errors():
			s.l.Errorf("error from channel: %v", err)
		case err := <-sub.Err():
			s.l.Errorf("error from subscription: %v", err)
		}
	}
}
