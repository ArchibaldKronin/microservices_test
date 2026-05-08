package decoder

import (
	"fmt"

	"github.com/ArchibaldKronin/microservices_test/assembly/internal/converter/kafka"
	"github.com/ArchibaldKronin/microservices_test/assembly/internal/model"
	events_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/events/v1"
	"google.golang.org/protobuf/proto"
)

type decoder struct{}

func NewOrderpaidDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.OrderPaidEvent, error) {
	var pb events_v1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return kafka.OrderPaidEventToModel(&pb), nil
}
