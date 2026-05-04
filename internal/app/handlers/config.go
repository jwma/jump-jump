package handlers

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/config"
	"github.com/jwma/jump-jump/internal/app/models"
)

// GetConfigAPI godoc
// @Summary 获取系统配置信息
// @Description 获取系统配置信息
// @Tags 系统配置
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=models.GetConfigAPIResponseData}
// @Failure 401 {object} nil
// @Router /config [get]
func GetConfigAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		tenantID := user.TenantID
		c.JSON(http.StatusOK, models.NewSuccessResponse(models.GetConfigAPIResponseData{
			Config: config.GetTenantConfig(tenantID),
		}))
	})
}

// UpdateIdLengthConfigAPI godoc
// @Summary 更新短链接 ID 设置
// @Description 更新短链接 ID 设置
// @Tags 系统配置
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body models.UpdateIdLengthRequest true "更新短链接 ID 设置请求"
// @Success 200 {object} models.Response{data=models.GetConfigAPIResponseData}
// @Failure 401 {object} nil
// @Router /config/id-length [patch]
func UpdateIdLengthConfigAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		if user.Role != models.RoleAdmin {
			c.JSON(http.StatusOK, models.NewErrorResponse("你无权修改设置"))
			return
		}

		p := &models.UpdateIdLengthRequest{}
		if err := c.ShouldBindJSON(p); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		if p.IdMinimumLength <= p.IdLength && p.IdLength <= p.IdMaximumLength &&
			p.IdMinimumLength > 0 && p.IdLength > 0 && p.IdMaximumLength > 0 {
			config.UpdateIdConfig(user.TenantID, &config.IdConfig{
				IdLength:        p.IdLength,
				IdMinimumLength: p.IdMinimumLength,
				IdMaximumLength: p.IdMaximumLength,
			})
			c.JSON(http.StatusOK, models.NewSuccessResponse(models.GetConfigAPIResponseData{
				Config: config.GetTenantConfig(user.TenantID),
			}))
			return
		}

		c.JSON(http.StatusOK, models.NewErrorResponse("最小长度 <= 默认长度 <= 最大长度，三个值均大于 0"))
	})
}

// UpdateShortLinkNotFoundConfigAPI godoc
// @Summary 更新短链接 404 设置
// @Description 更新短链接 404 设置
// @Tags 系统配置
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body models.UpdateNotFoundConfigRequest true "更新短链接 404 设置请求"
// @Success 200 {object} models.Response{data=models.GetConfigAPIResponseData}
// @Failure 401 {object} nil
// @Router /config/short-link-404-handling [patch]
func UpdateShortLinkNotFoundConfigAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		if user.Role != models.RoleAdmin {
			c.JSON(http.StatusOK, models.NewErrorResponse("你无权修改配置"))
			return
		}

		p := &models.UpdateNotFoundConfigRequest{}
		if err := c.ShouldBindJSON(p); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		if !slices.Contains([]string{config.ShortLinkNotFoundContentMode, config.ShortLinkNotFoundRedirectMode}, p.Mode) {
			c.JSON(http.StatusOK, models.NewErrorResponse("处理模式参数不正确"))
			return
		}

		config.UpdateShortLinkNotFoundConfig(user.TenantID, &config.ShortLinkNotFoundConfig{
			Mode:  p.Mode,
			Value: p.Value,
		})
		c.JSON(http.StatusOK, models.NewSuccessResponse(models.GetConfigAPIResponseData{
			Config: config.GetTenantConfig(user.TenantID),
		}))
	})
}
