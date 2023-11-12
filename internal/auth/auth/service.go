package auth

import (
	"blockchaincrawler/internal/auth/transport"
	"blockchaincrawler/internal/kafka"
	pb "blockchaincrawler/pkg/protobuf/userservice/gw"
	"blockchaincrawler/pkg/redis"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Service struct {
	secretString      string
	userGrpcTransport *transport.UserGrpcTransport
	codeProducer      *kafka.Producer
	redisdb           *redis.Redis
}

func New(secret string, userTransport *transport.UserGrpcTransport, codeProducer *kafka.Producer, redisdb *redis.Redis) UseCase {
	return &Service{
		secretString:      secret,
		userGrpcTransport: userTransport,
		redisdb:           redisdb,
		codeProducer:      codeProducer,
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

func (a *Service) CreateUser(ctx context.Context, email string, password string) (*uuid.UUID, error) {
	userId, err := a.userGrpcTransport.CreateUser(ctx, &pb.User{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return &uuid.Nil, fmt.Errorf("failed to create new user: %w", err)
	}

	confirmCode := generateRandomCode()
	a.redisdb.Client.Set(ctx, email, confirmCode, time.Minute*5)

	byteConfirmCode, err := json.Marshal(confirmCode)
	if err != nil {
		return &uuid.Nil, fmt.Errorf("failed to marshal confirm code: %w", err)
	}

	a.codeProducer.ProduceMessage(byteConfirmCode)
	log.Println("Produced message: " + string(byteConfirmCode))

	return userId, nil
}

func (a *Service) ConfirmUser(ctx context.Context, email string, code string) error {
	redisCmd := a.redisdb.Get(ctx, email)
	if redisCmd.Err() != nil {
		return fmt.Errorf("failed to get code from redis db: %w", redisCmd.Err())
	}
	storedCode, _ := redisCmd.Result()
	if strings.Compare(storedCode, code) != 0 {
		return errors.New("invalid confirmation code")
	}

	isConfirmed, err := a.userGrpcTransport.ConfirmUser(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to confirm user on userservice: %w", err)
	}
	if !isConfirmed {
		return errors.New("user didnot get confirmed")
	}

	return nil
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

func generateRandomCode() string {
	return fmt.Sprintf("%d%d%d%d", rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10))
}
