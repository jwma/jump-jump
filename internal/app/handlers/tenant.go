package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/config"
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
// @Description 获取租户详情（租户成员或超管可查看）
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
		tenantID := c.Param("id")
		if getMemberOrNil(ctx.User, tenantID) == nil {
			c.JSON(http.StatusOK, models.NewErrorResponse("你无权查看此租户"))
			return
		}

		repo := repository.GetTenantRepo(db.GetPostgresPool())
		t, err := repo.GetByID(tenantID)
		if err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(t))
	})
}

// ListTenantsAPI godoc
// @Summary 我的租户列表
// @Description 列出当前用户所属的租户
// @Tags 租户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=[]models.Tenant}
// @Failure 401 {object} nil
// @Router /tenant/ [get]
func ListTenantsAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, ctx *AuthContext) {
		repo := repository.GetTenantRepo(db.GetPostgresPool())
		tenants, err := repo.ListByUser(ctx.User.ID)
		if err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(tenants))
	})
}

// UpdateTenantAPI godoc
// @Summary 更新租户信息
// @Description 更新租户名称或 Slug（仅 admin 或超管）
// @Tags 租户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "租户 ID"
// @Param body body models.UpdateTenantRequest true "更新租户请求"
// @Success 200 {object} models.Response{data=models.Tenant}
// @Failure 401 {object} nil
// @Router /tenant/{id} [patch]
func UpdateTenantAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, ctx *AuthContext) {
		tenantID := c.Param("id")
		member, err := isTenantAdmin(ctx.User, tenantID)
		if err != nil || !member.IsAdmin() {
			c.JSON(http.StatusOK, models.NewErrorResponse("仅管理员可修改租户信息"))
			return
		}

		p := &models.UpdateTenantRequest{}
		if err := c.ShouldBindJSON(p); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse("参数错误"))
			return
		}

		repo := repository.GetTenantRepo(db.GetPostgresPool())
		t, err := repo.Update(tenantID, p)
		if err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(t))
	})
}

// AddDomainAPI godoc
// @Summary 添加域名
// @Description 为租户添加域名（租户 admin 或超管）
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
		tenantID := c.Param("id")
		member, err := isTenantAdmin(ctx.User, tenantID)
		if err != nil || !member.IsAdmin() {
			c.JSON(http.StatusOK, models.NewErrorResponse("你无权管理域名"))
			return
		}

		p := &models.AddDomainRequest{}
		if err := c.ShouldBindJSON(p); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse("参数错误"))
			return
		}

		repo := repository.GetTenantRepo(db.GetPostgresPool())
		if err := repo.AddDomain(tenantID, p.Domain, p.IsDefault); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		domains, _ := repo.ListDomains(tenantID)
		c.JSON(http.StatusOK, models.NewSuccessResponse(domains))
	})
}

// RemoveDomainAPI godoc
// @Summary 删除域名
// @Description 删除租户域名（租户 admin 或超管）
// @Tags 租户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "租户 ID"
// @Param domain path string true "域名"
// @Success 200 {object} models.Response{data=[]models.TenantDomain}
// @Failure 401 {object} nil
// @Router /tenant/{id}/domains/{domain} [delete]
func RemoveDomainAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, ctx *AuthContext) {
		tenantID := c.Param("id")
		member, err := isTenantAdmin(ctx.User, tenantID)
		if err != nil || !member.IsAdmin() {
			c.JSON(http.StatusOK, models.NewErrorResponse("你无权管理域名"))
			return
		}

		domain := c.Param("domain")
		if domain == "" {
			c.JSON(http.StatusOK, models.NewErrorResponse("domain 参数不能为空"))
			return
		}

		repo := repository.GetTenantRepo(db.GetPostgresPool())
		if err := repo.RemoveDomain(tenantID, domain); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		domains, _ := repo.ListDomains(tenantID)
		c.JSON(http.StatusOK, models.NewSuccessResponse(domains))
	})
}

// ListDomainsAPI godoc
// @Summary 获取租户域名列表
// @Description 获取租户域名列表（租户成员或超管可查看）
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
		tenantID := c.Param("id")
		if getMemberOrNil(ctx.User, tenantID) == nil {
			c.JSON(http.StatusOK, models.NewErrorResponse("你无权查看此租户域名"))
			return
		}

		repo := repository.GetTenantRepo(db.GetPostgresPool())
		domains, err := repo.ListDomains(tenantID)
		if err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(domains))
	})
}

// GetTenantConfigAPI godoc
// @Summary 获取租户配置
// @Description 获取租户配置（租户成员或超管可查看）
// @Tags 租户配置
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "租户 ID"
// @Success 200 {object} models.Response{data=models.GetConfigAPIResponseData}
// @Failure 401 {object} nil
// @Router /tenant/{id}/config [get]
func GetTenantConfigAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, ctx *AuthContext) {
		tenantID := c.Param("id")
		if getMemberOrNil(ctx.User, tenantID) == nil {
			c.JSON(http.StatusOK, models.NewErrorResponse("你无权查看此租户配置"))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(models.GetConfigAPIResponseData{
			Config: config.GetTenantConfig(tenantID),
		}))
	})
}

// UpdateTenantConfigAPI godoc
// @Summary 更新租户配置
// @Description 更新租户配置（租户 admin 或超管）
// @Tags 租户配置
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "租户 ID"
// @Param body body models.UpdateTenantConfigRequest true "更新配置请求"
// @Success 200 {object} models.Response{data=models.GetConfigAPIResponseData}
// @Failure 401 {object} nil
// @Router /tenant/{id}/config [patch]
func UpdateTenantConfigAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, ctx *AuthContext) {
		tenantID := c.Param("id")
		member, err := isTenantAdmin(ctx.User, tenantID)
		if err != nil || !member.IsAdmin() {
			c.JSON(http.StatusOK, models.NewErrorResponse("你无权修改此租户配置"))
			return
		}

		p := &models.UpdateTenantConfigRequest{}
		if err := c.ShouldBindJSON(p); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse("参数错误"))
			return
		}

		if p.IdLength != nil && p.IdMinimumLength != nil && p.IdMaximumLength != nil {
			if *p.IdMinimumLength <= *p.IdLength && *p.IdLength <= *p.IdMaximumLength &&
				*p.IdMinimumLength > 0 && *p.IdLength > 0 && *p.IdMaximumLength > 0 {
				config.UpdateIdConfig(tenantID, &config.IdConfig{
					IdLength:        *p.IdLength,
					IdMinimumLength: *p.IdMinimumLength,
					IdMaximumLength: *p.IdMaximumLength,
				})
			} else {
				c.JSON(http.StatusOK, models.NewErrorResponse("最小长度 <= 默认长度 <= 最大长度，三个值均大于 0"))
				return
			}
		}

		if p.NotFoundMode != nil && p.NotFoundValue != nil {
			if *p.NotFoundMode == config.ShortLinkNotFoundContentMode ||
				*p.NotFoundMode == config.ShortLinkNotFoundRedirectMode {
				config.UpdateShortLinkNotFoundConfig(tenantID, &config.ShortLinkNotFoundConfig{
					Mode:  *p.NotFoundMode,
					Value: *p.NotFoundValue,
				})
			} else {
				c.JSON(http.StatusOK, models.NewErrorResponse("处理模式参数不正确"))
				return
			}
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(models.GetConfigAPIResponseData{
			Config: config.GetTenantConfig(tenantID),
		}))
	})
}
