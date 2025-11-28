package repo

import (
	"context"
	"errors"

	models "github.com/fvaiiii/ordering_products/inventory/internal/models"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrProductAlreadyExists = errors.New("product already exists")
)

type Products interface {
	Create(ctx context.Context, product *models.Product) error
	Read(ctx context.Context, productID string) (*models.Product, error)
	Upsert(ctx context.Context, products []*models.Product) error
	List(ctx context.Context) ([]*models.Product, error)
	Delete(ctx context.Context, productID string) error
}
