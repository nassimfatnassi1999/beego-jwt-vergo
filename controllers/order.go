package controllers

import (
	"beego-jwt-vergo/models"
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
)

type OrderController struct {
	beego.Controller
}

// POST /orders
func (c *OrderController) CreateOrder() {
	var body struct {
		ProductId int64 `json:"product_id"`
		Quantity  int   `json:"quantity"`
	}

	json.Unmarshal(c.Ctx.Input.CopyBody(1<<20), &body)

	order, err := models.CreateOrder(body.ProductId, body.Quantity)
	if err != nil {
		c.CustomAbort(400, err.Error())
	}

	c.Data["json"] = order
	c.ServeJSON()
}

// GET /orders
func (c *OrderController) GetAllOrders() {
	orders, err := models.GetAllOrders()
	if err != nil {
		c.CustomAbort(500, err.Error())
	}

	c.Data["json"] = orders
	c.ServeJSON()
}

// GET /orders/:id
func (c *OrderController) GetOrder() {
	id, _ := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)

	order, err := models.GetOrderById(id)
	if err != nil {
		c.CustomAbort(404, err.Error())
	}

	c.Data["json"] = order
	c.ServeJSON()
}

// DELETE /orders/:id
func (c *OrderController) DeleteOrder() {
	id, _ := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)

	err := models.DeleteOrder(id)
	if err != nil {
		c.CustomAbort(400, err.Error())
	}

	c.Data["json"] = map[string]string{"message": "order deleted"}
	c.ServeJSON()
}
