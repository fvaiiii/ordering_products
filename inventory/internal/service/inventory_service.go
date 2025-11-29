package service

import (
	"context"

	"github.com/fvaiiii/ordering_products/inventory/internal/models"
	"github.com/fvaiiii/ordering_products/inventory/internal/port/repo"
)

type InventoryService struct {
	repo repo.Products
}

func NewInventoryService(repo repo.Products) *InventoryService {
	return &InventoryService{
		repo: repo,
	}
}

func (i *InventoryService) GetProduct(ctx context.Context, productID string) (*models.Product, error) {
	product, err := i.repo.GetProduct(ctx, productID)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (i *InventoryService) ListProducts(ctx context.Context, filter models.ProductsFilter) ([]*models.Product, error) {
	products, err := i.repo.ListProducts(ctx, filter)
	if err != nil {
		return nil, err
	}

	return products, nil

}
