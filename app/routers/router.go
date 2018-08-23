package routers

import (
	"github.com/jwma/jump-jump/app/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	// 根据不同短链接的slug进行跳转
	beego.Router("/:slug:string", &controllers.JumpController{})

	// 欢迎页面
	beego.Get("/", func(c *context.Context) {
		c.WriteString("码极工作室/MJ STUDIO")
	})

	// 根据不同短链接的slug进行跳转
	beego.Router("/admin/link", &controllers.LinkController{})
}
