package main

import (
	"log"
	"net"

	grpcapi "github.com/fvaiiii/ordering_products/payment/internal/api/grpc"
	"github.com/fvaiiii/ordering_products/payment/internal/service"
	v1 "github.com/fvaiiii/ordering_products/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
)

func main() {
	svc := service.NewPaymentService()
	handler := grpcapi.NewPaymentServer(svc)

	grpcServer := grpc.NewServer()
	v1.RegisterPaymentServiceServer(grpcServer, handler)

	lis, _ := net.Listen("tcp", ":50052")
	log.Printf("Starting PaymentService on :50052")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
