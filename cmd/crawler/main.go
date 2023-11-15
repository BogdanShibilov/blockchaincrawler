package main

import (
	"fmt"

	"github.com/bogdanshibilov/blockchaincrawler/internal/crawler/crawler"
	"github.com/bogdanshibilov/blockchaincrawler/pkg/logger"
)

func main() {
	l := logger.NewZap()
	c, err := crawler.NewCrawler("wss://ethereum-sepolia.publicnode.com", l)
	if err != nil {
		l.Panic(err)
	}

	blocks := make(chan *crawler.Result)
	c.CrawlNewBlocks(blocks)

	for b := range blocks {
		fmt.Println(b.Block.Number())
	}
}
