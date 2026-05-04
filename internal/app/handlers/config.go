package handlers

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/config"
	"github.com/jwma/jump-jump/internal/app/models"
)

func GetConfigAPI(c *gin.Context) {
	c.JSON(http.StatusOK, models.NewSuccessResponse(models.GetConfigAPIResponseData{Config: config.GetSystemConfig()}))
}

func UpdateLandingHostsAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		if user.Role != models.RoleAdmin {
			c.JSON(http.StatusOK, models.NewErrorResponse("你无权修改短链接域名"))
			return
		}

		p := &models.UpdateLandingHostsAPIRequest{}
		if err := c.ShouldBindJSON(p); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		config.UpdateLandingHosts(p.Hosts)
		c.JSON(http.StatusOK, models.NewSuccessResponse(models.GetConfigAPIResponseData{Config: config.GetSystemConfig()}))
	})
}

func UpdateIdLengthConfigAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		if user.Role != models.RoleAdmin {
			c.JSON(http.StatusOK, models.NewErrorResponse("你无权修改随机 ID 长度设置"))
			return
		}

		p := &config.IdConfig{}
		if err := c.ShouldBindJSON(p); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		if p.IdMinimumLength <= p.IdLength && p.IdLength <= p.IdMaximumLength &&
			p.IdMinimumLength > 0 && p.IdLength > 0 && p.IdMaximumLength > 0 {
			config.UpdateIdConfig(p)
			c.JSON(http.StatusOK, models.NewSuccessResponse(models.GetConfigAPIResponseData{Config: config.GetSystemConfig()}))
			return
		}

		c.JSON(http.StatusOK, models.NewErrorResponse("最小长度 <= 默认长度 <= 最大长度，三个值均大于 0"))
	})
}

func UpdateShortLinkNotFoundConfigAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		if user.Role != models.RoleAdmin {
			c.JSON(http.StatusOK, models.NewErrorResponse("你无权修改短链接 404 处理配置"))
			return
		}

		p := &config.ShortLinkNotFoundConfig{}
		if err := c.ShouldBindJSON(p); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		if !slices.Contains([]string{config.ShortLinkNotFoundContentMode, config.ShortLinkNotFoundRedirectMode}, p.Mode) {
			c.JSON(http.StatusOK, models.NewErrorResponse("处理模式参数不正确"))
			return
		}

		config.UpdateShortLinkNotFoundConfig(p)
		c.JSON(http.StatusOK, models.NewSuccessResponse(models.GetConfigAPIResponseData{Config: config.GetSystemConfig()}))
	})
}
