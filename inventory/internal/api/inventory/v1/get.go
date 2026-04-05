package v1

import (
	"context"
	"errors"

	"github.com/ArchibaldKronin/microservices_test/inventory/internal/converter"
	"github.com/ArchibaldKronin/microservices_test/inventory/internal/model"
	inventory_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) GetPart(ctx context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	id := req.GetUuid()
	part, err := a.inventoryService.GetPart(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrNotFound):
			return nil, status.Errorf(codes.NotFound, "запчасть с id: %s не найдена", id)
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &inventory_v1.GetPartResponse{
		Part: converter.PartToProto(part),
	}, nil
}
