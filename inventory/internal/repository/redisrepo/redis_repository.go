package redisrepo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/fvaiiii/ordering_products/inventory/internal/config"
	models "github.com/fvaiiii/ordering_products/inventory/internal/models"
	"github.com/fvaiiii/ordering_products/inventory/internal/port/repo"
	"github.com/redis/go-redis/v9"
)

var _ repo.Products = (*RedisRepository)(nil)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(cfg *config.RedisConfig) (*RedisRepository, error) {
	log.Printf("Connecting to Redis at %s (DB: %d)", cfg.Addr, cfg.DB)
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Printf("Successfully connected to Redis")
	return &RedisRepository{client: client}, nil
}

func (r *RedisRepository) AddProduct(product *models.Product) error {
	if product == nil || product.Uuid == "" {
		return errors.New("invalid product")
	}

	key := fmt.Sprintf("product:%s", product.Uuid)
	exists, err := r.client.Exists(context.Background(), key).Result()
	if err != nil {
		return fmt.Errorf("redis exists error: %w", err)
	}
	if exists > 0 {
		return repo.ErrAlreadyExists
	}

	data, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("marshal product error: %w", err)
	}

	if err := r.client.Set(context.Background(), key, data, 0).Err(); err != nil {
		return fmt.Errorf("redis set error: %w", err)
	}

	log.Printf("Product saved: %s", product.Uuid)
	return nil
}

func (r *RedisRepository) GetProduct(ctx context.Context, productID string) (*models.Product, error) {
	key := fmt.Sprintf("product:%s", productID)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, repo.ErrNotFound
		}
		return nil, fmt.Errorf("redis get error: %w", err)
	}

	var product models.Product
	if err := json.Unmarshal(data, &product); err != nil {
		return nil, fmt.Errorf("unmarshal product error: %w", err)
	}

	return &product, nil
}

func (r *RedisRepository) ListProducts(ctx context.Context, filter models.ProductsFilter) ([]*models.Product, error) {
	var (
		cursor   uint64
		products []*models.Product
	)

	for {
		keys, nextCursor, err := r.client.Scan(
			ctx,
			cursor,
			"product:*",
			100,
		).Result()
		if err != nil {
			return nil, fmt.Errorf("redis scan error: %w", err)
		}

		if len(keys) > 0 {
			pipe := r.client.Pipeline()
			cmds := make([]*redis.StringCmd, 0, len(keys))

			for _, key := range keys {
				cmds = append(cmds, pipe.Get(ctx, key))
			}

			if _, err := pipe.Exec(ctx); err != nil && err != redis.Nil {
				return nil, fmt.Errorf("pipeline exec error: %w", err)
			}

			for _, cmd := range cmds {
				data, err := cmd.Bytes()
				if err != nil {
					continue
				}

				var product models.Product
				if err := json.Unmarshal(data, &product); err != nil {
					continue
				}
				if matchesFilter(&product, filter) {
					products = append(products, &product)
				}
			}
		}

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}

	return products, nil
}

func matchesFilter(product *models.Product, filter models.ProductsFilter) bool {
	if len(filter.Uuids) > 0 && !contains(filter.Uuids, product.Uuid) {
		return false
	}

	if len(filter.Names) > 0 && !contains(filter.Names, product.Name) {
		return false
	}

	if len(filter.Categories) > 0 && !containsCategory(filter.Categories, product.Category) {
		return false
	}

	if len(filter.ManufacturerCountries) > 0 &&
		!contains(filter.ManufacturerCountries, product.Manufacturer.Country) {
		return false
	}

	return true
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func containsCategory(slice []models.Category, item models.Category) bool {
	for _, c := range slice {
		if c == item {
			return true
		}
	}
	return false
}
