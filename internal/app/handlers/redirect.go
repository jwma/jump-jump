package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/repository"
	"net/http"
)

func Redirect(c *gin.Context) {
	slRepo := repository.GetShortLinkRepo(db.GetRedisClient())
	s, err := slRepo.Get(c.Param("id"))
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	if !s.IsEnable {
		c.String(http.StatusOK, "你访问的页面不存在哦")
		return
	}

	// 保存短链接请求记录（IP、User-Agent）
	rhRepo := repository.GetRequestHistoryRepo(db.GetRedisClient())
	go rhRepo.Save(models.NewRequestHistory(s, c.ClientIP(), c.Request.UserAgent()))

	c.Redirect(http.StatusTemporaryRedirect, s.Url)
}
