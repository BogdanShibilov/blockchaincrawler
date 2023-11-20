package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bogdanshibilov/blockchaincrawler/internal/crawler/config"
	"github.com/bogdanshibilov/blockchaincrawler/internal/crawler/crawler"
	"github.com/bogdanshibilov/blockchaincrawler/internal/crawler/transport"
	"github.com/bogdanshibilov/blockchaincrawler/pkg/logger"
	"github.com/ethereum/go-ethereum/core/types"
)

type App struct {
	logger *logger.ZapLogger
	cfg    *config.Config
}

func New(l *logger.ZapLogger, cfg *config.Config) *App {
	return &App{
		logger: l,
		cfg:    cfg,
	}
}

func (a *App) Run() {
	var (
		l   = a.logger
		cfg = a.cfg
	)

	ctx, cancel := context.WithCancel(context.TODO())

	blockInfoTransport, err := transport.NewBlockInfo(cfg.Transport.BlockInfoTransport)
	if err != nil {
		l.Panicf("failed to connect to blockinfo: %v", err)
	}

	blockCrawler, err := crawler.NewCrawler(cfg.ExternalNode.Protocol+"://"+cfg.ExternalNode.Hostname, l, blockInfoTransport)
	if err != nil {
		l.Panicf("failed to create a new crawler: %v", err)
	}

	blocks := make(chan *types.Block)
	headers := make(chan *types.Header)
	errCh := make(chan error)
	sub, err := blockCrawler.Crawl(ctx, blocks, headers, errCh)
	if err != nil {
		l.Panicf("failed to crawl: %v", err)
	}
	defer sub.Unsubscribe()
	defer close(blocks)
	defer close(headers)
	defer close(errCh)
	go blockCrawler.WriteBlocks(ctx, blocks, errCh)

	go func() {
		for err := range errCh {
			l.Errorf("error occured: %v", err)
		}
	}()

	a.gracefulShutdown(cancel)
	l.Infoln("End")
}

func (a *App) gracefulShutdown(cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
	signal.Stop(ch)
	cancel()
}
