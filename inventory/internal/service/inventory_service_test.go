package service

import (
	"github.com/brianvoe/gofakeit"
	"github.com/fvaiiii/ordering_products/inventory/internal/models"
	"github.com/fvaiiii/ordering_products/inventory/internal/port/repo"
)

func (s *ServiceSuite) TestGetProduct() {
	productUuid := gofakeit.UUID()

	expectedManyfacturer := &models.Manufacturer{
		Name:    gofakeit.Name(),
		Country: gofakeit.Company(),
		Website: gofakeit.URL(),
	}

	expectedProduct := &models.Product{
		Uuid:          productUuid,
		Name:          gofakeit.Name(),
		Description:   gofakeit.Sentence(10),
		Price:         gofakeit.Price(100, 1000),
		StockQuantity: gofakeit.Int64(),
		Category:      models.Category(gofakeit.Int32()%4 + 1),
		Manufacturer:  *expectedManyfacturer,
		CreatedAt:     gofakeit.Date(),
		UpdatedAt:     gofakeit.Date(),
	}

	s.inventoryRepository.On("GetProduct", s.ctx, productUuid).Return(expectedProduct, nil).Once()

	product, err := s.service.GetProduct(s.ctx, productUuid)
	s.Require().NoError(err)
	s.Require().Equal(expectedProduct, product)

}

func (s *ServiceSuite) TestGetProductNotFound() {
	productUuid := gofakeit.UUID()
	s.inventoryRepository.On("GetProduct", s.ctx, productUuid).Return((*models.Product)(nil), repo.ErrNotFound).Once()
	product, err := s.service.GetProduct(s.ctx, productUuid)
	s.Require().Error(err, repo.ErrNotFound)
	s.Require().Nil(product)
}

func (s *ServiceSuite) TestListProducts() {
	productUuid := gofakeit.UUID()
	filter := models.ProductsFilter{
		Uuids: []string{productUuid},
	}
	expectedManyfacturer := &models.Manufacturer{
		Name:    gofakeit.Name(),
		Country: gofakeit.Country(),
		Website: gofakeit.URL(),
	}
	expectedProducts := []*models.Product{
		{
			Uuid:          productUuid,
			Name:          gofakeit.Name(),
			Description:   gofakeit.Sentence(10),
			Price:         gofakeit.Price(100, 1000),
			StockQuantity: gofakeit.Int64(),
			Category:      models.Category(gofakeit.Int32()%4 + 1),
			Manufacturer:  *expectedManyfacturer,
			CreatedAt:     gofakeit.Date(),
			UpdatedAt:     gofakeit.Date(),
		},
	}
	s.inventoryRepository.On("ListProducts", s.ctx, filter).Return(expectedProducts, nil).Once()
	products, err := s.service.ListProducts(s.ctx, filter)
	s.Require().NoError(err)
	s.Require().Equal(expectedProducts, products)

}
