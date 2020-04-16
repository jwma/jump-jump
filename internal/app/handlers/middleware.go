package handlers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/repository"
	"github.com/jwma/jump-jump/internal/app/utils"
	"github.com/thoas/go-funk"
	"log"
	"net/http"
	"os"
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
		repo := repository.GetUserRepo(db.GetRedisClient())
		u, err := repo.FindOneByUsername(claims["identifier"].(string))
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

// 检查当前请求的 Host 是否属于我们所设定的 Host 列表中的其中一个
// 如果不在设定列表中，则返回 HTTP Code 400 并中断后续逻辑的处理
func AllowedHostsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedHosts := os.Getenv("ALLOWED_HOSTS")

		if allowedHosts != "" && allowedHosts != "*" {
			h := strings.Split(c.Request.Host, ":")[0]

			if !funk.ContainsString(strings.Split(allowedHosts, ","), h) {
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
