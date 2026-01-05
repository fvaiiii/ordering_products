package clients

import (
	"context"
	"fmt"

	"github.com/fvaiiii/ordering_products/order/internal/models"
	"github.com/fvaiiii/ordering_products/order/internal/repo"
	paymentv1 "github.com/fvaiiii/ordering_products/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
)

var _ = (*repo.PaymentClient)(nil)

type PaymentClient struct {
	conn   *grpc.ClientConn
	client paymentv1.PaymentServiceClient
}

func NewPaymentClient(conn *grpc.ClientConn) PaymentClient {
	return PaymentClient{
		conn:   conn,
		client: paymentv1.NewPaymentServiceClient(conn),
	}
}
func (c *PaymentClient) Close() error {
	return c.conn.Close()
}

func (c *PaymentClient) PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod models.PaymentMethod) (string, error) {

	req := &paymentv1.PaymentRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: paymentMethod.ToProto(),
	}

	resp, err := c.client.PayOrder(ctx, req)
	if err != nil {
		return "", fmt.Errorf("payment service error: %w", err)
	}

	return resp.GetTransactionUuid(), nil

}
