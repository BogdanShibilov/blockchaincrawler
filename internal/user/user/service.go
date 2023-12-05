package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/bogdanshibilov/blockchaincrawler/internal/user/entity"
	"github.com/bogdanshibilov/blockchaincrawler/internal/user/repository"
	pb "github.com/bogdanshibilov/blockchaincrawler/pkg/protobuf/user/gw"
	"github.com/bogdanshibilov/blockchaincrawler/pkg/redis"
)

type Service struct {
	repo  repository.UserRepository
	cache *redis.Redis
}

func New(repo repository.UserRepository, cache *redis.Redis) UseCase {
	return &Service{
		repo:  repo,
		cache: cache,
	}
}

func (s *Service) CreateUser(ctx context.Context, user *entity.User) (uuid.UUID, error) {
	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("CreateUser() error: %w", err)
	}

	return id, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := s.getUserByEmailCached(ctx, email)
	if user == nil {
		user, err = s.repo.GetUserByEmail(ctx, email)
		if err != nil {
			return nil, fmt.Errorf("GetUserByEmail() error: %w", err)
		}

		err = s.setUserByEmailCache(ctx, email, user)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetUserById(ctx context.Context, id string) (*entity.User, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse uuid: %w", err)
	}

	user, err := s.repo.GetUserById(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from repo: %w", err)
	}

	return user, nil
}

func (s *Service) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	users, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetAllUsers() error: %w", err)
	}

	return users, nil
}

func (s *Service) DeleteUserById(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteUserById(ctx, id)
	if err != nil {
		return fmt.Errorf("DeleteUserById() error: %w", err)
	}

	return nil
}

func (s *Service) ConfirmUser(ctx context.Context, email string) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("GetUserByEmail() error: %w", err)
	}

	user.IsConfirmed = true
	err = s.repo.UpdateUserById(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update user, confirmed to true: %w", err)
	}

	return nil
}

func (s *Service) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) error {
	uuid, err := uuid.Parse(req.UserId)
	if err != nil {
		return err
	}
	profile := &entity.Profile{
		UserId:  uuid,
		Name:    req.Profile.Name,
		Surname: req.Profile.Surname,
		AboutMe: req.Profile.AboutMe,
	}

	err = s.repo.CreateOrUpdateProfile(ctx, profile)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetProfileById(ctx context.Context, id string) (*entity.Profile, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	profile, err := s.repo.GetProfileById(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
