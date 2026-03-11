# Go Hexagonal Architecture

Visit my [blog post](https://blog.ranchocooper.com/2025/03/20/go-hexagonal/) for more details.

![CI Status](https://github.com/RanchoCooper/go-hexagonal/actions/workflows/ci.yaml/badge.svg)
[![codecov](https://codecov.io/gh/RanchoCooper/go-hexagonal/branch/master/graph/badge.svg)](https://codecov.io/gh/RanchoCooper/go-hexagonal)

![Hexagonal Architecture](https://github.com/Sairyss/domain-driven-hexagon/raw/master/assets/images/DomainDrivenHexagon.png)

## Project Overview

This project is a Go microservice framework based on [Hexagonal Architecture](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software)) and [Domain-Driven Design](https://en.wikipedia.org/wiki/Domain-driven_design). It provides a clean project structure and design patterns to help developers build maintainable, testable, and scalable applications.

Hexagonal Architecture (also known as [Ports and Adapters Architecture](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software))) divides the application into internal and external parts, implementing [Separation of Concerns](https://en.wikipedia.org/wiki/Separation_of_concerns) and [Dependency Inversion Principle](https://en.wikipedia.org/wiki/Dependency_inversion_principle) through well-defined interfaces (ports) and implementations (adapters).

## Core Features

### Architecture Design
- **[Domain-Driven Design (DDD)](https://en.wikipedia.org/wiki/Domain-driven_design)** - Organize business logic through Aggregates, Entities, and Value Objects
- **[Hexagonal Architecture](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software))** - Divide the application into domain, application, and adapter layers
- **[Dependency Injection](https://en.wikipedia.org/wiki/Dependency_injection)** - Use [Wire](https://github.com/google/wire) for dependency injection
- **[Repository Pattern](https://en.wikipedia.org/wiki/Repository_pattern)** - Abstract data access layer with transaction support
- **[Domain Events](https://en.wikipedia.org/wiki/Domain-driven_design)** - Implement Event-Driven Architecture with sync/async support
- **[CQRS Pattern](https://en.wikipedia.org/wiki/Command_Query_Responsibility_Segregation)** - Command and Query Responsibility Segregation

### Technical Implementation
- **[RESTful API](https://en.wikipedia.org/wiki/Representational_state_transfer)** - HTTP API using [Gin](https://github.com/gin-gonic/gin)
- **Database Support** - [GORM](https://gorm.io) with MySQL support, configurable connection pool
- **Cache Support** - [Redis](https://en.wikipedia.org/wiki/Redis) caching with health checks
- **Enhanced Cache Features**:
  - Negative caching to prevent cache penetration
  - Distributed locking for cache consistency
  - Key tracking with LRU eviction
  - Pipeline support for batch operations
- **Logging System** - [Zap](https://go.uber.org/zap) for structured logging
- **Configuration Management** - [Viper](https://github.com/spf13/viper) for flexible configuration
- **Graceful Shutdown** - Support for graceful service startup and shutdown
- **Transaction Support** - Real transaction factory with rollback/commit
- **Cache-DB Transaction Consistency** - Cache operations within database transactions

### Monitoring
- **Prometheus Metrics** - HTTP requests, DB operations, cache performance, domain events
- **Health Checks** - Redis and database health monitoring

## Project Structure

```
.
├── adapter/                # Adapter Layer - External system interactions
│   ├── dependency/        # Dependency injection (Wire)
│   ├── job/               # Scheduled tasks
│   └── repository/        # Data repository adapters
│       ├── mysql/         # MySQL implementation
│       ├── postgre/       # PostgreSQL implementation
│       └── redis/         # Redis implementation
│           ├── client.go           # Redis client
│           ├── enhanced_cache.go  # Enhanced cache
│           └── example_cache.go   # Example cache repo
├── api/                   # API Layer
│   ├── dto/              # Data Transfer Objects
│   ├── error_code/       # Error codes
│   ├── grpc/             # gRPC handlers
│   ├── middleware/       # HTTP middleware
│   └── http/             # HTTP handlers
├── application/           # Application Layer
│   ├── core/            # Core interfaces
│   └── example/         # Use cases
├── cmd/                  # Entry points
├── config/               # Configuration
├── domain/               # Domain Layer
│   ├── aggregate/       # Aggregates
│   ├── dto/            # Domain DTOs
│   ├── event/          # Domain events
│   ├── model/          # Domain models
│   ├── repo/           # Repository interfaces
│   ├── service/        # Domain services
│   └── vo/             # Value Objects
├── tests/                # Integration tests
└── util/                 # Utilities
    ├── errors/          # Error handling
    ├── log/             # Logging
    ├── metrics/         # Prometheus metrics
    └── transaction/     # Transaction utilities
```

## Architecture Layers

### Domain Layer
- **Models**: Domain entities (Example)
- **Repository Interfaces**: IExampleRepo, IExampleCacheRepo
- **Services**: ExampleService with transaction support
- **Events**: ExampleCreatedEvent, ExampleUpdatedEvent, ExampleDeletedEvent

### Application Layer
- Use cases: Create, Get, Update, Delete, FindByName

### Adapter Layer
- **MySQL**: GORM-based implementation with transaction support
- **Redis**: Client with connection pooling, enhanced cache with LRU, locks, negative caching
- **Event Bus**: In-memory and async event processing

### API Layer
- HTTP handlers with request validation
- DTO transformation
- Error handling middleware

## Key Features

### Transaction Factory
```go
// Real transaction with automatic rollback
txFactory := repository.NewTransactionFactory(map[repository.StoreType]any{
    repository.MySQLStore: clients.MySQL,
})
service := service.NewExampleService(repo, cacheRepo, txFactory)
```

### Enhanced Cache with Pipeline
```go
// Pipeline for batch operations
pipe := client.Pipeline()
pipe.Set(ctx, key1, value1, ttl1)
pipe.Set(ctx, key2, value2, ttl2)
cmds, err := pipe.Exec(ctx)
```

### Cache Metrics
```go
metrics.RecordCacheHit("example_cache", "GetByID_hit")
metrics.RecordCacheHit("example_cache", "GetByID_miss")
```

### Cache with LRU Key Tracking
```go
cache := redis.NewEnhancedCache(client, redis.CacheOptions{
    EnableKeyTracking: true,
    MaxTrackedKeys:   100,
})
```

## Error Handling

```go
// Create domain error
return errors.New(errors.ErrorTypeNotFound, "entity not found")

// Wrap error with context
return errors.Wrapf(err, errors.ErrorTypePersistence, "failed to save entity %d", id)

// Use errors.Is for comparison
if errors.Is(err, redis.Nil) {
    return ErrCacheMiss
}
```

## Getting Started

### Prerequisites
- Go 1.21+
- MySQL
- Redis

### Configuration

Edit `config/config.yaml`:

```yaml
mysql:
  host: localhost
  port: 3306
  user: root
  password: password
  database: hex_demo
  max_idle_conns: 10
  max_open_conns: 100

redis:
  host: localhost
  port: 6379
  pool_size: 10
```

### Run

```bash
# Run tests
go test ./... -short

# Run application
go run cmd/main.go
```

## License

[Apache 2.0 License](LICENSE)
