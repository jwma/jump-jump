package controllers

import (
	"github.com/astaxie/beego"
	"github.com/jwma/jump-jump/app/db"
	"github.com/go-redis/redis"
	"golang.org/x/crypto/scrypt"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Post() {
	username := c.GetString("username")
	password := c.GetString("password")

	client := db.GetRedisClient()
	user, err := client.HGetAll("u:" + username).Result()
	if err != nil {
		switch err {
		case redis.Nil:
			c.Data["json"] = map[string]interface{}{"code": 4999, "msg": "用户名或密码错误"}
			c.ServeJSON()
			return
		default:
			beego.Error(err)
			c.Data["json"] = map[string]interface{}{"code": 4999, "msg": "服务繁忙，请稍后重试..."}
			c.ServeJSON()
			return
		}
	}

	dk, _ := scrypt.Key([]byte(password), []byte(user["Salt"]), 1<<15, 8, 1, 32)
	if user["Password"] != string(dk) {
		c.Data["json"] = map[string]interface{}{"code": 4999, "msg": "用户名或密码错误"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{"code": 0, "msg": "ok"}
	c.ServeJSON()
}
