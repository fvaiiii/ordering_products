package main

import (
	"context"
	"log"
	"net"
	"time"

	grpcapi "github.com/fvaiiii/ordering_products/inventory/internal/api/grpc"
	"github.com/fvaiiii/ordering_products/inventory/internal/config"
	"github.com/fvaiiii/ordering_products/inventory/internal/models"
	"github.com/fvaiiii/ordering_products/inventory/internal/repository/redisrepo"
	"github.com/fvaiiii/ordering_products/inventory/internal/service"
	v1 "github.com/fvaiiii/ordering_products/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SeedTestData(repo *redisrepo.RedisRepository) error {
	product1 := &models.Product{
		Uuid:          "123e4567-e89b-12d3-a456-426614174001",
		Name:          "Яблоки",
		Description:   "Свежие зеленые яблоки",
		Price:         120.50,
		StockQuantity: 100,
		Category:      models.FRUITS,
		Manufacturer: models.Manufacturer{
			Name:    "Фруктовая ферма",
			Country: "Россия",
			Website: "https://fruits.ru",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	product2 := &models.Product{
		Uuid:          "123e4567-e89b-12d3-a456-426614174002",
		Name:          "Морковь",
		Description:   "Свежая молодая морковь",
		Price:         80.25,
		StockQuantity: 50,
		Category:      models.VEGETABLES,
		Manufacturer: models.Manufacturer{
			Name:    "Овощной двор",
			Country: "Россия",
			Website: "https://vegetables.ru",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := repo.AddProduct(product1); err != nil {
		log.Printf("Failed to add product1: %w", err)
		return err
	}

	if err := repo.AddProduct(product2); err != nil {
		log.Printf("Failed to add product2: %w", err)
		return err
	}

	log.Printf("Test data added to Redis: %s, %s", product1.Name, product2.Name)
	return nil
}

func main() {
	cfg := config.Load()
	log.Printf("=== Starting Inventory Service ===")
	log.Printf("Port: %s", cfg.Server.Port)
	log.Printf("Redis: %s (DB: %d)", cfg.Redis.Addr, cfg.Redis.DB)
	log.Printf("Log Level: %s", cfg.Logging.Level)

	repo, err := redisrepo.NewRedisRepository(&cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %w", err)
	}

	ctx := context.Background()

	existingProducts, err := repo.ListProducts(ctx, models.ProductsFilter{})
	if err != nil {
		log.Printf("Warning: failed to check existing products: %v", err)
	} else {
		log.Printf("Found %d existing products in Redis", len(existingProducts))
	}

	if len(existingProducts) == 0 {
		log.Println("Redis is empty, seeding test data...")
		if err := SeedTestData(repo); err != nil {
			log.Printf("Warning: failed to seed test data: %v", err)
		} else {
			products, _ := repo.ListProducts(ctx, models.ProductsFilter{})
			log.Printf("Successfully seeded %d products", len(products))
		}
	}

	testProduct1, err := repo.GetProduct(ctx, "123e4567-e89b-12d3-a456-426614174001")
	if err != nil {
		log.Printf("Warning: test product 1 not found: %v", err)
	} else {
		log.Printf("Test product 1 available: %s (Stock: %d)",
			testProduct1.Name, testProduct1.StockQuantity)
	}

	testProduct2, err := repo.GetProduct(ctx, "123e4567-e89b-12d3-a456-426614174002")
	if err != nil {
		log.Printf("Warning: test product 2 not found: %v", err)
	} else {
		log.Printf("Test product 2 available: %s (Stock: %d)",
			testProduct2.Name, testProduct2.StockQuantity)
	}

	svc := service.NewInventoryService(repo)
	handler := grpcapi.NewServer(svc)

	grpcServer := grpc.NewServer()
	v1.RegisterInventoryServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":"+cfg.Server.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Server.Port, err)
	}

	log.Printf("InventoryService (Redis) listening on :%s", cfg.Server.Port)
	log.Printf("Ready to accept gRPC connections")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}
