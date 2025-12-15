package repository

import (
	"context"
	"sync"

	"github.com/fvaiiii/ordering_products/order/internal/models"
	"github.com/fvaiiii/ordering_products/order/internal/repo"
)

var _ = (*repo.Order)(nil)

type OrderRepo struct {
	orders map[string]*models.Order
	mu     *sync.RWMutex
}

func NewOrderRepo() *OrderRepo {
	return &OrderRepo{
		orders: make(map[string]*models.Order),
		mu:     new(sync.RWMutex),
	}
}

func (r *OrderRepo) Create(ctx context.Context, order *models.Order) error {

	if order == nil {
		return repo.ErrInvalidData
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.orders[order.OrderUuid]; exists {
		return repo.ErrAlreadyExists
	}

	r.orders[order.OrderUuid] = order

	return nil
}

func (r *OrderRepo) GetByUUID(ctx context.Context, uuid string) (*models.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[uuid]
	if !ok {
		return nil, repo.ErrNotFound
	}

	return order, nil
}
func (r *OrderRepo) Update(ctx context.Context, order *models.Order) error {

	if order == nil {
		return repo.ErrInvalidData
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.orders[order.OrderUuid]; !exists {
		return repo.ErrNotFound
	}
	r.orders[order.OrderUuid] = order
	return nil
}
func (r *OrderRepo) Exists(ctx context.Context, uuid string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.orders[uuid]
	return ok
}
