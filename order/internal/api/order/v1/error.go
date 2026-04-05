package v1

import (
	"context"
	"net/http"

	order_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/openapi/order/v1"
)

func (a *api) NewError(_ context.Context, err error) *order_v1.GenericErrorStatusCode {
	return &order_v1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: order_v1.GenericError{
			Code:    order_v1.NewOptInt(http.StatusInternalServerError),
			Message: order_v1.NewOptString(err.Error()),
		},
	}
}
