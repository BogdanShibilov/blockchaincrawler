package crawler

import "github.com/ethereum/go-ethereum/core/types"

type Result struct {
	Block *types.Block
	Err   error
}
