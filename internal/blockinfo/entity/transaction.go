package entity

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgtype"
)

type Transaction struct {
	Hash         string          `json:"hash" gorm:"primaryKey;"`
	BlockHash    string          `json:"blockHash"`
	ChainID      *pgtype.Numeric `json:"chainId,omitempty" gorm:"type:numeric;"`
	MaxFeePerGas *pgtype.Numeric `json:"maxFeePerGas" gorm:"type:numeric;"`
	Input        []byte          `json:"input"`
	Nonce        uint64          `json:"nonce"`
	GasPrice     *pgtype.Numeric `json:"gasPrice" gorm:"type:numeric;"`
	Gas          uint64          `json:"gas"`
	Value        *pgtype.Numeric `json:"value" gorm:"type:numeric;"`
	V            *pgtype.Numeric `json:"v" gorm:"type:numeric;"`
	R            *pgtype.Numeric `json:"r" gorm:"type:numeric;"`
	S            *pgtype.Numeric `json:"s" gorm:"type:numeric;"`
	To           string          `json:"to"`
}

func (t *Transaction) From(gethTx *types.Transaction) {
	t.Hash = gethTx.Hash().Hex()
	t.ChainID = &pgtype.Numeric{Int: gethTx.ChainId(), Exp: 0, Status: pgtype.Present}
	t.MaxFeePerGas = &pgtype.Numeric{Int: gethTx.GasFeeCap(), Exp: 0, Status: pgtype.Present}
	t.Input = gethTx.Data()
	t.Nonce = gethTx.Nonce()
	t.GasPrice = &pgtype.Numeric{Int: gethTx.GasPrice(), Exp: 0, Status: pgtype.Present}
	t.Gas = gethTx.Gas()
	t.Value = &pgtype.Numeric{Int: gethTx.Value(), Exp: 0, Status: pgtype.Present}
	v, r, s := gethTx.RawSignatureValues()
	t.V = &pgtype.Numeric{Int: v, Exp: 0, Status: pgtype.Present}
	t.R = &pgtype.Numeric{Int: r, Exp: 0, Status: pgtype.Present}
	t.S = &pgtype.Numeric{Int: s, Exp: 0, Status: pgtype.Present}
	t.To = gethTx.To().Hex()
}

func (t *Transaction) SetBlockHash(blockHash string) *Transaction {
	t.BlockHash = blockHash
	return t
}
