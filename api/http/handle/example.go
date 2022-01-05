package handle

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "go-hexagonal/api/http/dto"
    "go-hexagonal/internal/domain.model/service"
    "go-hexagonal/util/logger"
)

/**
 * @author Rancho
 * @date 2022/1/5
 */

func CreateExample(ctx *gin.Context) {
    body := dto.CreateExampleReq{}
    err := ctx.ShouldBindJSON(&body)
    if err != nil {
        logger.Log.Errorf(ctx, "create example failed.%v", err.Error())
        ctx.Abort()
        return
    }

    example, err := service.Service.Create(ctx, body)
    if err != nil {
        logger.Log.Errorf(ctx, "create example failed.%v", err.Error())
        ctx.Abort()
        return
    }
    ctx.JSON(http.StatusOK, example)
}
