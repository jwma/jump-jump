package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/config"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/repository"
	"log"
	"net/http"
)

func LandingHome(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "https://github.com/jwma/jump-jump")
}

func Redirect(c *gin.Context) {
	if c.Param("id") == "favicon.ico" {
		c.String(http.StatusNotFound, "")
		return
	}

	slRepo := repository.GetShortLinkRepo(db.GetRedisClient())
	s, err := slRepo.Get(c.Param("id"))

	if err != nil {
		log.Printf("查找短链接失败，error: %v\n", err)
		cc := config.GetShortLinkNotFoundConfig()

		switch cc.Mode {
		case config.ShortLinkNotFoundContentMode:
			c.String(http.StatusOK, cc.Value)
			break
		case config.ShortLinkNotFoundRedirectMode:
			c.Redirect(http.StatusTemporaryRedirect, cc.Value)
			break
		default:
			c.String(http.StatusOK, "你访问的页面不存在哦")
		}

		return
	}

	if !s.IsEnable {
		c.String(http.StatusOK, "你访问的页面不存在哦")
		return
	}

	// 保存短链接请求记录（IP、User-Agent），保存活跃链接记录
	rhRepo := repository.GetRequestHistoryRepo(db.GetRedisClient())
	alRepo := repository.GetActiveLinkRepo(db.GetRedisClient())
	go func() {
		rhRepo.Save(models.NewRequestHistory(s, c.ClientIP(), c.Request.UserAgent()))
		alRepo.Save(s.Id)
	}()

	c.Redirect(http.StatusTemporaryRedirect, s.Url)
}
