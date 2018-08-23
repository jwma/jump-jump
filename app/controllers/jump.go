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
	linkJson, err := client.Get("l:" + slug).Result()
	if err != nil {
		beego.Error(err)
		c.Ctx.WriteString("链接不存在")
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

	var link models.Link
	if err := json.Unmarshal([]byte(linkJson), &link); err != nil {
		beego.Error(err)
		c.Ctx.WriteString("链接不存在")
		return
	}
	if !link.IsEnabled {
		c.Ctx.WriteString("链接已失效")
		return
	}

	c.Redirect(link.Url, 302)
}
