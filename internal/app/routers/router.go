package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/handlers"
	"net/http"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	if gin.Mode() == gin.DebugMode { // 开发环境下，开启 CORS
		corsCfg := cors.DefaultConfig()
		corsCfg.AllowAllOrigins = true
		corsCfg.AddAllowHeaders("Authorization")
		r.Use(cors.New(corsCfg))
	}

	r.Use(handlers.AllowedHostsMiddleware())

	r.LoadHTMLFiles("./web/admin/index.html")
	r.StaticFS("/static", http.Dir("./web/admin/static"))
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.POST("/v1/user/login", handlers.Login)
	r.GET("/v1/user/info", handlers.JWTAuthenticatorMiddleware(), handlers.GetUserInfoAPI())
	r.POST("/v1/user/logout", handlers.JWTAuthenticatorMiddleware(), handlers.LogoutAPI())
	r.PATCH("/v1/user/change-password", handlers.JWTAuthenticatorMiddleware(), handlers.ChangePasswordAPI())

	r.GET("/v1/config", handlers.GetConfig)
	r.PATCH("/v1/config/landing-hosts", handlers.JWTAuthenticatorMiddleware(),
		handlers.UpdateLandingHostsAPI())
	r.PATCH("/v1/config/id-length", handlers.JWTAuthenticatorMiddleware(),
		handlers.UpdateIdLengthConfigAPI())
	r.PATCH("/v1/config/short-link-404-handling", handlers.JWTAuthenticatorMiddleware(),
		handlers.UpdateShortLinkNotFoundConfigAPI())

	shortLinkAPI := r.Group("/v1/short-link")
	shortLinkAPI.Use(handlers.JWTAuthenticatorMiddleware())
	shortLinkAPI.GET("/", handlers.ListShortLinksAPI())
	shortLinkAPI.GET("/:id", handlers.GetShortLinkAPI())
	shortLinkAPI.POST("/", handlers.CreateShortLinkAPI())
	shortLinkAPI.PATCH("/:id", handlers.UpdateShortLinkAPI())
	shortLinkAPI.DELETE("/:id", handlers.DeleteShortLinkAPI())
	shortLinkAPI.GET("/:id/*action", handlers.ShortLinkActionAPI())

	return r
}

func SetupLandingRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/:id", handlers.Redirect)

	return r
}
