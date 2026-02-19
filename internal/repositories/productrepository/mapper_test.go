package productrepository_test

import (
	"testing"
	"time"

	"go-product-restapi/internal/entity"
	"go-product-restapi/internal/repositories/productrepository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPostgresProductMapper_ToDTO(t *testing.T) {
	id := uuid.New()
	name := "Test Product"
	description := "Test Description"
	price := 10.0
	salePrice := 8.0

	product := entity.NewProduct(name, &description, price, &salePrice)
	product.SetId(&id)

	mapper := productrepository.NewPostgresProductMapper()

	type args struct {
		product entity.Product
	}

	tests := []struct {
		name     string
		args     args
		expected productrepository.ProductRow
	}{
		{
			name: "Should map entity to DTO correctly",
			args: args{
				product: *product,
			},
			expected: productrepository.ProductRow{
				ID:          &id,
				Name:        name,
				Description: &description,
				Price:       price,
				SalePrice:   &salePrice,
			},
		},
		{
			name: "Should map entity with nil fields to DTO correctly",
			args: args{
				product: entity.Product{
					ID:          &id,
					Name:        name,
					Price:       price,
					Description: nil,
					SalePrice:   nil,
				},
			},
			expected: productrepository.ProductRow{
				ID:          &id,
				Name:        name,
				Description: nil,
				Price:       price,
				SalePrice:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapper.ToDTO(tt.args.product)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestPostgresProductMapper_ToEntity(t *testing.T) {
	id := uuid.New()
	name := "Test Product"
	description := "Test Description"
	price := 10.0
	salePrice := 8.0
	createdAt := time.Now()
	updatedAt := time.Now()

	dto := productrepository.ProductRow{
		ID:          &id,
		Name:        name,
		Description: &description,
		Price:       price,
		SalePrice:   &salePrice,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	mapper := productrepository.NewPostgresProductMapper()

	type args struct {
		dto productrepository.ProductRow
	}

	tests := []struct {
		name     string
		args     args
		expected *entity.Product
	}{
		{
			name: "Should map DTO to entity correctly",
			args: args{
				dto: dto,
			},
			expected: func() *entity.Product {
				p := entity.NewProduct(name, &description, price, &salePrice)
				p.SetId(&id)
				return p
			}(),
		},
		{
			name: "Should map DTO with nil fields to entity correctly",
			args: args{
				dto: productrepository.ProductRow{
					ID:          &id,
					Name:        name,
					Description: nil,
					Price:       price,
					SalePrice:   nil,
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
				},
			},
			expected: func() *entity.Product {
				p := entity.NewProduct(name, nil, price, nil)
				p.SetId(&id)
				return p
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapper.ToEntity(tt.args.dto)

			assert.NotNil(t, got)
			assert.Equal(t, tt.expected, got)
		})
	}
}
