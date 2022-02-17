package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/repository"
	"github.com/jwma/jump-jump/internal/app/utils"
	"net/http"
	"strings"
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
	err := c.BindJSON(f)
	if err != nil {
		c.JSON(http.StatusOK, models.NewErrorResponse("用户名或密码错误"))
		return
	}

	repo := repository.GetUserRepo(db.GetRedisClient())
	u, err := repo.FindOneByUsername(strings.TrimSpace(f.Username))
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
		Token: utils.GenerateJWT(u.Username),
	}))
}

// GetUserInfoAPI godoc
// @Security ApiKeyAuth
// @Summary 获取账号信息
// @Description 获取账号信息
// @Tags 账号
// @Accept json
// @Produce json
// @Success 200 {object} models.Response{data=models.GetUserInfoAPIResponseData}
// @Failure 401
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
// @Security ApiKeyAuth
// @Summary 登出
// @Description 登出
// @Tags 账号
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Failure 401
// @Router /user/logout [post]
func LogoutAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		c.JSON(http.StatusOK, models.NewSuccessResponse(nil))
	})
}

// ChangePasswordAPI godoc
// @Security ApiKeyAuth
// @Summary 修改账号密码
// @Description 修改账号密码
// @Tags 账号
// @Accept json
// @Produce json
// @Param body body models.ChangePasswordAPIRequest true "修改密码请求"
// @Success 200 {object} models.Response
// @Failure 401
// @Router /user/change-password [post]
func ChangePasswordAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		if user.Username == "guest" {
			c.JSON(http.StatusOK, models.NewErrorResponse("该账号不支持修改密码"))
			return
		}

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

		repo := repository.GetUserRepo(db.GetRedisClient())
		err := repo.UpdatePassword(user)
		if err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(nil))
	})
}
