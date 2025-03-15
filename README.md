# Go Hexagonal Architecture

![Hexagonal Architecture](https://github.com/Sairyss/domain-driven-hexagon/raw/master/assets/images/DomainDrivenHexagon.png)

## 项目概述

本项目是一个基于六边形架构（Hexagonal Architecture）和领域驱动设计（Domain-Driven Design）的Go语言微服务框架。它提供了一个清晰的项目结构和设计模式，帮助开发者构建可维护、可测试和可扩展的应用程序。

六边形架构（也称为端口和适配器架构）将应用程序分为内部和外部两部分，通过定义明确的接口（端口）和实现（适配器）来实现关注点分离和依赖倒置原则。这种架构使得业务逻辑与技术实现细节解耦，便于单元测试和功能扩展。

## 核心特性

### 架构设计
- **领域驱动设计 (DDD)** - 通过聚合、实体、值对象等概念组织业务逻辑
- **六边形架构** - 将应用程序分为领域、应用和适配器层
- **依赖注入** - 使用Wire实现依赖注入，提高代码的可测试性和灵活性
- **仓储模式** - 抽象数据访问层，支持事务处理
- **领域事件** - 实现事件驱动架构，支持系统内部组件的松耦合通信
- **CQRS模式** - 命令和查询职责分离，优化读写操作

### 技术实现
- **RESTful API** - 使用Gin框架实现HTTP API
- **数据库支持** - 集成GORM，支持MySQL、PostgreSQL等数据库
- **缓存支持** - 集成Redis缓存
- **日志系统** - 使用Zap实现高性能日志记录
- **配置管理** - 使用Viper实现灵活的配置管理
- **优雅关闭** - 支持服务的优雅启动和关闭
- **单元测试** - 使用go-sqlmock和redismock实现数据库和缓存的模拟测试

### 开发工具链
- **代码质量** - 集成Golangci-lint进行代码质量检查
- **提交规范** - 使用Commitlint确保Git提交信息符合规范
- **预提交钩子** - 使用Pre-commit进行代码检查和格式化
- **CI/CD** - 集成GitHub Actions实现持续集成和部署

## 项目结构

```
.
├── adapter/                # 适配器层 - 实现与外部系统的交互
│   ├── amqp/               # 消息队列适配器
│   ├── dependency/         # 依赖注入配置
│   ├── job/                # 定时任务适配器
│   └── repository/         # 数据仓库适配器
│       ├── mysql/          # MySQL实现
│       ├── postgre/        # PostgreSQL实现
│       └── redis/          # Redis实现
├── api/                    # API层 - 处理HTTP请求和响应
│   ├── dto/                # 数据传输对象
│   ├── error_code/         # 错误码定义
│   ├── grpc/               # gRPC API
│   └── http/               # HTTP API
│       ├── handle/         # 请求处理器
│       ├── middleware/     # 中间件
│       ├── paginate/       # 分页处理
│       └── validator/      # 请求验证
├── application/            # 应用层 - 协调领域对象完成用例
│   ├── core/               # 核心接口和错误定义
│   └── example/            # 示例用例实现
├── cmd/                    # 命令行入口
│   └── http_server/        # HTTP服务器启动
├── config/                 # 配置文件和配置管理
├── domain/                 # 领域层 - 核心业务逻辑
│   ├── aggregate/          # 聚合
│   ├── event/              # 领域事件
│   ├── model/              # 领域模型
│   ├── repo/               # 仓库接口
│   ├── service/            # 领域服务
│   └── vo/                 # 值对象
├── tests/                  # 集成测试
└── util/                   # 工具函数
    ├── clean_arch/         # 架构检查工具
    └── log/                # 日志工具
```

## 架构分层

### 领域层 (Domain Layer)
领域层是应用程序的核心，包含业务逻辑和规则。它独立于其他层，不依赖于任何外部组件。

- **模型 (Model)**: 领域实体和值对象
- **仓库接口 (Repository Interface)**: 定义数据访问接口
- **领域服务 (Service)**: 处理跨实体的业务逻辑
- **领域事件 (Event)**: 定义领域内的事件

### 应用层 (Application Layer)
应用层协调领域对象完成特定的应用任务。它依赖于领域层，但不包含业务规则。

- **用例 (Use Cases)**: 定义应用程序的功能
- **命令和查询 (Commands & Queries)**: 实现CQRS模式
- **事件处理器 (Event Handlers)**: 处理领域事件

### 适配器层 (Adapter Layer)
适配器层实现与外部系统的交互，如数据库、消息队列等。

- **仓库实现 (Repository Implementation)**: 实现数据访问接口
- **消息队列适配器 (AMQP Adapter)**: 实现消息发布和订阅
- **定时任务 (Job)**: 实现定时任务

### API层 (API Layer)
API层处理HTTP请求和响应，是应用程序的入口点。

- **控制器 (Controllers)**: 处理HTTP请求
- **中间件 (Middleware)**: 实现横切关注点
- **数据传输对象 (DTOs)**: 定义请求和响应的数据结构

## 依赖注入

本项目使用Google Wire实现依赖注入，通过以下方式组织依赖关系：

```go
// 初始化服务
func InitializeServices(ctx context.Context) (*service.Services, error) {
    // 创建仓库
    entityExample := entity.NewExample()

    // 创建事件总线和处理器
    inMemoryEventBus := event.NewInMemoryEventBus()
    loggingHandler := event.NewLoggingEventHandler()
    exampleHandler := event.NewExampleEventHandler()

    // 注册事件处理器
    inMemoryEventBus.Subscribe(loggingHandler)
    inMemoryEventBus.Subscribe(exampleHandler)

    // 创建服务
    exampleService := service.NewExampleService(ctx)
    exampleService.Repository = entityExample
    exampleService.EventBus = inMemoryEventBus

    // 创建服务容器
    services := service.NewServices(exampleService, inMemoryEventBus)

    return services, nil
}
```

## 领域事件

领域事件用于在系统内部组件之间进行通信，实现松耦合的事件驱动架构：

```go
// 发布事件
evt := event.NewExampleCreatedEvent(example.Id, example.Name, example.Alias)
e.EventBus.Publish(ctx, evt)

// 处理事件
func (h *ExampleEventHandler) HandleEvent(ctx context.Context, event Event) error {
    switch event.EventName() {
    case ExampleCreatedEventName:
        return h.handleExampleCreated(ctx, event)
    // ...
    }
    return nil
}
```

## 应用层用例

应用层用例实现了命令和查询职责分离 (CQRS) 模式：

```go
// 创建示例用例
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

## 使用指南

### 环境准备

使用Docker启动MySQL：
```bash
docker run --name mysql-local \
  -e MYSQL_ROOT_PASSWORD=mysqlroot \
  -e MYSQL_DATABASE=go-hexagonal \
  -e MYSQL_USER=user \
  -e MYSQL_PASSWORD=mysqlroot \
  -p 3306:3306 \
  -d mysql:latest
```

### 开发工具安装

```bash
# 安装开发工具
make init && make precommit.rehook
```

或手动安装：

```bash
# 安装pre-commit
brew install pre-commit
# 安装golangci-lint
brew install golangci-lint
# 安装commitlint
npm install -g @commitlint/cli @commitlint/config-conventional
# 添加commitlint配置
echo "module.exports = {extends: ['@commitlint/config-conventional']}" > commitlint.config.js
# 添加pre-commit钩子
make precommit.rehook
```

### 运行项目

```bash
# 运行项目
go run cmd/main.go
```

### 测试

```bash
# 运行测试
go test ./...
```

## 扩展计划

- **依赖注入改进** - 完善Wire依赖注入配置
- **HTTP处理改进** - 优化HTTP请求处理实现
- **领域事件增强** - 完善领域事件机制
- **gRPC支持** - 添加gRPC服务实现
- **热加载配置** - 实现配置热加载
- **监控集成** - 集成Prometheus监控
- **消息队列集成** - 集成Kafka等消息队列

## 参考资料

- **架构**
  - [Freedom DDD Framework](https://github.com/8treenet/freedom)
  - [Hexagonal Architecture in Go](https://medium.com/@matiasvarela/hexagonal-architecture-in-go-cfd4e436faa3)
  - [Dependency Injection in A Nutshell](https://appliedgo.net/di/)
- **项目规范**
  - [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0)
  - [Improving Your Go Project With pre-commit hooks](https://goangle.medium.com/golang-improving-your-go-project-with-pre-commit-hooks-a265fad0e02f)
- **代码参考**
  - [Go CleanArch](https://github.com/roblaszczak/go-cleanarch)
