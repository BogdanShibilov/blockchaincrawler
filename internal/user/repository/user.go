package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/bogdanshibilov/blockchaincrawler/internal/user/database/postgres"
	"github.com/bogdanshibilov/blockchaincrawler/internal/user/entity"
)

type User struct {
	main *postgres.Pg
}

func New(db *postgres.Pg) UserRepository {
	return &User{db}
}
func (u *User) CreateUser(ctx context.Context, user *entity.User) (uuid.UUID, error) {
	err := u.main.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		userProfile := &entity.Profile{
			UserId: user.ID,
		}
		if err := tx.Create(&userProfile).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
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
	err := u.main.WithContext(ctx).Delete(&entity.User{ID: id}, id).Error
	if err != nil {
		return fmt.Errorf("failed to delete user by id: %w", err)
	}

	return nil
}

func (u *User) CreateOrUpdateProfile(ctx context.Context, p *entity.Profile) error {
	res := u.main.DB.WithContext(ctx).Where("user_id = ?", p.UserId).Save(&p)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (u *User) GetProfileById(ctx context.Context, id uuid.UUID) (p *entity.Profile, err error) {
	res := u.main.DB.WithContext(ctx).Where("user_id = ?", id).First(&p)
	if res.Error != nil {
		return nil, res.Error
	}

	return p, nil
}
