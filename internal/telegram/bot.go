package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"tester/internal/common/telegram/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const ServiceName = "TGBot"

type Bot struct {
	config config.Config

	bot *tgbotapi.BotAPI
	wg  *sync.WaitGroup

	shutdownCtx       context.Context
	cancelShutdownCtx func()
	alertsChan        chan alertInfo

	logger *slog.Logger
}

func New(conf config.Config, logger *slog.Logger) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(conf.BotToken)
	if err != nil {
		return nil, fmt.Errorf("can't create tg bot: %w", err)
	}

	ctx, cancelCtx := context.WithCancel(context.Background())
	b := &Bot{
		config:            conf,
		bot:               bot,
		wg:                new(sync.WaitGroup),
		shutdownCtx:       ctx,
		cancelShutdownCtx: cancelCtx,
		alertsChan:        make(chan alertInfo, 1),
		logger:            logger,
	}

	return b, nil
}

func (b *Bot) Start() {
	b.wg.Add(1)
	go b.processAlerts()
}

func (b *Bot) Shutdown() {
	b.cancelShutdownCtx()
	close(b.alertsChan)
	b.wg.Wait()
}
