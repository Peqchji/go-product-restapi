package productrepository

import (
	"context"
	"testing"
	"time"

	"go-product-restapi/internal/entity"
	"go-product-restapi/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func TestPostgresProductRepository_GetById(t *testing.T) {
	id := uuid.New()
	name := "Test Product"
	description := "Test Description"
	price := 10.0
	salePrice := 8.0
	createdAt := time.Now()
	updatedAt := time.Now()

	product := &entity.Product{
		ID:          &id,
		Name:        name,
		Description: &description,
		Price:       price,
		SalePrice:   &salePrice,
	}

	type args struct {
		ctx context.Context
		id  uuid.UUID
	}

	tests := []struct {
		name        string
		args        args
		mockSetup   func(m pgxmock.PgxPoolIface)
		expected    *entity.Product
		expectedErr bool
	}{

		{
			name: "Should return product when product exists",
			args: args{
				ctx: context.Background(),
				id:  id,
			},
			mockSetup: func(m pgxmock.PgxPoolIface) {
				rows := m.NewRows([]string{"id", "name", "description", "price", "sale_price", "created_at", "updated_at"}).
					AddRow(&id, name, &description, price, &salePrice, createdAt, updatedAt)
				m.ExpectQuery("^SELECT id, name, description, price, sale_price, created_at, updated_at FROM products WHERE id = \\$1$").
					WithArgs(id).
					WillReturnRows(rows)
			},
			expected:    product,
			expectedErr: false,
		},
		{
			name: "Should return error when product does not exist",
			args: args{
				ctx: context.Background(),
				id:  id,
			},

			mockSetup: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery("^SELECT id, name, description, price, sale_price, created_at, updated_at FROM products WHERE id = \\$1$").
					WithArgs(id).
					WillReturnError(pgx.ErrNoRows)
			},
			expected:    nil,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mock.Close()

			tt.mockSetup(mock)

			l := logger.NewGlobalLoggerFactory().GetGlobalLogger()
			mapper := NewPostgresProductMapper()
			r := NewPostgresProductRepository(mock, l, mapper)

			got, err := r.GetById(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.expectedErr {
				t.Errorf("PostgresProductRepository.GetById() error = %v, wantErr %v", err, tt.expectedErr)
				return
			}

			if !tt.expectedErr {
				assert.Equal(t, tt.expected, got)
			} else {
				assert.Nil(t, got)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPostgresProductRepository_Save(t *testing.T) {
	id := uuid.New()
	name := "Test Product"
	description := "Test Description"
	price := 10.0
	salePrice := 8.0

	product := entity.NewProduct(name, &description, price, &salePrice)
	product.SetId(&id)

	type args struct {
		ctx     context.Context
		product entity.Product
	}

	tests := []struct {
		name      string
		args      args
		mockSetup func(m pgxmock.PgxPoolIface)
		expected  *entity.Product
		wantErr   bool
	}{

		{
			name: "Should return saved product with ID when save is successful",
			args: args{
				ctx:     context.Background(),
				product: *product,
			},

			mockSetup: func(m pgxmock.PgxPoolIface) {
				returnedId := id
				returnedCreatedAt := time.Now()
				returnedUpdatedAt := time.Now()

				rows := m.NewRows([]string{"id", "name", "description", "price", "sale_price", "created_at", "updated_at"}).
					AddRow(&returnedId, name, &description, price, &salePrice, returnedCreatedAt, returnedUpdatedAt)

				m.ExpectQuery("^INSERT INTO products").
					WithArgs(product.ID, product.Name, product.Description, product.Price, product.SalePrice).
					WillReturnRows(rows)
			},
			expected: product,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mock.Close()

			tt.mockSetup(mock)

			l := logger.NewGlobalLoggerFactory().GetGlobalLogger()
			mapper := NewPostgresProductMapper()
			r := NewPostgresProductRepository(mock, l, mapper)

			got, err := r.Save(tt.args.ctx, tt.args.product)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostgresProductRepository.Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.Equal(t, tt.expected, got)
			} else {
				assert.Nil(t, got)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
