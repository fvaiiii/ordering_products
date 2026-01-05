package main

import (
	"log"
	"net/http"
	"os"
	"time"

	// Импортируем твои пакеты с алиасами
	httprouter "github.com/fvaiiii/ordering_products/order/internal/api/http" // ← РОУТЕР
	"github.com/fvaiiii/ordering_products/order/internal/api/http/handlers"   // ← HANDLERS
	"github.com/fvaiiii/ordering_products/order/internal/clients"
	"github.com/fvaiiii/ordering_products/order/internal/repository"
	"github.com/fvaiiii/ordering_products/order/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	logger := log.New(os.Stdout, "OrderService ", log.LstdFlags|log.Lshortfile)
	logger.Println("Starting OrderService...")

	logger.Println("Creating repository...")
	repo := repository.NewOrderRepo()

	logger.Println("Connecting to InventoryService...")
	inventoryConn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatal("Failed to connect to InventoryService:", err)
	}
	defer inventoryConn.Close()

	logger.Println("Connecting to PaymentService...")
	paymentConn, err := grpc.NewClient(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatal("Failed to connect to PaymentService:", err)
	}
	defer paymentConn.Close()

	inventoryClient := clients.NewInventoryClient(inventoryConn)
	paymentClient := clients.NewPaymentClient(paymentConn)

	orderService := service.NewOrderService(
		repo,
		inventoryClient,
		paymentClient,
	)

	orderHandler := handlers.NewOrderHandler(orderService)

	router := httprouter.NewRouter(orderHandler)

	logger.Println("Waiting for gRPC connections...")
	time.Sleep(500 * time.Millisecond)

	logger.Println("HTTP server listening on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Fatal("Server failed:", err)
	}
}
