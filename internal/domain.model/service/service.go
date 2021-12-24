package service

import (
    "context"
)

/**
 * @author Rancho
 * @date 2021/12/10
 */

type DomainServer struct {
    *ExampleService
}

type DomainServerOption func(srv *DomainServer)

func NewDomainService(ctx context.Context, opts ...DomainServerOption) {
    srv := &DomainServer{}

    for _, opt := range opts {
        opt(srv)
    }
}

func WithExampleService(ctx context.Context) DomainServerOption {
    return func(s *DomainServer) {
        s.ExampleService = NewExampleService(ctx)
    }
}
