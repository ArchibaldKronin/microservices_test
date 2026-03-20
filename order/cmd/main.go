package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	order_v1 "github.com/ArchibaldKronin/microservices_test/platform/pkg/openapi/order/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type OrderStatus string
type PaymentMethod string

const (
	httpPort          = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second

	OrderStatusPENDINGPAYMENT OrderStatus = "PENDING_PAYMENT"
	OrderStatusPAID           OrderStatus = "PAID"
	OrderStatusCANCELLED      OrderStatus = "CANCELLED"

	PaymentMethodUNKNOWN       PaymentMethod = "UNKNOWN"
	PaymentMethodCARD          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCREDITCARD    PaymentMethod = "CREDIT_CARD"
	PaymentMethodINVESTORMONEY PaymentMethod = "INVESTOR_MONEY"
)

// func mapOrderStatus(os order_v1.OrderStatus) OrderStatus {
// 	switch os {
// 	case order_v1.OrderStatusCANCELLED:
// 		return OrderStatusCANCELLED
// 	case order_v1.OrderStatusPAID:
// 		return OrderStatusPAID
// 	default:
// 		return OrderStatusPENDINGPAYMENT
// 	}
// }

func mapOrderStatusToDTO(os OrderStatus) order_v1.OrderStatus {
	switch os {
	case OrderStatusCANCELLED:
		return order_v1.OrderStatusCANCELLED
	case OrderStatusPAID:
		return order_v1.OrderStatusPAID
	default:
		return order_v1.OrderStatusPENDINGPAYMENT
	}
}

func mapPaymentMethod(pm order_v1.PaymentMethod) PaymentMethod {
	switch pm {
	case order_v1.PaymentMethodCARD:
		return PaymentMethodCARD
	case order_v1.PaymentMethodCREDITCARD:
		return PaymentMethodCREDITCARD
	case order_v1.PaymentMethodINVESTORMONEY:
		return PaymentMethodINVESTORMONEY
	case order_v1.PaymentMethodSBP:
		return PaymentMethodSBP
	default:
		return PaymentMethodUNKNOWN
	}
}
func mapPaymentMethodToDTO(pm PaymentMethod) order_v1.PaymentMethod {
	switch pm {
	case PaymentMethodCARD:
		return order_v1.PaymentMethodCARD
	case PaymentMethodCREDITCARD:
		return order_v1.PaymentMethodCREDITCARD
	case PaymentMethodINVESTORMONEY:
		return order_v1.PaymentMethodINVESTORMONEY
	case PaymentMethodSBP:
		return order_v1.PaymentMethodSBP
	default:
		return order_v1.PaymentMethodUNKNOWN
	}
}

// заказ
type Order struct {
	OrderId       string
	UserId        string
	PartIds       []string
	Total_price   float64
	TransactionID *string
	PaymentMethod *PaymentMethod
	Status        OrderStatus
}

func NewOrder(userId string, partIds []string) *Order {
	return &Order{
		OrderId: uuid.NewString(),
		UserId:  userId,
		PartIds: partIds,
		////////////////////////////
		Total_price: 100,
		////////////////////////////
		Status: OrderStatusPENDINGPAYMENT,
	}
}

// Маппинг из Order в DTO
func ToDTO(o *Order) (order_v1.OrderDto, error) {
	orderIdDTO, err := uuid.Parse(o.OrderId)
	if err != nil {
		return order_v1.OrderDto{}, err
	}
	userIdDTO, err := uuid.Parse(o.UserId)
	if err != nil {
		return order_v1.OrderDto{}, err
	}
	partsIdsDTO, err := convertStringsToUUIDs(o.PartIds)
	if err != nil {
		return order_v1.OrderDto{}, fmt.Errorf("error converting strings to UUIDs: %w", err)
	}

	var optTransactionId order_v1.OptUUID
	if o.TransactionID != nil {
		txId, err := uuid.Parse(*o.TransactionID)
		if err != nil {
			return order_v1.OrderDto{}, err
		}
		optTransactionId = order_v1.OptUUID{
			Value: txId,
			Set:   true,
		}
	}

	var optPaymentMethod order_v1.OptPaymentMethod
	if o.PaymentMethod != nil {
		optPaymentMethod = order_v1.OptPaymentMethod{
			Value: mapPaymentMethodToDTO(*o.PaymentMethod),
			Set:   true,
		}
	}

	return order_v1.OrderDto{
		OrderUUID:       orderIdDTO,
		UserUUID:        userIdDTO,
		PartUuids:       partsIdsDTO,
		TotalPrice:      float32(o.Total_price),
		TransactionUUID: optTransactionId,
		PaymentMethod:   optPaymentMethod,
		Status:          mapOrderStatusToDTO(o.Status),
	}, nil
}

