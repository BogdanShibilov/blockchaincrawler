package repository

import (
	"context"
	"fmt"

	"github.com/bogdanshibilov/blockchaincrawler/internal/user/database/postgres"
	"github.com/bogdanshibilov/blockchaincrawler/internal/user/entity"
	"github.com/google/uuid"
)

type User struct {
	main *postgres.Pg
}

func New(db *postgres.Pg) UserRepository {
	return &User{db}
}
func (u *User) CreateUser(ctx context.Context, user *entity.User) (uuid.UUID, error) {
	res := u.main.DB.WithContext(ctx).Create(user)
	if res.Error != nil {
		return uuid.Nil, fmt.Errorf("failed to create user: %w", res.Error)
	}

	return user.ID, nil
}

func (u *User) GetUserById(ctx context.Context, id uuid.UUID) (user *entity.User, err error) {
	res := u.main.DB.WithContext(ctx).First(&user, id)
	if res.Error != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", res.Error)
	}

	return user, nil
}

func (u *User) GetUserByEmail(ctx context.Context, email string) (user *entity.User, err error) {
	res := u.main.DB.Where("email = ?", email).WithContext(ctx).First(&user)
	if res.Error != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", res.Error)
	}

	return user, nil
}

func (u *User) GetAllUsers(ctx context.Context) (users []*entity.User, err error) {
	res := u.main.DB.WithContext(ctx).Find(&users)
	if res.Error != nil {
		return nil, fmt.Errorf("failed to get all users: %w", res.Error)
	}

	return users, nil
}

func (u *User) UpdateUserById(ctx context.Context, newUser *entity.User) error {
	res := u.main.DB.WithContext(ctx).Model(newUser).Updates(newUser)
	if res.Error != nil {
		return fmt.Errorf("failed to update user: %w", res.Error)
	}

	return nil
}

func (u *User) DeleteUserById(ctx context.Context, id uuid.UUID) error {
	res := u.main.WithContext(ctx).Delete(&entity.User{}, id)
	if res.Error != nil {
		return fmt.Errorf("failed to delete user by id: %w", res.Error)
	}

	return nil
}
