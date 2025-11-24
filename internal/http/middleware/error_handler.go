package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/L11D/avito-review-assign-service/internal/errors"
	"github.com/L11D/avito-review-assign-service/internal/http/dto"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

		if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            
            switch e := err.(type) {
            case *errors.AppError:
                errDTO := dto.ErrorDTO{
                    Code:    string(e.Code),
                    Message: e.Message,
                }
                c.JSON(e.StatusCode, gin.H{"error": errDTO})

                slog.WarnContext(c.Request.Context(), e.Message,
                    slog.String("path", c.Request.URL.Path),
                    slog.String("method", c.Request.Method),
                    slog.String("code", string(e.Code)),
                    slog.Int("status_code", e.StatusCode),
                )
                
            default:
				slog.ErrorContext(c.Request.Context(), e.Error(),
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