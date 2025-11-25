package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"runtime/debug"

	appErrors "github.com/L11D/avito-review-assign-service/internal/errors"
	"github.com/L11D/avito-review-assign-service/pkg/api/dto"
	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			var appError *appErrors.AppError

			if errors.As(err, &appError) {
				errDTO := dto.ErrorDTO{
					Code:    string(appError.Code),
					Message: appError.Message,
				}
				c.JSON(appError.StatusCode, gin.H{"error": errDTO})

				slog.WarnContext(c.Request.Context(), appError.Message,
					slog.String("path", c.Request.URL.Path),
					slog.String("method", c.Request.Method),
					slog.String("code", string(appError.Code)),
					slog.Int("status_code", appError.StatusCode),
				)
			} else {
				slog.ErrorContext(c.Request.Context(), err.Error(),
					slog.String("path", c.Request.URL.Path),
					slog.String("method", c.Request.Method),
				)
				slog.ErrorContext(c.Request.Context(), string(debug.Stack()))

				errDTO := dto.ErrorDTO{
					Code:    "INTERNAL_ERROR",
					Message: "Internal server error",
				}

				c.JSON(http.StatusInternalServerError, gin.H{"error": errDTO})
			}

			c.Abort()

			return
		}
	}
}
