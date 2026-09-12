// Package middleware provides custom Gin middlewares for logging and localization.
package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// I18nMiddleware detects language preference from query parameters ('lang') or the 'Accept-Language' HTTP header.
func I18nMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Check query parameter (e.g. ?lang=id or ?lang=en)
		lang := strings.ToLower(strings.TrimSpace(c.Query("lang")))

		// 2. If query parameter is empty, inspect HTTP Accept-Language header
		if lang == "" {
			acceptLang := strings.ToLower(c.GetHeader("Accept-Language"))
			if strings.Contains(acceptLang, "id") {
				lang = "id"
			} else if strings.Contains(acceptLang, "en") {
				lang = "en"
			}
		}

		// 3. Fallback default to 'en' if invalid or empty
		if lang != "id" && lang != "en" {
			lang = "en"
		}

		// Store active language code in Gin Context
		c.Set("lang", lang)

		c.Next()
	}
}

