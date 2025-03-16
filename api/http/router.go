package http

import (
	"github.com/gin-gonic/gin"

	"go-hexagonal/api/http/middleware"
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

	router := gin.New()

	// Apply middleware
	router.Use(gin.Recovery())
	router.Use(middleware.RequestID()) // Add request ID middleware
	router.Use(middleware.Cors())
	router.Use(middleware.RequestLogger()) // Add request logging middleware
	router.Use(middleware.Translations())

	// Health check
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// Debug tools
	if config.GlobalConfig.HTTPServer.Pprof {
		middleware.RegisterPprof(router)
	}

	// Unified API version
	api := router.Group("/api")
	{
		// Example API
		exampleHandler := NewExampleHandlerV2(useCaseFactory)
		examples := api.Group("/examples")
		{
			examples.POST("", exampleHandler.Create)
			examples.GET("/:id", exampleHandler.Get)
			// Add more endpoints as needed
		}
	}

	return router
}
