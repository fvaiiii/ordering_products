package service

import (
	"context"
	"log"

	v1 "github.com/fvaiiii/ordering_products/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentService struct{}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

func (s *PaymentService) PayOrder(ctx context.Context, req *v1.PaymentRequest) (*v1.PaymentResponse, error) {
	if req.GetOrderUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "order_uuid is required")
	}
	if req.GetUserUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_uuid is required")
	}
	if req.GetPaymentMethod() == v1.PaymentMethod_PAYMENT_METHOD_UNKNOWN {
		return nil, status.Error(codes.InvalidArgument, "payment_method is required")
	}

	transactionUuid := uuid.New().String()

	log.Printf("payment was successful, transaction_uuid: %s", transactionUuid)

	return &v1.PaymentResponse{
		TransactionUuid: transactionUuid,
	}, nil

}
