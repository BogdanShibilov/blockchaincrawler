package blockinfo

import (
	"context"
)

type UseCase interface {
	CreateHeader(ctx context.Context, header []byte) error
	CreateTransaction(ctx context.Context, txJson []byte, blockHash string) error
	CreateWithdrawal(ctx context.Context, withdrawalJson []byte, blockHash string) error
	GetHeaders(ctx context.Context, page int, pageSize int) (*PagedResult, error)
}
