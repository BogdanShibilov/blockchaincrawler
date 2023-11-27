package blockinfo

import (
	"context"
)

type UseCase interface {
	CreateHeader(ctx context.Context, header []byte) error
}
