package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/repository"
	"github.com/jwma/jump-jump/internal/app/utils"
	"slices"
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
		jwtStr, err := parseAuthorizationHeader(c.Request.Header.Get("Authorization"))
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}

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

		repo := repository.GetUserRepo(db.GetPostgresPool())
		u, err := repo.FindOneByUsername(claims["identifier"].(string))
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}

		c.Set("user", u)
	}
}

type AuthAPIFunc func(c *gin.Context, user *models.User)

func Authenticator(f AuthAPIFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, exists := c.Get("user")
		if !exists {
			log.Println("请求的 API Func 没有经过 JWTAuthenticatorMiddleware 处理，请修改路由设置")
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}
		user := u.(*models.User)
		f(c, user)
	}
}

func AllowedHostsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedHosts := os.Getenv("ALLOWED_HOSTS")

		if allowedHosts != "" && allowedHosts != "*" {
			h := strings.Split(c.Request.Host, ":")[0]

			if !slices.Contains(strings.Split(allowedHosts, ","), h) {
				output := ""

				if gin.Mode() == gin.DebugMode {
					output = fmt.Sprintf("You can see this message because GIN_MODE=debug.\n"+
						"Invalid HTTP_HOST header: '%s'. "+
						"You may need to add '%s' to ALLOWED_HOSTS environment variable.", c.Request.Host, h)
				}

				c.String(http.StatusBadRequest, output)
				c.Abort()
				return
			}
		}
	}
}
