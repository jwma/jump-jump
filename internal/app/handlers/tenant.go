package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/repository"
)

func CreateTenantAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		if user.Role != models.RoleAdmin {
			c.JSON(http.StatusOK, models.NewErrorResponse("你无权创建租户"))
			return
		}

		p := &models.CreateTenantRequest{}
		if err := c.ShouldBindJSON(p); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse("参数错误"))
			return
		}

		repo := repository.GetTenantRepo(db.GetPostgresPool())
		t, err := repo.Create(p)
		if err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(t))
	})
}

func GetTenantAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		if user.Role != models.RoleAdmin {
			c.JSON(http.StatusOK, models.NewErrorResponse("你无权查看租户"))
			return
		}

		repo := repository.GetTenantRepo(db.GetPostgresPool())
		t, err := repo.GetByID(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(t))
	})
}

func ListTenantsAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		if user.Role != models.RoleAdmin {
			c.JSON(http.StatusOK, models.NewErrorResponse("你无权查看租户列表"))
			return
		}

		repo := repository.GetTenantRepo(db.GetPostgresPool())
		tenants, err := repo.List()
		if err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(tenants))
	})
}

func AddDomainAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		if user.Role != models.RoleAdmin {
			c.JSON(http.StatusOK, models.NewErrorResponse("你无权管理域名"))
			return
		}

		p := &models.AddDomainRequest{}
		if err := c.ShouldBindJSON(p); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse("参数错误"))
			return
		}

		repo := repository.GetTenantRepo(db.GetPostgresPool())
		if err := repo.AddDomain(c.Param("id"), p.Domain, p.IsDefault); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		domains, _ := repo.ListDomains(c.Param("id"))
		c.JSON(http.StatusOK, models.NewSuccessResponse(domains))
	})
}

func RemoveDomainAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		if user.Role != models.RoleAdmin {
			c.JSON(http.StatusOK, models.NewErrorResponse("你无权管理域名"))
			return
		}

		domain := c.Query("domain")
		if domain == "" {
			c.JSON(http.StatusOK, models.NewErrorResponse("domain 参数不能为空"))
			return
		}

		repo := repository.GetTenantRepo(db.GetPostgresPool())
		if err := repo.RemoveDomain(c.Param("id"), domain); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		domains, _ := repo.ListDomains(c.Param("id"))
		c.JSON(http.StatusOK, models.NewSuccessResponse(domains))
	})
}

func ListDomainsAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		if user.Role != models.RoleAdmin {
			c.JSON(http.StatusOK, models.NewErrorResponse("你无权查看域名"))
			return
		}

		repo := repository.GetTenantRepo(db.GetPostgresPool())
		domains, err := repo.ListDomains(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(domains))
	})
}
