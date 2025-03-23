// Package dependency provides dependency injection configuration
//go:build wireinject
// +build wireinject

package dependency

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/gorm"

	"go-hexagonal/adapter/converter"
	"go-hexagonal/adapter/repository"
	"go-hexagonal/adapter/repository/mysql/entity"
	"go-hexagonal/config"
	"go-hexagonal/domain/event"
	"go-hexagonal/domain/repo"
	"go-hexagonal/domain/service"
)

// RepositoryOption defines an option for repository initialization
type RepositoryOption func(*repository.ClientContainer)

// WithMySQL returns an option to initialize MySQL
func WithMySQL() RepositoryOption {
	return func(c *repository.ClientContainer) {
		if c.MySQL == nil {
			mysql, err := ProvideMySQL()
			if err != nil {
				panic("Failed to initialize MySQL: " + err.Error())
			}
			c.MySQL = mysql
		}
	}
}

// WithRedis returns an option to initialize Redis
func WithRedis() RepositoryOption {
	return func(c *repository.ClientContainer) {
		if c.Redis == nil {
			redis, err := ProvideRedis()
			if err != nil {
				panic("Failed to initialize Redis: " + err.Error())
			}
			c.Redis = redis
		}
	}
}

// ServiceOption defines an option for service initialization
type ServiceOption func(*service.Services, event.EventBus)

// WithExampleService returns an option to initialize the Example service
func WithExampleService() ServiceOption {
	return func(s *service.Services, eventBus event.EventBus) {
		if s.ExampleService == nil {
			exampleRepo := entity.NewExample()
			s.ExampleService = provideExampleService(exampleRepo, eventBus)
		}
	}
}

// WithTransactionFactory returns an option to initialize the transaction factory
func WithTransactionFactory() RepositoryOption {
	return func(c *repository.ClientContainer) {
		// This is a no-op since the transaction factory doesn't need initialization,
		// but we include it for consistency and future extensions
	}
}

// WithExampleConverter returns an option to initialize the example converter
func WithExampleConverter() ServiceOption {
	return func(s *service.Services, _ event.EventBus) {
		if s.Converter == nil {
			s.Converter = provideExampleConverter()
		}
	}
}

// InitializeServices initializes services based on the provided options
func InitializeServices(ctx context.Context, opts ...ServiceOption) (*service.Services, error) {
	// Initialize services container
	eventBus := provideEventBus()
	services := &service.Services{
		EventBus: eventBus,
	}

	// Apply service options
	for _, opt := range opts {
		opt(services, eventBus)
	}

	return services, nil
}

// InitializeRepositories initializes repository clients with the given options
func InitializeRepositories(opts ...RepositoryOption) (*repository.ClientContainer, error) {
	container := &repository.ClientContainer{}
	for _, opt := range opts {
		opt(container)
	}
	return container, nil
}

// ProvideMySQL creates and initializes a MySQL client
func ProvideMySQL() (*repository.MySQL, error) {
	if config.GlobalConfig.MySQL == nil {
		return nil, repository.ErrMissingMySQLConfig
	}

	db, err := repository.OpenGormDB()
	if err != nil {
		return nil, err
	}

	return &repository.MySQL{DB: db}, nil
}

// ProvideRedis creates and initializes a Redis client
func ProvideRedis() (*repository.Redis, error) {
	if config.GlobalConfig.Redis == nil {
		return nil, repository.ErrMissingRedisConfig
	}

	client := repository.NewRedisConn()
	return &repository.Redis{DB: client}, nil
}

// ProvideTransactionFactory creates and initializes a transaction factory
func ProvideTransactionFactory() repo.TransactionFactory {
	return repository.NewTransactionFactory()
}

// ProvideExampleConverter creates and initializes an example converter
func ProvideExampleConverter() service.Converter {
	return converter.NewExampleConverter()
}

// ProvideRepositoryClients creates a clients struct with all repositories
func ProvideRepositoryClients(mysql *repository.MySQL, redis *repository.Redis) *repository.ClientContainer {
	return &repository.ClientContainer{
		MySQL: mysql,
		Redis: redis,
	}
}

// MySQLSet provides a Wire provider set for MySQL
var MySQLSet = wire.NewSet(
	ProvideMySQL,
	wire.Bind(new(MySQLRepository), new(*repository.MySQL)),
)

// RedisSet provides a Wire provider set for Redis
var RedisSet = wire.NewSet(
	ProvideRedis,
	wire.Bind(new(RedisRepository), new(*repository.Redis)),
)

// MySQLRepository defines the interface for MySQL operations
type MySQLRepository interface {
	GetDB(ctx context.Context) *gorm.DB
	Close(ctx context.Context) error
}

// RedisRepository defines the interface for Redis operations
type RedisRepository interface {
	GetClient() *redis.Client
	Close(ctx context.Context) error
}

// provideEventBus creates and configures the event bus
func provideEventBus() *event.InMemoryEventBus {
	eventBus := event.NewInMemoryEventBus()

	// Register event handlers
	loggingHandler := event.NewLoggingEventHandler()
	exampleHandler := event.NewExampleEventHandler()
	eventBus.Subscribe(loggingHandler)
	eventBus.Subscribe(exampleHandler)

	return eventBus
}

// provideExampleService creates and configures the example service
func provideExampleService(repo repo.IExampleRepo, eventBus event.EventBus) *service.ExampleService {
	exampleService := service.NewExampleService(repo)
	exampleService.EventBus = eventBus
	return exampleService
}

// provideExampleConverter creates and returns an example converter
func provideExampleConverter() service.Converter {
	return converter.NewExampleConverter()
}

// Deprecated: Use the new InitializeServices with options pattern instead
func provideServices(exampleService *service.ExampleService, eventBus event.EventBus) *service.Services {
	return service.NewServices(exampleService, eventBus)
}
