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

type UserRoutes struct {
	api apigateway.UserUseCase
	l   *zap.SugaredLogger
	cfg *config.Config
}

func NewUserRoutes(handler *gin.RouterGroup, api apigateway.UserUseCase, l *zap.SugaredLogger, cfg *config.Config) {
	r := &UserRoutes{
		api: api,
		l:   l,
		cfg: cfg,
	}

	userHandler := handler.Group("/user")
	{
		userHandler.Use(middleware.JwtVerify(&cfg.Jwt))
		userHandler.GET("/profile", r.GetProfile)
		userHandler.PUT("/profile", r.UpdateProfile)
	}
}

// UpdateProfile godoc
// @Summary Updates profile
// @Description Updates your profile with given profile data
// @Tags user
// @Accept json
// @Produce json
// @Param profile body dto.UserProfileDto true "Your profile data"
// @Param Authorization header string true "Bearer token"
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 500
// @Router /user/profile [put]
func (u *UserRoutes) UpdateProfile(ctx *gin.Context) {
	var newProfile *dto.UserProfileDto
	err := ctx.ShouldBindJSON(&newProfile)
	if err != nil {
		u.l.Errorf("failed to bind json: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "bad json file was given"})
		return
	}
	id := ctx.Keys["userId"].(string)

	err = u.api.UpdateProfile(ctx, id, newProfile)
	if err != nil {
		u.l.Errorf("failed to update profile: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to update profile"})
		return
	}

	ctx.Status(http.StatusOK)
}

// GetProfile godoc
// @Summary Gets your profile
// @Description Returns your profile data
// @Tags user
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200
// @Failure 401
// @Failure 403
// @Failure 500
// @Router /user/profile [get]
func (u *UserRoutes) GetProfile(ctx *gin.Context) {
	id := ctx.Keys["userId"].(string)

	profile, err := u.api.GetProfileById(ctx, id)
	if err != nil {
		u.l.Errorf("failed to get profile: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to get profile"})
		return
	}

	ctx.JSON(http.StatusOK, profile)
}
