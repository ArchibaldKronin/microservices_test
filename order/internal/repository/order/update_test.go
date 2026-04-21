package order

// import (
// 	"context"
// 	"testing"

// 	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
// 	"github.com/brianvoe/gofakeit/v7"
// 	"github.com/stretchr/testify/require"
// )

// func TestUpdateSuccess(t *testing.T) {
// 	var (
// 		ctxCreate = context.Background()
// 		ctxUpdate = context.Background()
// 		ctxGet    = context.Background()

// 		userId     = gofakeit.UUID()
// 		partsId    = []string{gofakeit.UUID()}
// 		totalPrice = 42.42
// 		order      = model.NewOrder(userId, partsId, totalPrice)

// 		repo = NewRepository()

// 		expected = &model.Order{
// 			OrderId:     order.OrderId,
// 			PartIds:     partsId,
// 			Total_price: 90.0,
// 			UserId:      userId,
// 			Status:      model.OrderStatusPENDINGPAYMENT,
// 		}
// 	)
// 	repo.CreateOrder(ctxCreate, order)

// 	repo.UpdateOrder(ctxUpdate, expected)

// 	res := repo.GetOrder(ctxGet, order.OrderId)

// 	require.NotNil(t, res)
// 	require.Equal(t, expected, res)
// }

// func TestUpdateNil(t *testing.T) {
// 	var (
// 		ctxUpdate = context.Background()

// 		userId     = gofakeit.UUID()
// 		partsId    = []string{gofakeit.UUID()}
// 		totalPrice = 42.42
// 		order      = model.NewOrder(userId, partsId, totalPrice)

// 		repo = NewRepository()
// 	)

// 	res := repo.UpdateOrder(ctxUpdate, order)

// 	require.Nil(t, res)
// }
