package telegram

import (
	"bytes"
	"context"
	"embed"
	"text/template"

	"github.com/ArchibaldKronin/microservices_test/notification/internal/client/http"
	"github.com/ArchibaldKronin/microservices_test/notification/internal/model"
	def "github.com/ArchibaldKronin/microservices_test/notification/internal/service"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"go.uber.org/zap"
)

//go:embed templates/paid_notification.tmpl
var orderPaidTemplateFS embed.FS

type orderPaidTemplateData struct {
	OrderUuid     string
	PaymentMethod model.PaymentMethod
}

var orderPaidTemplate = template.Must(template.ParseFS(orderPaidTemplateFS, "templates/paid_notification.tmpl"))

//go:embed templates/assembled_notification.tmpl
var orderAssembledTemplateFS embed.FS

type orderAssembledTemplateData struct {
	OrderUuid    string
	BuildTimeSec int64
}

var orderAssembledTemplate = template.Must(template.ParseFS(orderAssembledTemplateFS, "templates/assembled_notification.tmpl"))

const chatID int64 = 8212543139

var _ def.TelegramService = (*service)(nil)

type service struct {
	telegramClient http.TelegramClient
}

func NewService(client http.TelegramClient) *service {
	return &service{telegramClient: client}
}

func (s *service) SendOrderPaidNotification(ctx context.Context, orderPaid model.OrderPaidEvent) error {
	msg, err := s.buildOrderPaidMessage(orderPaid)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatID, msg)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message OrderPaid send to chat", zap.Int64("chatID", chatID), zap.String("message", msg))
	return nil
}

func (s *service) SendOrderAssembledNotification(ctx context.Context, orderAssembled model.OrderAssembledEvent) error {
	msg, err := s.buildOrderAssembledMessage(orderAssembled)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatID, msg)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message OrderAssembled send to chat", zap.Int64("chatID", chatID), zap.String("message", msg))
	return nil
}

func (*service) buildOrderPaidMessage(orderPaid model.OrderPaidEvent) (string, error) {
	data := orderPaidTemplateData{
		OrderUuid:     orderPaid.OrderUuid,
		PaymentMethod: orderPaid.PaymentMethod,
	}

	var buf bytes.Buffer
	err := orderPaidTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (*service) buildOrderAssembledMessage(orderAssembled model.OrderAssembledEvent) (string, error) {
	data := orderAssembledTemplateData{
		OrderUuid:    orderAssembled.OrderUuid,
		BuildTimeSec: orderAssembled.BuildTimeSec,
	}

	var buf bytes.Buffer
	err := orderAssembledTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
