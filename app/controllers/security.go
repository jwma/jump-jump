package controllers

import (
	"github.com/astaxie/beego"
	"github.com/jwma/jump-jump/app/db"
	"github.com/go-redis/redis"
	"github.com/jwma/jump-jump/app/utils"
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

	dk, _ := utils.EncodePassword([]byte(password), []byte(user["Salt"]))
	if user["Password"] != string(dk) {
		c.Data["json"] = map[string]interface{}{"code": 4999, "msg": "用户名或密码错误"}
		c.ServeJSON()
		return
	}

	token := utils.GenerateJWT(user["Username"])
	c.Data["json"] = map[string]interface{}{"code": 0, "msg": "ok", "token": token}
	c.ServeJSON()
}
