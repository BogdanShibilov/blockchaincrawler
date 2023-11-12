package auth

import (
	"context"

	"github.com/google/uuid"
)

type UseCase interface {
	GenerateJwtToken(ctx context.Context, email string, password string) (*JwtUserToken, error)
	RenewJwtToken(ctx context.Context, refreshToken string) (*JwtUserToken, error)
	CreateUser(ctx context.Context, email string, password string) (*uuid.UUID, error)
	ConfirmUser(ctx context.Context, email string, code string) error
}
