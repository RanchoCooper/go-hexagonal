# Go Hexagonal Architecture

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
- **Repository Interfaces**: Define data access interfaces
- **Domain Services**: Handle business logic across entities
- **Domain Events**: Define events within the domain

### Application Layer
The application layer coordinates domain objects to complete specific application tasks. It depends on the domain layer but does not contain business rules.

- **Use Cases**: Define application functionality
- **Commands and Queries**: Implement CQRS pattern
- **Event Handlers**: Process domain events

### Adapter Layer
The adapter layer implements interaction with external systems, such as databases and message queues.

- **Repository Implementation**: Implement data access interfaces
- **Message Queue Adapters**: Implement message publishing and subscription
- **Scheduled Tasks**: Implement scheduled tasks

### API Layer
The API layer handles HTTP requests and responses, serving as the entry point to the application.

- **Controllers**: Handle HTTP requests
- **Middleware**: Implement cross-cutting concerns
- **Data Transfer Objects (DTOs)**: Define request and response data structures

## Dependency Injection

This project uses Google Wire for dependency injection, organizing dependencies as follows:

```go
// Initialize services
func InitializeServices(ctx context.Context) (*service.Services, error) {
    // Create repositories
    entityExample := entity.NewExample()

    // Create event bus and handlers
    inMemoryEventBus := event.NewInMemoryEventBus()
    loggingHandler := event.NewLoggingEventHandler()
    exampleHandler := event.NewExampleEventHandler()

    // Register event handlers
    inMemoryEventBus.Subscribe(loggingHandler)
    inMemoryEventBus.Subscribe(exampleHandler)

    // Create services
    exampleService := service.NewExampleService(ctx)
    exampleService.Repository = entityExample
    exampleService.EventBus = inMemoryEventBus

    // Create services container
    services := service.NewServices(exampleService, inMemoryEventBus)

    return services, nil
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
