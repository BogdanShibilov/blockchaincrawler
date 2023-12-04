package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/apigateway"
	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/config"
	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/controller/http/middleware"
)

type AdminRoutes struct {
	api apigateway.UseCase
	l   *zap.SugaredLogger
	cfg *config.Config
}

func NewAdminRoutes(handler *gin.RouterGroup, api apigateway.UseCase, l *zap.SugaredLogger, cfg *config.Config) {
	r := &AdminRoutes{
		api: api,
		l:   l,
		cfg: cfg,
	}

	adminHandler := handler.Group("/admin")
	{
		adminHandler.Use(middleware.JwtVerify(&cfg.Jwt))
		adminHandler.Use(middleware.AdminOnly())
		adminHandler.GET("/user", r.GetAllUsers)
		adminHandler.DELETE("/user/:id", r.DeleteUserById)
	}
}

// GetAllUsers godoc
// @Summary Gets all users
// @Description Returns your all users' data. Admin only
// @Tags admin
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200
// @Failure 401
// @Failure 403
// @Failure 500
// @Router /admin/user [get]
func (r *AdminRoutes) GetAllUsers(ctx *gin.Context) {
	users, err := r.api.GetAllUsers(ctx)
	if err != nil {
		r.l.Errorf("failed to get all users: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to get all users"})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// DeleteUserById godoc
// @Summary Delete user by id
// @Description Deletes user with given id
// @Tags admin
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 204
// @Failure 401
// @Failure 403
// @Failure 500
// @Router /admin/user/{id} [delete]
func (r *AdminRoutes) DeleteUserById(ctx *gin.Context) {
	id := ctx.Param("id")
	err := r.api.DeleteUserById(ctx, id)
	if err != nil {
		r.l.Errorf("failed to get delete user: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to get delete user"})
		return
	}

	ctx.Status(http.StatusNoContent)
}
