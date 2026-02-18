package productusecase

import (
	"context"
	"errors"

	"go-product-restapi/internal/entity"
	"go-product-restapi/internal/repositories"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	SalePrice   *float64 `json:"sale_price,omitempty"`
	Price       float64  `json:"price"`
}

type UpdateProductRequest struct {
	ID                 uuid.UUID `json:"id"`
	Name               *string   `json:"name,omitempty"`
	Description        *string   `json:"description,omitempty"`
	SetDescriptionNull bool
	SalePrice          *float64 `json:"sale_price,omitempty"`
	SetSalePriceNull   bool
	Price              *float64 `json:"price,omitempty"`
}

type IProductUseCase interface {
	CreateProduct(ctx context.Context, req *CreateProductRequest) (*entity.Product, error)
	UpdateProduct(ctx context.Context, req *UpdateProductRequest) error
}

type ProductUseCase struct {
	productRepo repositories.ProductRepository
}

func NewProductUseCase(productRepo repositories.ProductRepository) *ProductUseCase {
	return &ProductUseCase{productRepo: productRepo}
}

func (p *ProductUseCase) CreateProduct(ctx context.Context, req *CreateProductRequest) (*entity.Product, error) {
	newProduct := entity.NewProduct(
		req.Name,
		*req.Description,
		req.Price,
		*req.SalePrice,
	)

	err := p.productRepo.Save(ctx, *newProduct)
	if err != nil {
		return nil, err
	}

	return newProduct, nil
}

func (p *ProductUseCase) UpdateProduct(ctx context.Context, req *UpdateProductRequest) error {
	product, err := p.productRepo.GetById(ctx, req.ID)
	if err != nil {
		return errors.New("product not found")
	}

	if req.Name != nil {
		product.SetName(*req.Name)
	}

	if req.SetDescriptionNull {
		product.RemoveDescription()
	} else if req.Description != nil {
		product.SetDescription(*req.Description)
	}

	if req.SetSalePriceNull {
		product.RemoveSalePrice()
	} else if req.SalePrice != nil {
		product.SetSalePrice(*req.SalePrice)
	}

	if req.Price != nil {
		product.SetPrice(*req.Price)
	}

	err = p.productRepo.Save(ctx, *product)
	if err != nil {
		return err
	}

	return nil
}
