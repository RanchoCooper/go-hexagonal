// Package middleware provides HTTP request processing middleware
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-hexagonal/api/error_code"
	"go-hexagonal/util/errors"
)

// ErrorHandlerMiddleware handles API layer error responses uniformly
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Return appropriate HTTP status code and error message based on error type
			switch {
			case errors.IsValidationError(err):
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    error_code.InvalidParamsCode,
					"message": err.Error(),
				})
				return

			case errors.IsNotFoundError(err):
				c.JSON(http.StatusNotFound, gin.H{
					"code":    error_code.NotFoundCode,
					"message": err.Error(),
				})
				return

			case errors.IsPersistenceError(err):
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    error_code.ServerErrorCode,
					"message": "Database operation failed",
				})
				return

			default:
				// Default server error
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    error_code.ServerErrorCode,
					"message": "Internal server error",
				})
			}
		}
	}
}
