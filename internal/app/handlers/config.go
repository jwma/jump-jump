package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/config"
	"github.com/jwma/jump-jump/internal/app/models"
	"net/http"
)

func GetConfig(c *gin.Context) {
	cfg := config.GetConfig()
	landingHosts := cfg.GetStringSliceValue("landingHosts", make([]string, 0))

	c.JSON(http.StatusOK, gin.H{
		"msg":  "ok",
		"code": 0,
		"data": gin.H{
			"config": gin.H{"landingHosts": landingHosts},
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
