package auth

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	pb "github.com/bogdanshibilov/blockchaincrawler/pkg/protobuf/user/gw"
)

func (s *Service) CreateUser(ctx context.Context, email string, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", fmt.Errorf("failed to generate hash from password: %w", err)
	}
	res, err := s.userTransport.CreateUser(ctx,
		&pb.CreateUserRequest{
			Email:          email,
			HashedPassword: string(hashedPassword),
		})
	if err != nil {
		return "", fmt.Errorf("failed to send proto create user req: %w", err)
	}

	return res.Id, nil
}
