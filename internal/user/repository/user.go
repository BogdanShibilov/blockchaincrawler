package repository

import (
	"blockchaincrawler/internal/user/database/postgres"
	"blockchaincrawler/internal/user/entity"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type User struct {
	main *postgres.Pg
}

func New(db *postgres.Pg) *User {
	return &User{db}
}

func (ur *User) CreateUser(ctx context.Context, user *entity.User) (uuid.UUID, error) {
	res := ur.main.DB.WithContext(ctx).Create(user)
	if res.Error != nil {
		return uuid.Nil, fmt.Errorf("failed to create user: %w", res.Error)
	}

	return user.ID, nil
}

func (ur *User) GetUserById(ctx context.Context, id uuid.UUID) (user *entity.User, err error) {
	res := ur.main.DB.WithContext(ctx).First(&user, id)
	if res.Error != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", res.Error)
	}

	return user, nil
}

func (ur *User) GetUserByEmail(ctx context.Context, email string) (user *entity.User, err error) {
	res := ur.main.DB.Where("email = ?", email).WithContext(ctx).First(&user)
	if res.Error != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", res.Error)
	}

	return user, nil
}

func (ur *User) UpdateUserById(ctx context.Context, newUser *entity.User) error {
	res := ur.main.DB.WithContext(ctx).Model(newUser).Updates(newUser)
	if res.Error != nil {
		return fmt.Errorf("failed to update user: %w", res.Error)
	}

	return nil
}

func (ur *User) DeleteUserById(ctx context.Context, id uuid.UUID) error {
	res := ur.main.WithContext(ctx).Delete(&entity.User{}, id)
	if res.Error != nil {
		return fmt.Errorf("failed to delete user by id: %w", res.Error)
	}

	return nil
}
