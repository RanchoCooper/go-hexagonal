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

### Hexagonal Architecture-Specific Naming

- Repository interfaces: Prefix with `I`, e.g., `IUserRepository`
- Service interfaces: Prefix with `I`, e.g., `IUserService`
- Use case implementations: Name after the action, e.g., `CreateUserUseCase`
- Controllers: Name after resource/functionality, e.g., `UserController`
- DTOs: Suffix with `DTO`, `Request`, or `Response`, e.g., `UserDTO`, `CreateUserRequest`

```go
// Domain layer
type IUserRepository interface {
    FindByID(ctx context.Context, id string) (*model.User, error)
}

// Application layer
type CreateUserUseCase struct {
    userRepo domain.IUserRepository
    eventBus event.EventBus
}

// API layer
type CreateUserRequest struct {
    Username string `json:"username" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
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

The project uses the `api/error_code` package for unified error handling:

```go
import "go-hexagonal/api/error_code"

// Using predefined errors
if input.Name == "" {
    return nil, error_code.InvalidParams.WithDetails("Name cannot be empty")
}

// Creating custom errors with details
if !isValid {
    return nil, error_code.NewError(40001, "Invalid data format").WithDetails("Field X should be Y format")
}

// Error type checking
if errors.Is(err, error_code.NotFound) {
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

### Hexagonal Architecture Testing Strategy

#### Domain Layer Tests

Domain layer tests should focus on business logic without external dependencies:

```go
func TestExampleService_Validate(t *testing.T) {
    // Test domain logic in isolation
    service := service.NewExampleService()
    err := service.Validate(example)

    assert.NoError(t, err)
}
```

#### Application Layer Tests

Application layer tests should mock external dependencies:

```go
func TestCreateExampleUseCase_Execute(t *testing.T) {
    // Setup mocks
    mockRepo := mocks.NewMockExampleRepo(t)
    mockEventBus := mocks.NewMockEventBus(t)

    // Set expectations
    mockRepo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)
    mockEventBus.EXPECT().Publish(mock.Anything, mock.Anything).Return(nil)

    // Create use case with mocks
    useCase := usecase.NewCreateExampleUseCase(mockRepo, mockEventBus)

    // Execute use case
    result, err := useCase.Execute(ctx, input)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

#### Adapter Layer Tests

Adapter layer tests should verify interactions with external systems:

```go
func TestMySQLExampleRepo_Save(t *testing.T) {
    // Setup DB mock
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create mock DB: %v", err)
    }
    defer db.Close()

    // Set expectations
    mock.ExpectBegin()
    mock.ExpectExec("INSERT INTO examples").WithArgs(...).WillReturnResult(sqlmock.NewResult(1, 1))
    mock.ExpectCommit()

    // Create repo with mock DB
    repo := repository.NewMySQLExampleRepo(db)

    // Test save operation
    err = repo.Save(ctx, example)

    // Assert
    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
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

### Domain Layer

1. **Keep Domain Models Pure**: Domain models should not have dependencies on external libraries or frameworks
2. **Encapsulate Business Logic**: Keep business logic within domain services and entities
3. **Use Value Objects**: Create value objects for concepts with identity based on their attributes
4. **Rich Domain Model**: Prefer rich domain models over anemic models
5. **Domain Events**: Use domain events to propagate state changes across bounded contexts

```go
// Rich domain model with behavior
func (u *User) ChangePassword(currentPassword, newPassword string) error {
    if !u.ValidatePassword(currentPassword) {
        return ErrInvalidPassword
    }

    if err := u.ValidatePasswordStrength(newPassword); err != nil {
        return err
    }

    hashedPassword, err := hashPassword(newPassword)
    if err != nil {
        return err
    }

    u.Password = hashedPassword
    u.LastPasswordChange = time.Now()

    return nil
}
```

### Application Layer

1. **One Use Case, One File**: Each use case should be in its own file
2. **Keep Use Cases Focused**: Use cases should coordinate domain objects to complete a specific task
3. **Use DTOs at Boundaries**: Use DTOs for inputs and outputs at application boundaries
4. **Input Validation**: Validate inputs as early as possible
5. **Transaction Management**: Handle transactions at the application layer

```go
// Focused use case
func (uc *UpdateUserProfileUseCase) Execute(ctx context.Context, input UpdateUserProfileInput) (*UpdateUserProfileOutput, error) {
    // Input validation
    if err := uc.validator.Validate(input); err != nil {
        return nil, err
    }

    // Transaction management
    return uc.txManager.WithTransaction(ctx, func(ctx context.Context) (*UpdateUserProfileOutput, error) {
        // Get user
        user, err := uc.userRepo.FindByID(ctx, input.UserID)
        if err != nil {
            return nil, err
        }

        // Update user
        user.UpdateProfile(input.Name, input.Bio, input.Location)

        // Save user
        if err := uc.userRepo.Save(ctx, user); err != nil {
            return nil, err
        }

        // Publish event
        event := event.NewUserProfileUpdatedEvent(user.ID, user.Name)
        if err := uc.eventBus.Publish(ctx, event); err != nil {
            return nil, err
        }

        // Return output
        return &UpdateUserProfileOutput{
            User: mapUserToDTO(user),
        }, nil
    })
}
```

### Adapter Layer

1. **Infrastructure Concerns Only**: Keep infrastructure concerns in the adapter layer
2. **Implement Interfaces**: Implement interfaces defined by the domain layer
3. **Map Between Models**: Map between domain models and infrastructure-specific models
4. **Handle Infrastructure Errors**: Map infrastructure errors to domain errors
5. **Keep Adapters Thin**: Adapters should be thin and focused on infrastructure concerns

```go
// Repository implementation in adapter layer
func (r *MySQLUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
    var entity UserEntity

    // Infrastructure-specific logic
    result := r.db.WithContext(ctx).First(&entity, "id = ?", id)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, error_code.NotFound.WithDetails("User not found")
        }
        return nil, error_code.ServerError.WithDetails(result.Error.Error())
    }

    // Map to domain model
    return mapEntityToDomain(entity), nil
}
```

### API Layer

1. **Thin Controllers**: Keep controllers thin and focused on HTTP concerns
2. **Use DTOs**: Use DTOs for request and response data
3. **Validate Requests**: Validate request data before passing to use cases
4. **Consistent Responses**: Use consistent response formats
5. **Document APIs**: Document APIs using Swagger/OpenAPI

```go
// Thin controller
func (c *UserController) UpdateProfile(ctx *gin.Context) {
    // Parse request
    var req UpdateProfileRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.Error(error_code.InvalidParams.WithDetails(err.Error()))
        return
    }

    // Map to use case input
    input := application.UpdateUserProfileInput{
        UserID:   ctx.Param("id"),
        Name:     req.Name,
        Bio:      req.Bio,
        Location: req.Location,
    }

    // Execute use case
    output, err := c.updateProfileUseCase.Execute(ctx, input)
    if err != nil {
        ctx.Error(err)
        return
    }

    // Return response
    ctx.JSON(http.StatusOK, output)
}
```

## Dependency Injection

1. **Use Wire for DI**: Use Google Wire for dependency injection
2. **Separate Wire Configuration**: Keep Wire configuration in the adapter/dependency directory
3. **Provider Functions**: Use provider functions to create instances
4. **Group Related Providers**: Group related providers together
5. **Clear Dependencies**: Make dependencies explicit in provider functions

```go
// Provider functions
func ProvideUserRepository(db *gorm.DB) domain.IUserRepository {
    return repository.NewMySQLUserRepository(db)
}

