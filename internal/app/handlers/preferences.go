package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/repository"
)

// GetUserPreferencesAPI godoc
// @Summary 获取用户偏好设置
// @Description 获取用户偏好设置
// @Tags 账号
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=map[string]interface{}} "data.preferences 为 []*models.UserPreference"
// @Failure 401 {object} nil
// @Router /user/preferences [get]
func GetUserPreferencesAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		repo := repository.GetUserPreferenceRepo(db.GetPostgresPool())
		prefs, err := repo.GetAll(user.TenantID, user.Username)
		if err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusOK, models.NewSuccessResponse(map[string]interface{}{
			"preferences": prefs,
		}))
	})
}

type updatePreferencesRequest struct {
	Preferences []*models.UserPreference `json:"preferences" binding:"required"`
}

// UpdateUserPreferencesAPI godoc
// @Summary 更新用户偏好设置
// @Description 更新用户偏好设置
// @Tags 账号
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body object true "更新偏好设置请求" example({"preferences":[{"key":"theme","value":"dark"}]})
// @Success 200 {object} models.Response{data=map[string]interface{}} "data.preferences 为 []*models.UserPreference"
// @Failure 401 {object} nil
// @Router /user/preferences [put]
func UpdateUserPreferencesAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		req := &updatePreferencesRequest{}
		if err := c.ShouldBindJSON(req); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse("参数错误"))
			return
		}

		repo := repository.GetUserPreferenceRepo(db.GetPostgresPool())
		if err := repo.Upsert(user.TenantID, user.Username, req.Preferences); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		prefs, _ := repo.GetAll(user.TenantID, user.Username)
		c.JSON(http.StatusOK, models.NewSuccessResponse(map[string]interface{}{
			"preferences": prefs,
		}))
	})
}
