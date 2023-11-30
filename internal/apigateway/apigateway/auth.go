package apigateway

import (
	"context"

	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/controller/http/v1/dto"
)

func (a *ApiGateway) GenerateJwtToken(ctx context.Context, jwtReq *dto.UserCreds) (*dto.JwtToken, error) {
	return a.auth.GenerateJwtToken(ctx, jwtReq)
}

func (a *ApiGateway) RenewJwtToken(ctx context.Context, refToken *dto.RenewTokenRequest) (*dto.JwtToken, error) {
	return a.auth.RenewJwtToken(ctx, refToken)
}

func (a *ApiGateway) CreateUser(ctx context.Context, creds *dto.UserCreds) (string, error) {
	return a.auth.CreateUser(ctx, creds)
}

func (a *ApiGateway) SendConfirmationCode(ctx context.Context, sendConfReq *dto.SendConfirmCodeRequest) error {
	return a.auth.SendConfirmationCode(ctx, sendConfReq)
}

func (a *ApiGateway) ConfirmUser(ctx context.Context, confReq *dto.ConfirmUserRequest) error {
	return a.auth.ConfirmUser(ctx, confReq)
}
