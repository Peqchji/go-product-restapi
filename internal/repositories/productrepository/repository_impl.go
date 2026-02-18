package productrepository

import (
	"context"
	"time"

	"go-product-restapi/internal/entity"
	"go-product-restapi/internal/repositories"
	"go-product-restapi/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
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
	db     *pgxpool.Pool
	logger logger.ILogger
	mapper repositories.ProductMapper[ProductRow]
}

func NewPostgresProductRepository(
	db *pgxpool.Pool,
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
	const sql = "SELECT * FROM products WHERE id = $1"
	row := r.db.QueryRow(ctx, sql, id)

	var product ProductRow
	if err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.SalePrice,
		&product.CreatedAt,
		&product.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return r.mapper.ToEntity(product), nil
}

func (r *PostgresProductRepository) Save(ctx context.Context, product entity.Product) error {
	const sql = `
		INSERT INTO products (id, name, description, price, sale_price, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		ON CONFLICT (id) 
		DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			price = EXCLUDED.price,
			sale_price = EXCLUDED.sale_price,
			updated_at = NOW()
	`

	productDTO := r.mapper.ToDTO(product)
	_, err := r.db.Exec(
		ctx,
		sql,
		productDTO.ID,
		productDTO.Name,
		productDTO.Description,
		productDTO.Price,
		productDTO.SalePrice,
	)

	if err != nil {
		return err
	}

	return nil
}
