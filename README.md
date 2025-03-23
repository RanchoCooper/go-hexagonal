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
- **MongoDB Support** - Integration with MongoDB for document storage
- **Logging System** - Use [Zap](https://go.uber.org/zap) for high-performance logging
- **Configuration Management** - Use [Viper](https://github.com/spf13/viper) for flexible configuration management
- **[Graceful Shutdown](https://en.wikipedia.org/wiki/Graceful_exit)** - Support graceful service startup and shutdown
- **[Unit Testing](https://en.wikipedia.org/wiki/Unit_testing)** - Use [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock), [redismock](https://github.com/go-redis/redismock), and [testify/mock](https://github.com/stretchr/testify) for comprehensive test coverage with enhanced HTTP testing utilities and improved DTO handling
- **Transaction Support** - Provide no-operation transaction implementation, simplifying service layer interaction with repository layer, complete with mock transaction implementation and lifecycle hooks (Begin, Commit, and Rollback) for testing

### Development Toolchain
- **Code Quality** - Integrate [Golangci-lint](https://github.com/golangci/golangci-lint) for code quality checks
- **Commit Standards** - Use [Commitlint](https://github.com/conventional-changelog/commitlint) to ensure Git commit messages follow conventions
- **Pre-commit Hooks** - Use [Pre-commit](https://pre-commit.com) for code checking and formatting
- **[CI/CD](https://en.wikipedia.org/wiki/CI/CD)** - Integrate [GitHub Actions](https://github.com/features/actions) for continuous integration and deployment

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
    └── log/                # Logging utilities
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

Application layer use cases implement the Command and Query Responsibility Segregation (CQRS) pattern and depend on domain interfaces rather than concrete implementations:

```go
// Use case with interface dependency
type CreateUseCase struct {
    exampleService service.IExampleService
}

// Create example use case
func (uc *CreateUseCase) Execute(ctx context.Context, input dto.CreateExampleReq) (*dto.CreateExampleResp, error) {
    // Create a real transaction for atomic operations
    tx, err := repository.NewTransaction(ctx, repository.MySQLStore, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create transaction: %w", err)
    }
    defer tx.Rollback()

    // Convert DTO to domain model
    example := &model.Example{
        Name:  input.Name,
        Alias: input.Alias,
    }

    // Call domain service through interface
    createdExample, err := uc.exampleService.Create(ctx, example)
    if err != nil {
        return nil, fmt.Errorf("failed to create example: %w", err)
    }

    // Commit transaction
    if err = tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %w", err)
    }

    // Convert domain model to DTO
    result := &dto.CreateExampleResp{
        Id:        uint(createdExample.Id),
        Name:      createdExample.Name,
        Alias:     createdExample.Alias,
        CreatedAt: createdExample.CreatedAt,
        UpdatedAt: createdExample.UpdatedAt,
    }

    return result, nil
}
```

## Unit Testing

The project implements comprehensive unit testing strategies:

- **Interface-Based Testing**: Test against interfaces rather than concrete implementations
- **Mock Objects**: Use testify/mock to create mock implementations of interfaces
- **Transaction Mocking**: Separate database operations from business logic by mocking transactions
- **Standardized Testing Pattern**: Follow a consistent pattern for all tests
  - Create mock services
  - Set up test data and expectations
  - Execute use case
  - Assert results
  - Verify mock expectations

Example of a unit test with mocked dependencies:

```go
// Mock implementation of IExampleService interface for testing
type mockExampleService struct {
    mock.Mock
}

func (m *mockExampleService) Create(ctx context.Context, example *model.Example) (*model.Example, error) {
    args := m.Called(ctx, example)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.Example), args.Error(1)
}

// Test for successful creation
func TestCreateUseCase_Execute_Success(t *testing.T) {
    // Create mock service
    mockService := new(mockExampleService)

    // Set up mock behavior
    now := time.Now()
    expectedExample := &model.Example{
        Id:        1,
        Name:      "Test Example",
        Alias:     "test",
        CreatedAt: now,
        UpdatedAt: now,
    }
    mockService.On("Create", mock.Anything, mock.Anything).Return(expectedExample, nil)

    // Create testable use case
    useCase := newTestableCreateUseCase(mockService)

    // Execute use case
    ctx := context.Background()
    createReq := dto.CreateExampleReq{
        Name:  "Test Example",
        Alias: "test",
    }
    result, err := useCase.Execute(ctx, createReq)

    // Assert results
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, uint(expectedExample.Id), result.Id)
    mockService.AssertExpectations(t)
}
```

## Coding Standards

This project follows unified coding standards to ensure code quality and consistency. For detailed guidelines, please refer to [CODING_STYLE.md](./CODING_STYLE.md).

Key standards include:

- Code format and style (using go fmt and golangci-lint)
- Naming conventions (package names, variable names, interfaces and structs)
- Import package ordering
- Comment standards
- Error handling standards (using util/errors package)
- Testing standards
- CI/CD standards

Developers should ensure compliance with these standards before submitting code. Use the following commands for verification:

```bash
# Format code
make fmt

# Code quality check
make lint

# Run tests
make test
```

## Transaction Management

This project implements transaction interfaces and no-operation transactions, supporting different transaction management strategies:

```go
// Transaction interface
type Transaction interface {
    Begin() error
    Commit() error
    Rollback() error
    Conn(ctx context.Context) any
}

// No-operation transaction implementation
type NoopTransaction struct {
    conn any
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

This project implements clear data mapping and transformation between different layers using the [jinzhu/copier](https://github.com/jinzhu/copier) library for efficient object copying:

```go
// Entity to model conversion using copier
func (e EntityExample) ToModel() *model.Example {
    model := &model.Example{}
    copier.Copy(model, e)
    return model
}

// Model to entity conversion using copier
func (e *EntityExample) FromModel(m *model.Example) {
    copier.Copy(e, m)
}

// Batch conversion from entities to models
func EntitiesToModels(entities []EntityExample) []*model.Example {
    result := make([]*model.Example, len(entities))
    for i, entity := range entities {
        result[i] = entity.ToModel()
    }
    return result
}

// DTO to model conversion
func (req *CreateExampleReq) ToModel() *model.Example {
    model := &model.Example{}
    copier.Copy(model, req)
    return model
}

// Model to response DTO conversion
func ModelToResponse(m *model.Example) *ExampleResponse {
    if m == nil {
        return nil
    }

    resp := &ExampleResponse{}
    copier.Copy(resp, m)

    // Format time fields after copying
    resp.CreatedAt = m.CreatedAt.Format(time.RFC3339)
    resp.UpdatedAt = m.UpdatedAt.Format(time.RFC3339)

    return resp
}
```

Benefits of using the copier library:
- Simplifies conversion between similar structs
- Automatically copies fields with the same name and compatible types
- Supports deep copying of nested structs
- Reduces boilerplate code for object transformations

These conversions maintain a clear separation between different layers:
- Database entities (in the adapter layer)
- Domain models (in the domain layer)
- Data Transfer Objects (in the API layer)

This approach allows each layer to have its own representation of the data, optimized for its specific responsibilities.

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
