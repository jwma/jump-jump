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

// TenantResolverMiddleware is deprecated — kept only for reference.
// Apiserver no longer resolves tenants by domain. Use TenantContextMiddleware instead.
func TenantResolverMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		userID, _ := claims["user_id"].(string)
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}

		userRepo := repository.GetUserRepo(db.GetPostgresPool())
		u, err := userRepo.FindByID(userID)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}

		c.Set("auth_context", &AuthContext{User: u, Member: nil})
		c.Next()
	}
}

// TenantContextMiddleware reads X-Tenant-ID, validates membership, sets tenant_id and auth_context.
func TenantContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ac, exists := c.Get("auth_context")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}
		ctx := ac.(*AuthContext)

		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			tenantID = c.Query("tenant_id")
		}
		if tenantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "X-Tenant-ID header is required"})
			c.Abort()
			return
		}

		var member *models.TenantMember
		if ctx.User.IsSuper {
			member = &models.TenantMember{
				TenantID: tenantID, UserID: ctx.User.ID, Role: models.RoleAdmin,
			}
		} else {
			memberRepo := repository.GetTenantMemberRepo(db.GetPostgresPool())
			m, err := memberRepo.Get(tenantID, ctx.User.ID)
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{"msg": "你不是该租户的成员"})
				c.Abort()
				return
			}
			member = m
		}

		c.Set("tenant_id", tenantID)
		c.Set("auth_context", &AuthContext{User: ctx.User, Member: member})
		c.Next()
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

func getTenantID(c *gin.Context) string {
	tid, _ := c.Get("tenant_id")
	s, _ := tid.(string)
	return s
}

func SuperAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ac, exists := c.Get("auth_context")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}
		ctx := ac.(*AuthContext)
		if !ctx.User.IsSuper {
			c.JSON(http.StatusOK, models.NewErrorResponse("仅超级管理员可执行此操作"))
			c.Abort()
			return
		}
		c.Next()
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
