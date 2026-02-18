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
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

func (p *PostgresProductMapper) ToEntity(productDTO ProductRow) entity.Product {
	return entity.Product{
		ID:          productDTO.ID,
		Name:        productDTO.Name,
		Description: productDTO.Description,
		Price:       productDTO.Price,
		SalePrice:   productDTO.SalePrice,
		CreatedAt:   productDTO.CreatedAt,
		UpdatedAt:   productDTO.UpdatedAt,
	}
}

