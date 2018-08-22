package controllers

import (
	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
)

func getRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
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
		c.Ctx.WriteString("访问的链接不存在")
		return
	}

	c.Redirect(targetUrl, 302)
}
