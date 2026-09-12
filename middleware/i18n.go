package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// I18nMiddleware mendeteksi preferensi bahasa dari query parameter atau header Accept-Language
func I18nMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Cek dari query parameter (misal: ?lang=id atau ?lang=en)
		lang := strings.ToLower(strings.TrimSpace(c.Query("lang")))

		// 2. Jika query kosong, baca dari header HTTP Accept-Language
		if lang == "" {
			acceptLang := strings.ToLower(c.GetHeader("Accept-Language"))
			if strings.Contains(acceptLang, "id") {
				lang = "id"
			} else if strings.Contains(acceptLang, "en") {
				lang = "en"
			}
		}

		// 3. Fallback default ke 'en' jika tidak dikenali atau kosong
		if lang != "id" && lang != "en" {
			lang = "en"
		}

		// Simpan preferensi bahasa ke context
		c.Set("lang", lang)

		c.Next()
	}
}
