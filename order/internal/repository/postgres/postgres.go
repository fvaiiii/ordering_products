package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fvaiiii/ordering_products/order/internal/models"
	"github.com/fvaiiii/ordering_products/order/internal/repo"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		pool: pool,
	}
}

func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	productUUIDsJSON, err := json.Marshal(order.ProductUuids)
	if err != nil {
		return fmt.Errorf("failed to marshal product UUIDs: %w", err)
	}
	query := `
		INSERT INTO orders 
		(order_uuid, user_uuid, product_uuids, total_price, transaction_uuid, payment_method, status)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = r.pool.Exec(ctx, query,
		order.OrderUuid,
		order.UserUuid,
		productUUIDsJSON,
		order.TotalPrice,
		order.TransactionUuid,
		order.PaymentMethod,
		order.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	return nil
}

func (r *OrderRepository) GetByUUID(ctx context.Context, uuid string) (*models.Order, error) {
	query := `
		SELECT order_uuid, user_uuid, product_uuids, total_price, transaction_uuid, payment_method, status 
		FROM orders
		WHERE order_uuid = $1 
	`

	var order models.Order
	var productUUIDsJSON []byte
	err := r.pool.QueryRow(ctx, query, uuid).Scan(
		&order.OrderUuid,
		&order.UserUuid,
		&productUUIDsJSON,
		&order.TotalPrice,
		&order.TransactionUuid,
		&order.PaymentMethod,
		&order.Status,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
	}

	if err := json.Unmarshal(productUUIDsJSON, &order.ProductUuids); err != nil {
		return nil, fmt.Errorf("failed to unmarshal product UUIDs: %w", err)
	}

	return &order, nil
}

func (r *OrderRepository) Update(ctx context.Context, order *models.Order) error {
	productUUIDsJSON, err := json.Marshal(order.ProductUuids)
	if err != nil {
		return fmt.Errorf("failed to marshal product UUIDs: %w", err)
	}
	query := `
		UPDATE orders
		SET 
			user_uuid = $2, 
			product_uuids = $3, 
			total_price = $4, 
			transaction_uuid = $5, 
			payment_method = $6, 
			status = $7
		WHERE order_uuid = $1
	`

	result, err := r.pool.Exec(ctx, query,
		order.OrderUuid,
		order.UserUuid,
		productUUIDsJSON,
		order.TotalPrice,
		order.TransactionUuid,
		order.PaymentMethod,
		order.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	rowsAffecred := result.RowsAffected()
	if rowsAffecred == 0 {
		return repo.ErrNotFound
	}

	return nil
}

func (r *OrderRepository) Exists(ctx context.Context, uuid string) bool {
	query := `
        SELECT EXISTS(
            SELECT 1 FROM orders WHERE order_uuid = $1
        )
    `

	var exists bool
	err := r.pool.QueryRow(ctx, query, uuid).Scan(&exists)

	if err != nil {
		return false
	}
	return exists
}
