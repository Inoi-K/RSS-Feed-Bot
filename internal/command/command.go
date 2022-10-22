package command

import (
	"context"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/util"
	"github.com/Inoi-K/RSS-Feed-Bot/pkg/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type ICommand interface {
	Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, cfg *util.Config) error
}

type Scream struct{}

func (c *Scream) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, cfg *util.Config) error {
	cfg.Screaming = true
	return nil
}

type Whisper struct{}

func (c *Whisper) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, cfg *util.Config) error {
	cfg.Screaming = false
	return nil
}

type Menu struct{}

func (c *Menu) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, cfg *util.Config) error {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, cfg.FirstMenu)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = cfg.FirstMenuMarkup
	_, err := bot.Send(msg)
	return err
}

type Start struct{}

func (c *Start) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, cfg *util.Config) error {
	usr := upd.SentFrom()

	db := database.GetDB()

	return db.AddUser(ctx, usr.ID, usr.LanguageCode)
}

type Subscribe struct{}

func (c *Subscribe) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, cfg *util.Config) error {
	usr := upd.SentFrom()

	db := database.GetDB()

	urls := strings.Split(upd.Message.CommandArguments(), " ")
	for _, url := range urls {
		err := db.AddSource(ctx, usr.ID, url)
		if err != nil {
			return err
		}
	}

	return nil
}

type Unsubscribe struct{}

func (c *Unsubscribe) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, cfg *util.Config) error {
	usr := upd.SentFrom()

	db := database.GetDB()

	return db.RemoveSource(ctx, usr.ID, usr.LanguageCode)
}
