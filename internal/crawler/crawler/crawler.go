package crawler

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Crawler struct {
	*ethclient.Client
}

func NewCrawler(c *ethclient.Client) *Crawler {
	return &Crawler{c}
}

func (c *Crawler) Crawl() {
	tick := time.NewTicker(time.Second * 5)
	done := make(chan bool)
	var prevBlock *types.Block

	go func() {
		for {
			b := c.getLatestBlock()
			if prevBlock == nil {
				prevBlock = b
			}
			if b.Hash().Cmp(prevBlock.Hash()) != 0 {
				prevBlock = b
				b = c.getLatestBlock()
				fmt.Printf("Crawled block: %v\n", b.Hash())
				fmt.Printf("Number: %v\n", b.Number())
				fmt.Printf("Extra: %v\n", string(b.Extra()))
			}

			select {
			case <-done:
				return
			case t := <-tick.C:
				fmt.Println("ticked at ", t)
			}
		}
	}()

	time.Sleep(time.Minute)

}

func (c *Crawler) getLatestBlock() *types.Block {
	block, _ := c.BlockByNumber(context.TODO(), nil)

	return block
}
