package repo

import (
	"context"

	"github.com/fvaiiii/ordering_products/order/internal/models"
<<<<<<< Updated upstream
	inventoryv1 "github.com/fvaiiii/ordering_products/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/fvaiiii/ordering_products/shared/pkg/proto/payment/v1"
=======
>>>>>>> Stashed changes
)

type Order interface {
	Create(ctx context.Context, order *models.Order) error
	GetByUUID(ctx context.Context, uuid string) (*models.Order, error)
	Update(ctx context.Context, order *models.Order) error
	Exists(ctx context.Context, uuid string) bool
}

type InventoryClient interface {
<<<<<<< Updated upstream
	ListParts(ctx context.Context, uuids []string) ([]*inventoryv1.Product, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod paymentv1.PaymentMethod) (string, error)
=======
	ListParts(ctx context.Context, uuids []string) ([]*models.Product, error)
	Close() error
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod models.PaymentMethod) (string, error)
	Close() error
>>>>>>> Stashed changes
}
