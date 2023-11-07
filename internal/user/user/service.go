package user

import (
	"blockchaincrawler/internal/user/entity"
	"blockchaincrawler/internal/user/repository"
	"context"
	"fmt"
)

type Service struct {
	repo repository.UserRepository
}

func New(repo repository.UserRepository) UseCase {
	return &Service{repo}
}

func (s *Service) CreateUser(ctx context.Context, user *entity.User) error {
	err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("CreateUser() error: %w", err)
	}

	return nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("GetUserByEmail() error: %w", err)
	}

	return user, nil
}

func (s *Service) DeleteUserById(ctx context.Context, id int) error {
	err := s.repo.DeleteUserById(ctx, id)
	if err != nil {
		return fmt.Errorf("DeleteUserById() error: %w", err)
	}

	return nil
}
