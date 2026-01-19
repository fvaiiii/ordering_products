package repo

import (
	"context"

	v1 "github.com/fvaiiii/ordering_products/shared/pkg/proto/payment/v1"
)

//go:generate mockery --name=Payment --filename=payment_mock.go --output=../../../mocks --inpackage=false
type Payment interface {
	PayOrder(ctx context.Context, req *v1.PaymentRequest) (*v1.PaymentResponse, error)
}
