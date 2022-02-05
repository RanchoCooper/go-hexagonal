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
    once           sync.Once
    ExampleService *ExampleSvc
)

func Init(ctx context.Context) {
    once.Do(func() {
        ExampleService = NewExampleService(ctx)
    })
}
