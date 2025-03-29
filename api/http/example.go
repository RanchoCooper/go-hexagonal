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
	"go-hexagonal/domain/model"
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

	// If converter exists, convert DTO to domain model first
	var example *model.Example
	var err error

	if converter != nil {
		example, err = converter.FromCreateRequest(&body)
		if err != nil {
			log.SugaredLogger.Errorf("CreateExample.FromCreateRequest errs: %v", err)
			response.ToErrorResponse(error_code.ServerError)
			return
		}
	}

	// Execute use case
	var result any
	if example != nil && services != nil && services.ExampleService != nil {
		// Use converted model to execute creation
		created, err := services.ExampleService.Create(ctx, example.Name, example.Alias)
		if err != nil {
			log.SugaredLogger.Errorf("CreateExample failed: %v", err.Error())
			response.ToErrorResponse(error_code.ServerError)
			return
		}

		// Convert return result
		result, err = converter.ToExampleResponse(created)
		if err != nil {
			log.SugaredLogger.Errorf("CreateExample.ToExampleResponse errs: %v", err)
			response.ToErrorResponse(error_code.ServerError)
			return
		}
	} else {
		// Fall back to original implementation
		result, err = appFactory.CreateExampleUseCase().Execute(ctx, body)
		if err != nil {
			log.SugaredLogger.Errorf("CreateExample failed: %v", err.Error())
			response.ToErrorResponse(error_code.ServerError)
			return
		}
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

	// Use direct service call or application factory
	var err error
	if services != nil && services.ExampleService != nil {
		// Direct service call for deletion
		err = services.ExampleService.Delete(ctx, param.Id)
		if err != nil {
			log.SugaredLogger.Errorf("DeleteExample failed: %v", err.Error())
			response.ToErrorResponse(error_code.ServerError)
			return
		}
	} else {
		// Fall back to original implementation
		_, err = appFactory.DeleteExampleUseCase().Execute(ctx, param.Id)
		if err != nil {
			log.SugaredLogger.Errorf("DeleteExample failed.%v", err.Error())
			response.ToErrorResponse(error_code.ServerError)
			return
		}
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

	// If converter exists, convert DTO to domain model first
	var err error

	if converter != nil && services != nil && services.ExampleService != nil {
		example, err := converter.FromUpdateRequest(&body)
		if err != nil {
			log.SugaredLogger.Errorf("UpdateExample.FromUpdateRequest errs: %v", err)
			response.ToErrorResponse(error_code.ServerError)
			return
		}

		// Use converted model to execute update
		err = services.ExampleService.Update(ctx, example.Id, example.Name, example.Alias)
		if err != nil {
			log.SugaredLogger.Errorf("UpdateExample failed: %v", err.Error())
			response.ToErrorResponse(error_code.ServerError)
			return
		}
	} else {
		// Fall back to original implementation
		_, err = appFactory.UpdateExampleUseCase().Execute(ctx, body)
		if err != nil {
			log.SugaredLogger.Errorf("UpdateExample failed.%v", err.Error())
			response.ToErrorResponse(error_code.ServerError)
			return
		}
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

	// Execute use case
	var result any
	var err error

	if converter != nil && services != nil && services.ExampleService != nil {
		// Direct service call
		example, err := services.ExampleService.Get(ctx, param.Id)
		if err != nil {
			log.SugaredLogger.Errorf("GetExample failed: %v", err.Error())
			response.ToErrorResponse(error_code.ServerError)
			return
		}

		// Convert return result
		result, err = converter.ToExampleResponse(example)
		if err != nil {
			log.SugaredLogger.Errorf("GetExample.ToExampleResponse errs: %v", err)
			response.ToErrorResponse(error_code.ServerError)
			return
		}
	} else {
		// Fall back to original implementation
		result, err = appFactory.GetExampleUseCase().Execute(ctx, param.Id)
		if err != nil {
			log.SugaredLogger.Errorf("GetExample failed.%v", err.Error())
			response.ToErrorResponse(error_code.ServerError)
			return
		}
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

	// Execute use case
	var result any
	var err error

	if converter != nil && services != nil && services.ExampleService != nil {
		// Direct service call
		example, err := services.ExampleService.FindByName(ctx, name)
		if err != nil {
			log.SugaredLogger.Errorf("FindExampleByName failed: %v", err.Error())
			if err.Error() == "record not found" ||
				err.Error() == "failed to find example: record not found" {
				response.ToErrorResponse(error_code.NotFound.WithDetails("example not found"))
				return
			}
			response.ToErrorResponse(error_code.ServerError)
			return
		}

		// Convert return result
		result, err = converter.ToExampleResponse(example)
		if err != nil {
			log.SugaredLogger.Errorf("FindExampleByName.ToExampleResponse errs: %v", err)
			response.ToErrorResponse(error_code.ServerError)
			return
		}
	} else {
		// Fall back to original implementation
		result, err = appFactory.FindExampleByNameUseCase().Execute(ctx, name)
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
	}

	response.ToResponse(result)
}
