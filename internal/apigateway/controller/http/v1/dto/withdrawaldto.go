package dto

type WithdrawalDto struct {
	BlockHash string `json:"blockHash"`
	Index     uint64 `json:"validatorIndex"`
	Address   string `json:"address"`
	Amount    uint64 `json:"amount"`
}
