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

func (a *ApiSuite) TestListPartsSuccess() {
	var (
		uuid = gofakeit.UUID()

		filter = inventory_v1.PartsFilter{
			Uuids: []string{uuid},
		}

		filterMock = converter.FilterToDomain(&filter)

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
			UpdatedAt: time.Now(),
		}

		expParts = []*inventory_v1.Part{converter.PartToProto(&part)}
	)

	a.inventoryService.EXPECT().ListParts(a.ctx, filterMock).Return([]*model.Part{&part}, nil).Once()

	res, err := a.api.ListParts(a.ctx, &inventory_v1.ListPartsRequest{
		Filter: &filter,
	})

	a.NoError(err)
	a.NotNil(res)
	a.NotNil(res.Parts)

	a.Equal(expParts, res.Parts)
}

func (a *ApiSuite) TestListPartsErrInvalidArgument() {
	_, err := a.api.ListParts(a.ctx, &inventory_v1.ListPartsRequest{
		Filter: nil,
	})

	a.Error(err)
	st, ok := status.FromError(err)
	a.True(ok)

	a.Equal(codes.InvalidArgument, st.Code())
}

func (a *ApiSuite) TestListPartsErrInternal() {
	var (
		filter = inventory_v1.PartsFilter{
			Uuids: []string{"sfrs"},
		}

		filterMock = converter.FilterToDomain(&filter)
	)

	a.inventoryService.EXPECT().ListParts(a.ctx, filterMock).Return(([]*model.Part)(nil), fmt.Errorf("generic error")).Once()

	_, err := a.api.ListParts(a.ctx, &inventory_v1.ListPartsRequest{
		Filter: &filter,
	})

	a.Error(err)
	st, ok := status.FromError(err)
	a.True(ok)

	a.Equal(codes.Internal, st.Code())
}
