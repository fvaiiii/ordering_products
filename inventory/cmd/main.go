package main

import (
	"context"
	"log"
	"net"
	"time"

	grpcapi "github.com/fvaiiii/ordering_products/inventory/internal/api/grpc"
	"github.com/fvaiiii/ordering_products/inventory/internal/models"
	"github.com/fvaiiii/ordering_products/inventory/internal/repository"
	"github.com/fvaiiii/ordering_products/inventory/internal/service"
	v1 "github.com/fvaiiii/ordering_products/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SeedTestData(repo *repository.ProductsRepo) {
	// ctx := context.Background()
	product1 := &models.Product{
		Uuid:          "test-1",
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
		Uuid:          "test-2",
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

	repo.AddProduct(product1)
	repo.AddProduct(product2)
	log.Printf("Тестовые данные подготовлены: %s, %s", product1.Name, product2.Name)
}

func main() {
	repo := repository.NewProductsRepo()
	SeedTestData(repo)

	products, _ := repo.ListProducts(context.Background(), models.ProductsFilter{})
	log.Printf("В репозитории: %d продуктов", len(products))

	svc := service.NewInventoryService(repo)
	handler := grpcapi.NewServer(svc)

	grpcServer := grpc.NewServer()
	v1.RegisterInventoryServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("unable to listen port 50051: %v", err)
	}
	log.Printf("Starting InventoryService on :50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("server crushed: %v", err)
	}

}
