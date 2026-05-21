package app

import (
	"context"
	"fmt"
	"log"

	"github.com/ArchibaldKronin/microservices_test/notification/internal/config"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/closer"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

type App struct {
	diContainer *diContainer
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	// errGroup, errCtx := errgroup.WithContext(ctx)

	// errGroup.Go(func() error {
	// 	err := a.runOrderPaidConsumer(errCtx)
	// 	if err != nil {
	// 		logger.Error(errCtx, "OrderPaid consumer crashed", zap.Error(err))
	// 		return fmt.Errorf("OrderPaid consumer crashed: %w", err)
	// 	}

	// 	logger.Info(ctx, "OrderPaid consumer finished work healthy")
	// 	return nil
	// })

	// errGroup.Go(func() error {
	// 	err := a.runOrderAssembledConsumer(errCtx)
	// 	if err != nil {
	// 		logger.Error(errCtx, "OrderAssembled consumer crashed", zap.Error(err))
	// 		return fmt.Errorf("OrderAssembled consumer crashed: %w", err)
	// 	}

	// 	logger.Info(ctx, "OrderAssembled consumer finished work healthy")
	// 	return nil
	// })

	// err := errGroup.Wait()
	// if err != nil {
	// 	logger.Error(ctx, "Application stopped", zap.Error(err))
	// 	return err
	// }

	// return nil

	errCh := make(chan (error), 2)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		err := a.runOrderPaidConsumer(ctx)
		if err != nil {
			errCh <- fmt.Errorf("OrderPaid consumer crashed: %w", err)
		}
	}()

	go func() {
		err := a.runOrderAssembledConsumer(ctx)
		if err != nil {
			errCh <- fmt.Errorf("OrderAssembled consumer crashed: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "Shutdown signal received")
	case err := <-errCh:
		logger.Error(ctx, "Component crushed, shutting down", zap.Error(err))
		cancel()
		<-ctx.Done()
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initLogger,
		a.initDi,
		a.initCloser,
		a.initTelegramBot,
	}

	for _, f := range inits {
		f := f
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	err := logger.Init(ctx,
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger,
	)
	if err != nil {
		log.Printf("❌ logger init failed: %v; using noop logger", err)
		logger.SetNopLogger()
		return err
	}

	closer.AddNamed("Logger", func(ctx context.Context) error {
		_ = logger.Sync()     //nolint:gosec
		_ = logger.Close(ctx) //nolint:gosec
		return nil
	})

	return nil
}

func (a *App) initDi(ctx context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initCloser(ctx context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initTelegramBot(ctx context.Context) error {
	telegramBot, err := a.diContainer.TelegramBot(ctx)
	if err != nil {
		logger.Error(ctx, "Faild to init telegram bot", zap.Error(err))
		return err
	}

	telegramBot.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/start",
		bot.MatchTypeExact,
		func(ctx context.Context, b *bot.Bot, update *models.Update) {
			logger.Info(ctx, "start message, chat id", zap.Int64("chat_id", update.Message.Chat.ID))

			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text: `📝 Notification bot активирован!
Теперь Вы будете получать уведомления об оплате и сборке заказов!.`,
			})
			if err != nil {
				logger.Error(ctx, "Failed to send activation message", zap.Error(err))
			}
		},
	)

	go func() {
		logger.Info(ctx, "🤖 Telegram bot started...")
		telegramBot.Start(ctx)
	}()

	return nil
}

func (a *App) runOrderPaidConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 💵 OrderPaid Kafka consumer running")

	c, err := a.diContainer.OrderPaidConsumerService(ctx)
	if err != nil {
		return err
	}

	err = c.RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runOrderAssembledConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 📦 OrderAssembled Kafka consumer running")

	c, err := a.diContainer.OrderAssembledConsumerService(ctx)
	if err != nil {
		return err
	}

	err = c.RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}
