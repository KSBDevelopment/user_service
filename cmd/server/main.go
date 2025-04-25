package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"user_service/internal/config"
	"user_service/pkg/logging"
)

func main() {
	logging.InitLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		logging.Instance.Error(err)
		return
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(logging.Middleware)

	logging.Instance.Info("Starting application or port :" + cfg.Port)
	err = r.Run(":" + cfg.Port)
	if err != nil {
		logging.Instance.Fatal(err)
		return
	}
}
