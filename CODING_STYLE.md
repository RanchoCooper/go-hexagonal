# Go Hexagonal Project Coding Standards

## Directory Structure

The project follows the Hexagonal Architecture structure design:

```
go-hexagonal/
├── adapter/            # Adapter Layer - Connecting domain with external infrastructure
│   ├── repository/     # Repository implementations
│   ├── dependency/     # Dependency injection
│   ├── job/            # Background tasks
│   └── amqp/           # Message queue
├── api/                # API Layer - Handling HTTP, gRPC requests
│   ├── http/           # HTTP API handlers
│   ├── grpc/           # gRPC API handlers
│   ├── error_code/     # Error code definitions
│   └── dto/            # Data Transfer Objects
├── application/        # Application Layer - Orchestrating business flows
│   ├── example/        # Example application services
│   └── core/           # Core interfaces and utilities
├── domain/             # Domain Layer - Core business logic
│   ├── service/        # Domain services
│   ├── repo/           # Repository interfaces
│   ├── event/          # Domain events
│   ├── vo/             # Value objects
│   ├── model/          # Domain models
│   └── aggregate/      # Aggregate roots
├── cmd/                # Application entry points
├── config/             # Configuration
├── tests/              # Tests
└── util/               # Utilities
```

## Naming Conventions

### Package Naming Conventions

- Use lowercase words, no underscores or mixed case
- Package names should be short, meaningful nouns
- Avoid using common variable names as package names

```go
// Correct
package repository
package service

// Incorrect
package Repository
package service_impl
```

### Variable Naming Conventions

- Local variables: Use camelCase, e.g., `userID` instead of `userid`
- Global variables: Use camelCase, capitalize first letter if exported
- Constants: Use all uppercase with underscores, e.g., `MAX_CONNECTIONS`

```go
// Local variables
func processUser() {
    userID := 123
    firstName := "John"
}

// Global variables
var (
    GlobalConfig Configuration
    maxRetryCount = 3
)

// Constants
const (
    MAX_CONNECTIONS = 100
    DEFAULT_TIMEOUT = 30
)
```

### Interface and Struct Naming Conventions

- Interface naming: Usually end with "er", e.g., `Reader`, `Writer`
- Structs implementing specific interfaces: Should be named after functionality rather than interface name
- Avoid abbreviations unless they are very common (like HTTP, URL)

```go
// Interface
type EventHandler interface {
    Handle(event Event) error
}

// Implementation
type LoggingEventHandler struct {
    logger Logger
}
```

## Code Format and Style

All code must pass `go fmt` and `golangci-lint` checks to ensure consistent style:

```bash
# Use make commands for checks
make fmt
make lint
```

### Import Package Ordering

Arrange imports in the following order:

1. Standard library
2. Third-party packages
3. Internal project packages

```go
import (
    // Standard library
    "context"
    "fmt"

    // Third-party packages
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"

    // Internal project packages
    "go-hexagonal/domain/model"
    "go-hexagonal/util/log"
)
```

### File Internal Structure

File content should be organized in the following order:

1. Package documentation comments
2. Package declaration
3. Import packages
4. Constants
5. Variables
6. Type definitions
7. Function definitions

## Comment Standards

### Package Comments

Each package should have package comments placed before the package statement:

```go
// Package repository provides data access implementations
// for the domain repositories.
package repository
```

### Exported Functions and Types Comments

All exported functions, types, constants, and variables should have comments:

```go
// ExampleService handles business logic for Example entities.
// It provides CRUD operations and domain-specific validations.
type ExampleService struct {
    // fields
}

// Create creates a new example entity with the given data.
// It validates the input and publishes an event on successful creation.
// Returns the created entity or an error if validation fails.
func (s *ExampleService) Create(ctx context.Context, example *model.Example) (*model.Example, error) {
    // implementation
}
```

## Error Handling Standards

### Using Unified Error Handling Library

The project uses the `util/errors` package for unified error handling:

```go
import "go-hexagonal/util/errors"

// Creating errors
if input.Name == "" {
    return nil, errors.NewValidationError("Name cannot be empty", nil)
}

// Wrapping errors
result, err := repository.Find(id)
if err != nil {
    return nil, errors.Wrap(err, errors.ErrorTypePersistence, "Failed to query record")
}

// Error type checking
if errors.IsNotFoundError(err) {
    // Handle resource not found case
}
```

### HTTP Layer Error Handling

The API layer uses a unified error handling middleware:

```go
// Middleware is configured in router.go
router.Use(middleware.ErrorHandlerMiddleware())
```

## Testing Standards

### Test Naming Conventions

- Test functions should be named `TestXxx`, where `Xxx` is the name of the function being tested
- Table-driven test variables should be named `tests` or `testCases`

```go
func TestExampleService_Create(t *testing.T) {
    tests := []struct {
        name        string
        input       *model.Example
        mockSetup   func(repo *mocks.MockExampleRepo)
        wantErr     bool
        expectedErr string
    }{
        // Test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## CI/CD Standards

The project uses GitHub Actions for continuous integration to ensure code quality:

- Every commit and PR will run code checks and tests
- Code must pass all checks and tests to be merged
- It's recommended to use pre-commit hooks to check code quality locally

```bash
# Install pre-commit hooks
make pre-commit.install
```

## Best Practices

1. **Dependency Injection**: Always use dependency injection, avoid global variables and singletons
2. **Context Passing**: Always pass context through function calls for cancellation and timeout control
3. **Error Handling**: Use unified error handling, don't discard errors, wrap errors appropriately
4. **Test Coverage**: Ensure critical code has sufficient test coverage, use table-driven tests
5. **Concurrency Safety**: Ensure data structures accessed concurrently are thread-safe
