package apigateway

import (
	"context"

	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/controller/http/v1/dto"
)

type UseCase interface {
	BlocksUseCase
	AuthUseCase
}

type BlocksUseCase interface {
	GetHeaders(ctx context.Context, page int, pageSize int) (*dto.PagedDto, error)
	GetTxsByBlockHash(ctx context.Context, hash string, page int, pageSize int) (*dto.PagedDto, error)
	GetWsByBlockHash(ctx context.Context, hash string, page int, pageSize int) (*dto.PagedDto, error)
}

type AuthUseCase interface {
	GenerateJwtToken(ctx context.Context, jwtReq *dto.UserCreds) (*dto.JwtToken, error)
	RenewJwtToken(ctx context.Context, refToken *dto.RenewTokenRequest) (*dto.JwtToken, error)
	CreateUser(ctx context.Context, creds *dto.UserCreds) (string, error)
	SendConfirmationCode(ctx context.Context, sendConfReq *dto.SendConfirmCodeRequest) error
	ConfirmUser(ctx context.Context, confReq *dto.ConfirmUserRequest) error
}
