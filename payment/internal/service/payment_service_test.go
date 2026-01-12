package service

import (
	"github.com/brianvoe/gofakeit"
	v1 "github.com/fvaiiii/ordering_products/shared/pkg/proto/payment/v1"
)

func (s *ServiceSuite) TestPayOrder() {
	req := &v1.PaymentRequest{
		OrderUuid:     gofakeit.UUID(),
		UserUuid:      gofakeit.UUID(),
		PaymentMethod: v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
	}

	transactionUuid := gofakeit.UUID()
	expectedResp := &v1.PaymentResponse{
		TransactionUuid: transactionUuid,
	}
	resp, err := s.service.PayOrder(s.ctx, req)

	s.Require().NoError(err)
	s.Require().NotNil(expectedResp)
	s.Require().NotEmpty(resp.TransactionUuid)
}

func (s *ServiceSuite) TestPayErrorInvalidOrderUuid() {
	req := &v1.PaymentRequest{
		OrderUuid:     "",
		UserUuid:      gofakeit.UUID(),
		PaymentMethod: v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
	}
	expectedResp, err := s.service.PayOrder(s.ctx, req)
	s.Require().Error(err)
	s.Require().Nil(expectedResp)

}
func (s *ServiceSuite) TestPayErrorInvalidUserUuid() {
	req := &v1.PaymentRequest{
		OrderUuid:     gofakeit.UUID(),
		UserUuid:      "",
		PaymentMethod: v1.PaymentMethod_PAYMENT_METHOD_CARD,
	}
	expectedResp, err := s.service.PayOrder(s.ctx, req)
	s.Require().Error(err)
	s.Require().Nil(expectedResp)
}
