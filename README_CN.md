# Go 六边形架构

欢迎访问我的[博客文章](https://blog.ranchocooper.com/2025/03/20/go-hexagonal/)了解更多详情。

![CI Status](https://github.com/RanchoCooper/go-hexagonal/actions/workflows/ci.yaml/badge.svg)
[![codecov](https://codecov.io/gh/RanchoCooper/go-hexagonal/branch/master/graph/badge.svg)](https://codecov.io/gh/RanchoCooper/go-hexagonal)

![六边形架构](https://github.com/Sairyss/domain-driven-hexagon/raw/master/assets/images/DomainDrivenHexagon.png)

## 项目概述

本项目是一个基于[六边形架构](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software))和[领域驱动设计](https://en.wikipedia.org/wiki/Domain-driven_design)的 Go 微服务框架。它提供了清晰的项目结构和设计模式，帮助开发者构建可维护、可测试和可扩展的应用程序。

六边形架构（也称为端口和适配器架构）将应用程序分为内部和外部部分，通过明确定义的接口（端口）和实现（适配器）实现关注点分离和依赖倒置原则。

## 核心特性

### 架构设计
- **[领域驱动设计 (DDD)](https://en.wikipedia.org/wiki/Domain-driven_design)** - 通过聚合、实体和值对象组织业务逻辑
- **[六边形架构](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software))** - 将应用分为领域层、应用层和适配器层
- **[依赖注入](https://en.wikipedia.org/wiki/Dependency_injection)** - 使用 [Wire](https://github.com/google/wire) 进行依赖注入
- **[仓储模式](https://en.wikipedia.org/wiki/Repository_pattern)** - 抽象数据访问层，支持事务
- **[领域事件](https://en.wikipedia.org/wiki/Domain-driven_design)** - 实现事件驱动架构，支持同步/异步处理
- **[CQRS 模式](https://en.wikipedia.org/wiki/Command_Query_Responsibility_Segregation)** - 命令查询职责分离

### 技术实现
- **[RESTful API](https://en.wikipedia.org/wiki/Representational_state_transfer)** - 使用 [Gin](https://github.com/gin-gonic/gin) 框架实现 HTTP API
- **数据库支持** - [GORM](https://gorm.io) 支持 MySQL，可配置连接池
- **缓存支持** - [Redis](https://en.wikipedia.org/wiki/Redis) 缓存，带健康检查
- **增强缓存功能**:
  - 负缓存防止缓存穿透
  - 分布式锁保证缓存一致性
  - LRU 淘汰的键追踪
  - 管道支持批量操作
- **日志系统** - [Zap](https://go.uber.org/zap) 结构化日志
- **配置管理** - [Viper](https://github.com/spf13/viper) 灵活配置
- **优雅关闭** - 支持服务优雅启动和关闭
- **事务支持** - 真实事务工厂，支持回滚/提交
- **缓存-数据库事务一致性** - 缓存操作纳入数据库事务

### 监控
- **Prometheus 指标** - HTTP 请求、数据库操作、缓存性能、领域事件
- **健康检查** - Redis 和数据库健康监控

## 项目结构

```
.
├── adapter/                # 适配器层 - 外部系统交互
│   ├── dependency/        # 依赖注入 (Wire)
│   ├── job/               # 定时任务
│   └── repository/        # 数据仓储适配器
│       ├── mysql/         # MySQL 实现
│       ├── postgre/       # PostgreSQL 实现
│       └── redis/         # Redis 实现
│           ├── client.go           # Redis 客户端
│           ├── enhanced_cache.go  # 增强缓存
│           └── example_cache.go   # 示例缓存仓储
├── api/                   # API 层
│   ├── dto/              # 数据传输对象
│   ├── error_code/       # 错误码
│   ├── grpc/             # gRPC 处理器
│   ├── middleware/       # HTTP 中间件
│   └── http/             # HTTP 处理器
├── application/           # 应用层
│   ├── core/            # 核心接口
│   └── example/         # 用例
├── cmd/                  # 入口点
├── config/               # 配置
├── domain/               # 领域层
│   ├── aggregate/       # 聚合
│   ├── dto/              # 领域 DTO
│   ├── event:            # 领域事件
│   ├── model/            # 领域模型
│   ├── repo/             # 仓储接口
│   ├── service/          # 领域服务
│   └── vo/               # 值对象
├── tests/                # 集成测试
└── util/                 # 工具
    ├── errors/           # 错误处理
    ├── log/               # 日志
    ├── metrics/           # Prometheus 指标
    └── transaction/       # 事务工具
```

## 架构分层

### 领域层
- **模型**: 领域实体 (Example)
- **仓储接口**: IExampleRepo, IExampleCacheRepo
- **服务**: 带事务支持的 ExampleService
- **事件**: ExampleCreatedEvent, ExampleUpdatedEvent, ExampleDeletedEvent

### 应用层
- 用例: 创建、获取、更新、删除、按名称查询

### 适配器层
- **MySQL**: 基于 GORM 的实现，支持事务
- **Redis**: 带连接池的客户端，支持 LRU、分布式锁、负缓存
- **事件总线**: 内存和异步事件处理

### API 层
- 带请求验证的 HTTP 处理器
- DTO 转换
- 错误处理中间件

## 核心功能

### 事务工厂
```go
// 真实事务带自动回滚
txFactory := repository.NewTransactionFactory(map[repository.StoreType]any{
    repository.MySQLStore: clients.MySQL,
})
service := service.NewExampleService(repo, cacheRepo, txFactory)
```

### 增强缓存管道
```go
// 批量操作的管道
pipe := client.Pipeline()
pipe.Set(ctx, key1, value1, ttl1)
pipe.Set(ctx, key2, value2, ttl2)
cmds, err := pipe.Exec(ctx)
```

### 缓存指标
```go
metrics.RecordCacheHit("example_cache", "GetByID_hit")
metrics.RecordCacheHit("example_cache", "GetByID_miss")
```

### LRU 键追踪缓存
```go
cache := redis.NewEnhancedCache(client, redis.CacheOptions{
    EnableKeyTracking: true,
    MaxTrackedKeys:   100,
})
```

## 错误处理

```go
// 创建领域错误
return errors.New(errors.ErrorTypeNotFound, "entity not found")

// 包装错误并添加上下文
return errors.Wrapf(err, errors.ErrorTypePersistence, "failed to save entity %d", id)

// 使用 errors.Is 进行比较
if errors.Is(err, redis.Nil) {
    return ErrCacheMiss
}
```

## 快速开始

### 前置条件
- Go 1.21+
- MySQL
- Redis

### 配置

编辑 `config/config.yaml`:

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

### 运行

```bash
# 运行测试
go test ./... -short

# 运行应用程序
go run cmd/main.go
```

## 许可证

[Apache 2.0 License](LICENSE)
