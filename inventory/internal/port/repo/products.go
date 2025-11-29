package repo

import (
	"context"
	"errors"

	models "github.com/fvaiiii/ordering_products/inventory/internal/models"
)

var (
	ErrProductNotFound      = errors.New("product not found")
	ErrProductAlreadyExists = errors.New("product already exists")
)

type Products interface {
	GetProduct(ctx context.Context, productID string) (*models.Product, error)
	ListProducts(ctx context.Context, filter models.ProductsFilter) ([]*models.Product, error)
}
