package http

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"go-hexagonal/api/error_code"
	"go-hexagonal/api/http/handle"
	"go-hexagonal/application"
	"go-hexagonal/application/core"
	"go-hexagonal/application/example"
	"go-hexagonal/util/log"
)

// ExampleHandlerV2 is the V2 version of example handler using application layer
type ExampleHandlerV2 struct {
	useCaseFactory *application.UseCaseFactory
}

// NewExampleHandlerV2 creates a new V2 example handler
func NewExampleHandlerV2(useCaseFactory *application.UseCaseFactory) *ExampleHandlerV2 {
	return &ExampleHandlerV2{
		useCaseFactory: useCaseFactory,
	}
}

// Create handles example creation
func (h *ExampleHandlerV2) Create(ctx *gin.Context) {
	response := handle.NewResponse(ctx)
	var input example.CreateExampleInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		log.SugaredLogger.Errorf("CreateExample.BindJSON errs: %v", err)
		response.ToErrorResponse(error_code.InvalidParams)
		return
	}

	useCase := h.useCaseFactory.CreateExampleUseCase()
	result, err := useCase.Execute(ctx, input)
	if err != nil {
		log.SugaredLogger.Errorf("CreateExample failed: %v", err)
		response.ToErrorResponse(error_code.ServerError)
		return
	}

	response.ToResponse(result)
}

// Get handles example retrieval
func (h *ExampleHandlerV2) Get(ctx *gin.Context) {
	response := handle.NewResponse(ctx)
	id := cast.ToInt(ctx.Param("id"))

	if id <= 0 {
		log.SugaredLogger.Errorf("GetExample.Invalid ID: %v", id)
		response.ToErrorResponse(error_code.InvalidParams)
		return
	}

	input := example.GetExampleInput{ID: id}
	useCase := h.useCaseFactory.GetExampleUseCase()
	result, err := useCase.Execute(ctx, input)

	if err != nil {
		if err == core.ErrNotFound {
			response.ToErrorResponse(error_code.NotFound)
			return
		}
		log.SugaredLogger.Errorf("GetExample failed: %v", err)
		response.ToErrorResponse(error_code.ServerError)
		return
	}

	response.ToResponse(result)
}
