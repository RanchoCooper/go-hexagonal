package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

/**
 * @author Rancho
 * @date 2022/12/29
 */

func TestTrigger(t *testing.T) {
	// it's an example unit test for middleware

	t.Run("Exp start", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, engine := gin.CreateTestContext(w)

		engine.GET("/test",
			Trigger(1),
			func(ctx *gin.Context) {
			})

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)
		ctx.Request.AddCookie(&http.Cookie{
			Name:  "token",
			Value: "root-token",
			Path:  "/",
		})
		engine.HandleContext(ctx)
		assert.Nil(t, ctx.Errors)
	})

}
