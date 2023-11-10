package v1

import (
	"blockchaincrawler/internal/auth/auth"
	authpb "blockchaincrawler/pkg/protobuf/authservice/gw"
	"context"
	"fmt"

	"go.uber.org/zap"
)

type Service struct {
	authpb.UnimplementedAuthServiceServer
	usecase auth.UseCase
	logger  *zap.SugaredLogger
}

func NewService(u auth.UseCase, l *zap.SugaredLogger) *Service {
	return &Service{
		usecase: u,
		logger:  l,
	}
}

func (s *Service) GenerateJwtToken(ctx context.Context, request *authpb.GenerateJwtTokenRequest) (*authpb.GenerateJwtTokenResponse, error) {
	token, err := s.usecase.GenerateJwtToken(ctx, request.Email, request.Password)
	if err != nil {
		s.logger.Errorf("failed to generate jwt token error: %v", err)
		return nil, fmt.Errorf("generateJwtToken err: %w", err)
	}

	return &authpb.GenerateJwtTokenResponse{
		Token:        token.Token,
		RefreshToken: token.RefreshToken,
	}, nil
}

func (s *Service) RenewJwtToken(ctx context.Context, request *authpb.RenewJwtTokenRequest) (*authpb.RenewJwtTokenResponse, error) {
	newToken, err := s.usecase.RenewJwtToken(ctx, request.RefreshToken)
	if err != nil {
		s.logger.Errorf("failed to renew jwt token error: %v", err)
		return nil, fmt.Errorf("renewJwtToken err: %w", err)
	}

	return &authpb.RenewJwtTokenResponse{
		Token:        newToken.Token,
		RefreshToken: newToken.RefreshToken,
	}, nil
}