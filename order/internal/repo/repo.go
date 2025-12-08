package repo

import (
	"context"

	"github.com/fvaiiii/ordering_products/order/internal/models"
	inventoryv1 "github.com/fvaiiii/ordering_products/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/fvaiiii/ordering_products/shared/pkg/proto/payment/v1"
)

type Order interface {
	Create(order *models.Order) error
	GetByUUID(uuid string) (*models.Order, error)
	Update(order *models.Order) error
	Exists(uuid string) (bool, error)
}

type InventoryClient interface {
	ListProducts(ctx context.Context, filter inventoryv1.ProductsFilter) ([]*inventoryv1.Product, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod paymentv1.PaymentMethod) (string, error)
}
