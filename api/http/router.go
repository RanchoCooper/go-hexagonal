package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"go-hexagonal/api/http/middleware"
	"go-hexagonal/api/http/validator/custom"
	"go-hexagonal/config"
)

// NewServerRoute creates and configures the HTTP server routes
func NewServerRoute() *gin.Engine {
	if config.GlobalConfig.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Register custom validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		custom.RegisterValidators(v)
	}

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
		examples := api.Group("/examples")
		{
			examples.POST("", CreateExample)
			examples.GET("/:id", GetExample)
			examples.PUT("/:id", UpdateExample)
			examples.DELETE("/:id", DeleteExample)
			examples.GET("/name/:name", FindExampleByName)
		}
	}

	return router
}
