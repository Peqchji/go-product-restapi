package productrepository

import (
	"context"
	"time"

	"go-product-restapi/internal/entity"
	"go-product-restapi/internal/repositories"
	"go-product-restapi/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ProductRow struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description *string   `db:"description"`
	Price       float64   `db:"price"`
	SalePrice   *float64  `db:"sale_price"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type PostgresProductRepository struct {
	db *pgx.Conn
	logger *logger.Logger
	mapper *repositories.ProductMapper[ProductRow]
}

func NewPostgresProductRepository(db *pgx.Conn) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}

func (r *PostgresProductRepository) GetById(ctx context.Context, id uuid.UUID) (entity.Product, error) {
	var product entity.Product
	return product, nil
}

func (r *PostgresProductRepository) Save(ctx context.Context, product entity.Product) error {
	return nil
}
