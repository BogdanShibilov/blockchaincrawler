package auth

import (
	"context"
)

type UseCase interface {
	User
	Jwt
}

type Jwt interface {
	GenerateJwtToken(ctx context.Context, email string, password string) (*JwtToken, error)
	RenewJwtToken(ctx context.Context, refreshToken string) (*JwtToken, error)
}

type User interface {
	CreateUser(ctx context.Context, email string, password string) (string, error)
	SendConfirmationCode(ctx context.Context, email string) error
	ConfirmUser(ctx context.Context, email string, code string) error
}
