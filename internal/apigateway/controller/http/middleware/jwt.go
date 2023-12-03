package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/config"
)

func JwtVerify(cfg *config.Jwt) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var tokenString string
		tokenHeader := ctx.Request.Header.Get("Authorization")
		tokenFields := strings.Fields(tokenHeader)
		if len(tokenFields) == 2 && tokenFields[0] == "Bearer" {
			tokenString = tokenFields[1]
		} else {
			ctx.AbortWithStatus(http.StatusForbidden)

			return
		}

		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(cfg.Secret), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusForbidden)

			return
		}

		if !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		ctx.Set("userId", claims["userId"])
		ctx.Set("userEmail", claims["userEmail"])
		ctx.Set("role", claims["role"])
		ctx.Set("isConfirmed", claims["isConfirmed"])

		ctx.Next()
	}
}

// Always must be after verify VerifyJwt
func AdminOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role := ctx.Keys["role"]
		if role != "admin" {
			ctx.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		ctx.Next()
	}
}
