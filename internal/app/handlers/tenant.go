package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/repository"
)

// CreateTenantAPI godoc
// @Summary 创建租户
// @Description 创建租户（已登录用户均可创建，创建者自动成为管理员）
// @Tags 租户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body models.CreateTenantRequest true "创建租户请求"
// @Success 200 {object} models.Response{data=models.Tenant}
// @Failure 401 {object} nil
// @Router /tenant/ [post]
func CreateTenantAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, ctx *AuthContext) {
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

		memberRepo := repository.GetTenantMemberRepo(db.GetPostgresPool())
		if err := memberRepo.Save(&models.TenantMember{
			TenantID: t.ID,
			UserID:   ctx.User.ID,
			Role:     models.RoleAdmin,
		}); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse("创建租户成员关系失败"))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(t))
	})
}

// GetTenantAPI godoc
// @Summary 获取租户详情
// @Description 获取租户详情（仅管理员）
// @Tags 租户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "租户 ID"
// @Success 200 {object} models.Response{data=models.Tenant}
// @Failure 401 {object} nil
// @Router /tenant/{id} [get]
func GetTenantAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, ctx *AuthContext) {
		if !ctx.User.IsSuper && ctx.Member.Role != models.RoleAdmin {
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

// ListTenantsAPI godoc
// @Summary 租户列表
// @Description 获取租户列表（仅管理员）
// @Tags 租户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=[]models.Tenant}
// @Failure 401 {object} nil
// @Router /tenant/ [get]
func ListTenantsAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, ctx *AuthContext) {
		if !ctx.User.IsSuper && ctx.Member.Role != models.RoleAdmin {
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

// AddDomainAPI godoc
// @Summary 添加域名
// @Description 为租户添加域名（仅管理员）
// @Tags 租户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "租户 ID"
// @Param body body models.AddDomainRequest true "添加域名请求"
// @Success 200 {object} models.Response{data=[]models.TenantDomain}
// @Failure 401 {object} nil
// @Router /tenant/{id}/domains [post]
func AddDomainAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, ctx *AuthContext) {
		if !ctx.User.IsSuper && ctx.Member.Role != models.RoleAdmin {
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

// RemoveDomainAPI godoc
// @Summary 删除域名
// @Description 删除租户域名（仅管理员）
// @Tags 租户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "租户 ID"
// @Param domain query string true "域名"
// @Success 200 {object} models.Response{data=[]models.TenantDomain}
// @Failure 401 {object} nil
// @Router /tenant/{id}/domains [delete]
func RemoveDomainAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, ctx *AuthContext) {
		if !ctx.User.IsSuper && ctx.Member.Role != models.RoleAdmin {
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

// ListDomainsAPI godoc
// @Summary 获取租户域名列表
// @Description 获取租户域名列表（仅管理员）
// @Tags 租户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "租户 ID"
// @Success 200 {object} models.Response{data=[]models.TenantDomain}
// @Failure 401 {object} nil
// @Router /tenant/{id}/domains [get]
func ListDomainsAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, ctx *AuthContext) {
		if !ctx.User.IsSuper && ctx.Member.Role != models.RoleAdmin {
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
