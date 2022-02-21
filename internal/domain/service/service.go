package service

import (
    "context"
    "sync"
)

/**
 * @author Rancho
 * @date 2021/12/10
 */

var (
    once       sync.Once
    ExampleSvc *ExampleService
)

func Init(ctx context.Context) {
    once.Do(func() {
        ExampleSvc = NewExampleService(ctx)
    })
}
