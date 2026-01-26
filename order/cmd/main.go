package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	httprouter "github.com/fvaiiii/ordering_products/order/internal/api/http"
	"github.com/fvaiiii/ordering_products/order/internal/api/http/handlers"
	"github.com/fvaiiii/ordering_products/order/internal/clients"
	repository "github.com/fvaiiii/ordering_products/order/internal/repository/postgres"
	"github.com/fvaiiii/ordering_products/order/internal/service"
	"github.com/fvaiiii/ordering_products/pkg/migrator"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	logger := log.New(os.Stdout, "OrderService ", log.LstdFlags|log.Lshortfile)

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://order_user:order_password@localhost:5435/order_db?sslmode=disable"
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		logger.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	logger.Println("Applying database migrations")
	migrator.Init(db, "./migrations")
	if err := migrator.Migrator().Up(); err != nil {
		logger.Fatal("Failed to run migrations:", err)
	}
	logger.Println("Migrations applied successfully")

	logger.Println("Creating database connection pool")
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		logger.Fatal("Failed to create connection pool:", err)
	}
	defer pool.Close()

	logger.Println("Starting OrderService...")
	logger.Println("Creating repository...")
	repo := repository.NewOrderRepository(pool)
	logger.Println("Connecting to InventoryService...")
	inventoryConn, err := grpc.NewClient(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		logger.Fatal("Failed to connect to InventoryService:", err)
	}
	defer inventoryConn.Close()

	logger.Println("Connecting to PaymentService...")
	paymentConn, err := grpc.NewClient(
		"localhost:50053",
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
