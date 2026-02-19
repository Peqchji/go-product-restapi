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

type IDBPool interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type ProductRow struct {
	ID          *uuid.UUID `db:"id,omitempty"`
	Name        string     `db:"name"`
	Description *string    `db:"description,omitempty"`
	Price       float64    `db:"price"`
	SalePrice   *float64   `db:"sale_price,omitempty"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

type PostgresProductRepository struct {
	db     IDBPool
	logger logger.ILogger
	mapper repositories.ProductMapper[ProductRow]
}

func NewPostgresProductRepository(
	db IDBPool,
	logger logger.ILogger,
	productMapper repositories.ProductMapper[ProductRow],
) *PostgresProductRepository {
	return &PostgresProductRepository{
		db:     db,
		logger: logger.Named("PostgresProductRepository"),
		mapper: productMapper,
	}
}

func (r *PostgresProductRepository) GetById(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	const sql = "SELECT id, name, description, price, sale_price, created_at, updated_at FROM products WHERE id = $1"
	row := r.db.QueryRow(ctx, sql, id)

	var productRow ProductRow
	if err := row.Scan(
		&productRow.ID,
		&productRow.Name,
		&productRow.Description,
		&productRow.Price,
		&productRow.SalePrice,
		&productRow.CreatedAt,
		&productRow.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return r.mapper.ToEntity(productRow), nil
}

func (r *PostgresProductRepository) Save(ctx context.Context, product entity.Product) (*entity.Product, error) {
	const sql = `
		INSERT INTO products (id, name, description, price, sale_price, created_at, updated_at) 
		VALUES (COALESCE($1, gen_random_uuid()), $2, $3, $4, $5, NOW(), NOW())
		ON CONFLICT (id) 
		DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			price = EXCLUDED.price,
			sale_price = EXCLUDED.sale_price,
			updated_at = NOW()
		RETURNING id, name, description, price, sale_price, created_at, updated_at
	`

	productDTO := r.mapper.ToDTO(product)

	var productRow ProductRow
	err := r.db.QueryRow(
		ctx,
		sql,
		productDTO.ID,
		productDTO.Name,
		productDTO.Description,
		productDTO.Price,
		productDTO.SalePrice,
	).Scan(
		&productRow.ID,
		&productRow.Name,
		&productRow.Description,
		&productRow.Price,
		&productRow.SalePrice,
		&productRow.CreatedAt,
		&productRow.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return r.mapper.ToEntity(productRow), nil
}
