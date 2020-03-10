package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/utils"
	"net/http"
)

type loginForm struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	f := &loginForm{}
	err := c.BindJSON(f)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg":  "用户名或密码错误",
			"data": nil,
		})
		return
	}

	u := &models.User{Username: f.Username}
	err = u.Get()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg":  "用户名或密码错误",
			"data": nil,
		})
		return
	}

	dk, _ := utils.EncodePassword([]byte(f.Password), u.Salt)
	if string(u.Password) != string(dk) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg":  "用户名或密码错误11",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
		"data": gin.H{
			"token":    utils.GenerateJWT(u.Username),
			"username": u.Username,
		},
	})
}
