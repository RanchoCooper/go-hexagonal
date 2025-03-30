package http

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	httpMiddleware "go-hexagonal/api/http/middleware"
	"go-hexagonal/api/http/validator/custom"
	metricsMiddleware "go-hexagonal/api/middleware"
	"go-hexagonal/application"
	"go-hexagonal/config"
	"go-hexagonal/domain/repo"
	"go-hexagonal/domain/service"
)

// Service instances for API handlers
var (
	services  *service.Services
	converter service.Converter
)

// RegisterServices registers service instances for API handlers
func RegisterServices(s *service.Services) {
	services = s
}

// RegisterConverter registers a converter instance for API handlers
// This is mainly used for testing
func RegisterConverter(c service.Converter) {
	converter = c
}

// InitAppFactory initializes the application factory and sets it for API handlers
func InitAppFactory(s *service.Services) {
	// Create transaction factory instance
	txFactory := repo.NewNoOpTransactionFactory()

	// Create application factory with necessary parameters
	factory := application.NewFactory(
		s.ExampleService,
		txFactory,
	)

	// Use the external SetAppFactory function defined in example.go
	SetAppFactory(factory)
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
	router.Use(httpMiddleware.RequestID()) // Add request ID middleware
	router.Use(httpMiddleware.Cors())
	router.Use(httpMiddleware.RequestLogger()) // Add request logging middleware
	router.Use(httpMiddleware.Translations())
	router.Use(httpMiddleware.ErrorHandlerMiddleware()) // Add unified error handling middleware

	// Add metrics middleware for each handler
	router.Use(func(c *gin.Context) {
		// Use the path as a label for the metrics
		handlerName := c.FullPath()
		if handlerName == "" {
			handlerName = "unknown"
		}

		// Record the start time
		start := time.Now()

		// Process the request
		c.Next()

		// Record metrics after request is processed
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		// Record request metrics
		metricsMiddleware.RecordHTTPMetrics(handlerName, c.Request.Method, statusCode, duration)
	})

	// Health check
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Debug tools
	if config.GlobalConfig.HTTPServer.Pprof {
		httpMiddleware.RegisterPprof(router)
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
