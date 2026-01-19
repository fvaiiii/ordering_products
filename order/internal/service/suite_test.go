package service

import (
	"context"
	"testing"

	mocksClient "github.com/fvaiiii/ordering_products/order/internal/clients/mocks"
	mocksOrder "github.com/fvaiiii/ordering_products/order/internal/service/mocks"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
	ctx                 context.Context
	orderRepository     *mocksOrder.Order
	inventoryRepository *mocksClient.InventoryClient
	paymentRepository   *mocksClient.PaymentClient
	orderService        OrderService
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.orderRepository = mocksOrder.NewOrder(s.T())
	s.inventoryRepository = mocksClient.NewInventoryClient(s.T())
	s.paymentRepository = mocksClient.NewPaymentClient(s.T())
	s.orderService = NewOrderService(s.orderRepository, s.inventoryRepository, s.paymentRepository)
}

func TestServiceRun(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
