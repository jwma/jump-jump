package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

var matcher = language.NewMatcher([]language.Tag{
	language.English,
	language.Chinese,
})

func LanguageMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.Query("lang")

		if lang != "" {
			tag, _, _ := matcher.Match(language.Make(lang))
			if tag == language.Chinese {
				c.Set("locale", "zh")
			} else {
				c.Set("locale", "en")
			}
		} else {
			accept := c.GetHeader("Accept-Language")
			if accept != "" {
				tags, _, _ := language.ParseAcceptLanguage(accept)
				tag, _, _ := matcher.Match(tags...)
				if tag == language.Chinese {
					c.Set("locale", "zh")
				} else {
					c.Set("locale", "en")
				}
			} else {
				c.Set("locale", "en")
			}
		}

		c.Next()
	}
}
