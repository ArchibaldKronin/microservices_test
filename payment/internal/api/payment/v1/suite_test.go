package v1

import (
	"context"
	"testing"

	"github.com/ArchibaldKronin/microservices_test/payment/internal/service/mocks"
	"github.com/stretchr/testify/suite"
)

type ApiSuite struct {
	suite.Suite

	ctx context.Context //nolint: containedctx

	paymentService mocks.PaymentService

	api *api
}

func (a *ApiSuite) SetupTest() {
	a.ctx = context.Background()

	a.paymentService = *mocks.NewPaymentService(a.T())

	a.api = NewAPI(&a.paymentService)
}

func TestApiSuite(t *testing.T) {
	suite.Run(t, new(ApiSuite))
}
