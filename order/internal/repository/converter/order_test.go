package converter

import (
	"testing"

	serviceModel "github.com/ArchibaldKronin/microservices_test/order/internal/model"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository/model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestOrderToDomain(t *testing.T) {
	t.Run("Order payment nil", func(t *testing.T) {
		var (
			OrderID    = gofakeit.UUID()
			UserID     = gofakeit.UUID()
			PartIDs    = []string{gofakeit.UUID(), gofakeit.UUID()}
			TotalPrice = gofakeit.Float64()
		)
		order := model.Order{
			OrderID:       OrderID,
			UserID:        UserID,
			PartIDs:       PartIDs,
			TotalPrice:    TotalPrice,
			TransactionID: nil,
			PaymentMethod: nil,
			Status:        model.OrderStatusPENDINGPAYMENT,
		}
		expect := serviceModel.Order{
			OrderID:       OrderID,
			UserID:        UserID,
			PartIDs:       PartIDs,
			TotalPrice:    TotalPrice,
			TransactionID: nil,
			PaymentMethod: nil,
			Status:        serviceModel.OrderStatusPENDINGPAYMENT,
		}

		res := OrderToDomain(order)
		require.Equal(t, expect, res)
	})

	t.Run("Order whith payment", func(t *testing.T) {
		var (
			OrderID       = gofakeit.UUID()
			UserID        = gofakeit.UUID()
			PartIDs       = []string{gofakeit.UUID(), gofakeit.UUID()}
			TotalPrice    = gofakeit.Float64()
			TransactionID = lo.ToPtr(gofakeit.UUID())
		)
		order := model.Order{
			OrderID:       OrderID,
			UserID:        UserID,
			PartIDs:       PartIDs,
			TotalPrice:    TotalPrice,
			TransactionID: TransactionID,
			PaymentMethod: lo.ToPtr(model.PaymentMethodSBP),
			Status:        model.OrderStatusPAID,
		}
		expect := serviceModel.Order{
			OrderID:       OrderID,
			UserID:        UserID,
			PartIDs:       PartIDs,
			TotalPrice:    TotalPrice,
			TransactionID: TransactionID,
			PaymentMethod: lo.ToPtr(serviceModel.PaymentMethodSBP),
			Status:        serviceModel.OrderStatusPAID,
		}

		res := OrderToDomain(order)
		require.Equal(t, expect, res)
	})
}

func TestOrderToRepo(t *testing.T) {
	t.Run("Order payment nil", func(t *testing.T) {
		var (
			OrderID    = gofakeit.UUID()
			UserID     = gofakeit.UUID()
			PartIDs    = []string{gofakeit.UUID(), gofakeit.UUID()}
			TotalPrice = gofakeit.Float64()
		)
		order := serviceModel.Order{
			OrderID:       OrderID,
			UserID:        UserID,
			PartIDs:       PartIDs,
			TotalPrice:    TotalPrice,
			TransactionID: nil,
			PaymentMethod: nil,
			Status:        serviceModel.OrderStatusPENDINGPAYMENT,
		}
		expect := model.Order{
			OrderID:       OrderID,
			UserID:        UserID,
			PartIDs:       PartIDs,
			TotalPrice:    TotalPrice,
			TransactionID: nil,
			PaymentMethod: nil,
			Status:        model.OrderStatusPENDINGPAYMENT,
		}

		res := OrderToRepo(order)
		require.Equal(t, expect, res)
	})

	t.Run("Order whith payment", func(t *testing.T) {
		var (
			OrderID       = gofakeit.UUID()
			UserID        = gofakeit.UUID()
			PartIDs       = []string{gofakeit.UUID(), gofakeit.UUID()}
			TotalPrice    = gofakeit.Float64()
			TransactionID = lo.ToPtr(gofakeit.UUID())
		)
		order := serviceModel.Order{
			OrderID:       OrderID,
			UserID:        UserID,
			PartIDs:       PartIDs,
			TotalPrice:    TotalPrice,
			TransactionID: TransactionID,
			PaymentMethod: lo.ToPtr(serviceModel.PaymentMethodSBP),
			Status:        serviceModel.OrderStatusPAID,
		}
		expect := model.Order{
			OrderID:       OrderID,
			UserID:        UserID,
			PartIDs:       PartIDs,
			TotalPrice:    TotalPrice,
			TransactionID: TransactionID,
			PaymentMethod: lo.ToPtr(model.PaymentMethodSBP),
			Status:        model.OrderStatusPAID,
		}

		res := OrderToRepo(order)
		require.Equal(t, expect, res)
	})
}
