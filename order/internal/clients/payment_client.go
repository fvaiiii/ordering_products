package clients

import (
	"context"
<<<<<<< Updated upstream
	"errors"
	"fmt"

=======
	"fmt"

	"github.com/fvaiiii/ordering_products/order/internal/models"
>>>>>>> Stashed changes
	"github.com/fvaiiii/ordering_products/order/internal/repo"
	paymentv1 "github.com/fvaiiii/ordering_products/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
)

var _ = (*repo.PaymentClient)(nil)

type PaymentClient struct {
	conn   *grpc.ClientConn
	client paymentv1.PaymentServiceClient
}

func NewPaymentClient(conn *grpc.ClientConn) *PaymentClient {
	return &PaymentClient{
		conn:   conn,
		client: paymentv1.NewPaymentServiceClient(conn),
	}
}
func (c *PaymentClient) Close() error {
	return c.conn.Close()
}

<<<<<<< Updated upstream
func (c *PaymentClient) PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod paymentv1.PaymentMethod) (string, error) {
=======
func (c *PaymentClient) PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod models.PaymentMethod) (string, error) {
>>>>>>> Stashed changes

	req := &paymentv1.PaymentRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
<<<<<<< Updated upstream
		PaymentMethod: paymentMethod,
=======
		PaymentMethod: paymentMethod.ToProto(),
>>>>>>> Stashed changes
	}

	resp, err := c.client.PayOrder(ctx, req)
	if err != nil {
		return "", fmt.Errorf("payment service error: %w", err)
	}

<<<<<<< Updated upstream
	if resp.GetTransactionUuid() == "" {
		return "", errors.New("payment service returned empty transaction uuid")
	}

=======
>>>>>>> Stashed changes
	return resp.GetTransactionUuid(), nil

}
