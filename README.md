# Go Hexagonal Architecture

Welcome to visit my [blog post](https://blog.ranchocooper.com/2025/03/20/go-hexagonal/)

![Hexagonal Architecture](https://github.com/Sairyss/domain-driven-hexagon/raw/master/assets/images/DomainDrivenHexagon.png)

## Project Overview

This project is a Go microservice framework based on [Hexagonal Architecture](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software)) and [Domain-Driven Design](https://en.wikipedia.org/wiki/Domain-driven_design). It provides a clear project structure and design patterns to help developers build maintainable, testable, and scalable applications.

Hexagonal Architecture (also known as [Ports and Adapters Architecture](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software))) divides the application into internal and external parts, implementing [Separation of Concerns](https://en.wikipedia.org/wiki/Separation_of_concerns) and [Dependency Inversion Principle](https://en.wikipedia.org/wiki/Dependency_inversion_principle) through well-defined interfaces (ports) and implementations (adapters). This architecture decouples business logic from technical implementation details, facilitating unit testing and feature extension.

## Core Features

### Architecture Design
- **[Domain-Driven Design (DDD)](https://en.wikipedia.org/wiki/Domain-driven_design)** - Organize business logic through concepts like [Aggregates](https://en.wikipedia.org/wiki/Domain-driven_design), [Entities](https://en.wikipedia.org/wiki/Entity), and [Value Objects](https://en.wikipedia.org/wiki/Value_object)
- **[Hexagonal Architecture](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software))** - Divide the application into domain, application, and adapter layers
- **[Dependency Injection](https://en.wikipedia.org/wiki/Dependency_injection)** - Use [Wire](https://github.com/google/wire) for dependency injection, improving code testability and flexibility
- **[Repository Pattern](https://en.wikipedia.org/wiki/Repository_pattern)** - Abstract data access layer with transaction support
- **[Domain Events](https://en.wikipedia.org/wiki/Domain-driven_design)** - Implement [Event-Driven Architecture](https://en.wikipedia.org/wiki/Event-driven_architecture), supporting loosely coupled communication between system components
- **[CQRS Pattern](https://en.wikipedia.org/wiki/Command_Query_Responsibility_Segregation)** - Command and Query Responsibility Segregation, optimizing read and write operations
- **[Interface-Driven Design](https://en.wikipedia.org/wiki/Interface-based_programming)** - Use interfaces to define service contracts, implementing Dependency Inversion Principle

### Technical Implementation
- **[RESTful API](https://en.wikipedia.org/wiki/Representational_state_transfer)** - Implement HTTP API using the [Gin](https://github.com/gin-gonic/gin) framework
- **Database Support** - Integrate [GORM](https://gorm.io) with support for [MySQL](https://en.wikipedia.org/wiki/MySQL), [PostgreSQL](https://en.wikipedia.org/wiki/PostgreSQL), and other databases
- **Cache Support** - Integrate [Redis](https://en.wikipedia.org/wiki/Redis) caching with comprehensive error handling, local error definitions for cache misses, and health check implementation for monitoring cache availability
- **Enhanced Cache** - Advanced cache features including negative caching to prevent cache penetration, distributed locking for cache consistency, and key tracking for improved hit rates
- **MongoDB Support** - Integration with MongoDB for document storage
- **Logging System** - Use [Zap](https://go.uber.org/zap) for high-performance logging with structured context support for tracing and debugging
- **Configuration Management** - Use [Viper](https://github.com/spf13/viper) for flexible configuration management
- **[Graceful Shutdown](https://en.wikipedia.org/wiki/Graceful_exit)** - Support graceful service startup and shutdown
- **[Unit Testing](https://en.wikipedia.org/wiki/Unit_testing)** - Use [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock), [redismock](https://github.com/go-redis/redismock), and [testify/mock](https://github.com/stretchr/testify) for comprehensive test coverage with enhanced HTTP testing utilities and improved DTO handling
- **Transaction Support** - Provide no-operation transaction implementation, simplifying service layer interaction with repository layer, complete with mock transaction implementation and lifecycle hooks (Begin, Commit, and Rollback) for testing
- **Asynchronous Event Processing** - Support for asynchronous event handling with worker pools, event persistence, and replay capabilities

### Development Toolchain
- **Code Quality** - Integrate [Golangci-lint](https://github.com/golangci/golangci-lint) for code quality checks
- **Commit Standards** - Use [Commitlint](https://github.com/conventional-changelog/commitlint) to ensure Git commit messages follow conventions
- **Pre-commit Hooks** - Use [Pre-commit](https://pre-commit.com) for code checking and formatting
- **[CI/CD](https://en.wikipedia.org/wiki/CI/CD)** - Integrate [GitHub Actions](https://github.com/features/actions) for continuous integration and deployment

## Recent Enhancements

### Unified Error Handling
- Extended error handling with consistent error types and error wrapping functions
- Support for structured error details and HTTP status code mapping
- Error comparison capabilities for more robust error checking

### Enhanced Structured Logging
- Context-aware logging with support for request IDs, user IDs, and trace IDs
- Consistent log formatting and level management
- Improved debugging capabilities with contextual information

### Asynchronous Event System
- Worker pool-based event processing for improved throughput
- Event persistence and replay capabilities for reliability
- Graceful shutdown support for event processing

### Advanced Caching Features
- Negative caching to protect against cache penetration
- Distributed locking to prevent cache stampede
- Key tracking for improved cache hit rates
- Cache consistency mechanisms for data integrity

## Project Structure

```
.
├── adapter/                # Adapter Layer - External system interactions
│   ├── amqp/               # Message queue adapters
│   ├── dependency/         # Dependency injection configuration
│   │   └── wire.go         # Wire DI setup with interface bindings
│   ├── job/                # Scheduled task adapters
│   └── repository/         # Data repository adapters
│       ├── mysql/          # MySQL implementation
│       │   └── entity/     # Database entities and repo implementations
│       ├── postgre/        # PostgreSQL implementation
│       ├── mongo/          # MongoDB implementation
│       └── redis/          # Redis implementation
│           └── enhanced_cache.go  # Enhanced cache with advanced features
├── api/                    # API Layer - HTTP requests and responses
│   ├── dto/                # Data Transfer Objects for API
│   ├── error_code/         # Error code definitions
│   ├── grpc/               # gRPC API handlers
│   └── http/               # HTTP API handlers
│       ├── handle/         # Request handlers using domain interfaces
│       ├── middleware/     # HTTP middleware
│       ├── paginate/       # Pagination handling
│       └── validator/      # Request validation
├── application/            # Application Layer - Use cases coordinating domain objects
│   ├── core/               # Core interfaces and base implementations
│   │   └── interfaces.go   # UseCase and UseCaseHandler interfaces
│   └── example/            # Example use case implementations
│       ├── create_example.go     # Create example use case
│       ├── delete_example.go     # Delete example use case
│       ├── get_example.go        # Get example use case
│       ├── update_example.go     # Update example use case
│       └── find_example_by_name.go # Find example by name use case
├── cmd/                    # Command-line entry points
│   └── http_server/        # HTTP server startup
├── config/                 # Configuration files and management
├── domain/                 # Domain Layer - Core business logic
│   ├── aggregate/          # Aggregates (DDD concept)
│   ├── event/              # Domain events and event bus interfaces
│   │   ├── event_bus.go    # EventBus interface
│   │   ├── async_event_bus.go # Asynchronous event bus implementation
│   │   └── handlers.go     # Event handler interfaces
│   ├── model/              # Domain models (pure business entities)
│   ├── repo/               # Repository interfaces
│   │   └── transaction.go  # Transaction interface
│   ├── service/            # Domain services with interfaces
│   │   ├── example.go      # ExampleService implementation
│   │   └── interfaces.go   # Service interfaces (IExampleService, etc.)
│   └── vo/                 # Value objects (DDD concept)
├── tests/                  # Integration tests
│   ├── migrations/         # Test database migrations
│   ├── mysql.go            # MySQL test helpers
│   ├── postgres.go         # PostgreSQL test helpers
│   ├── redis.go            # Redis test helpers
│   └── *_test.go           # Test files
└── util/                   # Utility functions
    ├── clean_arch/         # Architecture checking tools
    ├── errors/             # Enhanced error types and handling
    └── log/                # Enhanced logging with context support
```

### Key Architectural Elements

This structure enforces the Hexagonal Architecture principles:

1. **Interface-Implementation Separation**:
   - Domain layer defines interfaces (ports)
   - Adapter layer provides implementations (adapters)
   - Dependency flows inward, with outer layers depending on inner layers

2. **Dependency Inversion**:
   - High-level modules (domain/application) depend on abstractions
   - Low-level modules (adapters) implement these abstractions
   - All dependencies are injected through interfaces

3. **Domain-Centric Design**:
   - Domain models are pure business entities without technical concerns
   - Repository interfaces declare what the domain needs
   - Service interfaces define business operations

4. **Clean Boundaries**:
   - Each layer has clear responsibilities and dependencies
   - Data transformation occurs at layer boundaries
   - No leakage of implementation details between layers

## Architecture Layers

### Domain Layer
The domain layer is the core of the application, containing business logic and rules. It is independent of other layers and does not depend on any external components.

- **Models**: Domain entities and value objects
  - `Example`: Example entity, containing basic properties like ID, name, alias, etc.

- **Repository Interfaces**: Define data access interfaces
  - `IExampleRepo`: Example repository interface, defining operations like create, read, update, delete, etc.
  - `IExampleCacheRepo`: Example cache interface, defining health check methods
  - `Transaction`: Transaction interface, supporting transaction begin, commit, and rollback

- **Domain Services**: Handle business logic across entities
  - `IExampleService`: Service interface defining contracts for example-related operations
  - `ExampleService`: Implementation of the example service interface, handling business logic for example entities

- **Domain Events**: Define events within the domain
  - `ExampleCreatedEvent`: Example creation event
  - `ExampleUpdatedEvent`: Example update event
  - `ExampleDeletedEvent`: Example deletion event
  - `AsyncEventBus`: Asynchronous event processing with persistence

### Application Layer
The application layer coordinates domain objects to complete specific application tasks. It depends on domain interfaces but not on concrete implementations, following the Dependency Inversion Principle.

- **Use Cases**: Define application functionality
  - `CreateExampleUseCase`: Create example use case
  - `GetExampleUseCase`: Get example use case
  - `UpdateExampleUseCase`: Update example use case
  - `DeleteExampleUseCase`: Delete example use case
  - `FindExampleByNameUseCase`: Find example by name use case

- **Commands and Queries**: Implement CQRS pattern
  - Each use case defines Input and Output structures, representing command/query inputs and results

- **Event Handlers**: Process domain events
  - `LoggingEventHandler`: Logging event handler, recording all events
  - `ExampleEventHandler`: Example event handler, processing events related to examples

### Adapter Layer
The adapter layer implements interaction with external systems, such as databases and message queues.

- **Repository Implementation**: Implement data access interfaces
  - `EntityExample`: MySQL implementation of example repository
  - `NoopTransaction`: No-operation transaction implementation, simplifying testing
  - `MySQL`: MySQL connection and transaction management
  - `Redis`: Redis connection and basic operations
  - `EnhancedCache`: Advanced Redis caching with anti-penetration protection

- **Message Queue Adapters**: Implement message publishing and subscription
  - Support for Kafka and other message queue integrations

- **Scheduled Tasks**: Implement scheduled tasks
  - Cron-based task scheduling system

### API Layer
The API layer handles HTTP requests and responses, serving as the entry point to the application.

- **Controllers**: Handle HTTP requests
  - `CreateExample`: Create example API
  - `GetExample`: Get example API
  - `UpdateExample`: Update example API
  - `DeleteExample`: Delete example API
  - `FindExampleByName`: Find example by name API

- **Middleware**: Implement cross-cutting concerns
  - Internationalization support
  - CORS support
  - Request ID tracking
  - Request logging

- **Data Transfer Objects (DTOs)**: Define request and response data structures
  - `CreateExampleReq`: Create example request
  - `UpdateExampleReq`: Update example request
  - `DeleteExampleReq`: Delete example request
  - `GetExampleReq`: Get example request

## Dependency Injection

This project uses Google Wire for dependency injection, organizing dependencies as follows:

```go
// Initialize services
func InitializeServices(ctx context.Context) (*service.Services, error) {
    wire.Build(
        // Repository dependencies
        entity.NewExample,
        wire.Bind(new(repo.IExampleRepo), new(*entity.EntityExample)),

        // Event bus dependencies
        provideEventBus,
        wire.Bind(new(event.EventBus), new(*event.InMemoryEventBus)),

        // Service dependencies
        provideExampleService,
        wire.Bind(new(service.IExampleService), new(*service.ExampleService)),
        provideServices,
    )
    return nil, nil
}

// Provide event bus
func provideEventBus() *event.InMemoryEventBus {
    eventBus := event.NewInMemoryEventBus()

    // Register event handlers
    loggingHandler := event.NewLoggingEventHandler()
    exampleHandler := event.NewExampleEventHandler()
    eventBus.Subscribe(loggingHandler)
    eventBus.Subscribe(exampleHandler)

    return eventBus
}

// Provide example service
func provideExampleService(repo repo.IExampleRepo, eventBus event.EventBus) *service.ExampleService {
    exampleService := service.NewExampleService(repo)
    exampleService.EventBus = eventBus
    return exampleService
}

// Provide services container
func provideServices(exampleService service.IExampleService, eventBus event.EventBus) *service.Services {
    return service.NewServices(exampleService, eventBus)
}
```

## Domain Events

The project supports both synchronous and asynchronous event handling:

### Synchronous Event Handling
```go
// Publish an event synchronously
err := eventBus.Publish(ctx, event.NewExampleCreatedEvent(example.ID, example.Name))
```

### Asynchronous Event Handling
```go
// Configure asynchronous event bus
config := event.DefaultAsyncEventBusConfig()
config.QueueSize = 1000
config.WorkerCount = 10
asyncEventBus := event.NewAsyncEventBus(config)

// Publish an event asynchronously
err := asyncEventBus.Publish(ctx, event.NewExampleCreatedEvent(example.ID, example.Name))

// Graceful shutdown
err := asyncEventBus.Close(5 * time.Second)
```

## Enhanced Caching

The enhanced caching system provides advanced features for robust caching:

```go
// Create an enhanced cache with default options
cache := redis.NewEnhancedCache(redisClient, redis.DefaultCacheOptions())

// Try to get a value with auto-loading if missing
var result MyData
err := cache.TryGetSet(ctx, "key:123", &result, 30*time.Minute, func() (interface{}, error) {
    // This only executes if the key is not in cache
    return fetchDataFromDatabase()
})

// Use distributed lock to prevent concurrent operations
err := cache.WithLock(ctx, "lock:resource", func() error {
    // This code is protected by a distributed lock
    return updateSharedResource()
})
```

## Error Handling

The error system provides a consistent way to handle and propagate errors:

```go
// Create a domain error
if entity == nil {
    return errors.New(errors.ErrorTypeNotFound, "entity not found")
}

// Wrap an error with additional context
if err := repo.Save(entity); err != nil {
    return errors.Wrapf(err, errors.ErrorTypePersistence, "failed to save entity %d", entity.ID)
}

// Check error types
if errors.IsNotFoundError(err) {
    // Handle not found case
}
```

## Structured Logging

The logging system supports context-aware structured logging:

```go
// Create a log context
logCtx := log.NewLogContext().
    WithRequestID(requestID).
    WithUserID(userID).
    WithOperation("CreateEntity")

// Log with context
logger.InfoContext(logCtx, "Creating new entity",
    zap.Int("entity_id", entity.ID),
    zap.String("entity_name", entity.Name))
```

## Project Improvements

The project has recently undergone the following improvements:

### 1. Unified API Versions
- **Problem**: The project had both v1 and v2 API versions, causing code duplication and maintenance difficulties
- **Solution**:
  - Unified API routes, placing all APIs under the `/api` path
  - Retained the `/v2` path for backward compatibility
  - Used application layer use cases to handle all requests, phasing out direct domain service calls

### 2. Enhanced Dependency Injection
- **Problem**: Wire dependency injection configuration had duplicate binding issues, causing generation failures
- **Solution**:
  - Refactored the `wire.go` file, removing duplicate binding definitions
  - Used provider functions instead of direct bindings
  - Added event handler registration logic

### 3. Eliminated Global Variables
- **Problem**: The project used global variables to store service instances, violating dependency injection principles
- **Solution**:
  - Removed global service variables
  - Properly injected services via the Factory pattern
  - Improved testability by making dependencies explicit

### 4. Enhanced Architecture Validation
- **Problem**: Architecture validation was manual and prone to errors
- **Solution**:
  - Implemented automated layer dependency checking
  - Enforced strict architectural boundaries through code scanning
  - Added validation to CI pipeline

### 5. Graceful Shutdown
- **Problem**: Application didn't handle shutdown gracefully, potentially causing data loss
- **Solution**:
  - Implemented a graceful shutdown mechanism for the server, ensuring all in-flight requests are completed before shutting down
  - Added shutdown timeout settings to prevent the shutdown process from hanging indefinitely
  - Improved signal handling, supporting SIGINT and SIGTERM signals

### 6. Internationalization Support
- **Problem**: The application lacked proper internationalization support
- **Solution**:
  - Added translation middleware for multi-language validation error messages
  - Automatically selected appropriate language based on the Accept-Language header

### 7. CORS Support
- **Problem**: Cross-origin requests were not properly handled
- **Solution**:
  - Added CORS middleware to handle cross-origin requests
  - Configured allowed origins, methods, headers, and credentials

### 8. Debugging Tools
- **Problem**: Diagnosis of performance issues in production was difficult
- **Solution**:
  - Integrated pprof performance analysis tools for diagnosing performance issues in production environments
  - Can be enabled or disabled via configuration file

### 9. Advanced Redis Integration
- **Problem**: Redis implementation was limited and lacked proper connection management
- **Solution**:
  - Enhanced Redis client with proper connection pooling
  - Added comprehensive health checks and monitoring
  - Improved error handling and connection lifecycle management

### 10. Structured Request Logging
- **Problem**: API requests lacked proper logging, making debugging difficult
- **Solution**:
  - Implemented comprehensive request logging middleware
  - Added request ID tracking for correlating logs
  - Configured log levels based on status codes

### 11. Unified Error Response Format
- **Problem**: Error responses had inconsistent formats across the API
- **Solution**:
  - Standardized error response structure with code, message, and details
  - Added documentation references to errors
  - Implemented consistent HTTP status code mapping

These optimizations make the project more robust, maintainable, and provide a better development experience.

## Extension Plans

- **gRPC Support** - Add gRPC service implementation
- **Monitoring Integration** - Integrate Prometheus monitoring

## References

- **Architecture**
  - [Freedom DDD Framework](https://github.com/8treenet/freedom)
  - [Hexagonal Architecture in Go](https://medium.com/@matiasvarela/hexagonal-architecture-in-go-cfd4e436faa3)
  - [Dependency Injection in A Nutshell](https://appliedgo.net/di/)
- **Project Standards**
  - [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0)
  - [Improving Your Go Project With pre-commit hooks](https://goangle.medium.com/golang-improving-your-go-project-with-pre-commit-hooks-a265fad0e02f)
- **Code References**
  - [Go CleanArch](https://github.com/roblaszczak/go-cleanarch)
