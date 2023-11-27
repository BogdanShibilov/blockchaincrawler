package entity

import "github.com/ethereum/go-ethereum/core/types"

type Block struct {
	Hash   string `json:"hash" gorm:"primaryKey;"`
	Header Header `json:"header" gorm:"foreignKey:BlockHash;"`
}

type Header struct {
	BlockHash   string `json:"blockHash"`
	ParentHash  string `json:"parentHash"`
	UncleHash   string `json:"sha3Uncles"`
	Coinbase    string `json:"miner"`
	Root        string `json:"stateRoot"`
	TxHash      string `json:"transactionsRoot"`
	ReceiptHash string `json:"receiptsRoot"`
	Difficulty  uint64 `json:"difficulty"` // TODO: big.int might be given which is too big
	Number      uint64 `json:"number"`     //       big.int might be given which is too big
	GasLimit    uint64 `json:"gasLimit"`
	GasUsed     uint64 `json:"gasUsed"`
	Time        uint64 `json:"timestamp"`
	Extra       []byte `json:"extraData"`
	MixDigest   string `json:"mixHash"`
	Nonce       uint64 `json:"nonce"`
}

func (h *Header) From(gethHeader *types.Header) {
	h.BlockHash = gethHeader.Hash().Hex()
	h.ParentHash = gethHeader.ParentHash.Hex()
	h.UncleHash = gethHeader.UncleHash.Hex()
	h.Coinbase = gethHeader.Coinbase.Hex()
	h.Root = gethHeader.Root.Hex()
	h.TxHash = gethHeader.TxHash.Hex()
	h.ReceiptHash = gethHeader.ReceiptHash.Hex()
	h.Difficulty = gethHeader.Difficulty.Uint64()
	h.Number = gethHeader.Number.Uint64()
	h.GasLimit = gethHeader.GasUsed
	h.GasUsed = gethHeader.GasUsed
	h.Time = gethHeader.Time
	h.Extra = gethHeader.Extra
	h.MixDigest = gethHeader.MixDigest.Hex()
	h.Nonce = gethHeader.Nonce.Uint64()
}
