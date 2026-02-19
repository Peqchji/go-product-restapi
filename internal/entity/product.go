package entity

import (
	"github.com/google/uuid"
)

type Product struct {
	ID          *uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Price       float64   `json:"price"`
	SalePrice   *float64  `json:"sale_price"`
}

func NewProduct(name string, description *string, price float64, salePrice *float64) *Product {
	return &Product{
		Name:        name,
		Description: description,
		Price:       price,
		SalePrice:   salePrice,
	}
}

func (p *Product) SetId(id *uuid.UUID) {
	p.ID = id
}

func (p *Product) SetName(name string) {
	p.Name = name
}

func (p *Product) SetDescription(description string) {
	p.Description = &description
}

func (p *Product) SetPrice(price float64) {
	p.Price = price
}

func (p *Product) SetSalePrice(salePrice float64) {
	p.SalePrice = &salePrice
}

func (p *Product) RemoveDescription() {
	p.Description = nil
}

func (p *Product) RemoveSalePrice() {
	p.SalePrice = nil
}
