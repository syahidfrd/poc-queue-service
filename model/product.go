package model

import "time"

type (
	Product struct {
		ID       uint64 `gorm:"primary_key;auto_increment"`
		Name     string
		Quantity uint32
		Price    uint64

		CreatedAt time.Time
		UpdatedAt time.Time
	}

	ProductStore interface {
		Exist(product *Product) bool
		Save(product *Product) error
		FindOneBy(query map[string]interface{}) (*Product, error)
	}
)

func NewProduct(name string, quantity uint32, price uint64) (product *Product) {
	product = &Product{}
	product.SetName(name)
	product.SetQuantity(quantity)
	product.SetPrice(price)
	return
}

func (p *Product) GetID() uint64 {
	return p.ID
}

func (p *Product) GetName() string {
	return p.Name
}

func (p *Product) GetQuantity() uint32 {
	return p.Quantity
}

func (p *Product) GetPrice() uint64 {
	return p.Price
}

func (p *Product) GetCreatedAt() time.Time {
	return p.CreatedAt
}

func (p *Product) GetUpdatedAt() time.Time {
	return p.UpdatedAt
}

func (p *Product) SetName(name string) {
	p.Name = name
}

func (p *Product) SetQuantity(quantity uint32) {
	p.Quantity = quantity
}

func (p *Product) SetPrice(price uint64) {
	p.Price = price
}
