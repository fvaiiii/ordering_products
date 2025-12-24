package dto

type CreateOrderRequest struct {
	UserUuid     string   `json:"user_uuid" binding:"required"`
	ProductUuids []string `json:"products_uuids" binding:"required"`
}

type PayOrderRequest struct {
	PaymentMethod string `json:"payment_method" binding:"required,oneof=CARD SBP CREDIT_CARD"`
}

type CreateOrderResponse struct {
	OrderUUID  string  `json:"order_uuid"`
	TotalPrice float64 `json:"total_price"`
}

type PayOrderResponse struct {
	TransactionUUID string `json:"transaction_uuid"`
}

type GetOrderResponse struct {
	OrderUUID       string   `json:"order_uuid"`
	UserUUID        string   `json:"user_uuid"`
	ProductUUIDs    []string `json:"products_uuids"`
	TotalPrice      float64  `json:"total_price"`
	TransactionUUID string   `json:"transaction_uuid,omitempty"`
	PaymentMethod   string   `json:"payment_method,omitempty"`
	Status          string   `json:"status"`
}
