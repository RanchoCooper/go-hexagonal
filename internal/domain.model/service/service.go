package service

import (
    "context"
)

/**
 * @author Rancho
 * @date 2021/12/10
 */

type DomainService struct {
    *ExampleService
}

type DomainServiceOption func(srv *DomainService)

func NewDomainService(opts ...DomainServiceOption) *DomainService {
    srv := &DomainService{}

    for _, opt := range opts {
        opt(srv)
    }
    return srv
}

func WithExampleService(ctx context.Context) DomainServiceOption {
    return func(s *DomainService) {
        s.ExampleService = NewExampleService(ctx)
    }
}
