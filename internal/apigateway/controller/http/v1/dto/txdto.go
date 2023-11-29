package dto

import (
	"encoding/json"
	"math/big"

	"github.com/jackc/pgtype"
)

type TxDto struct {
	Hash         string   `json:"hash"`
	BlockHash    string   `json:"blockHash"`
	ChainID      *big.Int `json:"chainId,omitempty"`
	MaxFeePerGas *big.Int `json:"maxFeePerGas"`
	Input        []byte   `json:"input"`
	Nonce        uint64   `json:"nonce"`
	GasPrice     *big.Int `json:"gasPrice"`
	Gas          uint64   `json:"gas"`
	Value        *big.Int `json:"value"`
	V            *big.Int `json:"v"`
	R            *big.Int `json:"r"`
	S            *big.Int `json:"s"`
	To           string   `json:"to"`
}

func (t *TxDto) UnmarshalJSON(data []byte) error {
	aux := struct {
		ChainID      *pgtype.Numeric `json:"chainId,omitempty"`
		MaxFeePerGas *pgtype.Numeric `json:"maxFeePerGas"`
		GasPrice     *pgtype.Numeric `json:"gasPrice"`
		Value        *pgtype.Numeric `json:"value"`
		V            *pgtype.Numeric `json:"v"`
		R            *pgtype.Numeric `json:"r"`
		S            *pgtype.Numeric `json:"s"`

		Hash      string `json:"hash"`
		BlockHash string `json:"blockHash"`
		Input     []byte `json:"input"`
		Nonce     uint64 `json:"nonce"`
		Gas       uint64 `json:"gas"`
		To        string `json:"to"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	t.ChainID = new(big.Int).Set(aux.ChainID.Int)
	t.MaxFeePerGas = new(big.Int).Set(aux.MaxFeePerGas.Int)
	t.GasPrice = new(big.Int).Set(aux.GasPrice.Int)
	t.Value = new(big.Int).Set(aux.Value.Int)
	t.V = new(big.Int).Set(aux.V.Int)
	t.R = new(big.Int).Set(aux.R.Int)
	t.S = new(big.Int).Set(aux.S.Int)

	t.Hash = aux.Hash
	t.BlockHash = aux.BlockHash
	t.Input = aux.Input
	t.Nonce = aux.Nonce
	t.Gas = aux.Gas
	t.To = aux.To

	return nil
}
