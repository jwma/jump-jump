package controllers

import (
	"github.com/astaxie/beego"
	"github.com/jwma/jump-jump/app/models"
	"github.com/jwma/jump-jump/app/db"
	"time"
	"github.com/jwma/jump-jump/app/utils"
	"github.com/go-redis/redis"
	"encoding/json"
)

// 生成随机且唯一的slug
func generateUniqueSlug() (string, error) {
	client := db.GetRedisClient()

	for true {
		slug := utils.RandStringRunes(6)
		_, err := client.Get("l:" + slug).Result()
		if err != nil {
			switch err {
			case redis.Nil:
				return slug, nil
			default:
				beego.Error(err)
				return "", err
			}
		}
	}
	return "", nil
}

type LinkController struct {
	beego.Controller
}

func (c *LinkController) Post() {
	url := c.GetString("url")
	isEnabled, _ := c.GetBool("isEnabled")
	description := c.GetString("description")

	if url == "" {
		c.Data["json"] = map[string]interface{}{"code": 4999, "msg": "请填写URL"}
		c.ServeJSON()
		return
	}

	slug, err := generateUniqueSlug()
	if err != nil {
		c.Data["json"] = map[string]interface{}{"code": 4999, "msg": "服务繁忙，请稍后重试..."}
		c.ServeJSON()
		return
	}

	now := time.Now().Unix()
	link := &models.Link{
		Slug:        slug,
		Url:         url,
		IsEnabled:   isEnabled,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	client := db.GetRedisClient()
	linkJson, err := json.Marshal(link)
	if err != nil {
		beego.Error(err)
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": "服务繁忙，请稍后重试..."}
		c.ServeJSON()
		return
	}
	client.Set("l:"+slug, string(linkJson), 0)
	c.Data["json"] = map[string]interface{}{"code": 0, "msg": "ok", "slug": slug}
	c.ServeJSON()
}
