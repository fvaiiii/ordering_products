package clients

import (
	"context"
	"fmt"
	"log"

	"github.com/fvaiiii/ordering_products/order/internal/models"
	"github.com/fvaiiii/ordering_products/order/internal/repo"
	inventoryv1 "github.com/fvaiiii/ordering_products/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
)

var _ = (*repo.InventoryClient)(nil)

type InventoryClient struct {
	conn   *grpc.ClientConn
	client inventoryv1.InventoryServiceClient
}

func NewInventoryClient(conn *grpc.ClientConn) InventoryClient {
	return InventoryClient{
		conn:   conn,
		client: inventoryv1.NewInventoryServiceClient(conn),
	}
}

func (c *InventoryClient) Close() error {
	return c.conn.Close()
}

func (c *InventoryClient) ListParts(ctx context.Context, uuids []string) ([]*models.Product, error) {

	log.Printf("[InventoryClient] I am requesting products: %v", uuids)

	resp, err := c.client.ListProducts(ctx, &inventoryv1.ListProductsRequest{
		Filter: &inventoryv1.ProductsFilter{Uuids: uuids},
	})
	if err != nil {
		log.Printf("[InventoryClient] error from InventoryService: %v", err)
		return nil, fmt.Errorf("failed to get products from inventory: %w", err)
	}

	log.Printf("[InventoryClient] %d products received", len(resp.Products))

	if len(resp.Products) != len(uuids) {
		return nil, fmt.Errorf("some products not found: expected %d, got %d", len(uuids), len(resp.Products))
	}

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
}
