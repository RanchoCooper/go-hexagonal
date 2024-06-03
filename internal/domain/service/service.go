package service

import (
	"context"
	"sync"
)

var (
	once       sync.Once
	ExampleSvc *ExampleService
)

func Init(ctx context.Context) {
	once.Do(func() {
		ExampleSvc = NewExampleService(ctx)
	})
}
