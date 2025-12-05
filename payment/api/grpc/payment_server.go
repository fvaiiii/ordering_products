package grpc

import (
	"context"

	"github.com/fvaiiii/ordering_products/payment/internal/service"
	v1 "github.com/fvaiiii/ordering_products/shared/pkg/proto/payment/v1"
)

type Server struct {
	v1.UnimplementedPaymentServiceServer
	service *service.PaymentService
}

func NewPaymentServer(svc *service.PaymentService) *Server {
	return &Server{
		service: svc,
	}
}

func (s *Server) PayOrder(ctx context.Context, req *v1.PaymentRequest) (*v1.PaymentResponse, error) {
	return s.service.PayOrder(ctx, req)
}
