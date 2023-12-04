package apigateway

import (
	"context"

	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/controller/http/v1/dto"
)

type UseCase interface {
	BlocksUseCase
	AuthUseCase
	UserUseCase
}

type BlocksUseCase interface {
	GetHeaders(ctx context.Context, page int, pageSize int) (*dto.PagedDto, error)
	GetTxsByBlockHash(ctx context.Context, hash string, page int, pageSize int) (*dto.PagedDto, error)
	GetWsByBlockHash(ctx context.Context, hash string, page int, pageSize int) (*dto.PagedDto, error)
	GetLastNBlocks(ctx context.Context, count int) ([]*dto.BlockDto, error)
}

type AuthUseCase interface {
	GenerateJwtToken(ctx context.Context, jwtReq *dto.UserCreds) (*dto.JwtToken, error)
	RenewJwtToken(ctx context.Context, refToken *dto.RenewTokenRequest) (*dto.JwtToken, error)
	CreateUser(ctx context.Context, creds *dto.UserCreds) (string, error)
	SendConfirmationCode(ctx context.Context, sendConfReq *dto.SendConfirmCodeRequest) error
	ConfirmUser(ctx context.Context, email string, code string) error
}

type UserUseCase interface {
	GetAllUsers(ctx context.Context) ([]*dto.UserDto, error)
	DeleteUserById(ctx context.Context, id string) error
	UpdateProfile(ctx context.Context, id string, p *dto.UserProfileDto) error
	GetProfileById(ctx context.Context, id string) (*dto.UserProfileDto, error)
}
