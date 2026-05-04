package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/config"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/repository"
)

func LandingHome(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "https://github.com/jwma/jump-jump")
}

func Redirect(c *gin.Context) {
	if c.Param("id") == "favicon.ico" {
		c.String(http.StatusNotFound, "")
		return
	}

	slRepo := repository.GetShortLinkRepo(db.GetPostgresPool(), db.GetRedisClient())
	s, err := slRepo.Get(c.Param("id"))
	if err != nil {
		// Resolve tenant for config
		host := strings.Split(c.Request.Host, ":")[0]
		tenantID, _ := config.ResolveTenantID(host)
		if tenantID == "" {
			c.String(http.StatusOK, "你访问的页面不存在哦")
			return
		}

		cc := config.GetShortLinkNotFoundConfig(tenantID)
		switch cc.Mode {
		case config.ShortLinkNotFoundContentMode:
			c.String(http.StatusOK, cc.Value)
		case config.ShortLinkNotFoundRedirectMode:
			c.Redirect(http.StatusTemporaryRedirect, cc.Value)
		default:
			c.String(http.StatusOK, "你访问的页面不存在哦")
		}
		return
	}

	if !s.IsEnable {
		c.String(http.StatusOK, "你访问的页面不存在哦")
		return
	}

	rhRepo := repository.GetRequestHistoryRepo(db.GetRedisClient(), db.GetPostgresPool())
	alRepo := repository.GetActiveLinkRepo(db.GetRedisClient())
	go func() {
		rhRepo.Save(models.NewRequestHistory(s, c.ClientIP(), c.Request.UserAgent()))
		alRepo.Save(s.Id)
	}()

	c.Redirect(http.StatusTemporaryRedirect, s.Url)
}
