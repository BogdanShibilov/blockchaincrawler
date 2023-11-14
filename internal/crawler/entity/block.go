package entity

import "math/big"

type Block struct {
	Hash          string
	Number        uint64
	BaseFeePerGas *big.Int
	Difficulty    *big.Int
	ExtraData     string
	GasLimit      uint64
	GasUsed       uint64
	Miner         string
	MixHash       string
	Nonce         string
	ParentHash    string
	ReceiptsRoot  string
	Sha3Uncles    string
	Size          uint64
	StateRoot     string
	Timestamp     uint64
	// Transactions     []Transaction
	TransactionsRoot string
	// Withdrawals      []Withdrawal
	WithdrawalsRoot string
}
