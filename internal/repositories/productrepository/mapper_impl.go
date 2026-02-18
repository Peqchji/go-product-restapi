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
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		SalePrice:   product.SalePrice,
	}
}

func (p *PostgresProductMapper) ToEntity(productDTO ProductRow) *entity.Product {
	product := entity.NewProduct(
		productDTO.Name,
		productDTO.Description,
		productDTO.Price,
		productDTO.SalePrice,
	)

	product.SetId(productDTO.ID)

	return product
}

