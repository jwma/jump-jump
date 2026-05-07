package routers

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	docs "github.com/jwma/jump-jump/docs"
	"github.com/jwma/jump-jump/internal/app/handlers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	docs.SwaggerInfo.Description = "Jump Jump short link service"
	docs.SwaggerInfo.Version = "v1"
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	url := ginSwagger.URL("/swagger/doc.json")
	detectAPIDocHost()
	docsR := r.Group("/swagger", gin.BasicAuth(getAPIDocBasicAccounts()))
	{
		docsR.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}

	if gin.Mode() == gin.DebugMode {
		corsCfg := cors.DefaultConfig()
		corsCfg.AllowAllOrigins = true
		corsCfg.AddAllowHeaders("Authorization", "X-Tenant-ID")
		r.Use(cors.New(corsCfg))
	}

	r.Use(handlers.AllowedHostsMiddleware())
	r.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPaths([]string{"/v1/"})))

	r.LoadHTMLFiles("./web/admin/index.html")
	r.StaticFS("/static", http.Dir("./web/admin/static"))
	r.NoRoute(func(c *gin.Context) {
		// Skip API paths — let them return the default 404
		if strings.HasPrefix(c.Request.URL.Path, "/v1/") || strings.HasPrefix(c.Request.URL.Path, "/swagger/") {
			c.JSON(http.StatusNotFound, gin.H{"msg": "not found", "code": 404})
			return
		}
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// Auth routes — no TenantResolverMiddleware
	auth := r.Group("/v1/auth")
	{
		auth.POST("/login", handlers.LoginAPI)
		auth.GET("/info", handlers.JWTAuthenticatorMiddleware(), handlers.GetAuthInfoAPI())
	}

	// Tenant-resolved API routes
	v1 := r.Group("/v1")
	v1.Use(handlers.TenantResolverMiddleware())
	{
		v1.GET("/user/info", handlers.JWTAuthenticatorMiddleware(), handlers.GetUserInfoAPI())
		v1.POST("/user/logout", handlers.JWTAuthenticatorMiddleware(), handlers.LogoutAPI())
		v1.POST("/user/change-password", handlers.JWTAuthenticatorMiddleware(), handlers.ChangePasswordAPI())
		v1.GET("/user/preferences", handlers.JWTAuthenticatorMiddleware(), handlers.GetUserPreferencesAPI())
		v1.PUT("/user/preferences", handlers.JWTAuthenticatorMiddleware(), handlers.UpdateUserPreferencesAPI())

		v1.GET("/config", handlers.JWTAuthenticatorMiddleware(), handlers.GetConfigAPI())
		v1.PATCH("/config/id-length", handlers.JWTAuthenticatorMiddleware(), handlers.UpdateIdLengthConfigAPI())
		v1.PATCH("/config/short-link-404-handling", handlers.JWTAuthenticatorMiddleware(), handlers.UpdateShortLinkNotFoundConfigAPI())

		shortLinkAPI := v1.Group("/short-link")
		shortLinkAPI.Use(handlers.JWTAuthenticatorMiddleware())
		shortLinkAPI.GET("/", handlers.ListShortLinksAPI())
		shortLinkAPI.GET("/:id", handlers.GetShortLinkAPI())
		shortLinkAPI.POST("/", handlers.CreateShortLinkAPI())
		shortLinkAPI.PATCH("/:id", handlers.UpdateShortLinkAPI())
		shortLinkAPI.DELETE("/:id", handlers.DeleteShortLinkAPI())
		shortLinkAPI.GET("/:id/*action", handlers.ShortLinkActionAPI())

		tenantAPI := v1.Group("/tenant")
		tenantAPI.Use(handlers.JWTAuthenticatorMiddleware())
		tenantAPI.GET("/", handlers.ListTenantsAPI())
		tenantAPI.GET("/:id", handlers.GetTenantAPI())
		tenantAPI.POST("/", handlers.CreateTenantAPI())
		tenantAPI.GET("/:id/domains", handlers.ListDomainsAPI())
		tenantAPI.POST("/:id/domains", handlers.AddDomainAPI())
		tenantAPI.DELETE("/:id/domains", handlers.RemoveDomainAPI())
	}

	return r
}

func SetupLandingRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", handlers.LandingHome)
	r.GET("/:id", handlers.Redirect)

	return r
}
