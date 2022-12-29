package http

import (
	"github.com/gin-gonic/gin"

	"go-hexagonal/config"
)

func NewServerRoute() *gin.Engine {
	if config.Config.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	example := router.Group("/example")
	{
		example.POST("", CreateExample)
		example.DELETE("/:id", DeleteExample)
		example.PUT("/:id", UpdateExample)
		example.GET("/:id", GetExample)
	}

	return router
}
