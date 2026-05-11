package decoder

import (
	"github.com/ArchibaldKronin/microservices_test/notification/internal/converter"
	"github.com/ArchibaldKronin/microservices_test/notification/internal/model"
	events_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/events/v1"
	"google.golang.org/protobuf/proto"
)

type orderPaidDecoder struct{}

func NewOrderPaidDecoder() *orderPaidDecoder {
	return &orderPaidDecoder{}
}

func (*orderPaidDecoder) Decode(data []byte) (model.OrderPaidEvent, error) {
	var pb events_v1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidEvent{}, err
	}

	return model.OrderPaidEvent{
		EventUuid:       pb.EventUuid,
		OrderUuid:       pb.OrderUuid,
		UserUuid:        pb.UserUuid,
		PaymentMethod:   converter.PaymentMethodToModel(pb.PaymentMethod),
		TransactionUuid: pb.TransactionUuid,
	}, nil
}
