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
    "go-hexagonal/util/log"
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
        log.Logger.Sugar().Error("CreateExample.BindAndValid errs: %v", errs)
        err := error_code.InvalidParams.WithDetails(errs.Errors()...)
        response.ToErrorResponse(err)
        return
    }
    example := &entity.Example{}
    err := copier.Copy(example, body)
    if err != nil {
        log.Logger.Sugar().Error("CreateExample failed.%v", err.Error())
        response.ToErrorResponse(error_code.CopyError)
        return
    }
    example, err = service.ExampleSvc.Create(ctx, example)
    if err != nil {
        log.Logger.Sugar().Error("CreateExample failed.%v", err.Error())
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
        log.Logger.Sugar().Error("DeleteExample.BindAndValid errs: %v", errs)
        errResp := error_code.InvalidParams.WithDetails(errs.Errors()...)
        response.ToErrorResponse(errResp)
        return
    }

    err := service.ExampleSvc.Delete(ctx, param.Id)
    if err != nil {
        log.Logger.Sugar().Error("DeleteExample failed.%v", err.Error())
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
        log.Logger.Sugar().Error("UpdateExample.BindAndValid errs: %v", errs)
        errResp := error_code.InvalidParams.WithDetails(errs.Errors()...)
        response.ToErrorResponse(errResp)
        return
    }
    example := &entity.Example{}
    copier.Copy(example, body)
    err := service.ExampleSvc.Update(ctx, example)
    if err != nil {
        log.Logger.Sugar().Error("UpdateExample failed.%v", err.Error())
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
        log.Logger.Sugar().Error("GetExample.BindAndValid errs: %v", errs)
        errResp := error_code.InvalidParams.WithDetails(errs.Errors()...)
        response.ToErrorResponse(errResp)
        return
    }
    result, err := service.ExampleSvc.Get(ctx, param.Id)
    if err != nil {
        log.Logger.Sugar().Error("GetExample failed.%v", err.Error())
        response.ToErrorResponse(error_code.ServerError)
        return
    }
    response.ToResponse(*result)
}
