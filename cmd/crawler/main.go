package main

import (
	"blockchaincrawler/internal/crawler/crawler"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	c, _ := ethclient.Dial("https://ethereum-holesky.publicnode.com")

	crawler := crawler.NewCrawler(c)

	crawler.Crawl()
}
