package controllers

import (
	"beego-jwt-vergo/models"
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
)

type ProductController struct {
	beego.Controller
}

// POST /products
func (c *ProductController) CreateProduct() {
	var p models.Product
	json.Unmarshal(c.Ctx.Input.CopyBody(1<<20), &p)

	id, err := models.CreateProduct(&p)
	if err != nil {
		c.CustomAbort(400, err.Error())
	}

	p.Id = id
	c.Data["json"] = p
	c.ServeJSON()
}

// GET /products
func (c *ProductController) GetAllProducts() {
	products, err := models.GetAllProducts()
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.Data["json"] = products
	c.ServeJSON()
}

// GET /products/:id
func (c *ProductController) GetProduct() {
	id, _ := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	p, err := models.GetProductById(id)
	if err != nil {
		c.CustomAbort(404, err.Error())
	}
	c.Data["json"] = p
	c.ServeJSON()
}

// PUT /products/:id
func (c *ProductController) UpdateProduct() {
	id, _ := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	p, err := models.GetProductById(id)
	if err != nil {
		c.CustomAbort(404, "product not found")
	}

	json.Unmarshal(c.Ctx.Input.CopyBody(1<<20), p)

	err = models.UpdateProduct(p)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}

	c.Data["json"] = p
	c.ServeJSON()
}

// DELETE /products/:id
func (c *ProductController) DeleteProduct() {
	id, _ := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	err := models.DeleteProduct(id)
	if err != nil {
		c.CustomAbort(400, err.Error())
	}
	c.Data["json"] = map[string]string{"message": "product deleted"}
	c.ServeJSON()
}
