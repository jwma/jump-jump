package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/models"
	"net/http"
)

func GetShortLinkAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		l := &models.ShortLink{Id: c.Param("id")}
		err := l.Get()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": err.Error(),
			})
			return
		}

		if !user.IsAdmin() && user.Username != l.CreatedBy {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "你无权查看",
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

		if !user.IsAdmin() && user.Username != l.CreatedBy {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "你无权修改",
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

func DeleteShortLinkAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		l := &models.ShortLink{Id: c.Param("id")}
		err := l.Get()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": err.Error(),
			})
			return
		}

		if !user.IsAdmin() && user.Username != l.CreatedBy {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "你无权修改",
			})
			return
		}

		l.Delete()
		c.JSON(http.StatusOK, gin.H{})
	})
}

func ShortLinkActionAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		if c.Param("action") == "/history" {
			l := &models.ShortLink{Id: c.Param("id")}
			err := l.Get()
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"msg": err.Error(),
				})
				return
			}

			if !user.IsAdmin() && user.Username != l.CreatedBy {
				c.JSON(http.StatusForbidden, gin.H{
					"msg": "你无权查看",
				})
				return
			}

			rh := &models.RequestHistory{}
			rh.SetLink(l)
			histories, _ := rh.GetAll()

			c.JSON(http.StatusOK, gin.H{
				"msg":  "ok",
				"data": gin.H{"histories": histories},
			})
			return
		} else if c.Param("action") == "/" {
			GetShortLinkAPI()(c)
			return
		}
		c.JSON(http.StatusNotFound, gin.H{})
	})
}
