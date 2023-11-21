package v1

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/bogdanshibilov/blockchaincrawler/internal/auth/auth"
	pb "github.com/bogdanshibilov/blockchaincrawler/pkg/protobuf/auth/gw"
)

type Service struct {
	pb.UnimplementedAuthServiceServer
	usecase auth.UseCase
	logger  *zap.SugaredLogger
}

func NewService(u auth.UseCase, l *zap.SugaredLogger) *Service {
	return &Service{
		usecase: u,
		logger:  l,
	}
}

func (s *Service) GenerateJwtToken(ctx context.Context, req *pb.GenerateJwtTokenRequest) (*pb.GenerateJwtTokenResponse, error) {
	jwt, err := s.usecase.GenerateJwtToken(ctx, req.Email, req.Password)
	if err != nil {
		s.logger.Errorf("failed to generate jwt: %v", err)
		return nil, fmt.Errorf("failed to generate jwt: %w", err)
	}

	return &pb.GenerateJwtTokenResponse{
		AccessToken:  jwt.AccessToken,
		RefreshToken: jwt.RefreshToken,
	}, nil
}

func (s *Service) RenewJwtToken(ctx context.Context, req *pb.RenewJwtTokenRequest) (*pb.RenewJwtTokenResponse, error) {
	jwt, err := s.usecase.RenewJwtToken(ctx, req.RefreshToken)
	if err != nil {
		s.logger.Errorf("failed to renew jwt: %v", err)
		return nil, fmt.Errorf("failed to renew jwt: %w", err)
	}

	return &pb.RenewJwtTokenResponse{
		AccessToken:  jwt.AccessToken,
		RefreshToken: jwt.RefreshToken,
	}, nil
}

func (s *Service) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	userId, err := s.usecase.CreateUser(ctx, req.Email, req.Password)
	if err != nil {
		s.logger.Errorf("failed to create user: %v", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &pb.CreateUserResponse{
		UserId: userId,
	}, nil
}

func (s *Service) SendConfirmationCode(ctx context.Context, req *pb.SendConfirmationCodeRequest) (*pb.SendConfirmationCodeResponse, error) {
	err := s.usecase.SendConfirmationCode(ctx, req.Email)
	if err != nil {
		s.logger.Errorf("failed to send confirmation code: %v", err)
		return nil, fmt.Errorf("failed to send confirmation code: %w", err)
	}

	return &pb.SendConfirmationCodeResponse{}, nil
}
