package productrepository

import (
	"go-product-restapi/internal/repositories"
	"go-product-restapi/internal/entity"
)

type PostgresProductMapper struct{}

var _ repositories.ProductMapper[ProductRow] = (*PostgresProductMapper)(nil)

func NewPostgresProductMapper() *PostgresProductMapper {
	return &PostgresProductMapper{}
}

func (p *PostgresProductMapper) ToDTO(product entity.Product) ProductRow {
	return ProductRow{
		ID:          product.ID(),
		Name:        product.Name(),
		Description: product.Description(),
		Price:       product.Price(),
		SalePrice:   product.SalePrice(),
	}
}

func (p *PostgresProductMapper) ToEntity(productDTO ProductRow) *entity.Product {
	var desc string
	var salePrice float64

	if productDTO.Description != nil {
		desc = *productDTO.Description
	}

	if productDTO.SalePrice != nil {
		salePrice = *productDTO.SalePrice
	}

	return entity.NewProduct(
		productDTO.Name,
		desc,
		productDTO.Price,
		salePrice,
	)
}

