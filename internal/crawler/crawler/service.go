package crawler

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/bogdanshibilov/blockchaincrawler/internal/crawler/transport"
	"github.com/bogdanshibilov/blockchaincrawler/pkg/logger"
)

type Service struct {
	crawler   Crawler
	l         *logger.SugaredLogger
	blockinfo *transport.BlockInfo
	wg        *sync.WaitGroup
}

func NewService(c Crawler, l *logger.SugaredLogger, bt *transport.BlockInfo, wg *sync.WaitGroup) *Service {
	return &Service{
		crawler:   c,
		l:         l,
		blockinfo: bt,
		wg:        wg,
	}
}

func (s *Service) Run(ctx context.Context) {
	var c = s.crawler
	sub, _ := c.Sub()

	for {
		select {
		case header := <-c.Headers():
			s.l.Infof("Received header with number: %v", header.Number)
			go s.createHeaders(ctx, header)
			go c.GetBlockByHash(ctx, header.Hash())
		case block := <-c.Blocks():
			blockHash := block.Hash().Hex()
			go s.createTxs(ctx, block.Transactions(), blockHash)
			go s.createWs(ctx, block.Withdrawals(), blockHash)
		case err := <-c.Errors():
			s.l.Errorf("error from channel: %v", err)
		case err := <-sub.Err():
			s.l.Errorf("error from subscription: %v", err)
		case <-ctx.Done():
			s.crawler.Unsub()
			s.wg.Wait()
			return
		}
	}
}

func (s *Service) createHeaders(ctx context.Context, header *types.Header) {
	s.wg.Add(1)
	defer s.wg.Done()

	err := s.blockinfo.CreateHeader(ctx, header)
	if err != nil {
		s.crawler.Errors() <- err
	}

}

func (s *Service) createTxs(ctx context.Context, txs []*types.Transaction, blockHash string) {
	s.wg.Add(1)
	defer s.wg.Done()

	err := s.blockinfo.CreateTransactions(ctx, txs, blockHash)
	if err != nil {
		s.crawler.Errors() <- err
	}
}

func (s *Service) createWs(ctx context.Context, ws []*types.Withdrawal, blockHash string) {
	s.wg.Add(1)
	defer s.wg.Done()

	err := s.blockinfo.CreateWithdrawals(ctx, ws, blockHash)
	if err != nil {
		s.crawler.Errors() <- err
	}
}
