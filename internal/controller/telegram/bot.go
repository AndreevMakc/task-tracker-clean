package telegram

import (
	"context"
	"fmt"
	"log"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"task-tracker-clean/internal/usecase"
)

type Bot struct {
	bot      *telego.Bot
	handler  *th.BotHandler
	uc       usecase.TaskUsecase
}

func NewBot(uc usecase.TaskUsecase, token string) (*Bot, error) {
	if token == "" {
		log.Println("TELEGRAM_BOT_TOKEN not set, bot disabled")
		return nil, nil
	}

	bot, err := telego.NewBot(token)
	if err != nil {
		return nil, err
	}

	return &Bot{bot: bot, uc: uc}, nil
}

func (b *Bot) Run(ctx context.Context) error {
	if b == nil || b.bot == nil {
		return nil
	}

	botUser, err := b.bot.GetMe(ctx)
	if err != nil {
		return fmt.Errorf("failed to get bot info: %w", err)
	}

	log.Printf("starting telegram bot: %s", botUser.Username)

	updatesChan, err := b.bot.UpdatesViaLongPolling(ctx, nil)
	if err != nil {
		return err
	}

	bh, err := th.NewBotHandler(b.bot, updatesChan)
	if err != nil {
		return err
	}
	b.handler = bh

	handler := NewTaskHandler(b.uc)
	handler.RegisterBotHandler(bh)

	log.Println("telegram bot started")

	go func() {
		if err := bh.Start(); err != nil {
			log.Printf("telegram bot handler error: %v", err)
		}
	}()

	return nil
}

func (b *Bot) Stop() error {
	if b == nil || b.handler == nil {
		return nil
	}
	return b.handler.Stop()
}