func ProvideUserService(repo domain.IUserRepository, eventBus event.EventBus) domain.IUserService {
    return service.NewUserService(repo, eventBus)
}

func ProvideUpdateUserProfileUseCase(
    userRepo domain.IUserRepository,
    eventBus event.EventBus,
    txManager transaction.Manager,
    validator validator.Validator,
) *application.UpdateUserProfileUseCase {
    return application.NewUpdateUserProfileUseCase(userRepo, eventBus, txManager, validator)
}

// Wire set
var UserSet = wire.NewSet(
    ProvideUserRepository,
    ProvideUserService,
    ProvideUpdateUserProfileUseCase,
)
```

## Event-Driven Architecture

1. **Domain Events**: Define domain events in the domain layer
2. **Event Bus Interface**: Define event bus interface in the domain layer
3. **Event Handlers**: Implement event handlers in the application layer
4. **Event Persistence**: Consider event persistence for reliability
5. **Asynchronous Processing**: Use asynchronous processing for performance

```go
// Domain event
type UserCreatedEvent struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

func NewUserCreatedEvent(id, name, email string) *UserCreatedEvent {
    return &UserCreatedEvent{
        ID:        id,
        Name:      name,
        Email:     email,
        CreatedAt: time.Now(),
    }
}

// Event handler
type UserCreatedEventHandler struct {
    notificationService notification.Service
}

func (h *UserCreatedEventHandler) Handle(ctx context.Context, event event.Event) error {
    userCreatedEvent, ok := event.(*UserCreatedEvent)
    if !ok {
        return errors.New("invalid event type")
    }

    return h.notificationService.SendWelcomeEmail(ctx, userCreatedEvent.Email, userCreatedEvent.Name)
}

// Publishing event
func (s *UserService) CreateUser(ctx context.Context, user *model.User) error {
    if err := s.userRepo.Save(ctx, user); err != nil {
        return err
    }

    event := NewUserCreatedEvent(user.ID, user.Name, user.Email)
    return s.eventBus.Publish(ctx, event)
}
```

## CQRS Pattern

1. **Command Handlers**: Implement command handlers for write operations
2. **Query Handlers**: Implement query handlers for read operations
3. **Separate Models**: Use separate models for queries and commands
4. **Optimized Queries**: Optimize query models for specific read use cases
5. **Command Validation**: Validate commands before execution

```go
// Command
type CreateUserCommand struct {
    Username  string `json:"username" validate:"required"`
    Email     string `json:"email" validate:"required,email"`
    FirstName string `json:"first_name" validate:"required"`
    LastName  string `json:"last_name" validate:"required"`
}

// Command handler
func (h *CreateUserCommandHandler) Handle(ctx context.Context, cmd CreateUserCommand) error {
    // Validate command
    if err := h.validator.Validate(cmd); err != nil {
        return err
    }

    // Execute command
    user := model.NewUser(cmd.Username, cmd.Email, cmd.FirstName, cmd.LastName)
    return h.userRepo.Save(ctx, user)
}

// Query
type GetUserByIDQuery struct {
    ID string `json:"id" validate:"required,uuid"`
}

// Query handler
func (h *GetUserByIDQueryHandler) Handle(ctx context.Context, query GetUserByIDQuery) (*UserDTO, error) {
    // Validate query
    if err := h.validator.Validate(query); err != nil {
        return nil, err
    }

    // Execute query
    user, err := h.userRepo.FindByID(ctx, query.ID)
    if err != nil {
        return nil, err
    }

    // Map to DTO
    return mapUserToDTO(user), nil
}
```

## Commit Message Standards

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Types:
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to the build process or auxiliary tools

Examples:
```
feat(user): add user registration endpoint
fix(auth): correct token validation logic
docs(readme): update installation instructions
```
