package v1

import (
	"blockchaincrawler/internal/user/entity"
	"blockchaincrawler/internal/user/user"
	pb "blockchaincrawler/pkg/protobuf/userservice/gw"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	pb.UnimplementedUserServiceServer
	usecase user.UseCase
	logger  *zap.SugaredLogger
}

func NewService(u user.UseCase, l *zap.SugaredLogger) *Service {
	return &Service{
		usecase: u,
		logger:  l,
	}
}

func (s *Service) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	userRequest := &entity.User{
		Email:    request.User.Email,
		Password: request.User.Password,
	}

	userId, err := s.usecase.CreateUser(ctx, userRequest)
	if err != nil {
		s.logger.Errorf("failed to CreateUser err: %w", err)
		return nil, fmt.Errorf("CreateUser error: %w", err)
	}

	return &pb.CreateUserResponse{
		Result: userId.String(),
	}, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, request *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	user, err := s.usecase.GetUserByEmail(ctx, request.Email)
	if err != nil {
		s.logger.Errorf("failed to GetUserByEmail err: %w", err)
		return nil, fmt.Errorf("GetUserByEmail error: %w", err)
	}

	return &pb.GetUserByEmailResponse{
		Result: &pb.User{
			Id:          user.ID.String(),
			Email:       user.Email,
			IsConfirmed: user.IsConfirmed,
		},
	}, nil
}

func (s *Service) DeleteUserById(ctx context.Context, request *pb.DeleteUserByIdRequest) (*pb.DeleteUserByIdResponse, error) {
	uuid, err := uuid.Parse(request.Id)
	if err != nil {
		s.logger.Errorf("failed to parse uuid err: %w", err)
		return nil, fmt.Errorf("failed to parse uuid error: %w", err)
	}
	err = s.usecase.DeleteUserById(ctx, uuid)
	if err != nil {
		s.logger.Errorf("failed to DeleteUserById err: %w", err)
		return nil, fmt.Errorf("DeleteUserById error: %w", err)
	}

	return &pb.DeleteUserByIdResponse{
		Id: request.Id,
	}, nil
}

func (s *Service) IsValidLogin(ctx context.Context, request *pb.IsValidLoginRequest) (*pb.IsValidLoginResponse, error) {
	err := s.usecase.IsValidLogin(ctx, request.Email, request.Password)
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return &pb.IsValidLoginResponse{
			IsValid: false,
		}, nil
	} else if err != nil {
		s.logger.Errorf("isValidLogin() error: %w", err)
		return nil, fmt.Errorf("IsValidLogin error %w: ", err)
	}

	return &pb.IsValidLoginResponse{
		IsValid: false,
	}, nil
}
