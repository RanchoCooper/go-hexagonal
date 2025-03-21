package http

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"go-hexagonal/api/http/middleware"
	"go-hexagonal/api/http/validator/custom"
	"go-hexagonal/config"
	"go-hexagonal/domain/service"
)

// Service instances for API handlers
var (
	services *service.Services
)

// RegisterServices registers service instances for API handlers
func RegisterServices(s *service.Services) {
	services = s
}

// NewServerRoute creates and configures the HTTP server routes
func NewServerRoute() *gin.Engine {
	if config.GlobalConfig.Env.IsProd() {
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
	router.Use(middleware.ErrorHandlerMiddleware()) // Add unified error handling middleware

	// Health check
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Debug tools
	if config.GlobalConfig.HTTPServer.Pprof {
		middleware.RegisterPprof(router)
	}

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
