package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
	ctx     context.Context
	service *PaymentService
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.service = NewPaymentService()
}

func TestServiceRun(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