// Маппинг из DTO в Order
// func fromDTO(d order_v1.OrderDto) *Order {
// 	order := &Order{
// 		OrderId:     d.OrderUUID.String(),
// 		UserId:      d.UserUUID.String(),
// 		PartIds:     convertUUIDs(d.PartUuids),
// 		Total_price: float64(d.TotalPrice),
// 	}

// 	if d.TransactionUUID.Set {
// 		temp := d.TransactionUUID.Value.String()
// 		order.TransactionID = &temp
// 	}

// 	if d.PaymentMethod.Set {
// 		temp := mapPaymentMethod(d.PaymentMethod.Value)
// 		order.PaymentMethod = &temp
// 	}

// 	order.Status = mapOrderStatus(d.Status)

// 	return order
// }

// хранилище заказов
type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*Order),
	}
}

func (s *OrderStorage) GetOrder(id string) *Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[id]
	if !ok {
		return nil
	}
	return order
}

func (s *OrderStorage) UpdateOrder(o *Order) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[o.OrderId] = o
}

type OrderHandler struct {
	storage *OrderStorage
}

// Реализует интерфейс order_v1.Handler
func NewOrderHandler(storage *OrderStorage) *OrderHandler {
	return &OrderHandler{
		storage: storage,
	}
}

func (h *OrderHandler) GetOrderByUUID(_ context.Context, params order_v1.GetOrderByUUIDParams) (order_v1.GetOrderByUUIDRes, error) {
	id := params.OrderUUID.String()
	order := h.storage.GetOrder(id)

	if order == nil {
		return &order_v1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("Order for id %s not found", id),
		}, nil
	}
	orderDTO, err := ToDTO(order)
	if err != nil {
		return &order_v1.InternalServerError{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	return &orderDTO, nil
}

func (h *OrderHandler) CreateOrder(_ context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderRes, error) {
	order := NewOrder(req.UserUUID.String(), convertUUIDs(req.PartUuids))

	h.storage.UpdateOrder(order)

	orderId, _ := uuid.Parse(order.OrderId)
	resp := order_v1.CreateOrderResponse{
		OrderUUID:  orderId,
		TotalPrice: float32(order.Total_price),
	}

	return &resp, nil
}

func (h *OrderHandler) PayOrderByUUID(_ context.Context, req *order_v1.PayOrderRequest, params order_v1.PayOrderByUUIDParams) (order_v1.PayOrderByUUIDRes, error) {
	id := params.OrderUUID.String()
	order := h.storage.GetOrder(id)

	if order == nil {
		return &order_v1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("Order for id %s not found", id),
		}, nil
	}

	paymentMethod := mapPaymentMethod(req.PaymentMethod)
	order.PaymentMethod = &paymentMethod

	//////////////////////////////////////
	trxId := uuid.New()
	trxIdStr := trxId.String()
	//////////////////////////////////////
	order.TransactionID = &trxIdStr

	order.Status = OrderStatusPAID

	h.storage.UpdateOrder(order)

	return &order_v1.PayOrderResponse{
		TransactionUUID: trxId,
	}, nil
}

func (h *OrderHandler) CancelOrderByUUID(_ context.Context, params order_v1.CancelOrderByUUIDParams) (order_v1.CancelOrderByUUIDRes, error) {
	id := params.OrderUUID.String()
	order := h.storage.GetOrder(id)

	if order == nil {
		return &order_v1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("Order for id %s not found", id),
		}, nil
	}

	switch order.Status {
	case OrderStatusPENDINGPAYMENT:
		order.Status = OrderStatusCANCELLED
		h.storage.UpdateOrder(order)
		return &order_v1.CancelOrderByUUIDNoContent{}, nil
	case OrderStatusPAID:
		return &order_v1.ConflictError{
			Code:    409,
			Message: fmt.Sprintf("Order %s is already paid", id),
		}, nil
		// already cancelled
	default:
		return &order_v1.CancelOrderByUUIDNoContent{}, nil

	}
}

func (h *OrderHandler) NewError(_ context.Context, err error) *order_v1.GenericErrorStatusCode {
	return &order_v1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: order_v1.GenericError{
			Code:    order_v1.NewOptInt(http.StatusInternalServerError),
			Message: order_v1.NewOptString(err.Error()),
		},
	}
}

func convertUUIDs(in []uuid.UUID) []string {
	ids := make([]string, 0, len(in))
	for _, u := range in {
		ids = append(ids, u.String())
	}
	return ids
}

func convertStringsToUUIDs(in []string) ([]uuid.UUID, error) {
	uuids := make([]uuid.UUID, 0, len(in))
	for _, id := range in {
		if u, err := uuid.Parse(id); err == nil {
			uuids = append(uuids, u)
		} else {
			return nil, err
		}
	}
	return uuids, nil
}

func main() {
	storage := NewOrderStorage()

	orderHandler := NewOrderHandler(storage)

	orderServer, err := order_v1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("ошибка создания Order сервера OpenAPI: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		///////////////////////////
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		///////////////////////////
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
