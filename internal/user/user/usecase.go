package user

import (
	"blockchaincrawler/internal/user/entity"
	"context"
)

type UseCase interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	DeleteUserById(ctx context.Context, id int) error
}
