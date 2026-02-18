package repositories

import "go-product-restapi/internal/entity"

type ProductMapper[T any] interface {
	ToDTO(product entity.Product) T
	ToEntity(productDTO T) entity.Product
}