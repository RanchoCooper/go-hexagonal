package handle

import (
    "github.com/gin-gonic/gin"
    "github.com/spf13/cast"

    "go-hexagonal/api/http/dto"
    "go-hexagonal/api/http/errcode"
    "go-hexagonal/api/http/validator"
    "go-hexagonal/internal/domain.model/service"
    "go-hexagonal/util/logger"
)

/**
 * @author Rancho
 * @date 2022/1/5
 */

func CreateExample(ctx *gin.Context) {
    response := NewResponse(ctx)
    body := dto.CreateExampleReq{}

    valid, errs := validator.BindAndValid(ctx, &body, ctx.ShouldBindJSON)
    if !valid {
        logger.Log.Errorf(ctx, "CreateExample.BindAndValid errs: %v", errs)
        errResp := errcode.InvalidParams.WithDetails(errs.Errors()...)
        response.ToErrorResponse(errResp)
        return
    }
    example, err := service.Service.ExampleService.Create(ctx, body)
    if err != nil {
        logger.Log.Errorf(ctx, "CreateExample failed.%v", err.Error())
        ctx.Abort()
        return
    }
    response.ToResponse(example)
}

func DeleteExample(ctx *gin.Context) {
    response := NewResponse(ctx)
    param := dto.DeleteExampleReq{}

    valid, errs := validator.BindAndValid(ctx, &param, ctx.ShouldBindUri)
    if !valid {
        logger.Log.Errorf(ctx, "DeleteExample.BindAndValid errs: %v", errs)
        errResp := errcode.InvalidParams.WithDetails(errs.Errors()...)
        response.ToErrorResponse(errResp)
        return
    }

    err := service.Service.ExampleService.Delete(ctx, param)
    if err != nil {
        logger.Log.Errorf(ctx, "DeleteExample failed.%v", err.Error())
        ctx.Abort()
        return
    }
    response.ToResponse(gin.H{})
}

func UpdateExample(ctx *gin.Context) {
    response := NewResponse(ctx)
    body := dto.UpdateExampleReq{Id: cast.ToUint(ctx.Param("id"))}

    valid, errs := validator.BindAndValid(ctx, &body, ctx.ShouldBindJSON)
    if !valid {
        logger.Log.Errorf(ctx, "UpdateExample.BindAndValid errs: %v", errs)
        errResp := errcode.InvalidParams.WithDetails(errs.Errors()...)
        response.ToErrorResponse(errResp)
        return
    }
    err := service.Service.ExampleService.Update(ctx, body)
    if err != nil {
        logger.Log.Errorf(ctx, "UpdateExample failed.%v", err.Error())
        ctx.Abort()
        return
    }
    response.ToResponse(gin.H{})
}

func GetExample(ctx *gin.Context) {
    response := NewResponse(ctx)
    param := dto.GetExampleReq{}

    valid, errs := validator.BindAndValid(ctx, &param, ctx.ShouldBindUri)
    if !valid {
        logger.Log.Errorf(ctx, "GetExample.BindAndValid errs: %v", errs)
        errResp := errcode.InvalidParams.WithDetails(errs.Errors()...)
        response.ToErrorResponse(errResp)
        return
    }
    result, err := service.Service.ExampleService.Get(ctx, param.Id)
    if err != nil {
        logger.Log.Errorf(ctx, "GetExample failed.%v", err.Error())
        ctx.Abort()
        return
    }
    response.ToResponse(*result)
}
