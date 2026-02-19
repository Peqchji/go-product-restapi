package productusecase

import (
	"context"
	"errors"
	"testing"

	"go-product-restapi/internal/entity"
	"go-product-restapi/pkg/logger"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductRepository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) GetById(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *MockProductRepository) Save(ctx context.Context, product entity.Product) (*entity.Product, error) {
	args := m.Called(ctx, product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

func TestProductUseCase_GetById(t *testing.T) {
	id := uuid.New()
	product := entity.NewProduct("Test", nil, 10.0, nil)
	product.SetId(&id)

	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name        string
		args        args
		mockSetup   func(m *MockProductRepository)
		expected    *entity.Product
		wantErr     bool
		expectedErr error
	}{

		{
			name: "Should return product when repo returns product",
			args: args{
				ctx: context.Background(),
				id:  id,
			},

			mockSetup: func(m *MockProductRepository) {
				m.On("GetById", mock.Anything, id).Return(product, nil).Once()
			},
			expected: product,
			wantErr:  false,
		},

		{
			name: "Should return error when repo returns not found",
			args: args{
				ctx: context.Background(),
				id:  id,
			},

			mockSetup: func(m *MockProductRepository) {
				m.On("GetById", mock.Anything, id).Return(nil, errors.New("not found")).Once()
			},
			expected:    nil,
			wantErr:     true,
			expectedErr: ErrProductNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProductRepository)
			tt.mockSetup(mockRepo)

			l := logger.NewGlobalLoggerFactory().GetGlobalLogger()
			p := NewProductUseCase(l, mockRepo)

			got, err := p.GetById(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductUseCase.GetById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.expected, got)
			if tt.wantErr && tt.expectedErr != nil {
				assert.Equal(t, tt.expectedErr, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProductUseCase_CreateProduct(t *testing.T) {
	req := &CreateProductRequest{
		Name:  "New Product",
		Price: 20.0,
	}

	type args struct {
		ctx context.Context
		req *CreateProductRequest
	}
	tests := []struct {
		name        string
		args        args
		mockSetup   func(m *MockProductRepository)
		expected    *entity.Product
		wantErr     bool
		expectedErr error
	}{

		{
			name: "Should return created product when repo saves successfully",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			mockSetup: func(m *MockProductRepository) {
				m.On("Save", mock.Anything, mock.MatchedBy(func(p entity.Product) bool {

					return p.Name == req.Name && p.Price == req.Price
				})).Return(&entity.Product{Name: req.Name, Price: req.Price}, nil).Once()
			},
			expected: &entity.Product{Name: req.Name, Price: req.Price},
			wantErr:  false,
		},

		{
			name: "Should return error when repo save fails",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			mockSetup: func(m *MockProductRepository) {
				m.On("Save", mock.Anything, mock.MatchedBy(func(p entity.Product) bool {

					return p.Name == req.Name
				})).Return(nil, errors.New("db error")).Once()
			},
			expected:    nil,
			wantErr:     true,
			expectedErr: ErrCreateProductFailed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProductRepository)
			tt.mockSetup(mockRepo)

			l := logger.NewGlobalLoggerFactory().GetGlobalLogger()
			p := NewProductUseCase(l, mockRepo)

			got, err := p.CreateProduct(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductUseCase.CreateProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.expected, got)
			if tt.wantErr && tt.expectedErr != nil {
				assert.Equal(t, tt.expectedErr, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProductUseCase_UpdateProduct(t *testing.T) {
	id := uuid.New()
	req := &UpdateProductRequest{
		ID: id,
	}
	name := "Updated Name"
	req.Name = &name

	existingProduct := entity.NewProduct("Old Name", nil, 10.0, nil)
	existingProduct.SetId(&id)

	type args struct {
		ctx context.Context
		req *UpdateProductRequest
	}
	tests := []struct {
		name        string
		args        args
		mockSetup   func(m *MockProductRepository)
		wantErr     bool
		expectedErr error
	}{

		{
			name: "Should return nil when product is updated successfully",
			args: args{
				ctx: context.Background(),
				req: req,
			},

			mockSetup: func(m *MockProductRepository) {
				m.On("GetById", mock.Anything, id).Return(existingProduct, nil).Once()
				m.On("Save", mock.Anything, mock.MatchedBy(func(p entity.Product) bool {
					return p.Name == "Updated Name" && *p.ID == id
				})).Return(existingProduct, nil).Once()
			},
			wantErr: false,
		},

		{
			name: "Should return error when product to update is not found",
			args: args{
				ctx: context.Background(),
				req: req,
			},

			mockSetup: func(m *MockProductRepository) {
				m.On("GetById", mock.Anything, id).Return(nil, errors.New("not found")).Once()
			},
			wantErr:     true,
			expectedErr: ErrProductNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProductRepository)
			tt.mockSetup(mockRepo)

			l := logger.NewGlobalLoggerFactory().GetGlobalLogger()
			p := NewProductUseCase(l, mockRepo)

			if err := p.UpdateProduct(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("ProductUseCase.UpdateProduct() error = %v, wantErr %v", err, tt.wantErr)
			} else if tt.wantErr && tt.expectedErr != nil {
				assert.Equal(t, tt.expectedErr, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
