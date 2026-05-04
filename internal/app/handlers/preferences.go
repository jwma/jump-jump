package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/repository"
)

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
