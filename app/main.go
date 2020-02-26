package main

import (
	"flag"
	"github.com/astaxie/beego"
	"github.com/fatih/structs"
	"github.com/jwma/jump-jump/app/db"
	"github.com/jwma/jump-jump/app/models"
	_ "github.com/jwma/jump-jump/app/routers"
	"github.com/jwma/jump-jump/app/utils"
	"time"
)

func main() {
	// 如果传入用户名密码，则表示是通过命令行创建管理员，不启动beego
	username := flag.String("username", "", "A admin user username.")
	password := flag.String("password", "", "A admin user password.")
	flag.Parse()

	if *username != "" && *password != "" {
		salt, _ := utils.RandomSalt(32)
		dk, _ := utils.EncodePassword([]byte(*password), salt)
		user := models.User{
			Username:  *username,
			Password:  string(dk),
			Salt:      string(salt),
			CreatedAt: time.Now().Unix(),
		}
		client := db.GetRedisClient()
		client.HMSet("u:"+*username, structs.Map(user))
		return
	}

	beego.Run()
}
