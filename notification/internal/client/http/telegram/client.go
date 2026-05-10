package telegram

import (
	"context"
	"time"

	def "github.com/ArchibaldKronin/microservices_test/notification/internal/client/http"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"github.com/go-telegram/bot"
	"go.uber.org/zap"
)

var _ def.TelegramClient = (*telegramClient)(nil)

type telegramClient struct {
	bot *bot.Bot
}

func NewTelegramClient(bot *bot.Bot) *telegramClient {
	return &telegramClient{bot: bot}
}

func (c *telegramClient) SendMessage(ctx context.Context, chatID int64, text string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := c.bot.SendMessage(
		ctx,
		&bot.SendMessageParams{
			ChatID:    chatID,
			Text:      text,
			ParseMode: "Markdown",
		},
	)
	if err != nil {
		logger.Error(ctx, "Failed to send message", zap.Error(err))
		return err
	}

	return nil
}
