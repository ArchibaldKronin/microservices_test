package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	payment_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// type PaymentMethod int

// const (
// 	PaymentMethodUnknown PaymentMethod = iota
// 	PaymentMethodCard
// 	PaymentMethodSBP
// 	PaymentMethodCreditCard
// 	PaymentMethodInvestorMoney
// )

const grpcPort = 50052

type PaymentService struct {
	payment_v1.UnimplementedPaymentServiceServer
}

func (PaymentService) PayOrder(_ context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error) {
	return &payment_v1.PayOrderResponse{
		TransactionUuid: uuid.NewString(),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc.UnaryServerInterceptor(LoggerUUID()),
		),
	)

	payment_v1.RegisterPaymentServiceServer(server, &PaymentService{})

	reflection.Register(server)

	go func() {
		log.Printf("🚀 gRPC Payment server listening on %d\n", grpcPort)
		err := server.Serve(lis)
		if err != nil {
			log.Printf("failed to serve Payment: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC Payment server...")
	server.GracefulStop()
	log.Println("✅ Payment server stopped")

}

func LoggerUUID() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return resp, err
		}

		if info.FullMethod == "/payment.v1.PaymentService/PayOrder" {
			if v, ok := resp.(*payment_v1.PayOrderResponse); ok {
				fmt.Println(v.TransactionUuid)
			}
		}
		return resp, err
	}
}
