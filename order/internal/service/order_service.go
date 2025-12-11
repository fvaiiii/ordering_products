package service

import (
	"context"
	"fmt"

	"github.com/fvaiiii/ordering_products/order/internal/clients"
	"github.com/fvaiiii/ordering_products/order/internal/models"
	"github.com/fvaiiii/ordering_products/order/internal/repo"
	"github.com/fvaiiii/ordering_products/order/internal/repository"
	paymentv1 "github.com/fvaiiii/ordering_products/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
)

type OrderService struct {
	repo      repository.OrderRepo
	inventory clients.InventoryClient
	payment   clients.PaymentClient
}

func NewOrderService(repo repository.OrderRepo) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, userUUID string, productsUUIDs []string) (*models.Order, error) {
	if userUUID == "" {
		return nil, fmt.Errorf("user_uuid is required")
	}

	if len(productsUUIDs) == 0 {
		return nil, fmt.Errorf("at least one product is required")
	}

	products, err := s.inventory.ListParts(ctx, productsUUIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	if len(products) != len(productsUUIDs) {
		return nil, fmt.Errorf("some products are not exist")
	}
	var totalPrice float64
	for _, product := range products {
		totalPrice += product.Price
	}
	orderUUID := generateOrderUUID()

	order := &models.Order{
		OrderUuid:    orderUUID,
		UserUuid:     userUUID,
		ProductUuids: productsUUIDs,
		TotalPrice:   totalPrice,
		Status:       models.OrderStatusPendingPayment,
	}

	err = s.repo.Create(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	return order, nil
}

/*
- Находит заказ по `order_uuid`. Если не существует — возвращает 404 Not Found.
- Вызывает `PaymentService.PayOrder`, передаёт `user_uuid`, `order_uuid` и `payment_method`. Получает`transaction_uuid`.
- Обновляет заказ: статус → `PAID`, сохраняет `transaction_uuid`, `payment_method`.
*/
func (s *OrderService) PayOrder(ctx context.Context, orderUUID string, paymentMethod paymentv1.PaymentMethod) (string, error) {
	if orderUUID == "" {
		return "", fmt.Errorf("order_uuid is required")
	}
	order, err := s.repo.GetByUUID(ctx, orderUUID)
	if err != nil {
		return "", repo.ErrNotFound
	}

	if order.Status != models.OrderStatusPendingPayment {
		return "", repo.ErrInvalidData
	}

	transactionUuid, err := s.payment.PayOrder(ctx, orderUUID, order.UserUuid, paymentMethod)
	if err != nil {
		return "", fmt.Errorf("failed to pay order: %w", err)
	}

	order.Status = models.OrderStatusPaid
	order.TransactionUuid = &transactionUuid
	order.PaymentMethod = &paymentMethod

	s.repo.Update(ctx, order)

	return transactionUuid, nil

}

/*
- Ищет заказ по UUID.
- Если найден — возвращает.
- Если не найден — 404 Not Found.
---
*/
func (s *OrderService) GetOrderByUUID(ctx context.Context, orderUUID string) (*models.Order, error) {
	return s.repo.GetByUUID(ctx, orderUUID)
}

/*
- Проверяет статус заказа.
- Если `PENDING_PAYMENT` — меняет статус на `CANCELLED`.
- Если `PAID` — возвращает ошибку 409.
*/
func (s *OrderService) CancelOrder(ctx context.Context, orderUUID string) error {
	order, err := s.repo.GetByUUID(ctx, orderUUID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}
	if !order.Status.IsValid() {
		return repo.ErrInvalidData
	}
	if order.Status == models.OrderStatusPendingPayment {
		order.Status = models.OrderStatusCancelled
	}
	if order.Status == models.OrderStatusPaid { // После сохранения проверяем?
		return fmt.Errorf("conflict")
	}
	if err := s.repo.Update(ctx, order); err != nil { // ВСЕГДА сохраняем!
		return fmt.Errorf("failed to save: %w", err)
	}
	return nil
}

func generateOrderUUID() string {
	return uuid.New().String()
}
