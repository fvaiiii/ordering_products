package service

import (
	"fmt"

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

func (s *ServiceSuite) TestPayOrder() {
	orderUuid := gofakeit.UUID()
	productUuid := gofakeit.UUID()
	userUuid := gofakeit.UUID()

	transactionUuid := gofakeit.UUID()
	paymentMethod := models.PaymentMethodCard

	originalOrder := &models.Order{
		OrderUuid:       orderUuid,
		UserUuid:        userUuid,
		ProductUuids:    []string{productUuid},
		TotalPrice:      100.50,
		TransactionUuid: nil,
		PaymentMethod:   nil,
		Status:          models.OrderStatusPendingPayment,
	}

	s.orderRepository.On("GetByUUID", s.ctx, orderUuid).
		Return(originalOrder, nil).
		Once()

	s.paymentRepository.On("PayOrder", s.ctx, orderUuid, userUuid, paymentMethod).
		Return(transactionUuid, nil).
		Once()

	s.orderRepository.On("Update", s.ctx, mock.AnythingOfType("*models.Order")).
		Return(nil).
		Once()

	resultTransactionUuid, err := s.orderService.PayOrder(s.ctx, orderUuid, models.PaymentMethod(paymentMethod))

	s.Require().NoError(err)
	s.Require().Equal(transactionUuid, resultTransactionUuid)
}

func (s *ServiceSuite) TestGetOrderByUUID() {
	orderUuid := gofakeit.UUID()
	userUuid := gofakeit.UUID()
	productUuid := gofakeit.UUID()
	transactionUuid := gofakeit.UUID()
	paymentMethod := string(models.PaymentMethodCard)

	expectedOrder := &models.Order{
		OrderUuid:       orderUuid,
		UserUuid:        userUuid,
		ProductUuids:    []string{productUuid},
		TotalPrice:      100.50,
		TransactionUuid: &transactionUuid,
		PaymentMethod:   &paymentMethod,
		Status:          models.OrderStatusPendingPayment,
	}
	s.orderRepository.On("GetByUUID", s.ctx, orderUuid).Return(expectedOrder, nil).Once()

	resultOrder, err := s.orderService.GetOrderByUUID(s.ctx, orderUuid)
	s.Require().NoError(err)
	s.Require().Equal(expectedOrder, resultOrder)
}
func (s *ServiceSuite) TestGetOrderByUUID_NotFound() {
	orderUuid := gofakeit.UUID()
	s.orderRepository.On("GetByUUID", s.ctx, orderUuid).Return(nil, fmt.Errorf("order not found")).Once()
	order, err := s.orderService.GetOrderByUUID(s.ctx, orderUuid)

	s.Require().Error(err)
	s.Require().Contains(err.Error(), "order not found")
	s.Require().Nil(order)

}

func (s *ServiceSuite) TestCancelOrder() {
	orderUuid := gofakeit.UUID()
	productUuid := gofakeit.UUID()
	userUuid := gofakeit.UUID()
	transactionUuid := gofakeit.UUID()
	paymentMethod := string(models.PaymentMethodCard)

	originalOrder := &models.Order{
		OrderUuid:       orderUuid,
		UserUuid:        userUuid,
		ProductUuids:    []string{productUuid},
		TotalPrice:      100.50,
		TransactionUuid: &transactionUuid,
		PaymentMethod:   &paymentMethod,
		Status:          models.OrderStatusPendingPayment,
	}

	s.orderRepository.On("GetByUUID", s.ctx, orderUuid).Return(originalOrder, nil)

	s.orderRepository.On("Update", s.ctx, originalOrder).Return(nil).Once()

	err := s.orderService.CancelOrder(s.ctx, orderUuid)

	s.Require().NoError(err)
}

func (s *ServiceSuite) TestCancelOrder_OrderNotFound() {
	orderUuid := gofakeit.UUID()
	s.orderRepository.On("GetByUUID", s.ctx, orderUuid).Return(nil, fmt.Errorf("order not found"))

	err := s.orderService.CancelOrder(s.ctx, orderUuid)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "order not found")

}

func (s *ServiceSuite) TestCancelOrder_AlreadyPaid() {
	orderUuid := gofakeit.UUID()
	userUuid := gofakeit.UUID()

	originalOrder := &models.Order{
		OrderUuid: orderUuid,
		UserUuid:  userUuid,
		Status:    models.OrderStatusPaid,
	}
	s.orderRepository.On("GetByUUID", s.ctx, orderUuid).Return(originalOrder, nil).Once()

	err := s.orderService.CancelOrder(s.ctx, orderUuid)

	s.Require().Error(err)
	s.Require().Contains(err.Error(), "cannot cancel paid order")
}

func (s *ServiceSuite) TestCancelOrder_AlreadyCancelled() {
	orderUuid := gofakeit.UUID()
	userUuid := gofakeit.UUID()

	originalOrder := &models.Order{
		OrderUuid: orderUuid,
		UserUuid:  userUuid,
		Status:    models.OrderStatusCancelled,
	}
	s.orderRepository.On("GetByUUID", s.ctx, orderUuid).Return(originalOrder, nil).Once()

	err := s.orderService.CancelOrder(s.ctx, orderUuid)

	s.Require().Error(err)
	s.Require().Contains(err.Error(), "invalid order status")
}

func (s *ServiceSuite) TestCancelOrder_UpdateFails() {
	orderUuid := gofakeit.UUID()
	userUuid := gofakeit.UUID()

	originalOrder := &models.Order{
		OrderUuid: orderUuid,
		UserUuid:  userUuid,
		Status:    models.OrderStatusPendingPayment,
	}
	s.orderRepository.On("GetByUUID", s.ctx, orderUuid).Return(originalOrder, nil).Once()

	s.orderRepository.On("Update", s.ctx, mock.AnythingOfType("*models.Order")).Return(fmt.Errorf("failed to save"))

	err := s.orderService.CancelOrder(s.ctx, orderUuid)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "failed to save")
}
