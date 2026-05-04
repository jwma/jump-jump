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

// GetShortLinkAPI godoc
// @Summary 获取指定 ID 短链接
// @Description 获取指定 ID 短链接详情
// @Tags 短链接
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "短链接 ID"
// @Success 200 {object} models.Response{data=models.GetShortLinkAPIResponseData}
// @Failure 401 {object} nil
// @Router /short-link/{id} [get]
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

// CreateShortLinkAPI godoc
// @Summary 创建短链接
// @Description 创建短链接
// @Tags 短链接
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body models.CreateShortLinkAPIRequest true "创建短链接请求"
// @Success 200 {object} models.Response{data=models.CreateShortLinkAPIResponseData}
// @Failure 401 {object} nil
// @Router /short-link/ [post]
func CreateShortLinkAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		params := &models.CreateShortLinkAPIRequest{}
		if err := c.ShouldBindJSON(&params); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse("参数错误"))
			return
		}

		tenantID := user.TenantID
		s := models.NewShortLink(tenantID, user.Username, params)
		repo := repository.GetShortLinkRepo(db.GetPostgresPool(), db.GetRedisClient())
		idCfg := config.GetIdConfig(tenantID)
		idLen := idCfg.IdLength

		if user.Role == models.RoleUser {
			s.Id = ""
		}

		if s.Id != "" {
			checkShortLink, _ := repo.Get(s.Id)
			if checkShortLink != nil && checkShortLink.Id != "" {
				c.JSON(http.StatusOK, models.NewErrorResponse(fmt.Sprintf("%s 已被占用", s.Id)))
				return
			}
		} else {
			if idCfg.IdMinimumLength <= params.IdLength && params.IdLength <= idCfg.IdMaximumLength {
				idLen = params.IdLength
			}
			id, err := repo.GenerateId(idLen)
			if err != nil {
				log.Printf("generate id failed: %v", err)
				c.JSON(http.StatusOK, models.NewErrorResponse("服务器繁忙，请稍后再试"))
				return
			}
			s.Id = utils.TrimShortLinkId(id)
		}

		if s.Id == "" {
			c.JSON(http.StatusOK, models.NewErrorResponse("ID 错误"))
			return
		}

		if err := repo.Save(s); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(models.CreateShortLinkAPIResponseData{
			ShortLinkData: models.ToShortLinkData(s),
		}))
	})
}

// UpdateShortLinkAPI godoc
// @Summary 更新短链接
// @Description 更新短链接
// @Tags 短链接
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "短链接 ID"
// @Param body body models.UpdateShortLinkAPIRequest true "更新短链接请求"
// @Success 200 {object} models.Response{data=models.UpdateShortLinkAPIResponseData}
// @Failure 401 {object} nil
// @Router /short-link/{id} [patch]
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

		if err := slRepo.Update(s, updateShortLink); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(models.UpdateShortLinkAPIResponseData{
			ShortLinkData: models.ToShortLinkData(s),
		}))
	})
}

// DeleteShortLinkAPI godoc
// @Summary 删除短链接
// @Description 删除短链接
// @Tags 短链接
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "短链接 ID"
// @Success 200 {object} models.Response
// @Failure 401 {object} nil
// @Router /short-link/{id} [delete]
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

		slRepo.Delete(s)
		c.JSON(http.StatusOK, models.NewSuccessResponse(nil))
	})
}

// ListShortLinksAPI godoc
// @Summary 短链接列表
// @Description 短链接列表
// @Tags 短链接
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码"
// @Param pageSize query int false "每页条数"
// @Success 200 {object} models.Response{data=models.ListShortLinksAPIResponseData}
// @Failure 401 {object} nil
// @Router /short-link/ [get]
func ListShortLinksAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		var page = utils.GetIntQueryValue(c, "page", 1)
		var pageSize = utils.GetIntQueryValue(c, "pageSize", 20)
		start := int64((page - 1) * pageSize)

		slRepo := repository.GetShortLinkRepo(db.GetPostgresPool(), db.GetRedisClient())
		result, err := slRepo.List(user.TenantID, user.Username, user.IsAdmin(), start, int64(pageSize))
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

// ShortLinkActionAPI godoc
// @Summary 短链接访问数据
// @Description 可查询短链接某个日期范围内的访问数据
// @Tags 短链接
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "短链接 ID"
// @Param startDate query string true "开始日期 YYYY-mm-dd"
// @Param endDate query string true "结束日期 YYYY-mm-dd"
// @Success 200 {object} models.Response{data=models.ShortLinkDataAPIResponseData}
// @Failure 401 {object} nil
// @Router /short-link/{id}/data [get]
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

			startTime, _ := time.ParseInLocation("2006-01-02", startDate, time.Local)
			endTime, err := time.ParseInLocation("2006-01-02", endDate, time.Local)
			if err != nil {
				c.JSON(http.StatusOK, models.NewErrorResponse("日期参数错误"))
				return
			}

			endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 0, time.Local)
			rhRepo := repository.GetRequestHistoryRepo(db.GetRedisClient(), db.GetPostgresPool())
			rhs := rhRepo.FindByDateRange(s.Id, startTime, endTime)
			daily, osDist := rhRepo.GetAggregatedStats(s.Id, startTime, endTime)

			c.JSON(http.StatusOK, models.NewSuccessResponse(&models.ShortLinkDataAPIResponseData{
				Histories: rhs,
				Daily:     daily,
				OSDist:    osDist,
			}))
			return
		} else if c.Param("action") == "/" {
			GetShortLinkAPI()(c)
			return
		}

		c.JSON(http.StatusOK, models.NewErrorResponse("请求资源不存在"))
	})
}
