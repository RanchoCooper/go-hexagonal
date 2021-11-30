# Hexagonal Architecture Based On DDD
![](https://image.slidesharecdn.com/javadev-hexagonalarchitectureforjavaapplications-150202062634-conversion-gate01/95/hexagonal-architecture-for-java-applications-11-638.jpg?cb=1423245064)

# Project Layout

## /api
contains API based on HTTP/rpc/gRPC etc, which only depending on `/internal/port.adapter/service`

## /cmd
Main applications entry for current project.

The directory name for each application should match the name of the executable you want to have (e.g., /cmd/myapp).

Don't put a lot of code in the `cmd` directory.

If you think the code can be reused in other projects or package, then it should live in the `util` directory.

If the code is not reusable or if you don't want others to reuse it, put that code in the `/internal` directory. 

You'll be surprised what others will do, so be explicit about your intentions!

It's common to have a small main function that imports and invokes code from other packages and nothing else.

## /config
Configuration files or default configs.

Put your `yaml` or template files here.

## /internal
Private package and code. This is the code you don't want others importing in their applications or libraries.

Note that this layout pattern is enforced by the Go compiler itself. See the [Go 1.4 release notes](https://go.dev/doc/go1.4#internalpackages) for more details.

Note that you are not limited to the top level internal directory. You can have more than one internal directory at any level of your project tree.

### /internal/application

### /internal/domain.model
The domain layer is the core of the project. It only focuses on the business and does not pay attention to the technical implementation details, so it does not rely on any other layers.

Encapsulate the core business logic, and provide business entities and business logic calculations to the Application layer through the `Domain Service` and `Domain Entity` methods. 


### /internal/port.adapter

#### /internal/port.adapter/amqp

#### /internal/port.adapter/dependency

#### /internal/port.adapter/repository

#### /internal/port.adapter/service

## /util

# Reference
[六边形架构](https://juejin.cn/post/6844903569947099143)

[Hexagonal Architecture using Go](https://cgarciarosales97.medium.com/hexagonal-architecture-using-go-fiber-b2925fd677b5)

[Hexagonal Architecture by example - a hands-on introduction](https://blog.allegro.tech/2020/05/hexagonal-architecture-by-example.html)