package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	docs "github.com/jwma/jump-jump/docs"
	"github.com/jwma/jump-jump/internal/app/handlers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"
	"strings"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @contact.name MJ Ma
// @contact.url https://www.linkedin.com/in/mj-profile/
// @contact.email m.mjw.ma@gmail.com

// @license.name MIT
// @license.url https://github.com/jwma/jump-jump/blob/master/LICENSE

func getAPIDocBasicAccounts() gin.Accounts {
	defaultUsername := "apidoc"
	defaultPassword := "showmethedoc"

	u := os.Getenv("API_DOC_USERNAME")
	p := os.Getenv("API_DOC_PASSWORD")

	if u == "" || p == "" {
		return gin.Accounts{defaultUsername: defaultPassword}
	}

	return gin.Accounts{u: p}
}

func detectAPIDocHost() {
	// 如果通过该环境变量指明了 API 文档使用的 Host，则直接使用
	h := os.Getenv("API_DOC_HOST")
	if h != "" {
		docs.SwaggerInfo.Host = h
		return
	}

	if gin.Mode() == gin.DebugMode {
		docs.SwaggerInfo.Host = os.Getenv("J2_API_ADDR")
	} else {
		docs.SwaggerInfo.Host = strings.Split(os.Getenv("ALLOWED_HOSTS"), ",")[0]
	}
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Swagger
	docs.SwaggerInfo.Title = "Jump Jump API Documentation"
	docs.SwaggerInfo.Description = "🚀🚀🚀"
	docs.SwaggerInfo.Version = "v1"
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	url := ginSwagger.URL("/swagger/doc.json")
	detectAPIDocHost()
	docsR := r.Group("/swagger", gin.BasicAuth(getAPIDocBasicAccounts()))
	{
		docsR.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}

	if gin.Mode() == gin.DebugMode { // 开发环境下，开启 CORS
		corsCfg := cors.DefaultConfig()
		corsCfg.AllowAllOrigins = true
		corsCfg.AddAllowHeaders("Authorization")
		r.Use(cors.New(corsCfg))
	}

	r.Use(handlers.AllowedHostsMiddleware())
	r.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPaths([]string{"/v1/"})))

	// serve dashboard static resources
	r.LoadHTMLFiles("./web/admin/index.html")
	r.StaticFS("/static", http.Dir("./web/admin/static"))
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// v1 API's
	v1 := r.Group("/v1")
	{
		// account stuff
		v1.POST("/user/login", handlers.LoginAPI)
		v1.GET("/user/info", handlers.JWTAuthenticatorMiddleware(), handlers.GetUserInfoAPI())
		v1.POST("/user/logout", handlers.JWTAuthenticatorMiddleware(), handlers.LogoutAPI())
		v1.POST("/user/change-password", handlers.JWTAuthenticatorMiddleware(), handlers.ChangePasswordAPI())

		// system configuration stuff
		v1.GET("/config", handlers.JWTAuthenticatorMiddleware(), handlers.GetConfigAPI)
		v1.PATCH("/config/landing-hosts", handlers.JWTAuthenticatorMiddleware(), handlers.UpdateLandingHostsAPI())
		v1.PATCH("/config/id-length", handlers.JWTAuthenticatorMiddleware(), handlers.UpdateIdLengthConfigAPI())
		v1.PATCH("/config/short-link-404-handling", handlers.JWTAuthenticatorMiddleware(), handlers.UpdateShortLinkNotFoundConfigAPI())

		// short link stuff
		shortLinkAPI := v1.Group("/short-link")
		shortLinkAPI.Use(handlers.JWTAuthenticatorMiddleware())
		shortLinkAPI.GET("/", handlers.ListShortLinksAPI())
		shortLinkAPI.GET("/:id", handlers.GetShortLinkAPI())
		shortLinkAPI.POST("/", handlers.CreateShortLinkAPI())
		shortLinkAPI.PATCH("/:id", handlers.UpdateShortLinkAPI())
		shortLinkAPI.DELETE("/:id", handlers.DeleteShortLinkAPI())
		shortLinkAPI.GET("/:id/*action", handlers.ShortLinkActionAPI())
	}

	return r
}

func SetupLandingRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", handlers.LandingHome)
	r.GET("/:id", handlers.Redirect)

	return r
}
