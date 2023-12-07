package crawler

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Crawler interface {
	GetBlockByHash(ctx context.Context, hash common.Hash)
	Close()
	Headers() <-chan *types.Header
	Blocks() <-chan *types.Block
	Errors() chan error
	Sub() (ethereum.Subscription, error)
	Unsub()
}
