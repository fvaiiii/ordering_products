package repository

import (
	"context"
	"errors"
	"sync"

	models "github.com/fvaiiii/ordering_products/inventory/internal/models"
	"github.com/fvaiiii/ordering_products/inventory/internal/port/repo"
)

var _ repo.Products = (*ProductsRepo)(nil)

type ProductsRepo struct {
	Products map[string]*models.Product
	mu       *sync.RWMutex
}

func NewProductsRepo() *ProductsRepo {
	return &ProductsRepo{
		Products: make(map[string]*models.Product),
		mu:       new(sync.RWMutex),
	}
}

func (r *ProductsRepo) AddProduct(product *models.Product) error {
	if product == nil || product.Uuid == "" {
		return errors.New("invalid product")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.Products[product.Uuid]; exists {
		return errors.New("product already exists: " + product.Uuid)
	}

	r.Products[product.Uuid] = product
	return nil
}

func (r *ProductsRepo) GetProduct(ctx context.Context, productID string) (*models.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, ok := r.Products[productID]
	if !ok {
		return nil, repo.ErrNotFound
	}

	return product, nil
}

func (r *ProductsRepo) ListProducts(ctx context.Context, filter models.ProductsFilter) ([]*models.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := make([]*models.Product, 0)
	for _, product := range r.Products {
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
