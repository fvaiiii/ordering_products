package models

type PaymentMethod int32

const (
	UNKNOWN     PaymentMethod = 0
	CARD        PaymentMethod = 1
	SBP         PaymentMethod = 2
	CREDIT_CARD PaymentMethod = 3
)
