package auth

import (
	"context"
	"fmt"
	"time"

	pb "github.com/bogdanshibilov/blockchaincrawler/pkg/protobuf/user/gw"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) GenerateJwtToken(ctx context.Context, email string, password string) (*JwtToken, error) {
	res, err := s.userTransport.GetUserByEmail(ctx, &pb.GetUserByEmailRequest{Email: email})
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email from user service: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(res.User.HashedPassword), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("failed to verify password: %w", err)
	}

	accessToken, err := s.generateAccessToken(res.User)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	refreshToken, err := s.generateRefreshToken(res.User)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &JwtToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) RenewJwtToken(ctx context.Context, refreshToken string) (*JwtToken, error) {
	panic("not implemented") // TODO: Implement
}

func (s *Service) generateAccessToken(user *pb.User) (string, error) {
	timeNow := time.Now()
	tokenLiveDuration := time.Duration(s.cfg.AccessTokenDuration)
	claims := &AccessTokenClaims{
		UserId:      user.Id,
		UserEmail:   user.Email,
		IsConfirmed: user.IsConfirmed,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  &jwt.NumericDate{Time: timeNow},
			ExpiresAt: &jwt.NumericDate{Time: timeNow.Add(time.Minute * tokenLiveDuration)},
		},
	}

	secret := []byte(s.cfg.JwtSecretKey)
	accessToken, err := jwt.
		NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return accessToken, nil
}

func (s *Service) generateRefreshToken(user *pb.User) (string, error) {
	timeNow := time.Now()
	tokenLiveDuration := time.Duration(s.cfg.RefreshTokenDuration)
	claims := &RefreshTokenClaims{
		UserId: user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  &jwt.NumericDate{Time: timeNow},
			ExpiresAt: &jwt.NumericDate{Time: timeNow.Add(time.Minute * tokenLiveDuration)},
		},
	}

	secret := []byte(s.cfg.JwtSecretKey)
	refreshToken, err := jwt.
		NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return refreshToken, nil
}
