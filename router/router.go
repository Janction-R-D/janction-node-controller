package router

import (
	"janction/controller"
	"janction/logger"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(logger.GinLogger(), logger.GinRecovery(true), cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	v1 := r.Group("/api/controller/v1")
	{
		v1.GET("/job", controller.GetJobHandler)
		v1.POST("/register", controller.RegisterNodeHandler)
		v1.POST("/ping", controller.PingHandler)
	}

	return r
}
