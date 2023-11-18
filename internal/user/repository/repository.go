package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/bogdanshibilov/blockchaincrawler/internal/user/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) (uuid.UUID, error)
	GetUserById(ctx context.Context, id uuid.UUID) (user *entity.User, err error)
	GetUserByEmail(ctx context.Context, email string) (user *entity.User, err error)
	GetAllUsers(ctx context.Context) (users []*entity.User, err error)
	UpdateUserById(ctx context.Context, newUser *entity.User) error
	DeleteUserById(ctx context.Context, id uuid.UUID) error
}
