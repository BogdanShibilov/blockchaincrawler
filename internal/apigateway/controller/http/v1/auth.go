package v1

import (
	"net/http"

	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/apigateway"
	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/controller/http/v1/dto"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthRoutes struct {
	api apigateway.UseCase
	l   *zap.SugaredLogger
}

func NewAuthRoutes(handler *gin.RouterGroup, api apigateway.UseCase, l *zap.SugaredLogger) {
	r := &AuthRoutes{
		api: api,
		l:   l,
	}

	authHandler := handler.Group("/auth")
	{
		authHandler.POST("/signin", r.GenerateJwtToken)
		authHandler.POST("/refreshjwt", r.RenewJwtToken)
		authHandler.POST("/signup", r.CreateUser)
		authHandler.POST("/getconfirmation", r.SendConfirmationCode)
		authHandler.POST("/confirmuser", r.ConfirmUser)
	}
}

func (r *AuthRoutes) GenerateJwtToken(ctx *gin.Context) {
	var userCreds *dto.UserCreds
	err := ctx.ShouldBindJSON(&userCreds)
	if err != nil {
		r.l.Errorf("failed to bind json request to creds: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "bad json was given")
		return
	}

	jwtToken, err := r.api.GenerateJwtToken(ctx, userCreds)
	if err != nil {
		r.l.Errorf("failed to generate jwt token: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "failed to sign in")
		return
	}

	ctx.JSON(http.StatusOK, jwtToken)
}

func (r *AuthRoutes) RenewJwtToken(ctx *gin.Context) {
	var renewReq *dto.RenewTokenRequest
	err := ctx.ShouldBindJSON(&renewReq)
	if err != nil {
		r.l.Errorf("failed to bind json request to dto: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "bad json was given")
		return
	}

	jwtToken, err := r.api.RenewJwtToken(ctx, renewReq)
	if err != nil {
		r.l.Errorf("failed to renew jwt token: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "failed to refresh token")
		return
	}

	ctx.JSON(http.StatusOK, jwtToken)
}

func (r *AuthRoutes) CreateUser(ctx *gin.Context) {
	var userCreds *dto.UserCreds
	err := ctx.ShouldBindJSON(&userCreds)
	if err != nil {
		r.l.Errorf("failed to bind json request to creds: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "bad json was given")
		return
	}

	_, err = r.api.CreateUser(ctx, userCreds)
	if err != nil {
		r.l.Errorf("failed to create user: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "failed to sign up")
		return
	}

	ctx.JSON(http.StatusOK, &gin.H{"message": "Successfully signed up"})
}

func (r *AuthRoutes) SendConfirmationCode(ctx *gin.Context) {
	var sendConfReq *dto.SendConfirmCodeRequest
	err := ctx.ShouldBindJSON(&sendConfReq)
	if err != nil {
		r.l.Errorf("failed to bind json request to dto: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "bad json was given")
		return
	}

	err = r.api.SendConfirmationCode(ctx, sendConfReq)
	if err != nil {
		r.l.Errorf("failed to send confirmation code: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "failed to send confirmation code")
		return
	}

	ctx.JSON(http.StatusOK, &gin.H{"message": "Successfully sent code"})
}

func (r *AuthRoutes) ConfirmUser(ctx *gin.Context) {
	var confReq *dto.ConfirmUserRequest
	err := ctx.ShouldBindJSON(&confReq)
	if err != nil {
		r.l.Errorf("failed to bind json request to dto: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "bad json was given")
		return
	}

	err = r.api.ConfirmUser(ctx, confReq)
	if err != nil {
		r.l.Errorf("failed to confirm user: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "failed to confirm user")
		return
	}

	ctx.JSON(http.StatusOK, &gin.H{"message": "Succesfully confirmed user"})
}
