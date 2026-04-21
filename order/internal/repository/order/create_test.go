package order

// import (
// 	"context"
// 	"testing"

// 	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
// 	"github.com/brianvoe/gofakeit/v7"
// 	"github.com/stretchr/testify/require"
// )

// func TestCreate(t *testing.T) {
// 	var (
// 		ctxCreate = context.Background()
// 		ctxGet    = context.Background()

// 		userId     = gofakeit.UUID()
// 		partsId    = []string{gofakeit.UUID()}
// 		totalPrice = 42.42
// 		order      = model.NewOrder(userId, partsId, totalPrice)

// 		repo = NewRepository()
// 	)
// 	repo.CreateOrder(ctxCreate, order)

// 	res := repo.GetOrder(ctxGet, order.OrderId)

// 	require.NotNil(t, res)
// 	require.Equal(t, order, res)
// }
