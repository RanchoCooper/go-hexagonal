// Package core provides core application layer interfaces and types
package core

import "context"

// UseCase defines the application layer use case interface
type UseCase interface {
	Execute(ctx context.Context, input interface{}) (interface{}, error)
}

// UseCaseHandler defines the use case handler interface
type UseCaseHandler interface {
	Handle(ctx context.Context, input interface{}) (interface{}, error)
}

// BaseUseCase provides a base implementation for use cases
type BaseUseCase struct {
	Handler UseCaseHandler
}

// Execute executes the use case
func (u *BaseUseCase) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	return u.Handler.Handle(ctx, input)
}
