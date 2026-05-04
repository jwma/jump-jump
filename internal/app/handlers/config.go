package handlers

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/config"
	"github.com/jwma/jump-jump/internal/app/models"
)

func GetConfigAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		tenantID := user.TenantID
		c.JSON(http.StatusOK, models.NewSuccessResponse(models.GetConfigAPIResponseData{
			Config: config.GetTenantConfig(tenantID),
		}))
	})
}

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
