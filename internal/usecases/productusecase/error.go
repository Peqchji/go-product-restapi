package productusecase

import "errors"

var (
	ErrProductNotFound     = errors.New("product not found")
	ErrCreateProductFailed = errors.New("failed to create product")
	ErrUpdateProductFailed = errors.New("failed to update product")
)
