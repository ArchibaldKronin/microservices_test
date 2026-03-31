package v1

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/inventory/internal/converter"
	inventory_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) ListParts(ctx context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	filter := req.GetFilter()
	if filter == nil {
		return nil, status.Error(codes.InvalidArgument, "handler error: запрос обязан содержать поле filter")
	}
	filterModel, err := converter.FilterToDomain(filter)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "handler error: запрос обязан содержать поле filter")
	}

	parts, err := a.inventoryService.ListParts(ctx, filterModel)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "handler error: %v", err)
	}

	var partsProto []*inventory_v1.Part

	for _, part := range parts {
		partsProto = append(partsProto, converter.PartToProto(part))
	}

	return &inventory_v1.ListPartsResponse{
		Parts: partsProto,
	}, nil
}
