package http

import (
	"github.com/gin-gonic/gin"

	"go-hexagonal/application"
	"go-hexagonal/config"
)

// NewServerRoute creates and configures the HTTP server routes
func NewServerRoute(useCaseFactory *application.UseCaseFactory) *gin.Engine {
	if config.GlobalConfig.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Legacy API (v1)
	example := router.Group("/example")
	{
		example.POST("", CreateExample)
		example.DELETE("/:id", DeleteExample)
		example.PUT("/:id", UpdateExample)
		example.GET("/:id", GetExample)
	}

	// New API using application layer use cases (v2)
	if useCaseFactory != nil {
		exampleHandler := NewExampleHandlerV2(useCaseFactory)
		v2 := router.Group("/v2")
		{
			exampleV2 := v2.Group("/example")
			{
				exampleV2.POST("", exampleHandler.Create)
				exampleV2.GET("/:id", exampleHandler.Get)
			}
		}
	}

	return router
}
