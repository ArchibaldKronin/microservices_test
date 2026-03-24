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

	order_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/openapi/order/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type (
	OrderStatus   string
	PaymentMethod string
)

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

func mapPaymentMethodToOrderDTO(pm PaymentMethod) order_v1.PaymentMethod {
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

func NewOrder(userId string, partIds []string, totalPrice float64) *Order {
	return &Order{
		OrderId:     uuid.NewString(),
		UserId:      userId,
		PartIds:     partIds,
		Total_price: totalPrice,
		Status:      OrderStatusPENDINGPAYMENT,
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
			Value: mapPaymentMethodToOrderDTO(*o.PaymentMethod),
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

type OrderService struct {
	storage         *OrderStorage
	inventoryClient InventoryClient
	paymentClient   PaymentClient
}

func NewOrderService(storage *OrderStorage, incCl InventoryClient, paymentCl PaymentClient) *OrderService {
	return &OrderService{
		storage:         storage,
		inventoryClient: incCl,
		paymentClient:   paymentCl,
	}
}

func (serv *OrderService) getPartsInfo(ctx context.Context, part_ids []string) ([]string, float64, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	parts, err := serv.inventoryClient.ListPart(ctx, part_ids)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return nil, 0, fmt.Errorf("неудачный запрос, неизвестная ошибка: %w", err)
		}

		switch st.Code() {
		case codes.NotFound:
			return nil, 0, ErrNotFound
		case codes.Unavailable:
			return nil, 0, ErrUnavailable
		default:
			return nil, 0, err
		}
	}

	if len(part_ids) != 0 {
		if len(parts) != len(part_ids) {
			return nil, 0, ErrNotFound
		}
	}

	ids := make([]string, 0, len(parts))
	total_price := 0.0
	for _, part := range parts {
		ids = append(ids, part.Uuid)
		total_price += part.Price
	}
	return ids, total_price, nil
}

func (serv *OrderService) CreateOrder(ctx context.Context, userId string, partIds []string) (*Order, error) {
	ids, totalPrice, err := serv.getPartsInfo(ctx, partIds)
	if err != nil {
		return nil, err
	}

	order := NewOrder(userId, ids, totalPrice)
	serv.storage.UpdateOrder(order)

	return order, nil
}

func (serv *OrderService) GetOrder(id string) (*Order, error) {
	order := serv.storage.GetOrder(id)
	if order == nil {
		return nil, ErrNotFound
	}
	return order, nil
}

func (serv *OrderService) CancelOrder(id string) (OrderStatus, error) {
	order, err := serv.GetOrder(id)
	if err != nil {
		return "", err
	}

	switch order.Status {
	case OrderStatusPENDINGPAYMENT:
		order.Status = OrderStatusCANCELLED
		serv.storage.UpdateOrder(order)
		return OrderStatusCANCELLED, nil
	case OrderStatusPAID:
		return "", ErrOrderAlreadyPaid
		// already cancelled
	default:
		return OrderStatusCANCELLED, nil
	}
}

func (serv *OrderService) PayOrder(ctx context.Context, orderId string, pm PaymentMethod) (string, error) {
	order, err := serv.GetOrder(orderId)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	userId := order.UserId
	payOrderResp, err := serv.paymentClient.PayOrder(ctx, userId, orderId, pm)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return "", fmt.Errorf("неудачный запрос, неизвестная ошибка: %w", err)
		}
		switch st.Code() {
		case codes.NotFound:
			return "", ErrNotFound
		case codes.Unavailable:
			return "", ErrUnavailable
		default:
			return "", err
		}
	}

	order.Status = OrderStatusPAID
	order.PaymentMethod = &pm
	order.TransactionID = &payOrderResp.TransactionUuid

	serv.storage.UpdateOrder(order)

	return payOrderResp.TransactionUuid, nil
}

type OrderHandler struct {
	service *OrderService
}

// Реализует интерфейс order_v1.Handler
func NewOrderHandler(service *OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) GetOrderByUUID(_ context.Context, params order_v1.GetOrderByUUIDParams) (order_v1.GetOrderByUUIDRes, error) {
	id := params.OrderUUID.String()
	order, err := h.service.GetOrder(id)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return &order_v1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order for id %s not found", id),
			}, nil
		default:
			log.Printf("ошибка получения данных: %v\n", err)
			return nil, err
		}
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

