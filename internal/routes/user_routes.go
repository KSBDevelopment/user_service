package routes

import (
	"github.com/gin-gonic/gin"
	"user_service/internal/bootstrap"
	"user_service/internal/delivery"
	"user_service/internal/middleware"
	"user_service/internal/repository"
	"user_service/internal/service"
	"user_service/pkg/logging"
)

func SetupUserRoutes(router *gin.Engine, bs *bootstrap.Container) {

	ri, err := bs.GetRepository("user")
	if err != nil {
		logging.Instance.Error(err)
	}

	r, ok := ri.(*repository.UserRepositoryImpl)
	if !ok {
		logging.Instance.Error(err)
	}

	s := service.NewUserService(r)
	h := delivery.NewUserHandler(s)

	userRoutes := router.Group("/api/v1/user")

	publicRoutes := userRoutes.Group("/")
	{
		publicRoutes.POST("/", h.CreateUser)
		publicRoutes.GET("/", h.GetUsersPaginated)
		publicRoutes.GET("/:id", h.GetUserByID)
	}

	privateRoutes := userRoutes.Group("/")
	privateRoutes.Use(middleware.AuthMiddleware(bs.Config.JwtSecret))
	{
		privateRoutes.PUT("/:id", h.UpdateUser)
		privateRoutes.DELETE("/:id", h.DeleteUser)
	}
}
