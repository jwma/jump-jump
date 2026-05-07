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

	// Auth routes — no tenant context required
	auth := r.Group("/v1/auth")
	{
		auth.POST("/login", handlers.LoginAPI)
		auth.GET("/info", handlers.JWTAuthenticatorMiddleware(), handlers.GetAuthInfoAPI())
	}

	// Super admin routes — JWT + SuperAdmin middleware
	superAPI := r.Group("/v1/super")
	superAPI.Use(handlers.JWTAuthenticatorMiddleware(), handlers.SuperAdminMiddleware())
	{
		// User management
		superAPI.POST("/users", handlers.CreateUserAPI)
		superAPI.GET("/users", handlers.ListUsersAPI)
		superAPI.GET("/users/:id", handlers.GetUserAPI)
		superAPI.POST("/users/:id/reset-password", handlers.ResetPasswordAPI)
		superAPI.PATCH("/users/:id/status", handlers.UpdateUserStatusAPI)

		// Tenant management
		superAPI.GET("/tenants", handlers.ListAllTenantsAPI)
		superAPI.PATCH("/tenants/:id/status", handlers.UpdateTenantStatusAPI)
		superAPI.GET("/tenants/:id/short-links", handlers.SuperListTenantShortLinksAPI)
		superAPI.GET("/tenants/:id/members", handlers.SuperListTenantMembersAPI)
		superAPI.GET("/tenants/:id/domains", handlers.SuperListTenantDomainsAPI)
	}

	// Tenant CRUD routes — JWT only, tenant ID from path params
	tenantAPI := r.Group("/v1/tenant")
	tenantAPI.Use(handlers.JWTAuthenticatorMiddleware())
	{
		tenantAPI.POST("/", handlers.CreateTenantAPI())
		tenantAPI.GET("/", handlers.ListTenantsAPI())
		tenantAPI.GET("/:id", handlers.GetTenantAPI())
		tenantAPI.PATCH("/:id", handlers.UpdateTenantAPI())

		// Domain management — stored here for landingserver redirect resolution
		tenantAPI.GET("/:id/domains", handlers.ListDomainsAPI())
		tenantAPI.POST("/:id/domains", handlers.AddDomainAPI())
		tenantAPI.DELETE("/:id/domains/:domain", handlers.RemoveDomainAPI())

		// Config
		tenantAPI.GET("/:id/config", handlers.GetTenantConfigAPI())
		tenantAPI.PATCH("/:id/config", handlers.UpdateTenantConfigAPI())

		// Member management
		tenantAPI.GET("/:id/members", handlers.ListTenantMembersAPI())
		tenantAPI.POST("/:id/invitations", handlers.InviteUserAPI())
		tenantAPI.PATCH("/:id/members/:userId/role", handlers.UpdateMemberRoleAPI())
		tenantAPI.DELETE("/:id/members/:userId", handlers.RemoveMemberAPI())
		tenantAPI.POST("/:id/leave", handlers.LeaveTenantAPI())
	}

	// User routes — JWT + TenantContext (X-Tenant-ID required for role info)
	userAPI := r.Group("/v1/user")
	userAPI.Use(handlers.JWTAuthenticatorMiddleware(), handlers.TenantContextMiddleware())
	{
		userAPI.GET("/info", handlers.GetUserInfoAPI())
		userAPI.POST("/logout", handlers.LogoutAPI())
		userAPI.POST("/change-password", handlers.ChangePasswordAPI())
		userAPI.GET("/preferences", handlers.GetUserPreferencesAPI())
		userAPI.PUT("/preferences", handlers.UpdateUserPreferencesAPI())
	}

	// Short link routes — JWT + TenantContext (X-Tenant-ID required)
	shortLinkAPI := r.Group("/v1/short-link")
	shortLinkAPI.Use(handlers.JWTAuthenticatorMiddleware(), handlers.TenantContextMiddleware())
	{
		shortLinkAPI.GET("/", handlers.ListShortLinksAPI())
		shortLinkAPI.GET("/:id", handlers.GetShortLinkAPI())
		shortLinkAPI.POST("/", handlers.CreateShortLinkAPI())
		shortLinkAPI.PATCH("/:id", handlers.UpdateShortLinkAPI())
		shortLinkAPI.DELETE("/:id", handlers.DeleteShortLinkAPI())
		shortLinkAPI.GET("/:id/*action", handlers.ShortLinkActionAPI())
	}

	// Invitation routes — JWT only, tenant-agnostic
	invAPI := r.Group("/v1/invitations")
	invAPI.Use(handlers.JWTAuthenticatorMiddleware())
	{
		invAPI.GET("/", handlers.ListMyInvitationsAPI())
		invAPI.POST("/:id/accept", handlers.AcceptInvitationAPI())
		invAPI.POST("/:id/reject", handlers.RejectInvitationAPI())
	}

	return r
}

func SetupLandingRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", handlers.LandingHome)
	r.GET("/:id", handlers.Redirect)

	return r
}
