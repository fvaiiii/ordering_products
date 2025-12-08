package clients

import (
	"context"
	"fmt"

	"github.com/fvaiiii/ordering_products/order/internal/repo"
	inventoryv1 "github.com/fvaiiii/ordering_products/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
)

var _ = (*repo.InventoryClient)(nil)

type InventoryClient struct {
	conn   *grpc.ClientConn
	client inventoryv1.InventoryServiceClient
}

func NewInventoryClient(conn *grpc.ClientConn) *InventoryClient {
	return &InventoryClient{
		conn:   conn,
		client: inventoryv1.NewInventoryServiceClient(conn),
	}
}

func (c *InventoryClient) Close() error {
	return c.conn.Close()
}

func (c *InventoryClient) GetProducts(ctx context.Context, uuids []string) ([]*inventoryv1.Product, error) {
	resp, err := c.client.ListProducts(ctx, &inventoryv1.ListProductsRequest{
		Filter: &inventoryv1.ProductsFilter{Uuids: uuids},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get products from inventory: %w", err)
	}

	if len(resp.Products) != len(uuids) {
		return nil, fmt.Errorf("some products not found: expected %d, got %d", len(uuids), len(resp.Products))
	}
	return resp.Products, nil
}
