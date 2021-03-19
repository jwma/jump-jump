package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/config"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/repository"
	"github.com/jwma/jump-jump/internal/app/utils"
	"log"
	"net/http"
	"time"
)

// GetShortLinkAPI godoc
// @Security ApiKeyAuth
// @Summary 获取指定 ID 短链接
// @Description 获取系统配置信息
// @Tags 短链接
// @Accept json
// @Produce json
// @Param id path string true "短链接 ID"
// @Success 200 {object} models.Response{data=models.GetShortLinkAPIResponseData}
// @Failure 401
// @Router /short-link/{id} [get]
func GetShortLinkAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		slRepo := repository.GetShortLinkRepo(db.GetRedisClient())
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
// @Security ApiKeyAuth
// @Summary 创建短链接
// @Description 创建短链接
// @Tags 短链接
// @Accept json
// @Produce json
// @Param body body models.CreateShortLinkAPIRequest true "创建短链接请求"
// @Success 200 {object} models.Response{data=models.CreateShortLinkAPIResponseData}
// @Failure 401
// @Router /short-link/ [post]
func CreateShortLinkAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		var err error
		params := &models.CreateShortLinkAPIRequest{}

		if err := c.ShouldBindJSON(&params); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse("参数错误"))
			return
		}

		s := models.NewShortLink(user.Username, params)
		repo := repository.GetShortLinkRepo(db.GetRedisClient())
		idCfg := config.GetIdConfig()
		idLen := idCfg.IdLength

		if user.Role == models.RoleUser {
			s.Id = "" // 如果是普通用户，创建时不可以指定 ID
		}

		if s.Id != "" {
			// 如果管理员指定了 ID，则检查 ID 是否可用
			checkShortLink, _ := repo.Get(s.Id)

			if checkShortLink.Id != "" {
				c.JSON(http.StatusOK, models.NewErrorResponse(fmt.Sprintf("%s 已被占用，请使用其他 ID。", s.Id)))
				return
			}
		} else {
			// 如果管理员没有指定 ID，则计算随机 ID 的长度
			if idCfg.IdMinimumLength <= params.IdLength && params.IdLength <= idCfg.IdMaximumLength { // 检查是否在合法的范围内
				idLen = params.IdLength
			}

			id, err := repo.GenerateId(idLen) // 生成指定长度的随机 ID

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

// UpdateShortLinkAPI godoc
// @Security ApiKeyAuth
// @Summary 更新短链接
// @Description 更新短链接
// @Tags 短链接
// @Accept json
// @Produce json
// @Param id path string true "短链接 ID"
// @Param body body models.UpdateShortLinkAPIRequest true "更新短链接请求"
// @Success 200 {object} models.Response{data=models.UpdateShortLinkAPIResponseData}
// @Failure 401
// @Router /short-link/{id} [patch]
func UpdateShortLinkAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		slRepo := repository.GetShortLinkRepo(db.GetRedisClient())
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

		repo := repository.GetShortLinkRepo(db.GetRedisClient())
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

// DeleteShortLinkAPI godoc
// @Security ApiKeyAuth
// @Summary 删除短链接
// @Description 删除短链接
// @Tags 短链接
// @Accept json
// @Produce json
// @Param id path string true "短链接 ID"
// @Success 200 {object} models.Response
// @Failure 401
// @Router /short-link/{id} [delete]
func DeleteShortLinkAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		slRepo := repository.GetShortLinkRepo(db.GetRedisClient())
		s, err := slRepo.Get(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
			return
		}

		if !user.IsAdmin() && user.Username != s.CreatedBy {
			c.JSON(http.StatusOK, models.NewErrorResponse("你无权删除此短链接"))
			return
		}

		repo := repository.GetShortLinkRepo(db.GetRedisClient())
		repo.Delete(s)
		c.JSON(http.StatusOK, models.NewSuccessResponse(nil))
	})
}

// ListShortLinksAPI godoc
// @Security ApiKeyAuth
// @Summary 短链接列表
// @Description 短链接列表
// @Tags 短链接
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页条数"
// @Success 200 {object} models.Response{data=models.ListShortLinksAPIResponseData}
// @Failure 401
// @Router /short-link/ [get]
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

		slRepo := repository.GetShortLinkRepo(db.GetRedisClient())
		result, err := slRepo.List(key, start, stop)
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

// ShortLinkDataAPI godoc
// @Security ApiKeyAuth
// @Summary 短链接访问数据
// @Description 可查询短链接某个日期范围内的访问数据
// @Tags 短链接
// @Accept json
// @Produce json
// @Param id path string true "短链接 ID"
// @Param startDate query string true "开始日期 YYYY-mm-dd"
// @Param endDate query string true "结束日期 YYYY-mm-dd"
// @Success 200 {object} models.Response{data=models.ShortLinkDataAPIResponseData}
// @Failure 401
// @Router /short-link/{id}/data [get]
func ShortLinkActionAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {

		if c.Param("action") == "/data" {
			slRepo := repository.GetShortLinkRepo(db.GetRedisClient())
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
