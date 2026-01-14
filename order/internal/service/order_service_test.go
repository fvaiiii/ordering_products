package service

import (
	"github.com/brianvoe/gofakeit"
	"github.com/fvaiiii/ordering_products/order/internal/models"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuite) TestCreateOrder() {
	userUuid := gofakeit.UUID()
	productUuid := gofakeit.UUID()

	expectedProduct := &models.Product{
		UUID:  productUuid,
		Name:  gofakeit.Name(),
		Price: 100.50,
	}

	s.inventoryRepository.On("ListParts", s.ctx, []string{productUuid}).
		Return([]*models.Product{expectedProduct}, nil).
		Once()

	s.orderRepository.On("Create", s.ctx, mock.AnythingOfType("*models.Order")).
		Return(nil).
		Once()

	order, err := s.orderService.CreateOrder(s.ctx, userUuid, []string{productUuid})

	s.Require().NoError(err)
	s.Require().NotNil(order)

	s.Require().NotEmpty(order.OrderUuid)
	s.Require().Equal(userUuid, order.UserUuid)
	s.Require().Equal([]string{productUuid}, order.ProductUuids)
	s.Require().Equal(expectedProduct.Price, order.TotalPrice)
	s.Require().Equal(models.OrderStatusPendingPayment, order.Status)
}

// func (s *ServiceSuite) TestPayOrder() {

// }

// func (s *ServiceSuite) TestGetOrderByUUID() {

// }

// func (s *ServiceSuite) TestCancelOrder() {

// }
