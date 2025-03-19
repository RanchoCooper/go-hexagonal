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

### Technical Implementation
- **[RESTful API](https://en.wikipedia.org/wiki/Representational_state_transfer)** - Implement HTTP API using the [Gin](https://github.com/gin-gonic/gin) framework
- **Database Support** - Integrate [GORM](https://gorm.io) with support for [MySQL](https://en.wikipedia.org/wiki/MySQL), [PostgreSQL](https://en.wikipedia.org/wiki/PostgreSQL), and other databases
- **Cache Support** - Integrate [Redis](https://en.wikipedia.org/wiki/Redis) caching
- **Logging System** - Use [Zap](https://go.uber.org/zap) for high-performance logging
- **Configuration Management** - Use [Viper](https://github.com/spf13/viper) for flexible configuration management
- **[Graceful Shutdown](https://en.wikipedia.org/wiki/Graceful_exit)** - Support graceful service startup and shutdown
- **[Unit Testing](https://en.wikipedia.org/wiki/Unit_testing)** - Use [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) and [redismock](https://github.com/go-redis/redismock) for database and cache mocking
- **[NoopTransaction](NoopTransaction)** - Provide no-operation transaction implementation, simplifying service layer interaction with repository layer

### Development Toolchain
- **Code Quality** - Integrate [Golangci-lint](https://github.com/golangci/golangci-lint) for code quality checks
- **Commit Standards** - Use [Commitlint](https://github.com/conventional-changelog/commitlint) to ensure Git commit messages follow conventions
- **Pre-commit Hooks** - Use [Pre-commit](https://pre-commit.com) for code checking and formatting
- **[CI/CD](https://en.wikipedia.org/wiki/CI/CD)** - Integrate [GitHub Actions](https://github.com/features/actions) for continuous integration and deployment

## Project Structure

```
.
├── adapter/                # Adapter Layer - Interaction with external systems
│   ├── amqp/               # Message queue adapters
│   ├── dependency/         # Dependency injection configuration
│   ├── job/                # Scheduled task adapters
│   └── repository/         # Data repository adapters
│       ├── mysql/          # MySQL implementation
│       │   └── entity/     # Database entities
│       ├── postgre/        # PostgreSQL implementation
│       └── redis/          # Redis implementation
├── api/                    # API Layer - Handle HTTP requests and responses
│   ├── dto/                # Data Transfer Objects
│   ├── error_code/         # Error code definitions
│   ├── grpc/               # gRPC API
│   └── http/               # HTTP API
│       ├── handle/         # Request handlers
│       ├── middleware/     # Middleware
│       ├── paginate/       # Pagination handling
│       └── validator/      # Request validation
├── application/            # Application Layer - Coordinate domain objects for use cases
│   ├── core/               # Core interfaces and error definitions
│   └── example/            # Example use case implementations
├── cmd/                    # Command-line entry points
│   └── http_server/        # HTTP server startup
├── config/                 # Configuration files and management
├── domain/                 # Domain Layer - Core business logic
│   ├── aggregate/          # Aggregates
│   ├── event/              # Domain events
│   ├── model/              # Domain models
│   ├── repo/               # Repository interfaces
│   ├── service/            # Domain services
│   └── vo/                 # Value objects
├── tests/                  # Integration tests
└── util/                   # Utility functions
    ├── clean_arch/         # Architecture checking tools
    └── log/                # Logging utilities
```

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
  - `ExampleService`: Example service, handling business logic for example entities, interacting with repositories and event bus

- **Domain Events**: Define events within the domain
  - `ExampleCreatedEvent`: Example creation event
  - `ExampleUpdatedEvent`: Example update event
  - `ExampleDeletedEvent`: Example deletion event

### Application Layer
The application layer coordinates domain objects to complete specific application tasks. It depends on the domain layer but does not contain business rules.

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
func provideServices(exampleService *service.ExampleService, eventBus event.EventBus) *service.Services {
    return service.NewServices(exampleService, eventBus)
}
```

## Domain Events

Domain events are used for communication between system components, implementing a loosely coupled event-driven architecture:

```go
// Publish event
evt := event.NewExampleCreatedEvent(example.Id, example.Name, example.Alias)
e.EventBus.Publish(ctx, evt)

// Handle event
func (h *ExampleEventHandler) HandleEvent(ctx context.Context, event Event) error {
    switch event.EventName() {
    case ExampleCreatedEventName:
        return h.handleExampleCreated(ctx, event)
    // ...
    }
    return nil
}
```

## Application Layer Use Cases

Application layer use cases implement the Command and Query Responsibility Segregation (CQRS) pattern:

```go
// Create example use case
func (h *CreateExampleHandler) Handle(ctx context.Context, input interface{}) (interface{}, error) {
    createInput, ok := input.(CreateExampleInput)
    if !ok {
        return nil, core.ErrInvalidInput
    }

    example := &model.Example{
        Name:  createInput.Name,
        Alias: createInput.Alias,
    }

    createdExample, err := h.ExampleService.Create(ctx, example)
    if err != nil {
        return nil, err
    }

    return CreateExampleOutput{
        ID:    createdExample.Id,
        Name:  createdExample.Name,
        Alias: createdExample.Alias,
    }, nil
}
```

## Transaction Management

This project implements transaction interfaces and no-operation transactions, supporting different transaction management strategies:

```go
// Transaction interface
type Transaction interface {
    Begin() error
    Commit() error
    Rollback() error
    Conn(ctx context.Context) interface{}
}

// No-operation transaction implementation
type NoopTransaction struct {
    conn interface{}
}

// Using transactions in services
func (s *ExampleService) Create(ctx context.Context, example *model.Example) (*model.Example, error) {
    // Create a no-operation transaction
    tr := repo.NewNoopTransaction(s.Repository)

    createdExample, err := s.Repository.Create(ctx, tr, example)
    // ...
}
```

## Data Mapping and Transformation

This project implements clear data mapping and transformation between different layers:

```go
// Entity to model conversion
func (e EntityExample) ToModel() *model.Example {
    return &model.Example{
        Id:        e.ID,
        Name:      e.Name,
        Alias:     e.Alias,
        CreatedAt: e.CreatedAt,
        UpdatedAt: e.UpdatedAt,
    }
}

// Model to entity conversion
func (e *EntityExample) FromModel(m *model.Example) {
    e.ID = m.Id
    e.Name = m.Name
    e.Alias = m.Alias
    e.CreatedAt = m.CreatedAt
    e.UpdatedAt = m.UpdatedAt
}
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
  - Removed the use of global variables `service.ExampleSvc` and `service.EventBus`
  - Passed service instances through dependency injection
  - Initialized services using dependency injection when starting the HTTP server

### 4. Improved Application Layer Integration
- **Problem**: Application layer use cases were not fully utilized, with the HTTP server not enabling the application layer by default
- **Solution**:
  - Enabled application layer use cases by default
  - Used the use case factory to create and manage use cases
  - Implemented clearer error handling and response mapping

## Recent Optimizations

The project has recently undergone the following optimizations:

1. **Environment Variable Support**:
   - Added functionality for environment variable overrides for configuration files, making the application more flexible in containerized deployments
   - Used a unified prefix (APP_) and hierarchical structure (e.g., APP_MYSQL_HOST) to organize environment variables

2. **Unified Error Handling**:
   - Implemented an application-level error type system, supporting different types of errors (validation, not found, unauthorized, etc.)
   - Added unified error response handling, mapping internal errors to appropriate HTTP status codes
   - Improved error logging to ensure all unexpected errors are properly recorded

3. **Request Logging Middleware**:
   - Added detailed request logging middleware to record request methods, paths, status codes, latency, and other information
   - In debug mode, request and response bodies can be logged to help developers troubleshoot issues
   - Intelligently identifies content types to avoid logging binary content

4. **Request ID Tracking**:
   - Generated unique request IDs for each request, facilitating tracking in distributed systems
   - Returned request IDs in response headers for client reference
   - Included request IDs in logs to correlate multiple log entries for the same request

5. **Graceful Shutdown**:
   - Implemented a graceful shutdown mechanism for the server, ensuring all in-flight requests are completed before shutting down
   - Added shutdown timeout settings to prevent the shutdown process from hanging indefinitely
   - Improved signal handling, supporting SIGINT and SIGTERM signals

6. **Internationalization Support**:
   - Added translation middleware for multi-language validation error messages
   - Automatically selected appropriate language based on the Accept-Language header

7. **CORS Support**:
   - Added CORS middleware to handle cross-origin requests
   - Configured allowed origins, methods, headers, and credentials

8. **Debugging Tools**:
   - Integrated pprof performance analysis tools for diagnosing performance issues in production environments
   - Can be enabled or disabled via configuration file

These optimizations make the project more robust, maintainable, and provide a better development experience.

## Usage Guide

### Environment Preparation

Start MySQL using Docker:
```bash
docker run --name mysql-local \
  -e MYSQL_ROOT_PASSWORD=mysqlroot \
  -e MYSQL_DATABASE=go-hexagonal \
  -e MYSQL_USER=user \
  -e MYSQL_PASSWORD=mysqlroot \
  -p 3306:3306 \
  -d mysql:latest
```

### Development Tool Installation

```bash
# Install development tools
make init && make precommit.rehook
```

Or install manually:

```bash
# Install pre-commit
brew install pre-commit
# Install golangci-lint
brew install golangci-lint
# Install commitlint
npm install -g @commitlint/cli @commitlint/config-conventional
# Add commitlint configuration
echo "module.exports = {extends: ['@commitlint/config-conventional']}" > commitlint.config.js
# Add pre-commit hook
make precommit.rehook
```

### Running the Project

```bash
# Run the project
go run cmd/main.go
```

### Testing

```bash
# Run tests
go test ./...
```

## Extension Plans

- **Dependency Injection Improvements** - Enhance Wire dependency injection configuration
- **HTTP Handling Improvements** - Optimize HTTP request handling implementation
- **Domain Event Enhancements** - Improve domain event mechanisms
- **gRPC Support** - Add gRPC service implementation
- **Hot Reload Configuration** - Implement configuration hot reloading
- **Monitoring Integration** - Integrate Prometheus monitoring
- **Message Queue Integration** - Integrate Kafka and other message queues

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
