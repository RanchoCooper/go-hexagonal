package paginate

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"go-hexagonal/config"
)

func GetPage(c *gin.Context) int {
	page := cast.ToInt(c.Query("page"))
	if page <= 0 {
		return 1
	}

	return page
}

func GetPageSize(c *gin.Context) int {
	pageSize := cast.ToInt(c.Query("page_size"))
	if pageSize <= 0 {
		return config.GlobalConfig.HTTPServer.DefaultPageSize
	}
	if pageSize > config.GlobalConfig.HTTPServer.MaxPageSize {
		return config.GlobalConfig.HTTPServer.MaxPageSize
	}

	return pageSize
}

func GetPageOffset(page, pageSize int) int {
	return (page - 1) * pageSize
}
