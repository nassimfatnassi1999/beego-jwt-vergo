package routers

import (
	"beego-jwt-vergo/controllers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	// Root route for simple health-check
	beego.Get("/", func(ctx *context.Context) {
		_ = ctx.Output.Body([]byte("Beego JWT API is running"))
	})

	// Non-versioned routes
	beego.Router("/register", &controllers.UserController{}, "post:RegisterUser")
	beego.Router("/login", &controllers.UserController{}, "post:LoginUser")
	beego.Router("/users", &controllers.UserController{}, "get:IndexAll")

	// Versioned API v1
	apiV1 := beego.NewNamespace("/v1",
		beego.NSRouter("/register", &controllers.UserController{}, "post:RegisterUser"),
		beego.NSRouter("/login", &controllers.UserController{}, "post:LoginUser"),
		beego.NSRouter("/users", &controllers.UserController{}, "get:IndexAll"),
	)

	//produit
	beego.Router("/products", &controllers.ProductController{}, "post:CreateProduct")
	beego.Router("/products", &controllers.ProductController{}, "get:GetAllProducts")
	beego.Router("/products/:id", &controllers.ProductController{}, "get:GetProduct")
	beego.Router("/products/:id", &controllers.ProductController{}, "put:UpdateProduct")
	beego.Router("/products/:id", &controllers.ProductController{}, "delete:DeleteProduct")

	// ORDERS
	beego.Router("/orders", &controllers.OrderController{}, "post:CreateOrder")
	beego.Router("/orders", &controllers.OrderController{}, "get:GetAllOrders")
	beego.Router("/orders/:id", &controllers.OrderController{}, "get:GetOrder")
	beego.Router("/orders/:id", &controllers.OrderController{}, "delete:DeleteOrder")

	beego.AddNamespace(apiV1)
}
