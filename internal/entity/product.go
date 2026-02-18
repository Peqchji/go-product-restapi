package entity

import (
	"github.com/google/uuid"
)

type Product struct {
	id          uuid.UUID `json:"id"`
	name        string    `json:"name"`
	description *string   `json:"description"`
	price       float64   `json:"price"`
	salePrice   *float64  `json:"sale_price"`
}

func NewProduct(name string, description string, price float64, salePrice float64) *Product {
	return &Product{
		id:          uuid.New(),
		name:        name,
		description: &description,
		price:       price,
		salePrice:   &salePrice,
	}
}

func (p *Product) ID() uuid.UUID {
	return p.id
}

func (p *Product) Name() string {
	return p.name
}

func (p *Product) Description() *string {
	return p.description
}

func (p *Product) Price() float64 {
	return p.price
}

func (p *Product) SalePrice() *float64 {
	return p.salePrice
}

func (p *Product) SetName(name string) {
	p.name = name
}

func (p *Product) SetDescription(description string) {
	p.description = &description
}

func (p *Product) SetPrice(price float64) {
	p.price = price
}

func (p *Product) SetSalePrice(salePrice float64) {
	p.salePrice = &salePrice
}

func (p *Product) RemoveDescription() {
	p.description = nil
}

func (p *Product) RemoveSalePrice() {
	p.salePrice = nil
}
