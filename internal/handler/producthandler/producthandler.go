package producthandler

import (
	"net/http"

	"go-product-restapi/internal/handler"
	"go-product-restapi/internal/usecases/productusecase"

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

func (h *ProductHandler) PatchProduct() echo.HandlerFunc {
	return func(c *echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(
				http.StatusBadRequest,
				handler.BaseAPIResponse{
					Successful: false,
				},
			)
		}

		ctx := c.Request().Context()
		defer ctx.Done()

		var req productusecase.UpdateProductRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(
				http.StatusBadRequest,
				handler.BaseAPIResponse{
					Successful: false,
				},
			)
		}

		req.ID = id

		if err := h.productUseCase.UpdateProduct(ctx, &req); err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				handler.BaseAPIResponse{
					Successful: false,
				},
			)
		}

		return c.JSON(
			http.StatusOK,
			handler.BaseAPIResponse{
				Successful: true,
			},
		)
	}
}

func (h *ProductHandler) CreateProduct() echo.HandlerFunc {
	return func(c *echo.Context) error {
		return c.JSON(
			http.StatusCreated,
			handler.BaseAPIResponse{
				Successful: true,
				Data:       nil,
			},
		)
	}
}
