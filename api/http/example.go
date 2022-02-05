package http

import (
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/copier"
    "github.com/spf13/cast"

    "go-hexagonal/api/dto"
    "go-hexagonal/api/error_code"
    "go-hexagonal/api/http/handle"
    "go-hexagonal/api/http/validator"
    "go-hexagonal/internal/domain/entity"
    "go-hexagonal/internal/domain/service"
    "go-hexagonal/util/logger"
)

/**
 * @author Rancho
 * @date 2022/1/5
 */

func CreateExample(ctx *gin.Context) {
    response := handle.NewResponse(ctx)
    body := dto.CreateExampleReq{}

    valid, errs := validator.BindAndValid(ctx, &body, ctx.ShouldBindJSON)
    if !valid {
        logger.Log.Errorf(ctx, "CreateExample.BindAndValid errs: %v", errs)
        err := error_code.InvalidParams.WithDetails(errs.Errors()...)
        response.ToErrorResponse(err)
        return
    }
    example := &entity.Example{}
    err := copier.Copy(example, body)
    if err != nil {
        logger.Log.Errorf(ctx, "CreateExample failed.%v", err.Error())
        response.ToErrorResponse(error_code.CopyError)
        return
    }
    example, err = service.ExampleService.Create(ctx, example)
    if err != nil {
        logger.Log.Errorf(ctx, "CreateExample failed.%v", err.Error())
        response.ToErrorResponse(error_code.ServerError)
        return
    }
    response.ToResponse(example)
}

func DeleteExample(ctx *gin.Context) {
    response := handle.NewResponse(ctx)
    param := dto.DeleteExampleReq{}

    valid, errs := validator.BindAndValid(ctx, &param, ctx.ShouldBindUri)
    if !valid {
        logger.Log.Errorf(ctx, "DeleteExample.BindAndValid errs: %v", errs)
        errResp := error_code.InvalidParams.WithDetails(errs.Errors()...)
        response.ToErrorResponse(errResp)
        return
    }

    err := service.ExampleService.Delete(ctx, param.Id)
    if err != nil {
        logger.Log.Errorf(ctx, "DeleteExample failed.%v", err.Error())
        response.ToErrorResponse(error_code.ServerError)
        return
    }
    response.ToResponse(gin.H{})
}

func UpdateExample(ctx *gin.Context) {
    response := handle.NewResponse(ctx)
    body := dto.UpdateExampleReq{Id: cast.ToUint(ctx.Param("id"))}

    valid, errs := validator.BindAndValid(ctx, &body, ctx.ShouldBindJSON)
    if !valid {
        logger.Log.Errorf(ctx, "UpdateExample.BindAndValid errs: %v", errs)
        errResp := error_code.InvalidParams.WithDetails(errs.Errors()...)
        response.ToErrorResponse(errResp)
        return
    }
    example := &entity.Example{}
    copier.Copy(example, body)
    err := service.ExampleService.Update(ctx, example)
    if err != nil {
        logger.Log.Errorf(ctx, "UpdateExample failed.%v", err.Error())
        response.ToErrorResponse(error_code.ServerError)
        return
    }
    response.ToResponse(gin.H{})
}

func GetExample(ctx *gin.Context) {
    response := handle.NewResponse(ctx)
    param := dto.GetExampleReq{}

    valid, errs := validator.BindAndValid(ctx, &param, ctx.ShouldBindUri)
    if !valid {
        logger.Log.Errorf(ctx, "GetExample.BindAndValid errs: %v", errs)
        errResp := error_code.InvalidParams.WithDetails(errs.Errors()...)
        response.ToErrorResponse(errResp)
        return
    }
    result, err := service.ExampleService.Get(ctx, param.Id)
    if err != nil {
        logger.Log.Errorf(ctx, "GetExample failed.%v", err.Error())
        response.ToErrorResponse(error_code.ServerError)
        return
    }
    response.ToResponse(*result)
}
