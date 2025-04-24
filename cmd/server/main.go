package main

import (
	"github.com/gin-gonic/gin"
	"user_service/pkg/logging"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	logging.InitLogger()
	r := gin.New()
	r.Use(logging.Middleware)

	logging.Instance.Info("Starting application or port :8080")
	err := r.Run(":8080")
	if err != nil {
		logging.Instance.Fatal(err)
		return
	}
}
