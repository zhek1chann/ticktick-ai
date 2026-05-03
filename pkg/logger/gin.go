package logger

import (
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"
)

// GinLogger returns a Gin middleware that logs HTTP requests using slog.
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// startTime := time.Now()

		requestPath := c.Request.URL.Path
		if strings.HasPrefix(requestPath, "/swagger/") {
			c.Next()
			return
		}
		rawQuery := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// latency := time.Since(startTime)
		statusCode := c.Writer.Status()
		// clientIP := c.ClientIP()
		method := c.Request.Method
		bytesWritten := c.Writer.Size()
		if bytesWritten < 0 {
			bytesWritten = 0
		}

		ctx := c.Request.Context()

		// Log errors, if any, with preserved log context from wrapped errors
		for _, ginErr := range c.Errors {
			ctxWithErr := ErrorCtx(ctx, ginErr.Err)
			slog.ErrorContext(ctxWithErr, ginErr.Error())
		}

		// Build common attributes
		attrs := []slog.Attr{
			slog.Int("status", statusCode),
			slog.String("method", method),
			slog.String("path", requestPath),
			// slog.String("ip", clientIP),
			// slog.Duration("latency", latency),
		}
		if rawQuery != "" {
			attrs = append(attrs, slog.String("query", rawQuery))
		}
		switch {
		case statusCode >= 500:
			slog.LogAttrs(ctx, slog.LevelError, "", attrs...)
		case statusCode >= 400:
			slog.LogAttrs(ctx, slog.LevelWarn, "", attrs...)
		default:
			slog.LogAttrs(ctx, slog.LevelInfo, "", attrs...)
		}
	}
}
