package grpc

import (
	"time"

	"github.com/fvaiiii/ordering_products/inventory/internal/models"
	v1 "github.com/fvaiiii/ordering_products/shared/pkg/proto/inventory/v1"
)

func domainArrToProto(products []*models.Product) []*v1.Product {
	if products == nil {
		return nil
	}
	arr := make([]*v1.Product, 0)
	for _, product := range products {
		arr = append(arr, domainToProto(product))
	}

	return arr
}

func filterFromProto(filter *v1.ProductsFilter) models.ProductsFilter {
	return models.ProductsFilter{
		Uuids:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            categoryFromProto(filter.Categories),
		ManufacturerCountries: filter.ManufacturerCountries,
	}
}

func categoryFromProto(categoriesFilter []v1.Category) []models.Category {
	arr := make([]models.Category, 0)
	for _, categoryFilter := range categoriesFilter {
		switch categoryFilter {
		case v1.Category_CATEGORY_VEGETABLES:
			arr = append(arr, models.VEGETABLES)
		case v1.Category_CATEGORY_FRUITS:
			arr = append(arr, models.FRUITS)
		case v1.Category_CATEGORY_MEATS:
			arr = append(arr, models.MEATS)
		default:
			arr = append(arr, models.UNKNOWN)
		}
	}
	return arr
}

func domainToProto(product *models.Product) *v1.Product {
	if product == nil {
		return nil
	}
	return &v1.Product{
		Uuid:          product.Uuid,
		Name:          product.Name,
		Description:   product.Description,
		Price:         product.Price,
		StockQuantity: product.StockQuantity,
		Category:      categoryToProto(product.Category),
		Manufacturer: &v1.Manufacturer{
			Name:    product.Manufacturer.Name,
			Country: product.Manufacturer.Country,
			Website: product.Manufacturer.Website,
		},
		CreatedAt: formatTime(product.CreatedAt),
		UpdatedAt: formatTime(product.UpdatedAt),
	}
}

func categoryToProto(category models.Category) v1.Category {
	switch category {
	case models.VEGETABLES:
		return v1.Category_CATEGORY_VEGETABLES
	case models.FRUITS:
		return v1.Category_CATEGORY_FRUITS
	case models.MEATS:
		return v1.Category_CATEGORY_MEATS
	default:
		return v1.Category_CATEGORY_UNKNOWN
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}
