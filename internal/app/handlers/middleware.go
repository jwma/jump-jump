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

// AuthContext carries the authenticated user and their tenant membership for the current request.
type AuthContext struct {
	User   *models.User
	Member *models.TenantMember
}

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

		userRepo := repository.GetUserRepo(db.GetPostgresPool())
		u, err := userRepo.FindByUsername(username)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}

		memberRepo := repository.GetTenantMemberRepo(db.GetPostgresPool())
		m, err := memberRepo.Get(tenantID, u.ID)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}

		c.Set("auth_context", &AuthContext{User: u, Member: m})
		c.Set("tenant_id", tenantID)
	}
}

type AuthAPIFunc func(c *gin.Context, ctx *AuthContext)

func Authenticator(f AuthAPIFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ac, exists := c.Get("auth_context")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}
		f(c, ac.(*AuthContext))
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
