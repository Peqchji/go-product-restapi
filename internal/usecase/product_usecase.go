package usecase

import (
	"context"

	"go-product-restapi/internal/entity"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	SalePrice   *float64 `json:"sale_price,omitempty"`
	Price       float64  `json:"price"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	SalePrice   *float64 `json:"sale_price,omitempty"`
	Price       *float64 `json:"price,omitempty"`
}

type ProductUseCase interface {
	CreateProduct(ctx context.Context, req *CreateProductRequest) (*entity.Product, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req *UpdateProductRequest) error
}
