package auth

import "github.com/golang-jwt/jwt/v5"

type JwtToken struct {
	AccessToken  string
	RefreshToken string
}

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	UserId      string `json:"userId"`
	UserEmail   string `json:"userEmail"`
	IsConfirmed bool   `json:"isConfirmed"`
}

type RefreshTokenClaims struct {
	jwt.RegisteredClaims
	UserId string `json:"userId"`
}
