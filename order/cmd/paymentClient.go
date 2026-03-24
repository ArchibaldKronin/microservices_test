package main

import (
	"context"

	payment_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
)

func mapPaymentMethodToPaymentDTO(pm PaymentMethod) payment_v1.PaymentMethod {
	switch pm {
	case PaymentMethodCARD:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CARD
	case PaymentMethodCREDITCARD:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case PaymentMethodINVESTORMONEY:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	case PaymentMethodSBP:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_SBP
	default:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

type PaymentClient interface {
	PayOrder(ctx context.Context, userId string, orderId string, paymentMethod PaymentMethod) (*payment_v1.PayOrderResponse, error)
}

type paymentClient struct {
	client payment_v1.PaymentServiceClient
}

func NewPaymentClient(conn *grpc.ClientConn) PaymentClient {
	client := payment_v1.NewPaymentServiceClient(conn)

	return &paymentClient{
		client: client,
	}
}

func (c *paymentClient) PayOrder(ctx context.Context, userId string, orderId string, paymentMethod PaymentMethod) (*payment_v1.PayOrderResponse, error) {
	resp, err := c.client.PayOrder(ctx,
		&payment_v1.PayOrderRequest{
			UserUuid:      userId,
			OrderUuid:     orderId,
			PaymentMethod: mapPaymentMethodToPaymentDTO(paymentMethod),
		},
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
