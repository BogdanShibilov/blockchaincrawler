package auth

import "context"

type UseCase interface {
	GenerateJwtToken(ctx context.Context, email string, password string) (*JwtToken, error)
	RenewJwtToken(ctx context.Context, refreshToken string) (*JwtToken, error)
}
