package clients

import (
	"context"
	"fmt"

<<<<<<< Updated upstream
=======
	"github.com/fvaiiii/ordering_products/order/internal/models"
>>>>>>> Stashed changes
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

<<<<<<< Updated upstream
func (c *InventoryClient) ListParts(ctx context.Context, uuids []string) ([]*inventoryv1.Product, error) {
=======
func (c *InventoryClient) ListParts(ctx context.Context, uuids []string) ([]*models.Product, error) {
>>>>>>> Stashed changes
	resp, err := c.client.ListProducts(ctx, &inventoryv1.ListProductsRequest{
		Filter: &inventoryv1.ProductsFilter{Uuids: uuids},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get products from inventory: %w", err)
	}

<<<<<<< Updated upstream
	if len(resp.Products) != len(uuids) {
		return nil, fmt.Errorf("some products not found: expected %d, got %d", len(uuids), len(resp.Products))
	}
	return resp.Products, nil
=======
	products := make([]*models.Product, 0, len(resp.Products))
	for _, protoProduct := range resp.Products {
		if protoProduct == nil {
			continue
		}

		product := &models.Product{
			UUID:  protoProduct.Uuid,
			Price: protoProduct.Price,
			Name:  protoProduct.Name,
		}
		products = append(products, product)
	}

	return products, nil
>>>>>>> Stashed changes
}
