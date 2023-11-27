package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/bogdanshibilov/blockchaincrawler/internal/user/entity"
)

type UseCase interface {
	CreateUser(ctx context.Context, user *entity.User) (uuid.UUID, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserById(ctx context.Context, id string) (*entity.User, error)
	GetAllUsers(ctx context.Context) ([]*entity.User, error)
	DeleteUserById(ctx context.Context, id uuid.UUID) error
	ConfirmUser(ctx context.Context, email string) error
}