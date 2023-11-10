package auth

import (
	"blockchaincrawler/internal/auth/transport"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Service struct {
	secretString      string
	userGrpcTransport *transport.UserGrpcTransport
}

func New(secret string, userTransport *transport.UserGrpcTransport) UseCase {
	return &Service{
		secretString:      secret,
		userGrpcTransport: userTransport,
	}
}
func (a *Service) GenerateJwtToken(ctx context.Context, email string, password string) (*JwtUserToken, error) {
	isValid, err := a.userGrpcTransport.IsValidLogin(ctx, email, password)
	if err != nil {
		return nil, fmt.Errorf("failed to check validness of email and password: %w", err)
	}
	if !isValid {
		return nil, errors.New("invalid credentials")
	}

	user, _ := a.userGrpcTransport.GetUserByEmail(ctx, email)
	userId, err := uuid.Parse(user.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse uuid: %w", err)
	}

	tokenString, _ := a.createNewTokenString(userId, jwt.SigningMethodHS256, time.Minute*5)
	refreshTokenString, err := a.createNewTokenString(userId, jwt.SigningMethodHS256, time.Minute*10)
	if err != nil {
		return nil, fmt.Errorf("failed to create new jwt token string: %w", err)
	}

	return &JwtUserToken{
		Token:        tokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (a *Service) RenewJwtToken(ctx context.Context, refreshToken string) (*JwtUserToken, error) {
	claims, err := a.parseToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to parse refresh token: %w", err)
	}
	userId, err := uuid.Parse(claims["id"].(string))
	if err != nil {
		return nil, fmt.Errorf("failed to parse uuid from jwt claims: %w", err)
	}

	tokenString, _ := a.createNewTokenString(userId, jwt.SigningMethodHS256, time.Minute*5)
	refreshTokenString, err := a.createNewTokenString(userId, jwt.SigningMethodHS256, time.Minute*10)
	if err != nil {
		return nil, fmt.Errorf("failed to create new jwt token string: %w", err)
	}

	return &JwtUserToken{
		Token:        tokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (a *Service) parseToken(tokenString string) (jwt.MapClaims, error) {
	secretKey := []byte(a.secretString)
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secretKey, nil
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

func (a *Service) createNewTokenString(
	userId uuid.UUID,
	signingMethod jwt.SigningMethod,
	duration time.Duration,
) (string, error) {

	claims := &UserClaims{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(duration)},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		},
	}

	secret := []byte(a.secretString)
	claimToken := jwt.NewWithClaims(signingMethod, claims)
	tokenString, err := claimToken.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign claimToken with secret: %w", err)
	}

	return tokenString, nil
}
