package repository

import (
	"blockchaincrawler/internal/user/entity"
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) (uuid.UUID, error)
	GetUserById(ctx context.Context, id uuid.UUID) (user *entity.User, err error)
	GetUserByEmail(ctx context.Context, email string) (user *entity.User, err error)
	UpdateUserById(ctx context.Context, newUser *entity.User) error
	DeleteUserById(ctx context.Context, id uuid.UUID) error
}
