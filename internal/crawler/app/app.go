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
)

type App struct {
	cfg config.Config
	l   *logger.SugaredLogger
}

func New(cfg config.Config, l *logger.SugaredLogger) *App {
	return &App{
		cfg: cfg,
		l:   l,
	}
}

func (a *App) Run() {

	ctx, cancel := context.WithCancel(context.Background())

	blockInfoTransport, err := transport.NewBlockInfo(a.cfg.Transport.BlockInfoTransport)
	if err != nil {
		a.l.Panicf("failed to create block info transport: %v", err)
	}

	c, err := a.createCrawler(ctx)
	if err != nil {
		a.l.Panicf("failed to create crawler: %v", err)
	}
	defer c.Close()

	crawlerService := crawler.NewService(c, a.l, blockInfoTransport)
	a.l.Infof("Starting crawler")
	go crawlerService.Run(ctx)

	a.gracefulShutdown(cancel)
}

func (a *App) createCrawler(ctx context.Context) (crawler.Crawler, error) {
	p := crawler.Protocol(a.cfg.NodeUrl.Protocol)
	nodeConn := crawler.NewGethConnection(p, a.cfg.NodeUrl.Hostname)
	crawlerBuilder := crawler.NewGethCrawlerBuilder(*nodeConn)
	crawlerDirector := crawler.NewDirector(crawlerBuilder)
	return crawlerDirector.BuildCrawler(ctx)
}

func (a *App) gracefulShutdown(cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
	signal.Stop(ch)
	cancel()
}
