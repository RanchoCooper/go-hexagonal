# Hexagonal Architecture Based On DDD

## Architecture
![](https://github.com/Sairyss/domain-driven-hexagon/raw/master/assets/images/DomainDrivenHexagon.png)

## Overview
- **Archived**
    - **Essential**
        - [x] Domain Driven Design
        - [x] Hexagonal Architecture
        - [x] Repository Design (with transaction)
    - **Technical**
        - Mock UT with [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) and [Redis Mock](https://github.com/go-redis/redismock)
        - Clean Arch Detect/Check Tool
    - **Chore**
        - [x] [Travis CI](https://travis-ci.com)
        - [x] [Circle CI](https://travis-ci.com) (abandon)
        - [x] [Golangci-lint](https://github.com/golangci/golangci-lint)
        - [x] [Commit Lint](https://github.com/conventional-changelog/commitlint)
        - [x] [Pre-commit Hook](https://pre-commit.com/)
- **Roadmap**
    - **Essential**
        - [ ] Support Dependency Inversion/Dependency Injection
        - [ ] Improve HTTP Handle Implement
        - [ ] Support Domain Event
        - [ ] Add GRPC Example
    - **Technical**
        - [ ] Integrate [air](https://github.com/cosmtrek/air)
        - [ ] Integrate [Kafka](https://kafka.apache.org)
        - [ ] Integrate [Prometheus](https://prometheus.io)
        - [ ] Hot reloading configuration
- **Primary Module**
    - [Zap](https://github.com/uber-go/zap)
    - [Gin](https://gin-gonic.com)
    - [GORM](https://gorm.io)
    - [Cast](https://github.com/spf13/cast)
    - [Copier](https://github.com/jinzhu/copier)
    - [Structs](https://github.com/RanchoCooper/structs)
    - [Structtag](https://github.com/fatih/structtag)

## Usage

### Pre-commit Hook && Commitlint && Golangci-lint


manually install

```bash
# install pre-commit
brew install pre-commit
# install golangci-lint
brew install golangci-lint
# install commitlint
npm install -g @commitlint/cli @commitlint/config-conventional
# add commitlint config
echo "module.exports = {extends: ['@commitlint/config-conventional']}" > commitlint.config.js
# add pre-commit hook
make precommit.rehook
```

or just type

```bash
make init && make precommit.rehook
```

# Reference
- **Architecture**
    - [Freedom DDD Framework](https://github.com/8treenet/freedom)
    - [Hexagonal Architecture in Go](https://medium.com/@matiasvarela/hexagonal-architecture-in-go-cfd4e436faa3)
    - [Dependency Injection in A Nutshell](https://appliedgo.net/di/)
- **Project Conventional**
    - [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0)
    - [Improving Your Go Project With pre-commit hooks](https://goangle.medium.com/golang-improving-your-go-project-with-pre-commit-hooks-a265fad0e02f)
