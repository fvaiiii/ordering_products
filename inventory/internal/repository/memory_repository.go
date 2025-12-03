package repository

import (
	"context"
	"sync"

	models "github.com/fvaiiii/ordering_products/inventory/internal/models"
	"github.com/fvaiiii/ordering_products/inventory/internal/port/repo"
)

var _ repo.Products = (*ProductsRepo)(nil)

type ProductsRepo struct {
	products map[string]*models.Product
	mu       *sync.RWMutex
}

func NewProductsRepo() *ProductsRepo {
	return &ProductsRepo{
		products: make(map[string]*models.Product),
		mu:       new(sync.RWMutex),
	}
}

func (r *ProductsRepo) GetProduct(ctx context.Context, productID string) (*models.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, ok := r.products[productID]
	if !ok {
		return nil, repo.ErrNotFound
	}

	return product, nil
}

func (r *ProductsRepo) ListProducts(ctx context.Context, filter models.ProductsFilter) ([]*models.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := make([]*models.Product, 0)
	for _, product := range r.products {
		if matchesFilter(product, filter) {
			res = append(res, product)
		}
	}

	return res, nil
}

func matchesFilter(product *models.Product, filter models.ProductsFilter) bool {
	if len(filter.Uuids) > 0 && !contains(filter.Uuids, product.Uuid) {
		return false
	}

	if len(filter.Names) > 0 && !contains(filter.Names, product.Name) {
		return false
	}

	if len(filter.Categories) > 0 && !containsCategory(filter.Categories, product.Category) {
		return false
	}

	if len(filter.ManufacturerCountries) > 0 && !contains(filter.ManufacturerCountries, product.Manufacturer.Country) {
		return false
	}

	return true
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func containsCategory(slice []models.Category, item models.Category) bool {
	for _, c := range slice {
		if c == item {
			return true
		}
	}
	return false
}
