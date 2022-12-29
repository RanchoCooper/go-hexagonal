package handle

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-hexagonal/api/dto"
	"go-hexagonal/api/error_code"
	"go-hexagonal/api/http/paginate"
)

type Response struct {
	Ctx *gin.Context
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
	details := err.Details
	if len(details) > 0 {
		response["details"] = details
	}

	r.Ctx.JSON(err.StatusCode(), response)
}