func (h *OrderHandler) CreateOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderRes, error) {
	reqIds := convertUUIDs(req.PartUuids)
	order, err := h.service.CreateOrder(ctx, req.UserUUID.String(), reqIds)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return &order_v1.NotFoundError{
				Code:    404,
				Message: "some parts not found",
			}, nil
		default:
			log.Printf("ошибка получения данных о запчастях: %v\n", err)
			return nil, err
		}
	}

	orderId, err := uuid.Parse(order.OrderId)
	if err != nil {
		log.Printf("ошибка парсинга uuid: %v\n", err)
		return nil, err
	}
	resp := order_v1.CreateOrderResponse{
		OrderUUID:  orderId,
		TotalPrice: float32(order.Total_price),
	}

	return &resp, nil
}

func (h *OrderHandler) PayOrderByUUID(ctx context.Context, req *order_v1.PayOrderRequest, params order_v1.PayOrderByUUIDParams) (order_v1.PayOrderByUUIDRes, error) {
	orderId := params.OrderUUID.String()
	paymentMethod := mapPaymentMethod(req.PaymentMethod)

	transactionId, err := h.service.PayOrder(ctx, orderId, paymentMethod)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return &order_v1.NotFoundError{
				Code:    404,
				Message: "order not found",
			}, nil
		case errors.Is(err, ErrUnavailable):
			return &order_v1.InternalServerError{
				Code:    503,
				Message: "unavailable",
			}, nil
		default:
			log.Printf("ошибка получения данных о заказе: %v\n", err)
			return nil, err
		}
	}

	txUUID, err := uuid.Parse(transactionId)
	if err != nil {
		log.Printf("ошибка парсинга uuid: %v\n", err)
		return nil, fmt.Errorf("ошибка парсинга uuid: %w", err)
	}

	return &order_v1.PayOrderResponse{
		TransactionUUID: txUUID,
	}, nil
}

func (h *OrderHandler) CancelOrderByUUID(_ context.Context, params order_v1.CancelOrderByUUIDParams) (order_v1.CancelOrderByUUIDRes, error) {
	id := params.OrderUUID.String()
	orderStatus, err := h.service.CancelOrder(id)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return &order_v1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order for id %s not found", id),
			}, nil
		case errors.Is(err, ErrOrderAlreadyPaid):
			return &order_v1.ConflictError{
				Code:    409,
				Message: fmt.Sprintf("Order %s is already paid", id),
			}, nil
		default:
			log.Printf("ошибка получения данных: %v\n", err)
			return nil, err
		}
	}

	switch orderStatus {
	case OrderStatusCANCELLED:
		return &order_v1.CancelOrderByUUIDNoContent{}, nil
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

/////////////////////////////////////////////////////////////////////// CLIENTS

const (
	serverInventoryAddress = "localhost:50051"
	serverPaymentAddress   = "localhost:50052"
)

func main() {
	/////////////////////////////////////////////////////////////////////// CLIENTS
	connInventory, err := grpc.NewClient(
		serverInventoryAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("ошибка соединения с Inventory: %v", err)
	}
	defer func() {
		if cerr := connInventory.Close(); cerr != nil {
			log.Printf("failed to close connect to Inventory: %v", cerr)
		}
	}()

	inventoryCl := NewInventoryClient(connInventory)

	connPayment, err := grpc.NewClient(
		serverPaymentAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("ошибка соединения с Payment: %v", err)
	}
	defer func() {
		if cerr := connPayment.Close(); cerr != nil {
			log.Printf("failed to close connect to Payment: %v", cerr)
		}
	}()

	paymentCl := NewPaymentClient(connPayment)
	/////////////////////////////////////////////////////////////////////
	/////////////////////////////////////////////////////////////////////// SERVER

	storage := NewOrderStorage()

	service := NewOrderService(storage, inventoryCl, paymentCl)

	orderHandler := NewOrderHandler(service)

	orderServer, err := order_v1.NewServer(orderHandler)
	if err != nil {
		log.Printf("ошибка создания Order сервера OpenAPI: %v", err)
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
