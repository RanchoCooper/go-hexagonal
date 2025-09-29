package handle

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-hexagonal/api/dto"
	"go-hexagonal/api/error_code"
	"go-hexagonal/api/http/paginate"
	"go-hexagonal/util/log"
)

type Response struct {
	Ctx *gin.Context
}

// StandardResponse defines the standard API response structure
type StandardResponse struct {
	Code    int         `json:"code"`              // Status code
	Message string      `json:"message"`           // Response message
	Data    interface{} `json:"data,omitempty"`    // Response data
	DocRef  string      `json:"doc_ref,omitempty"` // Documentation reference
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{Ctx: ctx}
}

func (r *Response) ToResponse(data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	r.Ctx.JSON(http.StatusOK, StandardResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func (r *Response) ToResponseList(list interface{}, totalRows int) {
	r.Ctx.JSON(http.StatusOK, StandardResponse{
		Code:    0,
		Message: "success",
		Data: gin.H{
			"list": list,
			"pager": dto.Pager{
				Page:      paginate.GetPage(r.Ctx),
				PageSize:  paginate.GetPageSize(r.Ctx),
				TotalRows: totalRows,
			},
		},
	})
}

func (r *Response) ToErrorResponse(err *error_code.Error) {
	response := StandardResponse{
		Code:    err.Code,
		Message: err.Msg,
	}

	if len(err.Details) > 0 {
		response.Data = gin.H{"details": err.Details}
	}

	if err.DocRef != "" {
		response.DocRef = err.DocRef
	}

	r.Ctx.JSON(err.StatusCode(), response)
}

// Success returns a success response
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, StandardResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Error unified error handling
func Error(c *gin.Context, err error) {
	// Handle API error codes
	if apiErr, ok := err.(*error_code.Error); ok {
		c.JSON(apiErr.StatusCode(), StandardResponse{
			Code:    apiErr.Code,
			Message: apiErr.Msg,
			Data:    gin.H{"details": apiErr.Details},
			DocRef:  apiErr.DocRef,
		})
		return
	}

	// Handle util/errors AppError
	if appErr, ok := err.(interface {
		Error() string
		Code() string
		Message() string
		Details() []string
	}); ok {
		// Map AppError to appropriate HTTP status code
		statusCode := http.StatusInternalServerError
		errorCode := error_code.ServerErrorCode

		// Map error types to appropriate status codes
		switch appErr.Code() {
		case "validation", "invalid_input":
			statusCode = http.StatusBadRequest
			errorCode = error_code.InvalidParamsCode
		case "not_found", "resource_not_found":
			statusCode = http.StatusNotFound
			errorCode = error_code.NotFoundCode
		case "unauthorized", "forbidden":
			statusCode = http.StatusForbidden
			errorCode = error_code.UnauthorizedTokenErrorCode
		case "conflict", "already_exists":
			statusCode = http.StatusConflict
			errorCode = error_code.AccountExistErrorCode
		}

		response := StandardResponse{
			Code:    errorCode,
			Message: appErr.Message(),
		}

		if details := appErr.Details(); len(details) > 0 {
			response.Data = gin.H{"details": details}
		}

		c.JSON(statusCode, response)
		return
	}

	// Log unexpected errors
	log.SugaredLogger.Errorf("Unexpected error: %v", err)

	// Default error response
	c.JSON(http.StatusInternalServerError, StandardResponse{
		Code:    error_code.ServerErrorCode,
		Message: "Internal server error",
	})
}

// GetQueryInt gets an integer from query parameters with a default value
func GetQueryInt(c *gin.Context, key string, defaultValue int) int {
	value, exists := c.GetQuery(key)
	if !exists {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// GetQueryString gets a string from query parameters with a default value
func GetQueryString(c *gin.Context, key string, defaultValue string) string {
	value, exists := c.GetQuery(key)
	if !exists {
		return defaultValue
	}
	return value
}

// IsNil checks if an interface is nil or its underlying value is nil
func IsNil(i any) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
