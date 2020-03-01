package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/models"
	"net/http"
)

func GetShortLink(c *gin.Context) {
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
}

func CreateShortLink(c *gin.Context) {
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

	l.CreatedBy = "admin"

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
}

func UpdateShortLink(c *gin.Context) {
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
}

func ShortLinkActionHandler(c *gin.Context) {
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
}
