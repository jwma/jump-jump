package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/config"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/repository"
	"github.com/jwma/jump-jump/internal/app/utils"
)

func GetShortLinkAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		slRepo := repository.GetShortLinkRepo(db.GetPostgresPool(), db.GetRedisClient())
		s, err := slRepo.Get(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		if !user.IsAdmin() && user.Username != s.CreatedBy {
			c.JSON(http.StatusOK, models.NewErrorResponse("你无权查看"))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(&models.GetShortLinkAPIResponseData{
			ShortLinkData: models.ToShortLinkData(s),
		}))
	})
}

func CreateShortLinkAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		var err error
		params := &models.CreateShortLinkAPIRequest{}

		if err := c.ShouldBindJSON(&params); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse("参数错误"))
			return
		}

		s := models.NewShortLink(user.Username, params)
		repo := repository.GetShortLinkRepo(db.GetPostgresPool(), db.GetRedisClient())
		idCfg := config.GetIdConfig()
		idLen := idCfg.IdLength

		if user.Role == models.RoleUser {
			s.Id = ""
		}

		if s.Id != "" {
			checkShortLink, _ := repo.Get(s.Id)
			if checkShortLink.Id != "" {
				c.JSON(http.StatusOK, models.NewErrorResponse(fmt.Sprintf("%s 已被占用，请使用其他 ID。", s.Id)))
				return
			}
		} else {
			if idCfg.IdMinimumLength <= params.IdLength && params.IdLength <= idCfg.IdMaximumLength {
				idLen = params.IdLength
			}

			id, err := repo.GenerateId(idLen)
			if err != nil {
				log.Printf("generate id failed, error: %v\n", err)
				c.JSON(http.StatusOK, models.NewErrorResponse("服务器繁忙，请稍后再试"))
				return
			}

			s.Id = utils.TrimShortLinkId(id)
		}

		if s.Id == "" {
			log.Println("短链接 ID 为空")
			c.JSON(http.StatusOK, models.NewErrorResponse("ID 错误"))
		}

		err = repo.Save(s)
		if err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(models.CreateShortLinkAPIResponseData{
			ShortLinkData: models.ToShortLinkData(s),
		}))
	})
}

func UpdateShortLinkAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		slRepo := repository.GetShortLinkRepo(db.GetPostgresPool(), db.GetRedisClient())
		s, err := slRepo.Get(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		if !user.IsAdmin() && user.Username != s.CreatedBy {
			c.JSON(http.StatusOK, models.NewErrorResponse("你无权修改此短链接"))
			return
		}

		updateShortLink := &models.UpdateShortLinkAPIRequest{}
		if err := c.ShouldBindJSON(updateShortLink); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		repo := repository.GetShortLinkRepo(db.GetPostgresPool(), db.GetRedisClient())
		err = repo.Update(s, updateShortLink)
		if err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(models.UpdateShortLinkAPIResponseData{
			ShortLinkData: models.ToShortLinkData(s),
		}))
	})
}

func DeleteShortLinkAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		slRepo := repository.GetShortLinkRepo(db.GetPostgresPool(), db.GetRedisClient())
		s, err := slRepo.Get(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		if !user.IsAdmin() && user.Username != s.CreatedBy {
			c.JSON(http.StatusOK, models.NewErrorResponse("你无权删除此短链接"))
			return
		}

		repo := repository.GetShortLinkRepo(db.GetPostgresPool(), db.GetRedisClient())
		repo.Delete(s)
		c.JSON(http.StatusOK, models.NewSuccessResponse(nil))
	})
}

func ListShortLinksAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		var page = utils.GetIntQueryValue(c, "page", 1)
		var pageSize = utils.GetIntQueryValue(c, "pageSize", 20)
		start := int64((page - 1) * pageSize)

		slRepo := repository.GetShortLinkRepo(db.GetPostgresPool(), db.GetRedisClient())
		result, err := slRepo.List(user.Username, user.IsAdmin(), start, int64(pageSize))
		if err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(&models.ListShortLinksAPIResponseData{
			ShortLinks: models.ToShortLinkDataSlice(result.ShortLinks),
			Total:      result.Total,
		}))
	})
}

func ShortLinkActionAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		if c.Param("action") == "/data" {
			slRepo := repository.GetShortLinkRepo(db.GetPostgresPool(), db.GetRedisClient())
			s, err := slRepo.Get(c.Param("id"))

			if err != nil {
				c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
				return
			}

			if !user.IsAdmin() && user.Username != s.CreatedBy {
				c.JSON(http.StatusOK, models.NewErrorResponse("你无权查看"))
				return
			}

			startDate := c.Query("startDate")
			endDate := c.Query("endDate")

			if startDate == "" || endDate == "" {
				c.JSON(http.StatusOK, models.NewErrorResponse("参数错误"))
				return
			}

			startTime, err := time.ParseInLocation("2006-01-02", startDate, time.Local)
			endTime, err := time.ParseInLocation("2006-01-02", endDate, time.Local)

			if err != nil {
				c.JSON(http.StatusOK, models.NewErrorResponse("日期参数错误"))
				return
			}

			endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 0, time.Local)
			rhRepo := repository.GetRequestHistoryRepo(db.GetRedisClient())
			rhs := rhRepo.FindByDateRange(s.Id, startTime, endTime)

			c.JSON(http.StatusOK, models.NewSuccessResponse(&models.ShortLinkDataAPIResponseData{Histories: rhs}))
			return
		} else if c.Param("action") == "/" {
			GetShortLinkAPI()(c)
			return
		}

		c.JSON(http.StatusOK, models.NewErrorResponse("请求资源不存在"))
	})
}
