package transport

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/config"
	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/controller/http/v1/dto"
	pb "github.com/bogdanshibilov/blockchaincrawler/pkg/protobuf/auth/gw"
)

type Auth struct {
	client pb.AuthServiceClient
	cfg    *config.AuthTransport
}

func NewAuth(cfg *config.AuthTransport) (*Auth, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.Dial(cfg.Host+":"+cfg.Port, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial auth grpc server on %v:%v error: %w", cfg.Host, cfg.Port, err)
	}

	client := pb.NewAuthServiceClient(conn)

	return &Auth{
		client: client,
		cfg:    cfg,
	}, nil
}

func (a *Auth) GenerateJwtToken(ctx context.Context, jwtReq *dto.UserCreds) (*dto.JwtToken, error) {
	req := &pb.GenerateJwtTokenRequest{
		Email:    jwtReq.Email,
		Password: jwtReq.Password,
	}

	res, err := a.client.GenerateJwtToken(ctx, req)
	if err != nil {
		return nil, err
	}

	return &dto.JwtToken{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (a *Auth) RenewJwtToken(ctx context.Context, refToken *dto.RenewTokenRequest) (*dto.JwtToken, error) {
	req := &pb.RenewJwtTokenRequest{
		RefreshToken: refToken.RefreshToken,
	}

	res, err := a.client.RenewJwtToken(ctx, req)
	if err != nil {
		return nil, err
	}

	return &dto.JwtToken{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (a *Auth) CreateUser(ctx context.Context, creds *dto.UserCreds) (string, error) {
	req := &pb.CreateUserRequest{
		Email:    creds.Email,
		Password: creds.Password,
	}

	res, err := a.client.CreateUser(ctx, req)
	if err != nil {
		return "", err
	}

	return res.UserId, nil
}

func (a *Auth) SendConfirmationCode(ctx context.Context, sendConfReq *dto.SendConfirmCodeRequest) error {
	req := &pb.SendConfirmationCodeRequest{
		Email: sendConfReq.Email,
	}

	res, err := a.client.SendConfirmationCode(ctx, req)
	if err != nil {
		return err
	}
	_ = res

	return nil
}

func (a *Auth) ConfirmUser(ctx context.Context, confReq *dto.ConfirmUserRequest) error {
	req := &pb.ConfirmUserRequest{
		Email: confReq.Email,
		Code:  confReq.Code,
	}

	res, err := a.client.ConfirmUser(ctx, req)
	if err != nil {
		return err
	}
	_ = res

	return nil
}
