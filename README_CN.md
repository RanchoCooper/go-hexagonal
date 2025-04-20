# Go 六边形架构

欢迎访问我的[博客文章](https://blog.ranchocooper.com/2025/03/20/go-hexagonal/)

![六边形架构](https://github.com/Sairyss/domain-driven-hexagon/raw/master/assets/images/DomainDrivenHexagon.png)

## 项目概述

本项目是一个基于[六边形架构](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software))和[领域驱动设计](https://en.wikipedia.org/wiki/Domain-driven_design)的 Go 微服务框架。它提供了清晰的项目结构和设计模式，帮助开发者构建可维护、可测试和可扩展的应用程序。

六边形架构（也称为[端口和适配器架构](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software))）将应用程序分为内部和外部部分，通过明确定义的接口（端口）和实现（适配器）实现[关注点分离](https://en.wikipedia.org/wiki/Separation_of_concerns)和[依赖倒置原则](https://en.wikipedia.org/wiki/Dependency_inversion_principle)。这种架构将业务逻辑与技术实现细节解耦，便于单元测试和功能扩展。

## 核心特性

### 架构设计
- **[领域驱动设计 (DDD)](https://en.wikipedia.org/wiki/Domain-driven_design)** - 通过[聚合](https://en.wikipedia.org/wiki/Domain-driven_design)、[实体](https://en.wikipedia.org/wiki/Entity)和[值对象](https://en.wikipedia.org/wiki/Value_object)等概念组织业务逻辑
- **[六边形架构](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software))** - 将应用分为领域层、应用层和适配器层
- **[依赖注入](https://en.wikipedia.org/wiki/Dependency_injection)** - 使用 [Wire](https://github.com/google/wire) 进行依赖注入，提高代码可测试性和灵活性
- **[仓储模式](https://en.wikipedia.org/wiki/Repository_pattern)** - 抽象数据访问层，支持事务
- **[领域事件](https://en.wikipedia.org/wiki/Domain-driven_design)** - 实现[事件驱动架构](https://en.wikipedia.org/wiki/Event-driven_architecture)，支持系统组件间的松耦合通信
- **[CQRS 模式](https://en.wikipedia.org/wiki/Command_Query_Responsibility_Segregation)** - 命令查询职责分离，优化读写操作
- **[接口驱动设计](https://en.wikipedia.org/wiki/Interface-based_programming)** - 使用接口定义服务契约，实现依赖倒置原则

### 技术实现
- **[RESTful API](https://en.wikipedia.org/wiki/Representational_state_transfer)** - 使用 [Gin](https://github.com/gin-gonic/gin) 框架实现 HTTP API
- **数据库支持** - 集成 [GORM](https://gorm.io)，支持 [MySQL](https://en.wikipedia.org/wiki/MySQL)、[PostgreSQL](https://en.wikipedia.org/wiki/PostgreSQL) 等数据库
- **缓存支持** - 集成 [Redis](https://en.wikipedia.org/wiki/Redis) 缓存，包含完整的错误处理、本地错误定义和健康检查实现
- **增强缓存** - 高级缓存特性，包括防止缓存穿透的负缓存、分布式锁和键追踪
- **MongoDB 支持** - 集成 MongoDB 用于文档存储
- **日志系统** - 使用 [Zap](https://go.uber.org/zap) 进行高性能日志记录，支持结构化上下文
- **配置管理** - 使用 [Viper](https://github.com/spf13/viper) 进行灵活的配置管理
- **[优雅关闭](https://en.wikipedia.org/wiki/Graceful_exit)** - 支持优雅的服务启动和关闭
- **[单元测试](https://en.wikipedia.org/wiki/Unit_testing)** - 使用 [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)、[redismock](https://github.com/go-redis/redismock) 和 [testify/mock](https://github.com/stretchr/testify) 进行全面的测试覆盖
- **事务支持** - 提供无操作事务实现，简化服务层与仓储层的交互
- **异步事件处理** - 支持异步事件处理，包含工作池、事件持久化和重放功能
- **监控和可观察性** - 集成 [Prometheus](https://prometheus.io) 进行指标收集，测量HTTP请求、数据库操作、缓存性能和领域事件，内置中间件跟踪请求性能

### 开发工具链
- **代码质量** - 集成 [Golangci-lint](https://github.com/golangci/golangci-lint) 进行代码质量检查
- **提交规范** - 使用 [Commitlint](https://github.com/conventional-changelog/commitlint) 确保 Git 提交信息符合规范
- **预提交钩子** - 使用 [Pre-commit](https://pre-commit.com) 进行代码检查和格式化
- **[CI/CD](https://en.wikipedia.org/wiki/CI/CD)** - 集成 [GitHub Actions](https://github.com/features/actions) 进行持续集成和部署

## 最近增强

### 统一错误处理
- 扩展错误处理，包含一致的错误类型和错误包装函数
- 支持结构化错误详情和 HTTP 状态码映射
- 错误比较功能，用于更健壮的错误检查

### 增强结构化日志
- 支持请求 ID、用户 ID 和追踪 ID 的上下文感知日志
- 一致的日志格式和级别管理
- 改进的调试能力，包含上下文信息

### 异步事件系统
- 基于工作池的事件处理，提高吞吐量
- 事件持久化和重放功能，提高可靠性
- 事件处理的优雅关闭支持

### 高级缓存特性
- 负缓存防止缓存穿透
- 分布式锁防止缓存雪崩
- 键追踪提高缓存命中率
- 缓存一致性机制确保数据完整性

### 双模式 HTTP 处理器
- 灵活的 HTTP 处理器，可同时适用于应用层工厂和直接域服务调用
- 支持在测试和简单用例中直接调用服务
- 改进的可测试性，通过更好的转换器集成优化请求/响应转换
- 当直接服务模式不可用时优雅降级到应用工厂模式
- 通过简化的模拟设置增强测试能力

### 全面的监控和可观察性
- 基于Prometheus的全应用层次指标收集
- HTTP请求跟踪，包括持续时间、状态码和错误率
- 数据库操作监控，包括查询持续时间和错误计数
- 事务性能指标，支持操作跟踪
- 缓存性能监控，支持命中率/未命中率分析
- 领域事件监控，提供业务流程洞察
- 可定制的指标端点，支持健康检查

## 项目结构

```
.
├── adapter/                # 适配器层 - 外部系统交互
│   ├── amqp/               # 消息队列适配器
│   ├── dependency/         # 依赖注入配置
│   │   └── wire.go         # Wire DI 设置和接口绑定
│   ├── job/                # 定时任务适配器
│   └── repository/         # 数据仓储适配器
│       ├── mysql/          # MySQL 实现
│       │   └── entity/     # 数据库实体和仓储实现
│       ├── postgre/        # PostgreSQL 实现
│       ├── mongo/          # MongoDB 实现
│       └── redis/          # Redis 实现
│           └── enhanced_cache.go  # 增强缓存实现
├── api/                    # API 层 - HTTP 请求和响应
│   ├── dto/                # API 数据传输对象
│   ├── error_code/         # 错误码定义
│   ├── grpc/               # gRPC API 处理器
│   ├── middleware/         # 全局中间件，包括指标收集
│   └── http/               # HTTP API 处理器
│       ├── handle/         # 使用领域接口的请求处理器
│       ├── middleware/     # HTTP 中间件
│       ├── paginate/       # 分页处理
│       └── validator/      # 请求验证
├── application/            # 应用层 - 用例协调领域对象
│   ├── core/               # 核心接口和基础实现
│   │   └── interfaces.go   # UseCase 和 UseCaseHandler 接口
│   └── example/            # 示例用例实现
│       ├── create_example.go     # 创建示例用例
│       ├── delete_example.go     # 删除示例用例
│       ├── get_example.go        # 获取示例用例
│       ├── update_example.go     # 更新示例用例
│       └── find_example_by_name.go # 按名称查找示例用例
├── cmd/                    # 命令行入口点
│   └── main.go             # 主应用入口点
├── config/                 # 配置管理
│   ├── config.go           # 配置结构和加载
│   └── config.yaml         # 配置文件
├── domain/                 # 领域层 - 核心业务逻辑
│   ├── aggregate/          # 领域聚合
│   ├── dto/                # 领域数据传输对象
│   ├── event/              # 领域事件
│   ├── model/              # 领域模型
│   ├── repo/               # 仓储接口
│   ├── service/            # 领域服务
│   └── vo/                 # 值对象
└── tests/                  # 测试工具和示例
    ├── migrations/         # 测试数据库迁移
    ├── mysql.go            # MySQL 测试工具
    ├── postgresql.go       # PostgreSQL 测试工具
    └── redis.go            # Redis 测试工具
```

## 架构设计原则

### 分层设计
1. **领域层** (`domain/`)
   - 包含核心业务逻辑和规则
   - 定义领域模型、聚合和值对象
   - 声明仓储接口和领域服务
   - 独立于外部关注点

2. **应用层** (`application/`)
   - 实现用例并协调领域对象
   - 处理事务边界
   - 协调领域对象和外部系统
   - 不包含业务规则

3. **适配器层** (`adapter/`)
   - 实现领域层和应用层定义的接口
   - 处理外部关注点（数据库、HTTP、消息）
   - 提供端口的具体实现
   - 包含技术细节和框架

4. **API 层** (`api/`)
   - 处理 HTTP/gRPC 请求和响应
   - 管理 DTO 和领域对象之间的数据转换
   - 实现 API 特定的验证和错误处理
   - 提供 API 文档和版本控制

### 设计模式和原则

1. **依赖倒置**
   - 高层模块定义接口
   - 低层模块实现接口
   - 依赖指向领域层

2. **接口隔离**
   - 接口针对用例特定
   - 客户端只依赖其使用的方法
   - 防止接口污染

3. **单一职责**
   - 每个组件只有一个变更原因
   - 清晰的关注点分离
   - 专注和可维护的代码

4. **开闭原则**
   - 对扩展开放
   - 对修改关闭
   - 通过新实现添加功能

### 测试策略

1. **单元测试**
   - 隔离测试领域逻辑
   - 模拟外部依赖
   - 快速和可靠的测试

2. **集成测试**
   - 测试适配器实现
   - 验证外部系统交互
   - 数据库和缓存测试

3. **端到端测试**
   - 测试完整用例
   - 验证系统行为
   - API 契约测试

## 依赖注入

本项目使用 Google Wire 进行依赖注入，组织依赖关系如下：

```go
// 初始化服务
func InitializeServices(ctx context.Context) (*service.Services, error) {
    wire.Build(
        // 仓储依赖
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

本项目支持同步和异步事件处理：

### 同步事件处理
```go
// 同步发布事件
err := eventBus.Publish(ctx, event.NewExampleCreatedEvent(example.ID, example.Name))
```

### 异步事件处理
```go
// 配置异步事件总线
config := event.DefaultAsyncEventBusConfig()
config.QueueSize = 1000
config.WorkerCount = 10
asyncEventBus := event.NewAsyncEventBus(config)

// 异步发布事件
err := asyncEventBus.Publish(ctx, event.NewExampleCreatedEvent(example.ID, example.Name))

// 优雅关闭
err := asyncEventBus.Close(5 * time.Second)
```

## 增强缓存

增强缓存系统提供了强大的缓存功能：

```go
// 使用默认选项创建增强缓存
cache := redis.NewEnhancedCache(redisClient, redis.DefaultCacheOptions())

// 尝试获取值，如果缺失则自动加载
var result MyData
err := cache.TryGetSet(ctx, "key:123", &result, 30*time.Minute, func() (interface{}, error) {
    // 仅在缓存中没有键时执行
    return fetchDataFromDatabase()
})

// 使用分布式锁防止并发操作
err := cache.WithLock(ctx, "lock:resource", func() error {
    // 这段代码受分布式锁保护
    return updateSharedResource()
})
```

## 错误处理

错误系统提供了一致的方式处理和传播错误：

```go
// 创建领域错误
if entity == nil {
    return errors.New(errors.ErrorTypeNotFound, "找不到实体")
}

// 包装错误并添加上下文
if err := repo.Save(entity); err != nil {
    return errors.Wrapf(err, errors.ErrorTypePersistence, "保存实体 %d 失败", entity.ID)
}

// 检查错误类型
if errors.IsNotFoundError(err) {
    // 处理未找到情况
}
```

## 结构化日志

日志系统支持上下文感知的结构化日志：

```go
// 创建日志上下文
logCtx := log.NewLogContext().
    WithRequestID(requestID).
    WithUserID(userID).
    WithOperation("CreateEntity")

// 带上下文的日志记录
logger.InfoContext(logCtx, "正在创建新实体",
    zap.Int("entity_id", entity.ID),
    zap.String("entity_name", entity.Name))
```

## 项目改进

本项目最近进行了以下改进：

### 1. 统一 API 版本
- **问题**：项目同时存在 v1 和 v2 API 版本，导致代码重复和维护困难
- **解决方案**：
  - 统一 API 路由，将所有 API 置于 `/api` 路径下
  - 保留 `/v2` 路径以向后兼容
  - 使用应用层用例处理所有请求，逐步淘汰直接调用领域服务

### 2. 增强依赖注入
- **问题**：Wire 依赖注入配置存在重复绑定问题，导致生成失败
- **解决方案**：
  - 重构 `wire.go` 文件，移除重复绑定定义
  - 使用提供者函数代替直接绑定
  - 添加事件处理器注册逻辑

### 3. 消除全局变量
- **问题**：项目使用全局变量存储服务实例，违反依赖注入原则
- **解决方案**：
  - 移除全局服务变量
  - 通过工厂模式正确注入服务
  - 通过显式依赖提高可测试性

### 4. 增强架构验证
- **问题**：架构验证是手动的，容易出错
- **解决方案**：
  - 实现自动化层依赖检查
  - 通过代码扫描强制执行严格的架构边界
  - 在 CI 管道中添加验证

### 5. 优雅关闭
- **问题**：应用程序不能优雅关闭，可能导致数据丢失
- **解决方案**：
  - 为服务器实现优雅关闭机制，确保所有在途请求在关闭前完成
  - 添加关闭超时设置，防止关闭过程无限期挂起
  - 改进信号处理，支持 SIGINT 和 SIGTERM 信号

### 6. 国际化支持
- **问题**：应用程序缺乏适当的国际化支持
- **解决方案**：
  - 添加多语言验证错误消息的翻译中间件
  - 根据 Accept-Language 头自动选择适当的语言

### 7. CORS 支持
- **问题**：未正确处理跨源请求
- **解决方案**：
  - 添加 CORS 中间件处理跨源请求
  - 配置允许的源、方法、头和凭证

### 8. 调试工具
- **问题**：生产环境中性能问题诊断困难
- **解决方案**：
  - 集成 pprof 性能分析工具，用于诊断生产环境中的性能问题
  - 可通过配置文件启用或禁用

### 9. 高级 Redis 集成
- **问题**：Redis 实现有限，缺乏适当的连接管理
- **解决方案**：
  - 使用适当的连接池增强 Redis 客户端
  - 添加全面的健康检查和监控
  - 改进错误处理和连接生命周期管理

### 10. 结构化请求日志
- **问题**：API 请求缺乏适当的日志记录，使调试困难
- **解决方案**：
  - 实现全面的请求日志中间件
  - 添加请求 ID 追踪，用于关联日志
  - 根据状态码配置日志级别

### 11. 统一错误响应格式
- **问题**：API 错误响应格式不一致
- **解决方案**：
  - 标准化错误响应结构，包含代码、消息和详情
  - 添加错误文档引用
  - 实现一致的 HTTP 状态码映射

这些优化使项目更加健壮、可维护，并提供更好的开发体验。

## 快速开始

### 前置条件
- Go 1.21 或更高版本
- Docker（用于运行依赖服务）
- Homebrew（macOS 用户）
- Node.js 和 npm（用于提交消息规范检查）
- pre-commit（用于代码质量检查）
- golangci-lint（用于代码静态分析）

### 安装

#### 1. 克隆仓库
```bash
git clone https://github.com/RanchoCooper/go-hexagonal.git
cd go-hexagonal
```

#### 2. 初始化开发环境（macOS）
项目的 Makefile 中包含了便捷的 init 目标，用于设置所有必需的工具：

```bash
# 安装和配置所有必需的依赖
make init
```

这个命令会安装：
- Go（如果尚未安装）
- Node.js 和 npm（用于提交消息规范检查）
- pre-commit 钩子
- golangci-lint
- commitlint，用于确保提交消息符合标准

#### 3. 手动安装（非 macOS）
如果你不使用 macOS 或更喜欢手动设置：

```bash
# 安装 golangci-lint
# 参见 https://golangci-lint.run/usage/install/

# 安装 pre-commit
pip install pre-commit

# 安装 commitlint
npm install -g @commitlint/cli @commitlint/config-conventional

# 设置 pre-commit 钩子
make pre-commit.install
```

### 开发工作流

#### 代码格式化
```bash
# 根据 Go 标准格式化代码
make fmt
```

#### 运行测试
```bash
# 运行测试（包含竞态检测和覆盖率报告）
make test
```

#### 代码质量检查
```bash
# 运行静态代码分析检查代码质量
make ci.lint
```

#### 运行所有检查
```bash
# 运行格式化、静态分析和测试
make all
```

### 配置
1. 复制 `config/config.yaml.example` 到 `config/config.yaml`（如适用）
2. 根据需要调整配置值
3. 环境变量可以覆盖配置文件值

### Docker 设置（可选）
如果你的项目使用 Docker 进行本地开发：

```bash
# 启动所需服务（MySQL、Redis 等）
docker-compose up -d

# 完成后停止服务
docker-compose down
```

### Pre-commit 钩子管理

项目使用 pre-commit 钩子确保提交前的代码质量：

```bash
# 更新 pre-commit 钩子到最新版本
make precommit.rehook
```

### 运行应用程序

```bash
# 运行应用程序
go run cmd/main.go
```

## 贡献指南

1. Fork 仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'feat: add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 扩展计划

- **gRPC 支持** - 添加 gRPC 服务实现
- **监控集成** - 集成 Prometheus 监控

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

## 许可证

本项目采用 [Apache 2.0 License](LICENSE) 许可证 - 详见 [LICENSE](LICENSE) 文件。
