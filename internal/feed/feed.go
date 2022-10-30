package feed

import (
	"context"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/consts"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
)

var cancel func()
var canceled chan struct{}

func Begin(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64) {
	ticker := time.NewTicker(consts.FeedUpdateIntervalSeconds * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				tick(bot, chatID)
			case <-canceled:
				ticker.Stop()
				return
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()

	// local cancel function
	canceled = make(chan struct{})
	cancel = func() {
		select {
		case <-canceled:
		default:
			close(canceled)
		}
	}
}

func tick(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "tick")
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("couldn't send tick: %v", err)
	}
}

func End() {
	if cancel != nil {
		cancel()
	}
}
