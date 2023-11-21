package auth

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/bogdanshibilov/blockchaincrawler/internal/auth/codeproducer"
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

func (s *Service) SendConfirmationCode(ctx context.Context, email string) error {
	code := generateRandomCode()

	s.codeDb.Client.Set(ctx, email, code, time.Minute*5)

	err := s.codeProducer.ProduceCode(&codeproducer.ConfirmationCode{
		Email: email,
		Code:  code,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ConfirmUser(ctx context.Context, email string, code string) error {
	redisCmd := s.codeDb.Get(ctx, email)
	if redisCmd.Err() != nil {
		return fmt.Errorf("failed to get code from redis db: %w", redisCmd.Err())
	}
	storedCode, err := redisCmd.Result()
	if err != nil {
		return fmt.Errorf("failed to get code from redis db: %w", redisCmd.Err())
	}

	if strings.Compare(code, storedCode) != 0 {
		return errors.New("invalid confirmation code")
	}

	err = s.userTransport.ConfirmUser(ctx, &pb.ConfirmUserRequest{
		Email: email,
	})
	if err != nil {
		return fmt.Errorf("failed to confirm user: %w", redisCmd.Err())
	}

	return nil
}

func generateRandomCode() string {
	code := rand.Intn(8999) + 1000
	return strconv.Itoa(code)
}
