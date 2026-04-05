package v1

import (
	"fmt"

	"github.com/ArchibaldKronin/microservices_test/payment/internal/model"
	payment_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/payment/v1"
	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *ApiSuite) TestPaySuccess() {
	var (
		userId        = gofakeit.UUID()
		paymentMethod = model.PaymentMethodCreditCard
		orderId       = gofakeit.UUID()
		paymentId     = gofakeit.UUID()
	)

	a.paymentService.EXPECT().PayOrder(a.ctx, paymentMethod).Return(paymentId, nil).Once()

	res, err := a.api.PayOrder(a.ctx, &payment_v1.PayOrderRequest{
		OrderUuid:     orderId,
		UserUuid:      userId,
		PaymentMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
	})

	a.NoError(err)
	a.NotNil(res)
	a.Equal(paymentId, res.TransactionUuid)
}

func (a *ApiSuite) TestPayErrGeneric() {
	var (
		userId        = gofakeit.UUID()
		paymentMethod = model.PaymentMethodCreditCard
		orderId       = gofakeit.UUID()
	)

	a.paymentService.EXPECT().PayOrder(a.ctx, paymentMethod).Return("", fmt.Errorf("generic error")).Once()

	res, err := a.api.PayOrder(a.ctx, &payment_v1.PayOrderRequest{
		OrderUuid:     orderId,
		UserUuid:      userId,
		PaymentMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
	})

	a.Nil(res)

	a.Error(err)
	st, ok := status.FromError(err)
	a.True(ok)
	a.Equal(codes.Internal, st.Code())
	a.Contains(st.Message(), orderId)
}
