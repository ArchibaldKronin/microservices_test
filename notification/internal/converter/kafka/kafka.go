package kafka

import (
	"github.com/ArchibaldKronin/microservices_test/notification/internal/model"
)

type OrderAssembledDecoder interface {
	Decode(data []byte) (model.OrderAssembledEvent, error)
}

type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaidEvent, error)
}
