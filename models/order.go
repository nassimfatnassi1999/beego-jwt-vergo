package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
)

type Order struct {
	Id         int64     `json:"id" orm:"pk;auto"`
	Product    *Product  `json:"product" orm:"rel(fk)"`
	Quantity   int       `json:"quantity"`
	TotalPrice float64   `json:"total_price"`
	CreatedAt  time.Time `json:"created_at" orm:"auto_now_add;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Order))
}

func (o *Order) TableName() string {
	return "orders"
}

// CREATE ORDER
func CreateOrder(productId int64, qty int) (*Order, error) {
	if qty <= 0 {
		return nil, errors.New("quantity must be > 0")
	}

	// Retrieve product
	product, err := GetProductById(productId)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Check stock
	if product.Stock < qty {
		return nil, errors.New("not enough stock")
	}

	// Calculate price
	total := product.Price * float64(qty)

	order := &Order{
		Product:    product,
		Quantity:   qty,
		TotalPrice: total,
	}

	o := orm.NewOrm()
	_, err = o.Insert(order)
	if err != nil {
		return nil, err
	}

	// Decrease stock
	product.Stock -= qty
	_ = UpdateProduct(product)

	return order, nil
}

// GET ALL ORDERS
func GetAllOrders() ([]*Order, error) {
	o := orm.NewOrm()
	var orders []*Order
	_, err := o.QueryTable(new(Order)).RelatedSel().All(&orders)
	return orders, err
}

// GET ORDER BY ID
func GetOrderById(id int64) (*Order, error) {
	o := orm.NewOrm()
	order := Order{Id: id}
	if err := o.Read(&order); err != nil {
		return nil, errors.New("order not found")
	}
	o.LoadRelated(&order, "Product")
	return &order, nil
}

// DELETE ORDER
func DeleteOrder(id int64) error {
	o := orm.NewOrm()
	_, err := o.Delete(&Order{Id: id})
	return err
}
