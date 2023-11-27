package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bogdanshibilov/blockchaincrawler/internal/crawler/config"
	"github.com/bogdanshibilov/blockchaincrawler/internal/crawler/crawler"
)

type App struct {
	cfg config.Config
}

func New(cfg config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())

	c, err := a.createCrawler(ctx)
	if err != nil {
		log.Printf("Error while creating crawler: %v\n", err)
	}
	defer c.Close()

	crawlerService := crawler.NewService(c)
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
