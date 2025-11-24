package routers

import (
	"beego-jwt-vergo/controllers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {

	// Root for testing
	beego.Get("/", func(ctx *context.Context) {
		_ = ctx.Output.Body([]byte("Beego JWT API is running"))
	})

	// Routes sans namespace
	beego.Router("/register", &controllers.UserController{}, "post:RegisterUser")
	beego.Router("/login", &controllers.UserController{}, "post:LoginUser")
	beego.Router("/users", &controllers.UserController{}, "get:IndexAll")

	// Version API v1 (cleaner)
	apiV1 := beego.NewNamespace("/v1",

		// /v1/register
		beego.NSRouter("/register", &controllers.UserController{}, "post:RegisterUser"),

		// /v1/login
		beego.NSRouter("/login", &controllers.UserController{}, "post:LoginUser"),

		// /v1/users
		beego.NSRouter("/users", &controllers.UserController{}, "get:IndexAll"),
	)

	// enable namespace
	beego.AddNamespace(apiV1)
}
