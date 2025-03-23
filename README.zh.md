# Go Hexagonal Architecture

欢迎阅读我的[博客文章](https://blog.ranchocooper.com/2025/03/20/go-hexagonal/)

![Hexagonal Architecture](https://github.com/Sairyss/domain-driven-hexagon/raw/master/assets/images/DomainDrivenHexagon.png)

## 项目概述

本项目是基于[六边形架构](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software))（Hexagonal Architecture）和[领域驱动设计](https://en.wikipedia.org/wiki/Domain-driven_design)（Domain-Driven Design）的Go微服务框架。它提供了清晰的项目结构和设计模式，帮助开发者构建可维护、可测试和可扩展的应用程序。

六边形架构（也称为[端口与适配器架构](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software))）将应用程序分为内部和外部部分，通过明确定义的接口（端口）和实现（适配器）实现[关注点分离](https://en.wikipedia.org/wiki/Separation_of_concerns)和[依赖倒置原则](https://en.wikipedia.org/wiki/Dependency_inversion_principle)。这种架构将业务逻辑与技术实现细节解耦，便于单元测试和功能扩展。

## 核心特性

### 架构设计
- **[领域驱动设计 (DDD)](https://en.wikipedia.org/wiki/Domain-driven_design)** - 通过[聚合](https://en.wikipedia.org/wiki/Domain-driven_design)、[实体](https://en.wikipedia.org/wiki/Entity)和[值对象](https://en.wikipedia.org/wiki/Value_object)等概念组织业务逻辑
- **[六边形架构](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software))** - 将应用程序分为领域层、应用层和适配器层
- **[依赖注入](https://en.wikipedia.org/wiki/Dependency_injection)** - 使用[Wire](https://github.com/google/wire)进行依赖注入，提高代码可测试性和灵活性
- **[资源库模式](https://en.wikipedia.org/wiki/Repository_pattern)** - 抽象数据访问层，支持事务处理
- **[领域事件](https://en.wikipedia.org/wiki/Domain-driven_design)** - 实现[事件驱动架构](https://en.wikipedia.org/wiki/Event-driven_architecture)，支持系统组件之间的松耦合通信
- **[CQRS模式](https://en.wikipedia.org/wiki/Command_Query_Responsibility_Segregation)** - 命令和查询责任分离，优化读写操作
- **[接口驱动设计](https://en.wikipedia.org/wiki/Interface-based_programming)** - 使用接口定义服务契约，实现依赖倒置原则

### 技术实现
- **[RESTful API](https://en.wikipedia.org/wiki/Representational_state_transfer)** - 使用[Gin](https://github.com/gin-gonic/gin)框架实现HTTP API
- **数据库支持** - 集成[GORM](https://gorm.io)，支持[MySQL](https://en.wikipedia.org/wiki/MySQL)、[PostgreSQL](https://en.wikipedia.org/wiki/PostgreSQL)等数据库
- **缓存支持** - 集成[Redis](https://en.wikipedia.org/wiki/Redis)缓存，具有全面的错误处理机制和专用错误定义，支持健康检查功能监控缓存可用性
- **MongoDB支持** - 集成MongoDB文档存储
- **日志系统** - 使用[Zap](https://go.uber.org/zap)进行高性能日志记录
- **配置管理** - 使用[Viper](https://github.com/spf13/viper)进行灵活的配置管理
- **[优雅关闭](https://en.wikipedia.org/wiki/Graceful_exit)** - 支持服务优雅启动和关闭
- **[单元测试](https://en.wikipedia.org/wiki/Unit_testing)** - 使用[go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)、[redismock](https://github.com/go-redis/redismock)和[testify/mock](https://github.com/stretchr/testify)实现全面的测试覆盖，支持请求/响应数据传输对象的增强测试
- **事务支持** - 提供无操作事务实现，简化服务层与仓储层的交互，同时支持完整的模拟事务实现和生命周期钩子（Begin、Commit和Rollback），便于单元测试

### 开发工具链
- **代码质量** - 集成[Golangci-lint](https://github.com/golangci/golangci-lint)进行代码质量检查
- **提交标准** - 使用[Commitlint](https://github.com/conventional-changelog/commitlint)确保Git提交消息遵循约定
- **预提交钩子** - 使用[Pre-commit](https://pre-commit.com)进行代码检查和格式化
- **[CI/CD](https://en.wikipedia.org/wiki/CI/CD)** - 集成[GitHub Actions](https://github.com/features/actions)进行持续集成和部署

## 项目结构

```
.
├── adapter/                # 适配器层 - 与外部系统交互
│   ├── amqp/               # 消息队列适配器
│   ├── dependency/         # 依赖注入配置
│   │   └── wire.go         # Wire DI设置与接口绑定
│   ├── job/                # 定时任务适配器
│   └── repository/         # 数据存储适配器
│       ├── mysql/          # MySQL实现
│       │   └── entity/     # 数据库实体和仓储实现
│       ├── postgre/        # PostgreSQL实现
│       ├── mongo/          # MongoDB实现
│       └── redis/          # Redis实现
├── api/                    # API层 - 处理HTTP请求和响应
│   ├── dto/                # API的数据传输对象
│   ├── error_code/         # 错误码定义
│   ├── grpc/               # gRPC API处理器
│   └── http/               # HTTP API处理器
│       ├── handle/         # 使用领域接口的请求处理器
│       ├── middleware/     # HTTP中间件
│       ├── paginate/       # 分页处理
│       └── validator/      # 请求验证
├── application/            # 应用层 - 协调领域对象实现用例
│   ├── core/               # 核心接口和基础实现
│   │   └── interfaces.go   # UseCase和UseCaseHandler接口
│   └── example/            # 示例用例实现
│       ├── create_example.go     # 创建示例用例
│       ├── delete_example.go     # 删除示例用例
│       ├── get_example.go        # 获取示例用例
│       ├── update_example.go     # 更新示例用例
│       └── find_example_by_name.go # 按名称查找示例用例
├── cmd/                    # 命令行入口点
│   └── http_server/        # HTTP服务器启动
├── config/                 # 配置文件和管理
├── domain/                 # 领域层 - 核心业务逻辑
│   ├── aggregate/          # 聚合（DDD概念）
│   ├── event/              # 领域事件和事件总线接口
│   │   ├── event_bus.go    # 事件总线接口
│   │   └── handlers.go     # 事件处理器接口
│   ├── model/              # 领域模型（纯业务实体）
│   ├── repo/               # 资源库接口
│   │   └── transaction.go  # 事务接口
│   ├── service/            # 领域服务及其接口
│   │   ├── example.go      # ExampleService实现
│   │   └── interfaces.go   # 服务接口（IExampleService等）
│   └── vo/                 # 值对象（DDD概念）
├── tests/                  # 集成测试
│   ├── migrations/         # 测试数据库迁移
│   ├── mysql.go            # MySQL测试助手
│   ├── postgres.go         # PostgreSQL测试助手
│   ├── redis.go            # Redis测试助手
│   └── *_test.go           # 测试文件
└── util/                   # 工具函数
    ├── clean_arch/         # 架构检查工具
    └── log/                # 日志工具
```

### 关键架构元素

此结构强制执行六边形架构原则：

1. **接口-实现分离**：
   - 领域层定义接口（端口）
   - 适配器层提供实现（适配器）
   - 依赖流向内部，外层依赖于内层

2. **依赖倒置**：
   - 高层模块（领域/应用层）依赖于抽象
   - 低层模块（适配器）实现这些抽象
   - 所有依赖通过接口注入

3. **领域中心设计**：
   - 领域模型是不含技术关注点的纯业务实体
   - 资源库接口声明领域需要什么
   - 服务接口定义业务操作

4. **清晰边界**：
   - 每一层都有明确的职责和依赖
   - 数据转换发生在层边界
   - 实现细节不在层间泄漏

## 架构层次

### 领域层
领域层是应用程序的核心，包含业务逻辑和规则。它独立于其他层，不依赖于任何外部组件。

- **模型(Models)**: 领域实体和值对象
  - `Example`: 示例实体，包含基本属性如ID、名称、别名等

- **资源库接口(Repository Interfaces)**: 定义数据访问接口
  - `IExampleRepo`: 示例资源库接口，定义了创建、读取、更新、删除等操作
  - `IExampleCacheRepo`: 示例缓存接口，定义了健康检查方法
  - `Transaction`: 事务接口，支持事务的开始、提交和回滚

- **领域服务(Domain Services)**: 处理跨实体的业务逻辑
  - `IExampleService`: 服务接口，定义示例相关操作的契约
  - `ExampleService`: 示例服务接口的实现，处理示例实体的业务逻辑

- **领域事件(Domain Events)**: 定义领域内的事件
  - `ExampleCreatedEvent`: 示例创建事件
  - `ExampleUpdatedEvent`: 示例更新事件
  - `ExampleDeletedEvent`: 示例删除事件

### 应用层
应用层协调领域对象完成特定应用任务。它依赖于领域接口而非具体实现，遵循依赖倒置原则。

- **用例(Use Cases)**: 定义应用功能
  - `CreateExampleUseCase`: 创建示例用例
  - `GetExampleUseCase`: 获取示例用例
  - `UpdateExampleUseCase`: 更新示例用例
  - `DeleteExampleUseCase`: 删除示例用例
  - `FindExampleByNameUseCase`: 按名称查找示例用例

- **命令和查询(Commands and Queries)**: 实现CQRS模式
  - 每个用例定义了Input和Output结构，分别代表命令/查询输入和结果

- **事件处理器(Event Handlers)**: 处理领域事件
  - `LoggingEventHandler`: 日志事件处理器，记录所有事件
  - `ExampleEventHandler`: 示例事件处理器，处理与示例相关的事件

### 适配器层
适配器层实现与外部系统的交互，如数据库和消息队列。

- **资源库实现(Repository Implementation)**: 实现数据访问接口
  - `EntityExample`: MySQL实现的示例资源库
  - `NoopTransaction`: 无操作事务实现，简化测试
  - `MySQL`: MySQL连接和事务管理
  - `Redis`: Redis连接和基本操作

- **消息队列适配器(Message Queue Adapters)**: 实现消息发布和订阅
  - 支持Kafka等消息队列的集成

- **定时任务(Scheduled Tasks)**: 实现定时任务
  - 基于cron的任务调度系统

### API层
API层处理HTTP请求和响应，作为应用程序的入口点。

- **控制器(Controllers)**: 处理HTTP请求
  - `CreateExample`: 创建示例API
  - `GetExample`: 获取示例API
  - `UpdateExample`: 更新示例API
  - `DeleteExample`: 删除示例API
  - `FindExampleByName`: 按名称查找示例API

- **中间件(Middleware)**: 实现横切关注点
  - 国际化支持
  - CORS支持
  - 请求ID跟踪
  - 请求日志

- **数据传输对象(DTOs)**: 定义请求和响应数据结构
  - `CreateExampleReq`: 创建示例请求
  - `UpdateExampleReq`: 更新示例请求
  - `DeleteExampleReq`: 删除示例请求
  - `GetExampleReq`: 获取示例请求

## 依赖注入

本项目使用Google Wire进行依赖注入，组织依赖关系如下：

```go
// 初始化服务
func InitializeServices(ctx context.Context) (*service.Services, error) {
    wire.Build(
        // 资源库依赖
        entity.NewExample,
        wire.Bind(new(repo.IExampleRepo), new(*entity.EntityExample)),

        // 事件总线依赖
        provideEventBus,
        wire.Bind(new(event.EventBus), new(*event.InMemoryEventBus)),

        // 服务依赖
        provideExampleService,
        wire.Bind(new(service.IExampleService), new(*service.ExampleService)),
        provideServices,
    )
    return nil, nil
}

// 提供事件总线
func provideEventBus() *event.InMemoryEventBus {
    eventBus := event.NewInMemoryEventBus()

    // 注册事件处理器
    loggingHandler := event.NewLoggingEventHandler()
    exampleHandler := event.NewExampleEventHandler()
    eventBus.Subscribe(loggingHandler)
    eventBus.Subscribe(exampleHandler)

    return eventBus
}

// 提供示例服务
func provideExampleService(repo repo.IExampleRepo, eventBus event.EventBus) *service.ExampleService {
    exampleService := service.NewExampleService(repo)
    exampleService.EventBus = eventBus
    return exampleService
}

// 提供服务容器
func provideServices(exampleService service.IExampleService, eventBus event.EventBus) *service.Services {
    return service.NewServices(exampleService, eventBus)
}
```

## 领域事件

领域事件用于系统组件之间的通信，实现松耦合的事件驱动架构：

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

应用层用例实现命令和查询责任分离(CQRS)模式，并依赖领域接口而非具体实现：

```go
// 使用接口依赖的用例
type CreateUseCase struct {
    exampleService service.IExampleService
}

// 创建示例用例
func (uc *CreateUseCase) Execute(ctx context.Context, input dto.CreateExampleReq) (*dto.CreateExampleResp, error) {
    // 创建事务保证原子操作
    tx, err := repository.NewTransaction(ctx, repository.MySQLStore, nil)
    if err != nil {
        return nil, fmt.Errorf("创建事务失败: %w", err)
    }
    defer tx.Rollback()

    // 将DTO转换为领域模型
    example := &model.Example{
        Name:  input.Name,
        Alias: input.Alias,
    }

    // 通过接口调用领域服务
    createdExample, err := uc.exampleService.Create(ctx, example)
    if err != nil {
        return nil, fmt.Errorf("创建示例失败: %w", err)
    }

    // 提交事务
    if err = tx.Commit(); err != nil {
        return nil, fmt.Errorf("提交事务失败: %w", err)
    }

    // 将领域模型转换为DTO
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

## 事务管理

本项目实现了事务接口和无操作事务，支持不同的事务管理策略：

```go
// 事务接口
type Transaction interface {
    Begin() error
    Commit() error
    Rollback() error
    Conn(ctx context.Context) any
}

// 无操作事务实现
type NoopTransaction struct {
    conn any
}

// 在服务中使用事务
func (s *ExampleService) Create(ctx context.Context, example *model.Example) (*model.Example, error) {
    // 创建一个无操作事务
    tr := repo.NewNoopTransaction(s.Repository)

    createdExample, err := s.Repository.Create(ctx, tr, example)
    // ...
}
```

## 数据映射和转换

本项目在不同层次间实现了清晰的数据映射和转换，使用[jinzhu/copier](https://github.com/jinzhu/copier)库进行高效的对象复制：

```go
// 使用copier进行实体到模型的转换
func (e EntityExample) ToModel() *model.Example {
    model := &model.Example{}
    copier.Copy(model, e)
    return model
}

// 使用copier进行模型到实体的转换
func (e *EntityExample) FromModel(m *model.Example) {
    copier.Copy(e, m)
}

// 批量将实体转换为模型
func EntitiesToModels(entities []EntityExample) []*model.Example {
    result := make([]*model.Example, len(entities))
    for i, entity := range entities {
        result[i] = entity.ToModel()
    }
    return result
}

// 数据传输对象到模型的转换
func (req *CreateExampleReq) ToModel() *model.Example {
    model := &model.Example{}
    copier.Copy(model, req)
    return model
}

// 模型到响应数据传输对象的转换
func ModelToResponse(m *model.Example) *ExampleResponse {
    if m == nil {
        return nil
    }

    resp := &ExampleResponse{}
    copier.Copy(resp, m)

    // 复制后格式化时间字段
    resp.CreatedAt = m.CreatedAt.Format(time.RFC3339)
    resp.UpdatedAt = m.UpdatedAt.Format(time.RFC3339)

    return resp
}
```

使用copier库的优势：
- 简化相似结构体之间的转换
- 自动复制同名且类型兼容的字段
- 支持嵌套结构体的深度复制
- 减少对象转换的样板代码

这些转换保持了不同层次之间的清晰分离：
- 数据库实体（在适配器层）
- 领域模型（在领域层）
- 数据传输对象（在API层）

这种方法允许每一层拥有自己的数据表示形式，针对其特定职责进行优化。

## 项目优化

本项目最近经过以下优化：

### 1. 统一API版本
- **问题**：项目同时存在v1和v2 API版本，导致代码重复和维护困难
- **解决方案**：
  - 统一API路由，将所有API放在`/api`路径下
  - 保留`/v2`路径以向后兼容
  - 使用应用层用例处理所有请求，逐步淘汰直接调用领域服务

### 2. 增强依赖注入
- **问题**：Wire依赖注入配置存在重复绑定问题，导致生成失败
- **解决方案**：
  - 重构`wire.go`文件，移除重复绑定定义
  - 使用provider函数替代直接绑定
  - 添加事件处理器注册逻辑

### 3. 消除全局变量
- **问题**：项目使用全局变量存储服务实例，违反依赖注入原则
- **解决方案**：
  - 移除全局变量`service.ExampleSvc`和`service.EventBus`的使用
  - 通过依赖注入传递服务实例
  - 启动HTTP服务器时使用依赖注入初始化服务

### 4. 改进应用层集成
- **问题**：应用层用例未充分利用，HTTP服务器默认不启用应用层
- **解决方案**：
  - 默认启用应用层用例
  - 使用用例工厂创建和管理用例
  - 实现更清晰的错误处理和响应映射

## 最近的优化

本项目最近经过以下优化：

1. **环境变量支持**：
   - 添加配置文件的环境变量覆盖功能，使应用在容器化部署中更灵活
   - 使用统一前缀(APP_)和层次结构(如APP_MYSQL_HOST)组织环境变量

2. **统一错误处理**：
   - 实现应用级错误类型系统，支持不同类型的错误(验证错误、未找到、未授权等)
   - 添加统一错误响应处理，将内部错误映射到合适的HTTP状态码
   - 改进错误日志记录，确保所有意外错误被正确记录

3. **请求日志中间件**：
   - 添加详细的请求日志中间件，记录请求方法、路径、状态码、延迟等信息
   - 在调试模式下，可以记录请求和响应体，帮助开发者排查问题
   - 智能识别内容类型，避免记录二进制内容

4. **请求ID跟踪**：
   - 为每个请求生成唯一的请求ID，便于在分布式系统中跟踪
   - 在响应头中返回请求ID，供客户端参考
   - 在日志中包含请求ID，关联同一请求的多个日志条目

5. **优雅关闭**：
   - 实现服务器的优雅关闭机制，确保所有运行中的请求在关闭前完成
   - 添加关闭超时设置，防止关闭过程无限期挂起
   - 改进信号处理，支持SIGINT和SIGTERM信号

6. **国际化支持**：
   - 添加翻译中间件，支持多语言验证错误消息
   - 根据Accept-Language头自动选择适当的语言

7. **CORS支持**：
   - 添加CORS中间件，处理跨域请求
   - 配置允许的来源、方法、头和凭据

8. **调试工具**：
   - 集成pprof性能分析工具，用于诊断生产环境中的性能问题
   - 可通过配置文件启用或禁用

这些优化使项目更加健壮、可维护，并提供更好的开发体验。

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

- **gRPC支持** - 添加gRPC服务实现
- **监控集成** - 集成Prometheus监控

## 参考资料

- **架构**
  - [Freedom DDD Framework](https://github.com/8treenet/freedom)
  - [Hexagonal Architecture in Go](https://medium.com/@matiasvarela/hexagonal-architecture-in-go-cfd4e436faa3)
  - [Dependency Injection in A Nutshell](https://appliedgo.net/di/)
- **项目标准**
  - [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0)
  - [Improving Your Go Project With pre-commit hooks](https://goangle.medium.com/golang-improving-your-go-project-with-pre-commit-hooks-a265fad0e02f)
- **代码参考**
  - [Go CleanArch](https://github.com/roblaszczak/go-cleanarch)
