package productusecase

import (
	"context"

	"go-product-restapi/internal/entity"
	"go-product-restapi/internal/repositories"
	"go-product-restapi/pkg/logger"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	SalePrice   *float64 `json:"sale_price,omitempty"`
	Price       float64  `json:"price"`
}

type UpdateProductRequest struct {
	ID          uuid.UUID `json:"id"`
	Name        *string   `json:"name,omitempty"`
	Description *string   `json:"description,omitempty"`
	SalePrice   *float64 `json:"sale_price,omitempty"`
	Price       *float64 `json:"price,omitempty"`

	IsSetDescription bool `json:"-"`
	IsSetSalePrice   bool `json:"-"`
}

type IProductUseCase interface {
	GetById(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	CreateProduct(ctx context.Context, req *CreateProductRequest) (*entity.Product, error)
	UpdateProduct(ctx context.Context, req *UpdateProductRequest) error
}

type ProductUseCase struct {
	logger      logger.ILogger
	productRepo repositories.ProductRepository
}

func NewProductUseCase(logger logger.ILogger, productRepo repositories.ProductRepository) *ProductUseCase {
	usecaseLogger := logger.Named("ProductUseCase")

	return &ProductUseCase{logger: usecaseLogger, productRepo: productRepo}
}

func (p *ProductUseCase) GetById(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	p.logger.Info("Getting product by id", id)
	product, err := p.productRepo.GetById(ctx, id)
	if err != nil {
		p.logger.Error("Failed to get product by id", id, err)

		return nil, ErrProductNotFound
	}

	return product, nil
}

func (p *ProductUseCase) CreateProduct(ctx context.Context, req *CreateProductRequest) (*entity.Product, error) {
	p.logger.Info("Creating product", req)
	newProduct := entity.NewProduct(
		req.Name,
		req.Description,
		req.Price,
		req.SalePrice,
	)

	createdProduct, err := p.productRepo.Save(ctx, *newProduct)
	if err != nil {
		p.logger.Error("Failed to create product", err)

		return nil, ErrCreateProductFailed
	}

	return createdProduct, nil
}

func (p *ProductUseCase) UpdateProduct(ctx context.Context, req *UpdateProductRequest) error {
	p.logger.Info("Updating product", req.ID)
	p.logger.Debug("UpdateProductRequest", req)

	product, err := p.productRepo.GetById(ctx, req.ID)
	if err != nil {
		p.logger.Error("Product not found", "id", req.ID, "error", err)

		return ErrProductNotFound
	}

	if req.Name != nil {
		product.SetName(*req.Name)
	}

	if req.IsSetDescription {
		if req.Description != nil {
			product.SetDescription(*req.Description)
		} else {
			product.RemoveDescription()
		}
	}

	if req.IsSetSalePrice {
		if req.SalePrice != nil {
			product.SetSalePrice(*req.SalePrice)
		} else {
			product.RemoveSalePrice()
		}
	}

	if req.Price != nil {
		product.SetPrice(*req.Price)
	}

	p.logger.Debug("Product updated", product)

	_, err = p.productRepo.Save(ctx, *product)
	if err != nil {
		p.logger.Error("Failed to save product", "id", req.ID, "error", err)
		return ErrUpdateProductFailed
	}

	p.logger.Info("Product updated successfully", "id", req.ID)
	return nil
}
