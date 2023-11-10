package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtUserToken struct {
	Token        string
	RefreshToken string
}

type UserClaims struct {
	Id uuid.UUID `json:"id"`
	jwt.RegisteredClaims
}
