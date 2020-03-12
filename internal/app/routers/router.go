package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/handlers"
	"net/http"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "码极工作室/MJ Studio")
	})

	r.POST("/v1/user/login", handlers.Login)

	shortLinkAPI := r.Group("/v1/short-link")
	shortLinkAPI.Use(handlers.JWTAuthenticatorMiddleware())
	shortLinkAPI.GET("/:id", handlers.GetShortLinkAPI())
	shortLinkAPI.POST("/", handlers.CreateShortLinkAPI())
	shortLinkAPI.PATCH("/:id", handlers.UpdateShortLinkAPI())
	shortLinkAPI.GET("/:id/*action", handlers.ShortLinkActionAPI())

	return r
}

func SetupLandingRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/:id", handlers.Redirect)

	return r
}
