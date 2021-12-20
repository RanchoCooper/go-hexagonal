# Hexagonal Architecture Based On DDD

## Diagram
![](https://github.com/Sairyss/domain-driven-hexagon/raw/master/assets/images/DomainDrivenHexagon.png)

## Architecture

Mainly based on:

- Domain-Driven Design (DDD)
- Hexagonal (Ports and Adapters) Architecture
- Secure by Design
- Clean Architecture
- Onion Architecture
- SOLID Principles
- Software Design Patterns

And many other sources (more links below in every chapter).

Before we begin, here are the PROS and CONS of using a complete architecture like this:

### Pros

- Independent of external frameworks, technologies, databases, etc. Frameworks and external resources can be plugged/unplugged with much less effort.
- Easily testable and scalable.
- More secure. Some security principles are baked in design itself.
- The solution can be worked on and maintained by different teams, without stepping on each other's toes.
- Easier to add new features. As the system grows over time, the difficulty in adding new features remains constant and relatively small.
- If the solution is properly broken apart along bounded context lines, it becomes easy to convert pieces of it into microservices if needed.

### Cons

- This is a sophisticated architecture which requires a firm understanding of quality software principles, such as SOLID, Clean/Hexagonal Architecture, Domain-Driven Design, etc. Any team implementing such a solution will almost certainly require an expert to drive the solution and keep it from evolving the wrong way and accumulating technical debt.
- Some of the practices presented here are not recommended for small-medium sized applications with not a lot of business logic. There is added up-front complexity to support all those building blocks and layers, boilerplate code, abstractions, data mapping etc. thus implementing a complete architecture like this is generally ill-suited to simple CRUD applications and could over-complicate such solutions. Some of the described below principles can be used in a smaller sized applications but must be implemented only after analyzing and understanding all pros and cons.


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
[Domain Driven Hexagonal](https://github.com/Sairyss/domain-driven-hexagon)

[Hexagonal Architecture using Go](https://cgarciarosales97.medium.com/hexagonal-architecture-using-go-fiber-b2925fd677b5)

[Hexagonal Architecture by example - a hands-on introduction](https://blog.allegro.tech/2020/05/hexagonal-architecture-by-example.html)