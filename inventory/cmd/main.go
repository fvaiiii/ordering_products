package main

import (
	"log"
	"net"

	grpcapi "github.com/fvaiiii/ordering_products/inventory/internal/api/grpc"
	"github.com/fvaiiii/ordering_products/inventory/internal/repository"
	"github.com/fvaiiii/ordering_products/inventory/internal/service"
	v1 "github.com/fvaiiii/ordering_products/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
)

func main() {
	repo := repository.NewProductsRepo()
	svc := service.NewInventoryService(repo)
	handler := grpcapi.NewServer(svc)

	grpcServer := grpc.NewServer()
	v1.RegisterInventoryServiceServer(grpcServer, handler)

	lis, _ := net.Listen("tcp", ":50051")
	log.Printf("Starting InventoryService on :50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}

}
