package entity

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgtype"
)

type Block struct {
	Hash         string        `json:"hash" gorm:"primaryKey;"`
	Header       Header        `json:"header" gorm:"foreignKey:BlockHash;"`
	Transactions []Transaction `json:"transactions" gorm:"foreignKey:BlockHash;"`
	Withdrawals  []Withdrawal  `json:"withdrawals" gorm:"foreignKey:BlockHash;"`
}

type Header struct {
	BlockHash   string          `json:"blockHash"`
	ParentHash  string          `json:"parentHash"`
	UncleHash   string          `json:"sha3Uncles"`
	Coinbase    string          `json:"miner"`
	Root        string          `json:"stateRoot"`
	TxHash      string          `json:"transactionsRoot"`
	ReceiptHash string          `json:"receiptsRoot"`
	Difficulty  *pgtype.Numeric `json:"difficulty" gorm:"type:numeric;"`
	Number      *pgtype.Numeric `json:"number" gorm:"type:numeric;"`
	GasLimit    uint64          `json:"gasLimit"`
	GasUsed     uint64          `json:"gasUsed"`
	Time        uint64          `json:"timestamp"`
	Extra       []byte          `json:"extraData"`
	MixDigest   string          `json:"mixHash"`
	Nonce       uint64          `json:"nonce"`
}

func (h *Header) From(gethHeader *types.Header) {
	h.BlockHash = gethHeader.Hash().Hex()
	h.ParentHash = gethHeader.ParentHash.Hex()
	h.UncleHash = gethHeader.UncleHash.Hex()
	h.Coinbase = gethHeader.Coinbase.Hex()
	h.Root = gethHeader.Root.Hex()
	h.TxHash = gethHeader.TxHash.Hex()
	h.ReceiptHash = gethHeader.ReceiptHash.Hex()
	h.Difficulty = &pgtype.Numeric{Int: gethHeader.Difficulty, Exp: 0, Status: pgtype.Present}
	h.Number = &pgtype.Numeric{Int: gethHeader.Number, Exp: 0, Status: pgtype.Present}
	h.GasLimit = gethHeader.GasUsed
	h.GasUsed = gethHeader.GasUsed
	h.Time = gethHeader.Time
	h.Extra = gethHeader.Extra
	h.MixDigest = gethHeader.MixDigest.Hex()
	h.Nonce = gethHeader.Nonce.Uint64()
}
