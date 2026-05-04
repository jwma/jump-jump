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
// @Router /user/login [post]
func LoginAPI(c *gin.Context) {
	f := &models.LoginAPIRequest{}
	if err := c.BindJSON(f); err != nil {
		c.JSON(http.StatusOK, models.NewErrorResponse("用户名或密码错误"))
		return
	}

	tenantID, _ := c.Get("tenant_id")
	tid, _ := tenantID.(string)
	if tid == "" {
		c.JSON(http.StatusOK, models.NewErrorResponse("无法识别租户"))
		return
	}

	repo := repository.GetUserRepo(db.GetPostgresPool())
	u, err := repo.FindOneByUsername(tid, strings.TrimSpace(f.Username))
	if err != nil {
		c.JSON(http.StatusOK, models.NewErrorResponse("用户名或密码错误"))
		return
	}

	dk, _ := utils.EncodePassword([]byte(f.Password), u.Salt)
	if string(u.Password) != string(dk) {
		c.JSON(http.StatusOK, models.NewErrorResponse("用户名或密码错误"))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(models.LoginAPIResponseData{
		Token: utils.GenerateJWT(u.Username, u.TenantID),
	}))
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
	return Authenticator(func(c *gin.Context, user *models.User) {
		c.JSON(http.StatusOK, models.NewSuccessResponse(models.GetUserInfoAPIResponseData{
			Username: user.Username,
			Role:     user.Role,
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
	return Authenticator(func(c *gin.Context, user *models.User) {
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
	return Authenticator(func(c *gin.Context, user *models.User) {
		p := &models.ChangePasswordAPIRequest{}
		if err := c.ShouldBindJSON(p); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse("请填写原密码和新密码"))
			return
		}

		dk, _ := utils.EncodePassword([]byte(p.Password), user.Salt)
		if string(user.Password) != string(dk) {
			c.JSON(http.StatusOK, models.NewErrorResponse("原密码错误"))
			return
		}

		user.RawPassword = p.NewPassword
		repo := repository.GetUserRepo(db.GetPostgresPool())
		if err := repo.UpdatePassword(user); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(nil))
	})
}
