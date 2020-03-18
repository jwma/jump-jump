package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/repository"
	"github.com/jwma/jump-jump/internal/app/utils"
	"net/http"
)

func GetShortLinkAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		slRepo := repository.GetShortLinkRepo()
		s, err := slRepo.Get(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg":  err.Error(),
				"code": 4999,
				"data": nil,
			})
			return
		}

		if !user.IsAdmin() && user.Username != s.CreatedBy {
			c.JSON(http.StatusOK, gin.H{
				"msg":  "你无权查看",
				"code": 4999,
				"data": nil,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg":  "ok",
			"code": 0,
			"data": gin.H{
				"shortLink": s,
			},
		})
	})
}

func CreateShortLinkAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		s := &models.ShortLink{CreatedBy: user.Username}

		if err := c.BindJSON(&s); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg":  "参数错误",
				"code": 4999,
				"data": nil,
			})
			return
		}

		repo := repository.GetShortLinkRepo()
		err := repo.Save(s)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg":  err.Error(),
				"code": 4999,
				"data": nil,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg":  "ok",
			"code": 0,
			"data": gin.H{"shortLink": s},
		})
	})
}

func UpdateShortLinkAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		slRepo := repository.GetShortLinkRepo()
		s, err := slRepo.Get(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg":  err.Error(),
				"code": 4999,
				"data": nil,
			})
			return
		}

		if !user.IsAdmin() && user.Username != s.CreatedBy {
			c.JSON(http.StatusOK, gin.H{
				"msg":  "你无权修改",
				"code": 4999,
				"data": nil,
			})
			return
		}

		updateShortLink := &models.UpdateShortLinkParameter{}
		if err := c.ShouldBindJSON(updateShortLink); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg":  err.Error(),
				"code": 4999,
				"data": nil,
			})
			return
		}

		repo := repository.GetShortLinkRepo()
		err = repo.Update(s, updateShortLink)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg":  err.Error(),
				"code": 4999,
				"data": nil,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"msg":  "ok",
			"code": 0,
			"data": gin.H{"shortLink": s},
		})
	})
}

func DeleteShortLinkAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		slRepo := repository.GetShortLinkRepo()
		s, err := slRepo.Get(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg":  err.Error(),
				"code": 4999,
				"data": nil,
			})
			return
		}

		if !user.IsAdmin() && user.Username != s.CreatedBy {
			c.JSON(http.StatusOK, gin.H{
				"msg":  "你无权修改",
				"code": 4999,
				"data": nil,
			})
			return
		}

		repo := repository.GetShortLinkRepo()
		repo.Delete(s)
		c.JSON(http.StatusOK, gin.H{
			"msg":  "ok",
			"code": 0,
			"data": nil,
		})
	})
}

func ShortLinkActionAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		if c.Param("action") == "/latest-request-history" {
			slRepo := repository.GetShortLinkRepo()
			s, err := slRepo.Get(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"msg":  err.Error(),
					"code": 4999,
					"data": nil,
				})
				return
			}

			if !user.IsAdmin() && user.Username != s.CreatedBy {
				c.JSON(http.StatusOK, gin.H{
					"msg":  "你无权查看",
					"code": 4999,
					"data": nil,
				})
				return
			}

			repo := repository.GetRequestHistoryRepo()
			r, err := repo.FindLatest(s.Id, 20)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"msg":  err.Error(),
					"code": 4999,
					"data": nil,
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"msg":  "ok",
				"code": 0,
				"data": r,
			})
			return
		} else if c.Param("action") == "/" {
			GetShortLinkAPI()(c)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg":  "请求资源不存在",
			"code": 4999,
			"data": nil,
		})
	})
}

func ListShortLinksAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		var page = utils.GetIntQueryValue(c, "page", 1)
		var pageSize = utils.GetIntQueryValue(c, "pageSize", 20)
		start := int64((page - 1) * pageSize)
		stop := start - 1 + int64(pageSize)

		var key string
		if user.IsAdmin() {
			key = utils.GetShortLinksKey()
		} else {
			key = utils.GetUserShortLinksKey(user.Username)
		}

		slRepo := repository.GetShortLinkRepo()
		result, err := slRepo.List(key, start, stop)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg":  err.Error(),
				"code": 4999,
				"data": nil,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg":  "",
			"code": 0,
			"data": result,
		})
	})
}
