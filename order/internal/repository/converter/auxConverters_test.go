package converter

import (
	"testing"

	serviceModel "github.com/ArchibaldKronin/microservices_test/order/internal/model"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository/model"
	"github.com/stretchr/testify/require"
)

func TestPaymentMethodToDomain(t *testing.T) {
	t.Run("Card", func(t *testing.T) {
		pm := model.PaymentMethodCARD
		expected := serviceModel.PaymentMethodCARD

		res := PaymentMethodToDomain(pm)
		require.Equal(t, expected, res)
	})
	t.Run("Credit card", func(t *testing.T) {
		pm := model.PaymentMethodCREDITCARD
		expected := serviceModel.PaymentMethodCREDITCARD

		res := PaymentMethodToDomain(pm)
		require.Equal(t, expected, res)
	})
	t.Run("Investor money", func(t *testing.T) {
		pm := model.PaymentMethodINVESTORMONEY
		expected := serviceModel.PaymentMethodINVESTORMONEY

		res := PaymentMethodToDomain(pm)
		require.Equal(t, expected, res)
	})
	t.Run("SBP", func(t *testing.T) {
		pm := model.PaymentMethodSBP
		expected := serviceModel.PaymentMethodSBP

		res := PaymentMethodToDomain(pm)
		require.Equal(t, expected, res)
	})
	t.Run("Unknown", func(t *testing.T) {
		pm := model.PaymentMethodUNKNOWN
		expected := serviceModel.PaymentMethodUNKNOWN

		res := PaymentMethodToDomain(pm)
		require.Equal(t, expected, res)
	})
}

func TestPaymentMethodToRepo(t *testing.T) {
	t.Run("Card", func(t *testing.T) {
		pm := serviceModel.PaymentMethodCARD
		expected := model.PaymentMethodCARD

		res := PaymentMethodToRepo(pm)
		require.Equal(t, expected, res)
	})
	t.Run("Credit card", func(t *testing.T) {
		pm := serviceModel.PaymentMethodCREDITCARD
		expected := model.PaymentMethodCREDITCARD

		res := PaymentMethodToRepo(pm)
		require.Equal(t, expected, res)
	})
	t.Run("Investor money", func(t *testing.T) {
		pm := serviceModel.PaymentMethodINVESTORMONEY
		expected := model.PaymentMethodINVESTORMONEY

		res := PaymentMethodToRepo(pm)
		require.Equal(t, expected, res)
	})
	t.Run("SBP", func(t *testing.T) {
		pm := serviceModel.PaymentMethodSBP
		expected := model.PaymentMethodSBP

		res := PaymentMethodToRepo(pm)
		require.Equal(t, expected, res)
	})
	t.Run("Unknown", func(t *testing.T) {
		pm := serviceModel.PaymentMethodUNKNOWN
		expected := model.PaymentMethodUNKNOWN

		res := PaymentMethodToRepo(pm)
		require.Equal(t, expected, res)
	})
}

func TestStatusToDomain(t *testing.T) {
	t.Run("Pending", func(t *testing.T) {
		s := model.OrderStatusPENDINGPAYMENT
		expected := serviceModel.OrderStatusPENDINGPAYMENT

		res := StatusToDomain(s)
		require.Equal(t, expected, res)
	})
	t.Run("Paid", func(t *testing.T) {
		s := model.OrderStatusPAID
		expected := serviceModel.OrderStatusPAID

		res := StatusToDomain(s)
		require.Equal(t, expected, res)
	})
	t.Run("Cancelled", func(t *testing.T) {
		s := model.OrderStatusCANCELLED
		expected := serviceModel.OrderStatusCANCELLED

		res := StatusToDomain(s)
		require.Equal(t, expected, res)
	})
}

func TestStatusToRepo(t *testing.T) {
	t.Run("Pending", func(t *testing.T) {
		s := serviceModel.OrderStatusPENDINGPAYMENT
		expected := model.OrderStatusPENDINGPAYMENT

		res := StatusToRepo(s)
		require.Equal(t, expected, res)
	})
	t.Run("Paid", func(t *testing.T) {
		s := serviceModel.OrderStatusPAID
		expected := model.OrderStatusPAID

		res := StatusToRepo(s)
		require.Equal(t, expected, res)
	})
	t.Run("Cancelled", func(t *testing.T) {
		s := serviceModel.OrderStatusCANCELLED
		expected := model.OrderStatusCANCELLED

		res := StatusToRepo(s)
		require.Equal(t, expected, res)
	})
}
