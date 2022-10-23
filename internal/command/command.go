package command

import (
	"context"
	"fmt"
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

	args := upd.Message.CommandArguments()
	// Remove urls if args are specified
	// Otherwise display inline buttons with sources
	if len(args) > 0 {
		urls := strings.Split(args, " ")
		for _, url := range urls {
			err := db.RemoveSource(ctx, usr.ID, url)
			if err != nil {
				return err
			}
		}
	} else {
		// TODO return inline buttons with user's sources
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "**Please choose a subscription you'd like to unsubscribe from:**\n*test*")
		msg.ParseMode = tgbotapi.ModeMarkdownV2

		sourcesTitleURL, err := db.GetUserSourcesTitleURL(ctx, usr.ID)
		if err != nil {
			return err
		}

		var buttons [][]tgbotapi.InlineKeyboardButton
		for _, sourceTitleURL := range sourcesTitleURL {
			row := tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(sourceTitleURL[0], fmt.Sprintf("%v %v", cfg.UnsubscribeCallbackData, sourceTitleURL[1])),
			)
			buttons = append(buttons, row)
		}
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)

		_, err = bot.Send(msg)
		if err != nil {
			return err
		}
	}

	return nil
}
