package payment

import (
	"context"
	"testing"

	"github.com/ArchibaldKronin/microservices_test/payment/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPaySuccess(t *testing.T) {
	ctx := context.Background()
	svc := NewService()

	res, err := svc.PayOrder(ctx, model.PaymentMethodCard)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	_, errParse := uuid.Parse(res)

	require.NoError(t, errParse)
}
