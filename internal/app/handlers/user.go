package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/repository"
	"github.com/jwma/jump-jump/internal/app/utils"
)

// LoginAPI godoc
// @Summary 账号登入
// @Description 账号密码登入
// @Tags 账号
// @Accept json
// @Produce json
// @Param body body models.LoginAPIRequest true "登入请求"
// @Success 200 {object} models.Response{data=models.LoginAPIResponseData}
// @Router /auth/login [post]
func LoginAPI(c *gin.Context) {
	f := &models.LoginAPIRequest{}
	if err := c.BindJSON(f); err != nil {
		c.JSON(http.StatusOK, models.NewErrorResponse("用户名或密码错误"))
		return
	}

	userRepo := repository.GetUserRepo(db.GetPostgresPool())
	u, err := userRepo.FindByUsername(strings.TrimSpace(f.Username))
	if err != nil {
		c.JSON(http.StatusOK, models.NewErrorResponse("用户名或密码错误"))
		return
	}

	dk, _ := utils.EncodePassword([]byte(f.Password), u.Salt)
	if string(u.Password) != string(dk) {
		c.JSON(http.StatusOK, models.NewErrorResponse("用户名或密码错误"))
		return
	}

	if !u.IsActive {
		c.JSON(http.StatusOK, models.NewErrorResponse("账号已被禁用"))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(models.LoginAPIResponseData{
		Token: utils.GenerateJWT(u.ID, u.IsSuper),
	}))
}

// GetAuthInfoAPI godoc
// @Summary 获取当前用户信息及租户列表
// @Description 获取当前用户信息及租户列表
// @Tags 账号
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=models.AuthInfoResponseData}
// @Failure 401 {object} nil
// @Router /auth/info [get]
func GetAuthInfoAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, ctx *AuthContext) {
		userRepo := repository.GetUserRepo(db.GetPostgresPool())
		tenants, err := userRepo.GetUserTenants(ctx.User)
		if err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(models.AuthInfoResponseData{
			User: &models.AuthInfoUser{
				ID:       ctx.User.ID,
				Username: ctx.User.Username,
				IsSuper:  ctx.User.IsSuper,
			},
			Tenants: tenants,
		}))
	})
}

// GetUserInfoAPI godoc
// @Summary 获取账号信息
// @Description 获取账号信息
// @Tags 账号
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=models.GetUserInfoAPIResponseData}
// @Failure 401 {object} nil
// @Router /user/info [get]
func GetUserInfoAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, ctx *AuthContext) {
		role := ""
		if ctx.Member != nil {
			role = ctx.Member.Role
		}
		if ctx.User.IsSuper {
			role = models.RoleAdmin
		}
		c.JSON(http.StatusOK, models.NewSuccessResponse(models.GetUserInfoAPIResponseData{
			Username: ctx.User.Username,
			Role:     role,
			IsSuper:  ctx.User.IsSuper,
		}))
	})
}

// LogoutAPI godoc
// @Summary 登出
// @Description 登出
// @Tags 账号
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 401 {object} nil
// @Router /user/logout [post]
func LogoutAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, ctx *AuthContext) {
		c.JSON(http.StatusOK, models.NewSuccessResponse(nil))
	})
}

// ChangePasswordAPI godoc
// @Summary 修改账号密码
// @Description 修改账号密码
// @Tags 账号
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body models.ChangePasswordAPIRequest true "修改密码请求"
// @Success 200 {object} models.Response
// @Failure 401 {object} nil
// @Router /user/change-password [post]
func ChangePasswordAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, ctx *AuthContext) {
		p := &models.ChangePasswordAPIRequest{}
		if err := c.ShouldBindJSON(p); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse("请填写原密码和新密码"))
			return
		}

		dk, _ := utils.EncodePassword([]byte(p.Password), ctx.User.Salt)
		if string(ctx.User.Password) != string(dk) {
			c.JSON(http.StatusOK, models.NewErrorResponse("原密码错误"))
			return
		}

		ctx.User.RawPassword = p.NewPassword
		repo := repository.GetUserRepo(db.GetPostgresPool())
		if err := repo.UpdatePassword(ctx.User); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(nil))
	})
}
