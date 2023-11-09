package auth

import "context"

type UseCase interface {
	GenerateJwtToken(ctx context.Context, req GenerateJwtTokenRequest) (*JwtUserToken, error)
	RenewJwtToken(ctx context.Context, refreshToken string) (*JwtUserToken, error)
}
