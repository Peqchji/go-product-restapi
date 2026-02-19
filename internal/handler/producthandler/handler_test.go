package producthandler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-product-restapi/internal/entity"
	"go-product-restapi/internal/usecases/productusecase"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductUseCase
type MockProductUseCase struct {
	mock.Mock
}

func (m *MockProductUseCase) GetById(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *MockProductUseCase) CreateProduct(ctx context.Context, req *productusecase.CreateProductRequest) (*entity.Product, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *MockProductUseCase) UpdateProduct(ctx context.Context, req *productusecase.UpdateProductRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func TestProductHandler_GetProduct(t *testing.T) {
	id := uuid.New()
	product := entity.NewProduct("Test Product", nil, 10.0, nil)
	product.SetId(&id)

	tests := []struct {
		name           string
		idParam        string
		mockSetup      func(m *MockProductUseCase)
		expectedStatus int
		verifyBody     func(t *testing.T, body []byte)
	}{

		{
			name:    "Should return 200 and product data when product exists",
			idParam: id.String(),
			mockSetup: func(m *MockProductUseCase) {

				m.On("GetById", mock.Anything, id).Return(product, nil).Once()
			},
			expectedStatus: http.StatusOK,
			verifyBody: func(t *testing.T, body []byte) {
				var resp map[string]interface{}
				json.Unmarshal(body, &resp)
				assert.Equal(t, true, resp["successful"])
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, id.String(), data["id"])
			},
		},

		{
			name:    "Should return 500 when product does not exist",
			idParam: id.String(),
			mockSetup: func(m *MockProductUseCase) {
				m.On("GetById", mock.Anything, id).Return(nil, productusecase.ErrProductNotFound).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			verifyBody:     func(t *testing.T, body []byte) {},
		},

		{
			name:    "Should return 400 when product ID is invalid",
			idParam: "invalid-uuid",
			mockSetup: func(m *MockProductUseCase) {
				// mock should not be called
			},
			expectedStatus: http.StatusBadRequest,
			verifyBody:     func(t *testing.T, body []byte) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockProductUseCase)
			if tt.mockSetup != nil {
				tt.mockSetup(mockUseCase)
			}
			h := NewProductHandler(mockUseCase)
			e := echo.New()
			e.GET("/products/:id", h.GetProduct())

			req := httptest.NewRequest(http.MethodGet, "/products/"+tt.idParam, nil)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.verifyBody != nil {
				tt.verifyBody(t, rec.Body.Bytes())
			}
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestProductHandler_CreateProduct(t *testing.T) {
	createdProduct := entity.NewProduct("New Product", nil, 20.0, nil)
	id := uuid.New()
	createdProduct.SetId(&id)

	tests := []struct {
		name           string
		reqBody        string
		mockSetup      func(m *MockProductUseCase)
		expectedStatus int
		verifyBody     func(t *testing.T, body []byte)
	}{

		{
			name:    "Should return 200 and created product when payload is valid",
			reqBody: `{"name":"New Product","price":20.0}`,
			mockSetup: func(m *MockProductUseCase) {
				m.On("CreateProduct", mock.Anything, mock.MatchedBy(func(r *productusecase.CreateProductRequest) bool {
					return r.Name == "New Product" && r.Price == 20.0
				})).Return(createdProduct, nil).Once()
			},
			expectedStatus: http.StatusOK,
			verifyBody: func(t *testing.T, body []byte) {
				var resp map[string]interface{}
				json.Unmarshal(body, &resp)
				assert.Equal(t, true, resp["successful"])
			},
		},

		{
			name:    "Should return 400 when payload is invalid",
			reqBody: `{"name":"","price":-10.0}`,
			mockSetup: func(m *MockProductUseCase) {
				// mock should not be called
			},
			expectedStatus: http.StatusBadRequest,
			verifyBody:     func(t *testing.T, body []byte) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockProductUseCase)
			if tt.mockSetup != nil {
				tt.mockSetup(mockUseCase)
			}
			h := NewProductHandler(mockUseCase)
			e := echo.New()
			e.POST("/products", h.CreateProduct())

			req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(tt.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.verifyBody != nil {
				tt.verifyBody(t, rec.Body.Bytes())
			}
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestProductHandler_PatchProduct(t *testing.T) {
	id := uuid.New()

	tests := []struct {
		name           string
		idParam        string
		reqBody        string
		mockSetup      func(m *MockProductUseCase)
		expectedStatus int
	}{

		{
			name:    "Should return 200 when update is successful",
			idParam: id.String(),
			reqBody: `{"name":"Updated Name"}`,
			mockSetup: func(m *MockProductUseCase) {
				m.On("UpdateProduct", mock.Anything, mock.MatchedBy(func(r *productusecase.UpdateProductRequest) bool {
					return *r.Name == "Updated Name" && r.ID == id
				})).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
		},

		{
			name:    "Should return 400 when product ID is invalid",
			idParam: "invalid-uuid",
			reqBody: `{"name":"Updated Name"}`,
			mockSetup: func(m *MockProductUseCase) {
				// mock should not be called
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockProductUseCase)
			if tt.mockSetup != nil {
				tt.mockSetup(mockUseCase)
			}
			h := NewProductHandler(mockUseCase)
			e := echo.New()
			e.PATCH("/products/:id", h.PatchProduct())

			req := httptest.NewRequest(http.MethodPatch, "/products/"+tt.idParam, strings.NewReader(tt.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			mockUseCase.AssertExpectations(t)
		})
	}
}
