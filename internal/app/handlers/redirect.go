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
		c.JSON(http.StatusNotFound, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.Redirect(http.StatusMovedPermanently, l.Url)
}
