package service

import (
    "testing"

    "github.com/stretchr/testify/assert"

    "go-hexagonal/config"
    "go-hexagonal/internal/domain/entity"
)

/**
 * @author Rancho
 * @date 2021/12/25
 */

func TestExampleService_Create(t *testing.T) {
    if config.Config.App.Env == config.EnvGithub {
        t.SkipNow()
    }
    srv := NewExampleService(ctx)
    assert.NotNil(t, srv)
    assert.NotNil(t, srv.Repository)
    resp, err := srv.Create(ctx, &entity.Example{
        Name:  "RanchoCooper",
        Alias: "Rancho",
    })
    assert.Nil(t, err)
    assert.NotNil(t, resp)
}
