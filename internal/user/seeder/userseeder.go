package seeder

import (
	"fmt"

	"github.com/bogdanshibilov/blockchaincrawler/internal/user/database/postgres"
	"github.com/bogdanshibilov/blockchaincrawler/internal/user/entity"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserSeeder struct {
	pg *postgres.Pg
}

func NewUserSeeder(pg *postgres.Pg) *UserSeeder {
	return &UserSeeder{
		pg: pg,
	}
}

func (s *UserSeeder) Seed() {
	var u []*entity.User
	_ = s.pg.DB.Find(&u).Error
	if len(u) == 0 {
		data := getUsersData()
		for _, u := range data {
			s.createUser(u)
		}
	}
}

func (s *UserSeeder) createUser(user *entity.User) (uuid.UUID, error) {
	err := s.pg.DB.Transaction(func(tx *gorm.DB) error {
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

func getUsersData() []*entity.User {
	data := []*entity.User{
		{
			Role:        "admin",
			Email:       "admin@gmail.com",
			Password:    generateHash("admin"),
			IsConfirmed: true,
		},
		{
			Role:        "John",
			Email:       "John@gmail.com",
			Password:    generateHash("123456abc"),
			IsConfirmed: true,
		},
		{
			Role:        "Eve",
			Email:       "Eve@gmail.com",
			Password:    generateHash("123456abc"),
			IsConfirmed: false,
		},
		{
			Role:        "Smithy",
			Email:       "Smithy@gmail.com",
			Password:    generateHash("123456abc"),
			IsConfirmed: true,
		},
		{
			Role:        "Alice",
			Email:       "Alice@gmail.com",
			Password:    generateHash("123456abc"),
			IsConfirmed: true,
		},
		{
			Role:        "Carl",
			Email:       "Carl@gmail.com",
			Password:    generateHash("123456abc"),
			IsConfirmed: false,
		},
	}

	return data
}

func generateHash(password string) string {
	pass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(pass)
}
