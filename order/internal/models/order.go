package models

type Order struct {
	OrderUuid       string
	UserUuid        string
	ProductUuids    []string
	TotalPrice      float64
	TransactionUuid *string
	PaymentMethod   *string
	Status          OrderStatus
}

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusPendingPayment, OrderStatusPaid, OrderStatusCancelled:
		return true
	default:
		return false
	}
}
