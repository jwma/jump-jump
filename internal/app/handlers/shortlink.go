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

func GetShortLinkAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		slRepo := repository.GetShortLinkRepo(db.GetRedisClient())
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
		params := &models.CreateShortLinkParameter{}
		params.ShortLink = s

		if err := c.ShouldBindJSON(&params); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg":  "参数错误",
				"code": 4999,
				"data": nil,
			})
			return
		}

		repo := repository.GetShortLinkRepo(db.GetRedisClient())
		cfg := config.GetConfig()
		idLen := config.GetConfig().GetIntValue("idLength", 6)

		if user.Role == models.RoleUser {
			s.Id = "" // 如果是普通用户，创建时不可以指定 ID
		} else {
			if s.Id != "" { // 如果管理员指定了 ID，则检查 ID 是否可用
				checkShortLink, _ := repo.Get(s.Id)

				if checkShortLink.Id != "" {
					c.JSON(http.StatusOK, gin.H{
						"msg":  fmt.Sprintf("%s 已被占用，请使用其他 ID。", s.Id),
						"code": 4999,
						"data": nil,
					})
					return
				}
			} else { // 如果管理员没有指定 ID，则计算随机 ID 的长度
				idMinimumLength := cfg.GetIntValue("idMinimumLength", 2)
				idMaximumLength := cfg.GetIntValue("idMaximumLength", 10)

				if idMinimumLength <= params.IdLength && params.IdLength <= idMaximumLength { // 检查是否在合法的范围内
					idLen = params.IdLength
				}
			}
		}

		id, err := repo.GenerateId(idLen) // 生成 ID

		if err != nil {
			log.Printf("generate id failed, error: %v\n", err)
			c.JSON(http.StatusOK, gin.H{
				"msg":  "服务器繁忙，请稍后再试",
				"code": 4999,
				"data": nil,
			})
			return
		}

		s.Id = utils.TrimShortLinkId(id)

		if s.Id == "" {
			log.Println("短链接 ID 为空")
			c.JSON(http.StatusOK, gin.H{
				"msg":  "ID 错误",
				"code": 4999,
				"data": nil,
			})
		}

		err = repo.Save(s)
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
		slRepo := repository.GetShortLinkRepo(db.GetRedisClient())
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

		repo := repository.GetShortLinkRepo(db.GetRedisClient())
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
		slRepo := repository.GetShortLinkRepo(db.GetRedisClient())
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

		repo := repository.GetShortLinkRepo(db.GetRedisClient())
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

		if c.Param("action") == "/data" {
			slRepo := repository.GetShortLinkRepo(db.GetRedisClient())
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

			startDate := c.Query("startDate")
			endDate := c.Query("endDate")

			if startDate == "" || endDate == "" {
				c.JSON(http.StatusOK, gin.H{
					"msg":  "参数错误",
					"code": 4999,
					"data": nil,
				})
				return
			}

			startTime, err := time.ParseInLocation("2006-01-02", startDate, time.Local)
			endTime, err := time.ParseInLocation("2006-01-02", endDate, time.Local)

			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"msg":  "日期参数错误",
					"code": 4999,
					"data": nil,
				})
				return
			}

			endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 0, time.Local)
			rhRepo := repository.GetRequestHistoryRepo(db.GetRedisClient())
			rhs := rhRepo.FindByDateRange(s.Id, startTime, endTime)

			c.JSON(http.StatusOK, gin.H{
				"msg":  "ok",
				"code": 0,
				"data": gin.H{
					"histories": rhs,
				},
			})
			return
		} else if c.Param("action") == "/latest-request-history" {
			slRepo := repository.GetShortLinkRepo(db.GetRedisClient())
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

			repo := repository.GetRequestHistoryRepo(db.GetRedisClient())
			size := utils.GetIntQueryValue(c, "size", 20)

			r, err := repo.FindLatest(s.Id, int64(size))

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

		slRepo := repository.GetShortLinkRepo(db.GetRedisClient())
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
