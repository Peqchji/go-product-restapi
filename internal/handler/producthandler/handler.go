package producthandler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"go-product-restapi/internal/handler"
	"go-product-restapi/internal/usecases/productusecase"
	"go-product-restapi/pkg/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type ProductHandler struct {
	productUseCase productusecase.IProductUseCase
}

func NewProductHandler(productUseCase productusecase.IProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
	}
}

// @Summary Get Product
// @Description Get a product by ID
// @Tags products
// @ID get-product
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID (UUID)"
// @Success 200 {object} producthandler.ProductSuccessResponse
// @Failure 400 {object} handler.BaseAPIResponse
// @Failure 500 {object} handler.BaseAPIResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct() echo.HandlerFunc {
	return func(c *echo.Context) error {
		idStr := c.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			return c.JSON(
				http.StatusBadRequest,
				handler.BaseAPIResponse{
					Successful: false,
					ErrorCode:  utils.ToPointer(ErrCodeInvalidRequest),
				},
			)
		}

		ctx := c.Request().Context()
		defer ctx.Done()

		product, err := h.productUseCase.GetById(ctx, id)
		if err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				handler.BaseAPIResponse{
					Successful: false,
					ErrorCode:  utils.ToPointer(ErrProductNotFound),
				},
			)
		}

		return c.JSON(
			http.StatusOK,
			handler.BaseAPIResponse{
				Successful: true,
				Data: ProductResponse{
					ID:          idStr,
					Name:        product.Name,
					Description: product.Description,
					Price:       product.Price,
					SalePrice:   product.SalePrice,
				},
			},
		)
	}
}

// @Summary Patch Product
// @Description Patch a product by ID
// @Tags products
// @ID patch-product
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID (UUID)"
// @Param product body productusecase.UpdateProductRequest true "Product Request"
// @Success 200 {object} handler.BaseAPIResponse
// @Failure 400 {object} handler.BaseAPIResponse
// @Failure 500 {object} handler.BaseAPIResponse
// @Router /products/{id} [patch]
func (h *ProductHandler) PatchProduct() echo.HandlerFunc {
	return func(c *echo.Context) error {
		idStr := c.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			return c.JSON(
				http.StatusBadRequest,
				handler.BaseAPIResponse{
					Successful: false,
					ErrorCode:  utils.ToPointer(ErrCodeInvalidRequest),
				})
		}

		bodyBytes, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return c.JSON(
				http.StatusBadRequest,
				handler.BaseAPIResponse{
					Successful: false,
					ErrorCode:  utils.ToPointer(ErrCodeInvalidRequest),
				},
			)
		}

		c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var req productusecase.UpdateProductRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(
				http.StatusBadRequest,
				handler.BaseAPIResponse{
					Successful: false,
					ErrorCode:  utils.ToPointer(ErrCodeInvalidRequest),
				},
			)
		}

		rawBody := make(map[string]interface{})
		if err := json.Unmarshal(bodyBytes, &rawBody); err != nil {
			return c.JSON(
				http.StatusBadRequest,
				handler.BaseAPIResponse{
					Successful: false,
					ErrorCode:  utils.ToPointer(ErrCodeInvalidRequest),
				},
			)
		}

		_, isDescriptionPresent := rawBody["description"]
		_, isSalePricePresent := rawBody["sale_price"]

		req.ID = id
		req.IsSetDescription = isDescriptionPresent
		req.IsSetSalePrice = isSalePricePresent

		if err := h.validateUpdateProductRequest(&req); err != nil {
			return c.JSON(http.StatusBadRequest, handler.BaseAPIResponse{
				Successful: false,
				ErrorCode:  utils.ToPointer(ErrCodeInvalidRequest),
			})
		}

		ctx := c.Request().Context()
		if err := h.productUseCase.UpdateProduct(ctx, &req); err != nil {
			errCode := ErrCodeUpdateProductFailed
			if errors.Is(err, productusecase.ErrProductNotFound) {
				errCode = ErrProductNotFound
			}

			return c.JSON(http.StatusInternalServerError, handler.BaseAPIResponse{
				Successful: false,
				ErrorCode:  utils.ToPointer(errCode),
			})
		}

		return c.JSON(
			http.StatusOK,
			handler.BaseAPIResponse{
				Successful: true,
			},
		)
	}
}

func (h *ProductHandler) validateUpdateProductRequest(req *productusecase.UpdateProductRequest) error {
	if req.Price != nil && *req.Price < 0 {
		return errors.New("price must be greater or equal to 0")
	}

	if req.SalePrice != nil && *req.SalePrice < 0 {
		return errors.New("sale price must be greater or equal to 0")
	}

	return nil
}

// @Summary Create Product
// @Description Create a new product
// @Tags products
// @ID create-product
// @Accept  json
// @Produce  json
// @Param product body productusecase.CreateProductRequest true "Product Request"
// @Success 200 {object} producthandler.ProductSuccessResponse
// @Failure 400 {object} handler.BaseAPIResponse
// @Failure 500 {object} handler.BaseAPIResponse
// @Router /products [post]
func (h *ProductHandler) CreateProduct() echo.HandlerFunc {
	return func(c *echo.Context) error {
		ctx := c.Request().Context()
		defer ctx.Done()

		var req productusecase.CreateProductRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(
				http.StatusBadRequest,
				handler.BaseAPIResponse{
					Successful: false,
					ErrorCode:  utils.ToPointer(ErrCodeInvalidRequest),
				},
			)
		}

		if err := h.validateCreateProductRequest(&req); err != nil {
			return c.JSON(
				http.StatusBadRequest,
				handler.BaseAPIResponse{
					Successful: false,
					ErrorCode:  utils.ToPointer(ErrCodeInvalidRequest),
				},
			)
		}

		product, err := h.productUseCase.CreateProduct(ctx, &req)
		if err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				handler.BaseAPIResponse{
					Successful: false,
					ErrorCode:  utils.ToPointer(ErrCodeCreateProductFailed),
				},
			)
		}

		return c.JSON(
			http.StatusOK,
			handler.BaseAPIResponse{
				Successful: true,
				Data: ProductResponse{
					ID:          product.ID.String(),
					Name:        product.Name,
					Description: product.Description,
					Price:       product.Price,
					SalePrice:   product.SalePrice,
				},
			},
		)
	}
}

func (h *ProductHandler) validateCreateProductRequest(req *productusecase.CreateProductRequest) error {
	if req.Price < 0 {
		return errors.New("price must be greater or equal to 0")
	}

	if req.SalePrice != nil && *req.SalePrice < 0 {
		return errors.New("sale price must be greater or equal to 0")
	}

	return nil
}
