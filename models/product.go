package models

import (
	"errors"

	"github.com/astaxie/beego/orm"
)

type Product struct {
	Id          int64   `json:"id" orm:"pk;auto"`
	Name        string  `json:"name" orm:"size(150)"`
	Description string  `json:"description" orm:"type(text)"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

func init() {
	orm.RegisterModel(new(Product))
}

func (p *Product) TableName() string {
	return "products"
}

// CREATE
func CreateProduct(p *Product) (int64, error) {
	if p.Name == "" || p.Price <= 0 {
		return 0, errors.New("invalid product data")
	}
	o := orm.NewOrm()
	return o.Insert(p)
}

// READ ALL
func GetAllProducts() ([]Product, error) {
	o := orm.NewOrm()
	var products []Product
	_, err := o.QueryTable(new(Product)).All(&products)
	return products, err
}

// READ ONE
func GetProductById(id int64) (*Product, error) {
	o := orm.NewOrm()
	p := Product{Id: id}
	if err := o.Read(&p); err != nil {
		return nil, errors.New("product not found")
	}
	return &p, nil
}

// UPDATE
func UpdateProduct(p *Product) error {
	o := orm.NewOrm()
	_, err := o.Update(p)
	return err
}

// DELETE
func DeleteProduct(id int64) error {
	o := orm.NewOrm()
	_, err := o.Delete(&Product{Id: id})
	return err
}
