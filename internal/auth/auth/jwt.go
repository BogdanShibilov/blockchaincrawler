package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	pb "github.com/bogdanshibilov/blockchaincrawler/pkg/protobuf/user/gw"
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
	claims, err := s.parseToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to parse refresh token: %w", err)
	}
	userId := claims["userId"].(string)

	res, err := s.userTransport.GetUserById(ctx, &pb.GetUserByIdRequest{Id: userId})
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	accessToken, err := s.generateAccessToken(res.User)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	refreshToken, err = s.generateRefreshToken(res.User)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &JwtToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) generateAccessToken(user *pb.User) (string, error) {
	timeNow := time.Now()
	tokenLiveDuration := time.Duration(s.cfg.AccessTokenDuration)
	claims := &AccessTokenClaims{
		UserId:      user.Id,
		Role:        user.Role,
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

func (s *Service) parseToken(tokenString string) (jwt.MapClaims, error) {
	secret := []byte(s.cfg.JwtSecretKey)
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token string: %w", err)
	}

	if token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			return claims, nil
		} else {
			return nil, fmt.Errorf("failed to extract claims: %w", err)
		}
	}
	return nil, errors.New("invalid token")
}
