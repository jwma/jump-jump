package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/models"
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
	h := models.NewRequestHistory(l, c.Request.RemoteAddr, c.Request.UserAgent())
	go h.Save() // 因为需要继续处理重定向，所以保存请求记录失败不做处理

	c.Redirect(http.StatusTemporaryRedirect, l.Url)
}
