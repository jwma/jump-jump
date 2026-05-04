package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jwma/jump-jump/internal/app/config"
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

// TenantResolverMiddleware resolves the tenant from the Host header.
func TenantResolverMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		host := strings.Split(c.Request.Host, ":")[0]

		tenantID, err := config.ResolveTenantID(host)
		if err != nil || tenantID == "" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"msg": "unknown domain"})
			return
		}

		c.Set("tenant_id", tenantID)
		c.Next()
	}
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

		tenantID, _ := claims["tenant_id"].(string)
		username, _ := claims["identifier"].(string)
		if tenantID == "" || username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}

		repo := repository.GetUserRepo(db.GetPostgresPool())
		u, err := repo.FindOneByUsername(tenantID, username)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}

		c.Set("user", u)
		c.Set("tenant_id", tenantID)
	}
}

type AuthAPIFunc func(c *gin.Context, user *models.User)

func Authenticator(f AuthAPIFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, exists := c.Get("user")
		if !exists {
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
				c.String(http.StatusBadRequest, "")
				c.Abort()
				return
			}
		}
	}
}
