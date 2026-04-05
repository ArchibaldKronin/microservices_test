package payment

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/payment/internal/model"
	"github.com/google/uuid"
)

func (*service) PayOrder(_ context.Context, pm model.PaymentMethod) (string, error) {
	return uuid.NewString(), nil
}
