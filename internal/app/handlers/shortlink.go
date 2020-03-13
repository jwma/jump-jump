package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/models"
	"log"
	"net/http"
	"strconv"
)

func GetShortLinkAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		l := &models.ShortLink{Id: c.Param("id")}
		err := l.Get()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": err.Error(),
			})
			return
		}

		if !user.IsAdmin() && user.Username != l.CreatedBy {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "你无权查看",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg": "ok",
			"data": gin.H{
				"shortLink": l,
			},
		})
	})
}

func CreateShortLinkAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		l := models.ShortLink{}

		if err := c.BindJSON(&l); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "参数错误",
			})
			return
		}

		err := l.GenerateId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		l.CreatedBy = user.Username
		err = l.Save()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg":  "ok",
			"data": gin.H{"shortLink": l},
		})
	})
}

func UpdateShortLinkAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		l := &models.ShortLink{Id: c.Param("id")}
		err := l.Get()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": err.Error(),
			})
			return
		}

		if !user.IsAdmin() && user.Username != l.CreatedBy {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "你无权修改",
			})
			return
		}

		updateShortLink := &models.UpdateShortLinkParameter{}
		if err := c.ShouldBindJSON(updateShortLink); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}

		err = l.Update(updateShortLink)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{
			"msg":  "ok",
			"data": gin.H{"shortLink": l},
		})
	})
}

func DeleteShortLinkAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		l := &models.ShortLink{Id: c.Param("id")}
		err := l.Get()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": err.Error(),
			})
			return
		}

		if !user.IsAdmin() && user.Username != l.CreatedBy {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "你无权修改",
			})
			return
		}

		l.Delete()
		c.JSON(http.StatusOK, gin.H{})
	})
}

func ShortLinkActionAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		if c.Param("action") == "/history" {
			l := &models.ShortLink{Id: c.Param("id")}
			err := l.Get()
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"msg": err.Error(),
				})
				return
			}

			if !user.IsAdmin() && user.Username != l.CreatedBy {
				c.JSON(http.StatusForbidden, gin.H{
					"msg": "你无权查看",
				})
				return
			}

			rh := &models.RequestHistory{}
			rh.SetLink(l)
			histories, _ := rh.GetAll()

			c.JSON(http.StatusOK, gin.H{
				"msg":  "ok",
				"data": gin.H{"histories": histories},
			})
			return
		} else if c.Param("action") == "/" {
			GetShortLinkAPI()(c)
			return
		}
		c.JSON(http.StatusNotFound, gin.H{})
	})
}

func ListShortLinksAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, user *models.User) {
		var page = 1
		var pageSize int64 = 20
		var err error

		if c.Query("page") != "" {
			page, err = strconv.Atoi(c.Query("page"))
			if err != nil {
				page = 1
			}
		}
		start := int64(page-1) * pageSize
		stop := start - 1 + pageSize

		client := db.GetRedisClient()
		var key string
		if user.IsAdmin() {
			key = "links"
		} else {
			key = fmt.Sprintf("links:%s", user.Username)
		}

		ids, err := client.ZRevRange(key, start, stop).Result()
		if err != nil {
			log.Printf("fail to list short links, err: %v\n", err)
			c.JSON(http.StatusOK, gin.H{
				"code": 499,
				"msg":  "系统繁忙请稍后再试...",
			})
			return
		}

		if len(ids) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"data": []string{},
			})
			return
		}

		linkRs := make([]*redis.StringCmd, 0)
		p := client.Pipeline()
		for _, id := range ids {
			r := p.Get(fmt.Sprintf("link:%s", id))
			linkRs = append(linkRs, r)

		}
		_, _ = p.Exec()

		links := make([]*models.ShortLink, 0)
		for _, cmd := range linkRs {
			l := &models.ShortLink{}
			err = json.Unmarshal([]byte(cmd.Val()), l)
			links = append(links, l)
		}

		total, _ := client.ZCard(key).Result()

		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"total":      total,
				"shortLinks": links,
			},
		})
	})
}
