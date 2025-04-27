package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"user_service/internal/bootstrap"
	"user_service/internal/routes"
	"user_service/pkg/logging"
)

func main() {
	logging.InitLogger()

	bs, err := bootstrap.Init()
	if err != nil {
		logging.Instance.Error(err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(logging.Middleware)

	routes.SetupRoutes(r, bs)
	logging.Instance.Info("Starting application or port :" + bs.Config.Port)
	err = r.Run(":" + bs.Config.Port)
	if err != nil {
		logging.Instance.Fatal(err)
		return
	}
}
