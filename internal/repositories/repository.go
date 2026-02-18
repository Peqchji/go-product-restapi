package repositories

import (
	"context"

	"go-product-restapi/internal/entity"

	"github.com/google/uuid"
)


type ProductRepository interface {
	GetById(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	Save(ctx context.Context, product entity.Product) (*entity.Product, error)
}