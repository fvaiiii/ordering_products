package repo

import (
	"context"

	"github.com/fvaiiii/ordering_products/order/internal/models"
)

type Order interface {
	Create(ctx context.Context, order *models.Order) error
	GetByUUID(ctx context.Context, uuid string) (*models.Order, error)
	Update(ctx context.Context, order *models.Order) error
	Exists(ctx context.Context, uuid string) bool
}

type InventoryClient interface {
	ListParts(ctx context.Context, uuids []string) ([]*models.Product, error)
	Close() error
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod models.PaymentMethod) (string, error)
	Close() error
}
