package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/apigateway"
	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/config"
	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/controller/http/middleware"
	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/controller/http/v1/dto"
)

type AuthRoutes struct {
	api apigateway.UseCase
	l   *zap.SugaredLogger
	cfg *config.Config
}

func NewAuthRoutes(handler *gin.RouterGroup, api apigateway.UseCase, l *zap.SugaredLogger, cfg *config.Config) {
	r := &AuthRoutes{
		api: api,
		l:   l,
		cfg: cfg,
	}

	authHandler := handler.Group("/auth")
	{
		authHandler.POST("/signin", r.GenerateJwtToken)
		authHandler.POST("/signup", r.CreateUser)
		authHandler.POST("/refreshjwt", r.RenewJwtToken)
		authHandler.Use(middleware.JwtVerify(&r.cfg.Jwt))
		authHandler.POST("/getconfirmation", r.SendConfirmationCode)
		authHandler.POST("/confirmuser", r.ConfirmUser)
	}
}

// GenerateJwtToken godoc
// @Summary Generates jwt token
// @Description Returns jwt token which contains access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.UserCreds true "Email and password of account"
// @Success 200 {object} dto.JwtToken
// @Failure 400
// @Failure 500
// @Router /auth/signin [post]
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

// RenewJwtToken godoc
// @Summary Refreshes jwt token
// @Description Returns new jwt token which contains access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param refreshToken body dto.RenewTokenRequest true "Contains refresh token string"
// @Success 200 {object} dto.JwtToken
// @Failure 400
// @Failure 500
// @Router /auth/refreshjwt [post]
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

// CreateUser godoc
// @Summary Signs up
// @Description Create a new user in user database with given credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.UserCreds true "Email and password of account"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /auth/signup [post]
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

// SendConfirmationCode godoc
// @Summary Sends code
// @Description Sends code which can be used to confirm user account
// @Tags auth
// @Accept json
// @Produce json
// @Param email body dto.SendConfirmCodeRequest true "Email where code would be sent"
// @Param Authorization header string true "Bearer token"
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 500
// @Router /auth/getconfirmation [post]
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

	ctx.JSON(http.StatusOK, "Successfully sent code")
}

// ConfirmUser godoc
// @Summary Confrims user account
// @Description Tries to confirm user account with given code
// @Tags auth
// @Accept json
// @Produce json
// @Param code body dto.ConfirmUserRequest true "Code used to confirm user"
// @Param Authorization header string true "Bearer token"
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 403
// @Router /auth/confirmuser [post]
func (r *AuthRoutes) ConfirmUser(ctx *gin.Context) {
	var confReq *dto.ConfirmUserRequest
	err := ctx.ShouldBindJSON(&confReq)
	if err != nil {
		r.l.Errorf("failed to bind json request to dto: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "bad json was given")
		return
	}
	userEmail := ctx.Keys["userEmail"].(string)

	err = r.api.ConfirmUser(ctx, userEmail, confReq.Code)
	if err != nil {
		r.l.Errorf("failed to confirm user: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "failed to confirm user")
		return
	}

	ctx.JSON(http.StatusOK, &gin.H{"message": "Succesfully confirmed user"})
}
