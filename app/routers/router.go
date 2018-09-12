package routers

import (
	"github.com/jwma/jump-jump/app/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/dgrijalva/jwt-go"
)

func init() {
	// 根据不同短链接的slug进行跳转
	beego.Router("/:slug:string", &controllers.JumpController{})

	// 欢迎页面
	beego.Get("/", func(c *context.Context) {
		c.WriteString("码极工作室/MJ STUDIO")
	})

	// 后台API
	beego.Router("/admin/link", &controllers.LinkController{})

	// 登录态过滤器
	beego.InsertFilter("/admin/*", beego.BeforeRouter, func(ctx *context.Context) {
		authorization := ctx.Input.Header("Authorization")
		needLoginJson := map[string]interface{}{"code": 40001, "msg": "请登陆"}
		if authorization == "" {
			ctx.Output.JSON(needLoginJson, false, true)
		}

		secretKey := beego.AppConfig.String("secret_key")
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(authorization, claims, func(*jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil {
			ctx.Output.JSON(needLoginJson, false, true)
		}
		ctx.Input.SetData("username", claims["username"])
	})

	// 登录
	beego.Router("/login", &controllers.LoginController{})

	// 检查登录态
	beego.Post("/check-login", func(ctx *context.Context) {
		authorization := ctx.Input.Header("Authorization")
		needLoginJson := map[string]interface{}{"code": 40001, "msg": "请登陆"}

		if authorization == "" {
			ctx.Output.JSON(needLoginJson, false, true)
		}

		secretKey := beego.AppConfig.String("secret_key")
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(authorization, claims, func(*jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil {
			ctx.Output.JSON(needLoginJson, false, true)
		} else {
			ctx.Output.JSON(map[string]interface{}{"code": 0, "msg": "ok"}, false, true)
		}
	})
}
