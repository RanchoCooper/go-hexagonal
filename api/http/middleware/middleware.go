package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"go-hexagonal/util/log"
)

func Trigger(types ...int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(types) == 0 {
			return
		}

		path := ctx.FullPath()
		var body []byte
		// cache JSON data and rewrite to Request.Body
		if ctx.Request.Method == http.MethodPut && strings.HasPrefix(path, "/specific_url") {
			body, _ = ctx.GetRawData()
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		ctx.Next()

		// check if any business errors
		if ctx.Errors != nil {
			log.SugaredLogger.Errorf("NotifyTrigger fail due to business error, err: %s", ctx.Errors.String())

			return
		}

		for _, notifyType := range types {
			// logic handle
			log.SugaredLogger.Infof("handle type %d", notifyType)
		}
	}
}
