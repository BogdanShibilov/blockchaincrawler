package entity

type Bloom [256]byte
type Address [20]byte

type Block struct {
	Hash      string `json:"hash" gorm:"primaryKey;"`
	Number    uint64 `json:"number"`
	Timestamp uint64 `json:"timestamp"`
}

type Header struct {
	Block     Block  `json:"block" gorm:"foreignKey:BlockHash;"`
	BlockHash string `json:"blockHash"`

	ParentHash  string `json:"parentHash"`
	UncleHash   string `json:"sha3Uncles"`
	Miner       string `json:"miner"`
	Root        string `json:"stateRoot"`
	TxHash      string `json:"transactionsRoot"`
	ReceiptHash string `json:"receiptsRoot"`
	Bloom       string `json:"logsBloom"`
	Difficulty  uint64 `json:"difficulty"`
	GasLimit    uint64 `json:"gasLimit"`
	GasUsed     uint64 `json:"gasUsed"`
	Extra       string `json:"extraData"`
	MixDigest   string `json:"mixHash"`
	Nonce       uint64 `json:"nonce"`
}
