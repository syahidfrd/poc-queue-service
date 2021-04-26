package model

import "time"

type (
	Order struct {
		ID       uint64 `gorm:"primary_key;auto_increment"`
		Quantity uint32

		ProductID uint64 `gorm:"index"`
		Product   *Product

		CreatedAt time.Time
		UpdatedAt time.Time
	}

	OrderStore interface {
		Exist(product *Order) bool
		Save(product *Order) error
		FindOneBy(query map[string]interface{}) (*Order, error)
		FindAll() ([]*Order, error)
	}
)

func NewOrder(quantity uint32, productID uint64) (order *Order) {
	order = &Order{}
	order.SetQuantity(quantity)
	order.SetProductID(productID)
	return
}

func (o *Order) GetID() uint64 {
	return o.ID
}

func (o *Order) GetQuantity() uint32 {
	return o.Quantity
}

func (o *Order) GetProductID() uint64 {
	return o.ProductID
}

func (o *Order) GetProduct() *Product {
	return o.Product
}

func (o *Order) GetCreatedAt() time.Time {
	return o.CreatedAt
}

func (o *Order) GetUpdatedAt() time.Time {
	return o.UpdatedAt
}

func (o *Order) SetQuantity(quantity uint32) {
	o.Quantity = quantity
}

func (o *Order) SetProductID(productID uint64) {
	o.ProductID = productID
}
