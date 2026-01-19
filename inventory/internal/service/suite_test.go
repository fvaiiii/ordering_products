package service

import (
	"context"
	"testing"

	"github.com/fvaiiii/ordering_products/inventory/mocks"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
	ctx                 context.Context
	inventoryRepository *mocks.Products
	service             *InventoryService
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.inventoryRepository = mocks.NewProducts(s.T())
	s.service = NewInventoryService(s.inventoryRepository)
}

func TestServiceRun(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
