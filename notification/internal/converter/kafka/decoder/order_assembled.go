package decoder

import (
	"github.com/ArchibaldKronin/microservices_test/notification/internal/model"
	events_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/events/v1"
	"google.golang.org/protobuf/proto"
)

type orderAssemledDecoder struct{}

func NewOrderAssembledDecoder() *orderAssemledDecoder {
	return &orderAssemledDecoder{}
}

func (*orderAssemledDecoder) Decode(data []byte) (model.OrderAssembledEvent, error) {
	var pb events_v1.ShipAssembled
	err := proto.Unmarshal(data, &pb)
	if err != nil {
		return model.OrderAssembledEvent{}, err
	}

	return model.OrderAssembledEvent{
		EventUuid:    pb.EventUuid,
		OrderUuid:    pb.OrderUuid,
		UserUuid:     pb.UserUuid,
		BuildTimeSec: pb.BuildTimeSec,
	}, nil
}
