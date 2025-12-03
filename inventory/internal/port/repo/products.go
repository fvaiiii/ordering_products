package repo

import (
	"context"

	models "github.com/fvaiiii/ordering_products/inventory/internal/models"
)

type Products interface {
	GetProduct(ctx context.Context, productID string) (*models.Product, error)
	ListProducts(ctx context.Context, filter models.ProductsFilter) ([]*models.Product, error)
}
