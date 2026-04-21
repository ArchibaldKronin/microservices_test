package v1

import (
	"fmt"
	"time"

	"github.com/ArchibaldKronin/microservices_test/inventory/internal/converter"
	"github.com/ArchibaldKronin/microservices_test/inventory/internal/model"
	inventory_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/inventory/v1"
	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *ApiSuite) TestGetSuccess() {
	var (
		uuid = gofakeit.UUID()

		part = model.Part{
			Uuid:          uuid,
			Name:          "Turbo Engine X1",
			Description:   "High performance aircraft engine",
			Category:      model.CategoryEngine,
			Price:         125000.50,
			StockQuantity: 5,
			Dimensions: model.Dimensions{
				Length: 2.5,
				Width:  1.2,
				Height: 1.5,
				Weight: 850,
			},
			Manufacturer: model.Manufacturer{
				Name:    "JetCorp",
				Country: "USA",
				Website: "https://jetcorp.example.com",
			},
			Tags: []string{"engine", "turbo"},
			Metadata: map[string]model.Value{
				"horsepower":   model.Int64Value{Value: 4500},
				"fuel_type":    model.StringValue{Value: "Jet A-1"},
				"is_certified": model.BoolValue{Value: true},
			},
			CreatedAt: time.Now(),
			UpdatedAt: nil,
		}

		exp = converter.PartToProto(&part)
	)

	a.inventoryService.EXPECT().GetPart(a.ctx, uuid).Return(&part, nil).Once()

	res, err := a.api.GetPart(a.ctx, &inventory_v1.GetPartRequest{
		Uuid: uuid,
	})
	a.NoError(err)
	a.Equal(exp, res.Part)
}

func (a *ApiSuite) TestGetErrNotFound() {
	uuid := gofakeit.UUID()

	a.inventoryService.EXPECT().GetPart(a.ctx, uuid).Return(nil, model.ErrNotFound).Once()

	_, err := a.api.GetPart(a.ctx, &inventory_v1.GetPartRequest{
		Uuid: uuid,
	})

	a.Error(err)
	st, ok := status.FromError(err)
	a.True(ok)

	a.Equal(codes.NotFound, st.Code())
	a.Contains(st.Message(), uuid)
}

func (a *ApiSuite) TestGetErrInternal() {
	uuid := "sdgrtdsefade"

	a.inventoryService.EXPECT().GetPart(a.ctx, uuid).Return(nil, fmt.Errorf("generic error")).Once()

	_, err := a.api.GetPart(a.ctx, &inventory_v1.GetPartRequest{
		Uuid: uuid,
	})

	a.Error(err)
	st, ok := status.FromError(err)
	a.True(ok)

	a.Equal(codes.Internal, st.Code())
}
