package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/config"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/thoas/go-funk"
	"net/http"
)

func GetConfig(c *gin.Context) {
	cfg := config.GetConfig()
	landingHosts := cfg.GetStringSliceValue("landingHosts", make([]string, 0))
	idMinimumLength := cfg.GetIntValue("idMinimumLength", 2)
	idLength := cfg.GetIntValue("idLength", 6)
	idMaximumLength := cfg.GetIntValue("idMaximumLength", 10)
	shortLinkNotFoundConfig := cfg.GetStringStringMapValue("shortLinkNotFoundConfig",
		config.GetDefaultShortLinkNotFoundConfig())

	c.JSON(http.StatusOK, gin.H{
		"msg":  "ok",
		"code": 0,
		"data": gin.H{
			"config": gin.H{
				"landingHosts": landingHosts,
				"idConfig": gin.H{
					"idMinimumLength": idMinimumLength,
					"idLength":        idLength,
					"idMaximumLength": idMaximumLength,
				},
				"shortLinkNotFoundConfig": shortLinkNotFoundConfig,
			},
		},
	})
}

type LandingHostsParameter struct {
	Hosts []string `json:"hosts"`
}

func UpdateLandingHostsAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		if user.Role != models.RoleAdmin {
			c.JSON(http.StatusOK, gin.H{
				"msg":  "你无权修改短链接域名",
				"code": 4999,
				"data": nil,
			})
			return
		}

		p := &LandingHostsParameter{}
		if err := c.ShouldBindJSON(p); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg":  err.Error(),
				"code": 4999,
				"data": nil,
			})
			return
		}

		cfg := config.GetConfig()
		cfg.SetValue("landingHosts", p.Hosts)
		cfg.Persist()

		c.JSON(http.StatusOK, gin.H{
			"msg":  "",
			"code": 0,
			"data": gin.H{
				"config": gin.H{"landingHosts": cfg.GetStringSliceValue("landingHosts", make([]string, 0))},
			},
		})
	})
}

type idLengthParameter struct {
	IdMinimumLength int `json:"idMinimumLength"`
	IdLength        int `json:"idLength"`
	IdMaximumLength int `json:"idMaximumLength"`
}

func UpdateIdLengthConfigAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		if user.Role != models.RoleAdmin {
			c.JSON(http.StatusOK, gin.H{
				"msg":  "你无权修改随机 ID 长度设置",
				"code": 4999,
				"data": nil,
			})
			return
		}

		p := &idLengthParameter{}

		if err := c.ShouldBindJSON(p); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg":  err.Error(),
				"code": 4999,
				"data": nil,
			})
			return
		}

		if p.IdMinimumLength <= p.IdLength && p.IdLength <= p.IdMaximumLength &&
			p.IdMinimumLength > 0 && p.IdLength > 0 && p.IdMaximumLength > 0 {
			cfg := config.GetConfig()
			cfg.SetValue("idMinimumLength", p.IdMinimumLength)
			cfg.SetValue("idLength", p.IdLength)
			cfg.SetValue("idMaximumLength", p.IdMaximumLength)
			cfg.Persist()

			c.JSON(http.StatusOK, gin.H{
				"msg":  "",
				"code": 0,
				"data": gin.H{
					"config": gin.H{
						"idConfig": gin.H{
							"idMinimumLength": p.IdMinimumLength,
							"idLength":        p.IdLength,
							"idMaximumLength": p.IdMaximumLength,
						},
					},
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg":  "最小长度 <= 默认长度 <= 最大长度，三个值均大于 0",
			"code": 4999,
			"data": nil,
		})
	})
}

type shortLinkNotFoundConfigParameter struct {
	Mode  string `json:"mode" binding:"required"`
	Value string `json:"value" binding:"required"`
}

func UpdateShortLinkNotFoundConfigAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		if user.Role != models.RoleAdmin {
			c.JSON(http.StatusOK, gin.H{
				"msg":  "你无权修改短链接 404 处理配置",
				"code": 4999,
				"data": nil,
			})
			return
		}

		p := &shortLinkNotFoundConfigParameter{}

		if err := c.ShouldBindJSON(p); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg":  err.Error(),
				"code": 4999,
				"data": nil,
			})
			return
		}

		if !funk.ContainsString([]string{config.ShortLinkNotFoundContentMode, config.ShortLinkNotFoundRedirectMode},
			p.Mode) {
			c.JSON(http.StatusOK, gin.H{
				"msg":  "处理模式参数不正确",
				"code": 4999,
				"data": nil,
			})
			return
		}

		cc := config.NewShortLinkNotFoundConfig(p.Mode, p.Value)
		cfg := config.GetConfig()
		cfg.SetValue("shortLinkNotFoundConfig", cc)
		cfg.Persist()

		c.JSON(http.StatusOK, gin.H{
			"msg":  "",
			"code": 0,
			"data": gin.H{
				"config": gin.H{
					"shortLinkNotFoundConfig": cc,
				},
			},
		})
	})
}
