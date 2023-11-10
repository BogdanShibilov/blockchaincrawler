package auth

import "context"

type UseCase interface {
	GenerateJwtToken(ctx context.Context, email string, password string) (*JwtUserToken, error)
	RenewJwtToken(ctx context.Context, refreshToken string) (*JwtUserToken, error)
}
