package routes

import (
	"github.com/gin-gonic/gin"
	"user_service/internal/bootstrap"
)

func SetupRoutes(r *gin.Engine, bs *bootstrap.Container) {
	SetupUserRoutes(r, bs)
}
