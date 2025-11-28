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

func (r *ProductsRepo) Create(ctx context.Context, product *models.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if product == nil {
		return repo.ErrProductNotFound
	}

	if _, exists := r.products[product.Uuid]; exists {
		return repo.ErrProductAlreadyExists
	}
	r.products[product.Uuid] = product

	return nil
}

func (r *ProductsRepo) Read(ctx context.Context, productID string) (*models.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, ok := r.products[productID]
	if !ok {
		return nil, repo.ErrProductNotFound
	}

	return product, nil
}

func (r *ProductsRepo) Upsert(ctx context.Context, products []*models.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, p := range products {
		if p == nil {
			continue
		}
		r.products[p.Uuid] = p
	}

	return nil
}

func (r *ProductsRepo) List(ctx context.Context) ([]*models.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]*models.Product, 0, len(r.products))
	for _, p := range r.products {
		out = append(out, p)
	}

	return out, nil
}

func (r *ProductsRepo) Delete(ctx context.Context, productID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.products[productID]
	if !ok {
		return repo.ErrProductNotFound
	}

	delete(r.products, productID)
	return nil

}
