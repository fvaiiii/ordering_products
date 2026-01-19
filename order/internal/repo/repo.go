package repo

import (
	"context"

	"github.com/fvaiiii/ordering_products/order/internal/models"
)

//go:generate mockery --name=Order --filename=order_mock.go --output=../service/mocks --inpackage=false
type Order interface {
	Create(ctx context.Context, order *models.Order) error
	GetByUUID(ctx context.Context, uuid string) (*models.Order, error)
	Update(ctx context.Context, order *models.Order) error
	Exists(ctx context.Context, uuid string) bool
}

//go:generate mockery --name=InventoryClient --filename=inventory_client_mock.go --output=../clients/mocks --inpackage=false
type InventoryClient interface {
	ListParts(ctx context.Context, uuids []string) ([]*models.Product, error)
	Close() error
}

//go:generate mockery --name=PaymentClient --filename=payment_client_mock.go --output=../clients/mocks --inpackage=false
type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod models.PaymentMethod) (string, error)
	Close() error
}
