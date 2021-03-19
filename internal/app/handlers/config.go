package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/config"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/thoas/go-funk"
	"net/http"
)

// GetConfigAPI godoc
// @Security ApiKeyAuth
// @Summary 获取系统配置信息
// @Description 获取系统配置信息
// @Tags 系统配置
// @Accept json
// @Produce json
// @Success 200 {object} models.Response{data=models.GetConfigAPIResponseData}
// @Failure 401
// @Router /config [get]
func GetConfigAPI(c *gin.Context) {
	c.JSON(http.StatusOK, models.NewSuccessResponse(models.GetConfigAPIResponseData{Config: config.GetSystemConfig()}))
}

// UpdateLandingHostsAPI godoc
// @Security ApiKeyAuth
// @Summary 更新落地页 Hosts
// @Description 更新落地页 Hosts
// @Tags 系统配置
// @Accept json
// @Produce json
// @Param body body models.UpdateLandingHostsAPIRequest true "更新落地页 Hosts 请求"
// @Success 200 {object} models.Response{data=models.GetConfigAPIResponseData}
// @Failure 401
// @Router /config/landing-hosts [patch]
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

		config.UpdateLandingHosts(p.Hosts) // 更新
		c.JSON(http.StatusOK, models.NewSuccessResponse(models.GetConfigAPIResponseData{Config: config.GetSystemConfig()}))
	})
}

// UpdateIdLengthConfigAPI godoc
// @Security ApiKeyAuth
// @Summary 更新短链接 ID 设置
// @Description 更新短链接 ID 设置
// @Tags 系统配置
// @Accept json
// @Produce json
// @Param body body config.IdConfig true "更新短链接 ID 设置请求"
// @Success 200 {object} models.Response{data=models.GetConfigAPIResponseData}
// @Failure 401
// @Router /config/id-length [patch]
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

// UpdateShortLinkNotFoundConfigAPI godoc
// @Security ApiKeyAuth
// @Summary 更新短链接 404 设置
// @Description 更新短链接 404 设置
// @Tags 系统配置
// @Accept json
// @Produce json
// @Param body body config.ShortLinkNotFoundConfig true "更新短链接 404 设置请求"
// @Success 200 {object} models.Response{data=models.GetConfigAPIResponseData}
// @Failure 401
// @Router /config/short-link-404-handling [patch]
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

		if !funk.ContainsString([]string{config.ShortLinkNotFoundContentMode, config.ShortLinkNotFoundRedirectMode},
			p.Mode) {
			c.JSON(http.StatusOK, models.NewErrorResponse("处理模式参数不正确"))
			return
		}

		config.UpdateShortLinkNotFoundConfig(p)
		c.JSON(http.StatusOK, models.NewSuccessResponse(models.GetConfigAPIResponseData{Config: config.GetSystemConfig()}))
	})
}
