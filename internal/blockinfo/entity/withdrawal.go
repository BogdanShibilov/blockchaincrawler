package entity

import "github.com/ethereum/go-ethereum/core/types"

type Withdrawal struct {
	BlockHash string `json:"blockHash"`
	Index     uint64 `json:"validatorIndex"`
	Address   string `json:"address"`
	Amount    uint64 `json:"amount"`
}

func (w *Withdrawal) From(gethW *types.Withdrawal) {
	w.Index = gethW.Index
	w.Address = gethW.Address.Hex()
	w.Amount = gethW.Amount
}

func (w *Withdrawal) SetBlockHash(blockHash string) *Withdrawal {
	w.BlockHash = blockHash
	return w
}
