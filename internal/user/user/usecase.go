package user

import (
	"blockchaincrawler/internal/user/entity"
	"context"

	"github.com/google/uuid"
)

type UseCase interface {
	CreateUser(ctx context.Context, user *entity.User) (uuid.UUID, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	DeleteUserById(ctx context.Context, id uuid.UUID) error
}
