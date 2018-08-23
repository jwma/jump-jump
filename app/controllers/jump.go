package controllers

import (
	"github.com/astaxie/beego"
	"time"
	"encoding/json"
	"github.com/jwma/jump-jump/app/models"
	"github.com/jwma/jump-jump/app/db"
)

type JumpController struct {
	beego.Controller
}

func (c *JumpController) Get() {
	slug := c.Ctx.Input.Param(":slug")

	client := db.GetRedisClient()
	targetUrl, err := client.Get(slug).Result()
	if err != nil {
		beego.Error(err)
		c.Ctx.WriteString("访问的链接不存在")
		return
	}

	requestRecord := &models.RequestRecord{
		RemoteAddr: c.Ctx.Request.RemoteAddr,
		UserAgent:  c.Ctx.Request.UserAgent(),
		RequestAt:  time.Now().Unix(),
	}
	requestRecordJson, err := json.Marshal(requestRecord)
	if err != nil {
		beego.Error(err)
	}

	// 记录数据
	pipe := client.Pipeline()
	pipe.Incr("c:" + slug)
	pipe.LPush("r:"+slug, string(requestRecordJson))
	_, err = pipe.Exec()
	if err != nil {
		beego.Error(err)
	}

	c.Redirect(targetUrl, 302)
}
