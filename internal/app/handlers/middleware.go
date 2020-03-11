package handlers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/utils"
	"log"
	"net/http"
	"strings"
)

func parseAuthorizationHeader(a string) (string, error) {
	if a == "" {
		return "", fmt.Errorf("authorization 为空字符串")
	}
	t := strings.Split(a, " ")
	if len(t) < 2 {
		return "", fmt.Errorf("authorization 格式不正确")
	}
	return t[1], nil
}

func JWTAuthenticatorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization 提取 JWT
		jwtStr, err := parseAuthorizationHeader(c.Request.Header.Get("Authorization"))
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}

		// 校验 JWT
		claims := jwt.MapClaims{}
		_, err = jwt.ParseWithClaims(jwtStr, claims, func(*jwt.Token) (interface{}, error) {
			return []byte(utils.SecretKey), nil
		})
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}

		// 获取用户
		u := &models.User{Username: claims["identifier"].(string)}
		err = u.Get()
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}

		// 把当前请求用户保存到请求的上下文中
		c.Set("user", u)
	}
}
