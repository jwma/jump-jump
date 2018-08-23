package controllers

import (
	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
	"time"
	"encoding/json"
)

// 短链接请求记录结构
type RequestRecord struct {
	RemoteAddr string
	UserAgent  string
	RequestAt  int64
}

func getRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:        "localhost:6379",
		Password:    "",
		DB:          0,
		IdleTimeout: -1,
	})
	return client
}

type JumpController struct {
	beego.Controller
}

func (c *JumpController) Get() {
	slug := c.Ctx.Input.Param(":slug")

	client := getRedisClient()
	targetUrl, err := client.Get(slug).Result()
	if err != nil {
		beego.Warn(err)
		c.Ctx.WriteString("访问的链接不存在")
		return
	}

	requestRecord := &RequestRecord{
		c.Ctx.Request.RemoteAddr,
		c.Ctx.Request.UserAgent(),
		time.Now().Unix(),
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
