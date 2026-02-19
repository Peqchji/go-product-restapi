package component

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-product-restapi/internal/entity"
	"go-product-restapi/internal/handler/producthandler"
	"go-product-restapi/internal/usecases/productusecase"
	"go-product-restapi/pkg/logger"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
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

func TestProductComponent_GetProductFlow(t *testing.T) {
	id := uuid.New()
	product := entity.NewProduct("Integration Product", nil, 50.0, nil)
	product.SetId(&id)

	tests := []struct {
		name           string
		idParam        string
		mockSetup      func(m *MockProductRepository)
		expectedStatus int
		verifyBody     func(t *testing.T, body []byte)
	}{

		{
			name:    "Should return 200 and product when flow is successful",
			idParam: id.String(),
			mockSetup: func(m *MockProductRepository) {
				m.On("GetById", mock.Anything, id).Return(product, nil).Once()
			},
			expectedStatus: http.StatusOK,
			verifyBody: func(t *testing.T, body []byte) {
				var resp map[string]interface{}
				json.Unmarshal(body, &resp)
				assert.Equal(t, true, resp["successful"])
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, id.String(), data["id"])
				assert.Equal(t, "Integration Product", data["name"])
			},
		},

		{
			name:    "Should return 500 when product is not found in repository",
			idParam: id.String(),
			mockSetup: func(m *MockProductRepository) {
				m.On("GetById", mock.Anything, id).Return(nil, errors.New("not found")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			verifyBody:     func(t *testing.T, body []byte) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProductRepository)
			tt.mockSetup(mockRepo)

			l := logger.NewGlobalLoggerFactory().GetGlobalLogger()
			useCase := productusecase.NewProductUseCase(l, mockRepo)
			handler := producthandler.NewProductHandler(useCase)

			e := echo.New()
			e.GET("/products/:id", handler.GetProduct())

			req := httptest.NewRequest(http.MethodGet, "/products/"+tt.idParam, nil)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.verifyBody != nil {
				tt.verifyBody(t, rec.Body.Bytes())
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
