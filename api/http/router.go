package http

import (
    "github.com/gin-gonic/gin"

    "go-hexagonal/api/http/handle"
    "go-hexagonal/config"
)

func NewServerRoute() *gin.Engine {
    if config.Config.App.Debug {
        gin.SetMode(gin.DebugMode)
    } else {
        gin.SetMode(gin.ReleaseMode)
    }

    router := gin.Default()
    example := router.Group("/example")
    {
        example.POST("", handle.CreateExample)
        // example.DELETE("", handle.DeleteExample)
        // example.PUT("", handle.UpdateExample)
        // example.GET("", handle.GetExample)
    }

    return router
}
