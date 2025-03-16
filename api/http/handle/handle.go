package handle

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-hexagonal/api/dto"
	"go-hexagonal/api/error_code"
	"go-hexagonal/api/http/paginate"
	"go-hexagonal/application/core"
	"go-hexagonal/util/log"
)

type Response struct {
	Ctx *gin.Context
}

// StandardResponse standard response structure
type StandardResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{Ctx: ctx}
}

func (r *Response) ToResponse(data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	r.Ctx.JSON(http.StatusOK, data)
}

func (r *Response) ToResponseList(list interface{}, totalRows int) {
	r.Ctx.JSON(http.StatusOK, gin.H{
		"list": list,
		"pager": dto.Pager{
			Page:      paginate.GetPage(r.Ctx),
			PageSize:  paginate.GetPageSize(r.Ctx),
			TotalRows: totalRows,
		},
	})
}

func (r *Response) ToErrorResponse(err *error_code.Error) {
	response := gin.H{
		"code": err.Code,
		"msg":  err.Msg,
	}
	if details := err.Details; len(details) > 0 {
		response["details"] = details
	}

	r.Ctx.JSON(err.StatusCode(), response)
}

// Success returns a success response
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, StandardResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Error unified error handling
func Error(c *gin.Context, err error) {
	// Handle application errors
	if appErr, ok := err.(*core.Error); ok {
		c.JSON(getHttpStatusFromAppError(appErr), StandardResponse{
			Code:    appErr.Code,
			Message: appErr.Message,
			Data:    nil,
		})
		return
	}

	// Handle API error codes
	if apiErr, ok := err.(*error_code.Error); ok {
		c.JSON(http.StatusBadRequest, StandardResponse{
			Code:    apiErr.Code,
			Message: apiErr.Msg,
			Data:    nil,
		})
		return
	}

	// Log unexpected errors
	log.SugaredLogger.Errorf("Unexpected error: %v", err)

	// Default error response
	c.JSON(http.StatusInternalServerError, StandardResponse{
		Code:    error_code.ServerErrorCode,
		Message: "Internal server error",
		Data:    nil,
	})
}

// getHttpStatusFromAppError maps application error types to HTTP status codes
func getHttpStatusFromAppError(err *core.Error) int {
	switch err.Type {
	case core.ErrorTypeValidation:
		return http.StatusBadRequest
	case core.ErrorTypeNotFound:
		return http.StatusNotFound
	case core.ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case core.ErrorTypeForbidden:
		return http.StatusForbidden
	case core.ErrorTypeConflict:
		return http.StatusConflict
	case core.ErrorTypeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
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
func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
