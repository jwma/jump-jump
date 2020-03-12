package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/models"
	"log"
	"net/http"
)

func GetShortLinkAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		log.Println("当前请求用户：" + user.Username)
		l := &models.ShortLink{Id: c.Param("id")}
		err := l.Get()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": "ok",
			"data": gin.H{
				"shortLink": l,
			},
		})
	})
}

func CreateShortLinkAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		l := models.ShortLink{}

		if err := c.BindJSON(&l); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "参数错误",
			})
			return
		}

		err := l.GenerateId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		l.CreatedBy = user.Username

		err = l.Save()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg":  "ok",
			"data": gin.H{"shortLink": l},
		})
	})
}

func UpdateShortLinkAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		l := &models.ShortLink{Id: c.Param("id")}
		err := l.Get()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": err.Error(),
			})
			return
		}

		updateShortLink := &models.UpdateShortLinkParameter{}
		if err := c.ShouldBindJSON(updateShortLink); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}

		err = l.Update(updateShortLink)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{
			"msg":  "ok",
			"data": gin.H{"shortLink": l},
		})
	})
}

func ShortLinkActionAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		if c.Param("action") == "/history" {
			rh := &models.RequestHistory{}
			rh.SetLink(&models.ShortLink{Id: c.Param("id")})
			histories, _ := rh.GetAll()

			c.JSON(http.StatusOK, gin.H{
				"msg":  "ok",
				"data": gin.H{"histories": histories},
			})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{})
	})
}
