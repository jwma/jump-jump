package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/repository"
	"net/http"
)

func Redirect(c *gin.Context) {
	l := &models.ShortLink{Id: c.Param("id")}
	err := l.Get()
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	if !l.IsEnable {
		c.String(http.StatusOK, "你访问的页面不存在哦")
		return
	}

	// 保存短链接请求记录（IP、User-Agent）
	repo := repository.NewRequestHistoryRepository(db.GetRedisClient())
	go repo.Save(models.NewRequestHistory(l, c.Request.RemoteAddr, c.Request.UserAgent()))

	c.Redirect(http.StatusTemporaryRedirect, l.Url)
}
