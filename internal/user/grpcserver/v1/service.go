package v1

import (
	"context"
	"fmt"
	"time"

	"github.com/bogdanshibilov/blockchaincrawler/internal/user/entity"
	"github.com/bogdanshibilov/blockchaincrawler/internal/user/user"
	pb "github.com/bogdanshibilov/blockchaincrawler/pkg/protobuf/user/gw"
	"github.com/google/uuid"
	"go.uber.org/zap"
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

func (s *Service) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	createdAt := time.Now()
	newUser := &entity.User{
		Email:       req.Email,
		Password:    req.HashedPassword,
		IsConfirmed: false,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}

	id, err := s.usecase.CreateUser(ctx, newUser)
	if err != nil {
		s.logger.Errorf("failed to create user: %v", err)
		return nil, fmt.Errorf("faield to create user: %w", err)
	}

	return &pb.CreateUserResponse{
		Id: id.String(),
	}, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	user, err := s.usecase.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Errorf("failed to get user by email: %v", err)
		return nil, fmt.Errorf("faield to get user by email: %w", err)
	}

	return &pb.GetUserByEmailResponse{
		User: &pb.User{
			Id:             user.ID.String(),
			Email:          user.Email,
			HashedPassword: user.Password,
			IsConfirmed:    user.IsConfirmed,
		},
	}, nil
}

func (s *Service) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	users, err := s.usecase.GetAllUsers(ctx)
	if err != nil {
		s.logger.Errorf("failed to get all users: %v", err)
		return nil, fmt.Errorf("faield to get all users: %w", err)
	}

	var pbUsers []*pb.User
	for _, user := range users {
		pbUsers = append(pbUsers, &pb.User{
			Id:             user.ID.String(),
			Email:          user.Email,
			HashedPassword: user.Password,
			IsConfirmed:    user.IsConfirmed,
		})
	}

	return &pb.GetAllUsersResponse{
		Users: pbUsers,
	}, nil
}

func (s *Service) DeleteUserById(ctx context.Context, req *pb.DeleteUserByIdRequest) (*pb.DeleteUserByIdResponse, error) {
	uuid, err := uuid.Parse(req.Id)
	if err != nil {
		s.logger.Errorf("failed to parse uuid from request err: %v", err)
		return nil, fmt.Errorf("failed to parse uuid from request err: %w", err)
	}

	err = s.usecase.DeleteUserById(ctx, uuid)
	if err != nil {
		s.logger.Errorf("failed to delete user by id err: %v", err)
		return nil, fmt.Errorf("failed to delete user by id err: %w", err)
	}

	return &pb.DeleteUserByIdResponse{}, nil
}

func (s *Service) ConfirmUser(ctx context.Context, req *pb.ConfirmUserRequest) (*pb.ConfirmUserResponse, error) {
	err := s.usecase.ConfirmUser(ctx, req.Email)
	if err != nil {
		s.logger.Errorf("failed to confirm user with email %v error %w", req.Email, err)
		return nil, fmt.Errorf("failed to confirm user with email %v error %w", req.Email, err)
	}

	return &pb.ConfirmUserResponse{}, nil
}
