package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"go-hexagonal/api/dto"
	"go-hexagonal/api/error_code"
	"go-hexagonal/api/http/handle"
	"go-hexagonal/api/http/validator"
	"go-hexagonal/application"
	"go-hexagonal/util/log"
)

// appFactory is the application factory instance
var appFactory *application.Factory

// SetAppFactory sets the application factory
func SetAppFactory(factory *application.Factory) {
	appFactory = factory
}

func CreateExample(ctx *gin.Context) {
	response := handle.NewResponse(ctx)
	body := dto.CreateExampleReq{}

	valid, errs := validator.BindAndValid(ctx, &body, ctx.ShouldBindJSON)
	if !valid {
		log.SugaredLogger.Errorf("CreateExample.BindAndValid errs: %v", errs)
		err := error_code.InvalidParams.WithDetails(errs.Errors()...)
		response.ToErrorResponse(err)
		return
	}

	// Execute the use case
	result, err := appFactory.CreateExampleUseCase().Execute(ctx, body)
	if err != nil {
		log.SugaredLogger.Errorf("CreateExample failed: %v", err.Error())
		response.ToErrorResponse(error_code.ServerError)
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

func DeleteExample(ctx *gin.Context) {
	response := handle.NewResponse(ctx)
	param := dto.DeleteExampleReq{}

	valid, errs := validator.BindAndValid(ctx, &param, ctx.ShouldBindUri)
	if !valid {
		log.SugaredLogger.Errorf("DeleteExample.BindAndValid errs: %v", errs)
		errResp := error_code.InvalidParams.WithDetails(errs.Errors()...)
		response.ToErrorResponse(errResp)
		return
	}

	// Execute the use case
	err := appFactory.DeleteExampleUseCase().Execute(ctx, param.Id)
	if err != nil {
		log.SugaredLogger.Errorf("DeleteExample failed.%v", err.Error())
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
		log.SugaredLogger.Errorf("UpdateExample.BindAndValid errs: %v", errs)
		errResp := error_code.InvalidParams.WithDetails(errs.Errors()...)
		response.ToErrorResponse(errResp)
		return
	}

	// Execute the use case
	err := appFactory.UpdateExampleUseCase().Execute(ctx, body)
	if err != nil {
		log.SugaredLogger.Errorf("UpdateExample failed.%v", err.Error())
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
		log.SugaredLogger.Errorf("GetExample.BindAndValid errs: %v", errs)
		errResp := error_code.InvalidParams.WithDetails(errs.Errors()...)
		response.ToErrorResponse(errResp)
		return
	}

	// Execute the use case
	result, err := appFactory.GetExampleUseCase().Execute(ctx, param.Id)
	if err != nil {
		log.SugaredLogger.Errorf("GetExample failed.%v", err.Error())
		response.ToErrorResponse(error_code.ServerError)
		return
	}
	response.ToResponse(result)
}

func FindExampleByName(ctx *gin.Context) {
	response := handle.NewResponse(ctx)
	name := ctx.Param("name")

	if name == "" {
		log.SugaredLogger.Errorf("FindExampleByName. Name parameter is empty")
		response.ToErrorResponse(error_code.InvalidParams.WithDetails("name parameter is required"))
		return
	}

	// Execute the use case
	result, err := appFactory.FindExampleByNameUseCase().Execute(ctx, name)
	if err != nil {
		log.SugaredLogger.Errorf("FindExampleByName failed.%v", err.Error())
		if err.Error() == "record not found" ||
			err.Error() == "failed to find example: record not found" {
			response.ToErrorResponse(error_code.NotFound.WithDetails("example not found"))
			return
		}
		response.ToErrorResponse(error_code.ServerError)
		return
	}

	response.ToResponse(result)
}
