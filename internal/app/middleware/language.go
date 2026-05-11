package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func LanguageMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.Query("lang")

		if lang == "" {
			accept := c.GetHeader("Accept-Language")
			if strings.HasPrefix(accept, "zh") {
				lang = "zh"
			}
		}

		if lang != "zh" {
			lang = "en"
		}

		c.Set("locale", lang)
		c.Next()
	}
}
