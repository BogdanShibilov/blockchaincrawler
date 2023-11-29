package dto

import (
	"encoding/json"
	"math/big"

	"github.com/jackc/pgtype"
)

type HeaderDTO struct {
	BlockHash   string   `json:"blockHash"`
	ParentHash  string   `json:"parentHash"`
	UncleHash   string   `json:"sha3Uncles"`
	Coinbase    string   `json:"miner"`
	Root        string   `json:"stateRoot"`
	TxHash      string   `json:"transactionsRoot"`
	ReceiptHash string   `json:"receiptsRoot"`
	Difficulty  *big.Int `json:"difficulty"`
	Number      *big.Int `json:"number"`
	GasLimit    uint64   `json:"gasLimit"`
	GasUsed     uint64   `json:"gasUsed"`
	Time        uint64   `json:"timestamp"`
	Extra       []byte   `json:"extraData"`
	MixDigest   string   `json:"mixHash"`
	Nonce       uint64   `json:"nonce"`
}

func (h *HeaderDTO) UnmarshalJSON(data []byte) error {
	aux := struct {
		Difficulty *pgtype.Numeric `json:"difficulty"`
		Number     *pgtype.Numeric `json:"number"`

		BlockHash   string `json:"blockHash"`
		ParentHash  string `json:"parentHash"`
		UncleHash   string `json:"sha3Uncles"`
		Coinbase    string `json:"miner"`
		Root        string `json:"stateRoot"`
		TxHash      string `json:"transactionsRoot"`
		ReceiptHash string `json:"receiptsRoot"`
		GasLimit    uint64 `json:"gasLimit"`
		GasUsed     uint64 `json:"gasUsed"`
		Time        uint64 `json:"timestamp"`
		Extra       []byte `json:"extraData"`
		MixDigest   string `json:"mixHash"`
		Nonce       uint64 `json:"nonce"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	h.Difficulty = new(big.Int)
	h.Number = new(big.Int)
	if aux.Difficulty != nil && aux.Difficulty.Status == pgtype.Present {
		h.Difficulty.Set(aux.Difficulty.Int)
	}
	if aux.Number != nil && aux.Number.Status == pgtype.Present {
		h.Number.Set(aux.Number.Int)
	}

	h.BlockHash = aux.BlockHash
	h.ParentHash = aux.ParentHash
	h.UncleHash = aux.UncleHash
	h.Coinbase = aux.Coinbase
	h.Root = aux.Root
	h.TxHash = aux.TxHash
	h.ReceiptHash = aux.ReceiptHash
	h.GasLimit = aux.GasLimit
	h.GasUsed = aux.GasUsed
	h.Time = aux.Time
	h.Extra = aux.Extra
	h.MixDigest = aux.MixDigest
	h.Nonce = aux.Nonce

	return nil
}
