package dto

type Order struct {
	OrderUuid       string   `json:"order_uuid"`
	UserUuid        string   `json:"user_uuid"`
	ProductUuids    []string `json:"products_uuids"`
	TotalPrice      float64  `json:"total_price"`
	TransactionUuid *string  `json:"transaction_uuid,omitempty"`
	PaymentMethod   *string  `json:"payment_method,omitempty"`
	Status          string   `json:"status"`
}

type CreateOrderRequest struct {
	UserUuid     string   `json:"user_uuid"`
	ProductUuids []string `json:"products_uuids"`
}
type CreateOrderResponse struct {
	OrderUUID  string  `json:"order_uuid"`
	TotalPrice float64 `json:"total_price"`
}

type CreatePaymentRequest struct {
	PaymentMethod string `json:"paymentMethod"`
}
type CreatePaymentResponse struct {
	TransactionUuid string `json:"transaction_uuid"`
}

type GetPaymentRequest struct {
	OrderUuid string `json:"order_uuid"`
}

type GetPaymentRequestResponse struct {
	OrderUuid       string   `json:"order_uuid"`
	UserUuid        string   `json:"user_uuid"`
	ProductUuids    []string `json:"products_uuids"`
	TotalPrice      float64  `json:"total_price"`
	TransactionUuid string   `json:"transaction_uuid"`
	PaymentMethod   string   `json:"payment_method"`
	Status          string   `json:"status"`
}
