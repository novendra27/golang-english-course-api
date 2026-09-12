// Package middleware provides custom Gin middlewares for logging and localization.
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// HTTPLogger is a Gin middleware that records structured Zerolog logs for every incoming HTTP request.
func HTTPLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process downstream handlers
		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		event := log.Info()
		if statusCode >= 400 && statusCode < 500 {
			event = log.Warn()
		} else if statusCode >= 500 {
			event = log.Error()
		}

		event.
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Dur("latency", latency).
			Str("ip", clientIP).
			Str("user_agent", c.Request.UserAgent()).
			Str("error", errorMessage).
			Msg("HTTP Request")
	}
}

