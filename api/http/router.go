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

	// Check if use case factory is provided
	if useCaseFactory == nil {
		// If no use case factory is provided, use traditional API (not recommended)
		example := router.Group("/example")
		{
			example.POST("", CreateExample)
			example.DELETE("/:id", DeleteExample)
			example.PUT("/:id", UpdateExample)
			example.GET("/:id", GetExample)
		}
		return router
	}

	// Use API with application layer use cases (recommended)
	exampleHandler := NewExampleHandlerV2(useCaseFactory)

	// API routes
	api := router.Group("/api")
	{
		// Example resource
		exampleApi := api.Group("/example")
		{
			exampleApi.POST("", exampleHandler.Create)
			exampleApi.GET("/:id", exampleHandler.Get)
			// Add more endpoints here
		}
	}

	// Keep v2 path for backward compatibility
	v2 := router.Group("/v2")
	{
		exampleV2 := v2.Group("/example")
		{
			exampleV2.POST("", exampleHandler.Create)
			exampleV2.GET("/:id", exampleHandler.Get)
		}
	}

	return router
}
