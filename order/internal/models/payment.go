package models

import paymentv1 "github.com/fvaiiii/ordering_products/shared/pkg/proto/payment/v1"

type PaymentMethod string

const (
	PaymentMethodUnknown    PaymentMethod = "UNKNOWN"
	PaymentMethodCard       PaymentMethod = "CARD"
	PaymentMethodSBP        PaymentMethod = "SBP"
	PaymentMethodCreditCard PaymentMethod = "CREDIT_CARD"
)

func (pm PaymentMethod) ToProto() paymentv1.PaymentMethod {
	switch pm {
	case PaymentMethodCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case PaymentMethodSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	case PaymentMethodCreditCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNKNOWN
	}
}
