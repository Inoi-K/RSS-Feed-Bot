package command

import (
	"context"
	"fmt"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/consts"
	"github.com/Inoi-K/RSS-Feed-Bot/pkg/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type ICommand interface {
	Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error
}

type Menu struct{}

func (c *Menu) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, consts.FirstMenu)
	msg.ParseMode = consts.ParseMode
	msg.ReplyMarkup = consts.FirstMenuMarkup
	_, err := bot.Send(msg)
	return err
}

type Start struct{}

func (c *Start) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	usr := upd.SentFrom()

	db := database.GetDB()

	return db.AddUser(ctx, usr.ID, usr.LanguageCode)
}

type Subscribe struct{}

func (c *Subscribe) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	usr := upd.SentFrom()

	db := database.GetDB()

	urls := strings.Split(args, consts.ArgumentsSeparator)
	for _, url := range urls {
		err := db.AddSource(ctx, usr.ID, url)
		if err != nil {
			return err
		}
	}

	return nil
}

type Unsubscribe struct{}

func (c *Unsubscribe) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	usr := upd.SentFrom()

	db := database.GetDB()

	// Remove urls if args are specified
	// Otherwise display inline buttons with sources
	if len(args) > 0 {
		urls := strings.Split(args, consts.ArgumentsSeparator)
		for _, url := range urls {
			err := db.RemoveSource(ctx, usr.ID, url)
			if err != nil {
				return err
			}
		}
	} else {
		// TODO return inline buttons with user's sources
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "**Please choose a subscription you'd like to unsubscribe from:**\n*test*")
		msg.ParseMode = consts.ParseMode

		sourcesTitleURL, err := db.GetUserSourcesTitleURL(ctx, usr.ID)
		if err != nil {
			return err
		}

		var buttons [][]tgbotapi.InlineKeyboardButton
		for _, sourceTitleURL := range sourcesTitleURL {
			row := tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(sourceTitleURL[0], fmt.Sprintf("%v %v", consts.UnsubscribeButton, sourceTitleURL[1])),
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

// region Buttons

type UnsubscribeButton struct{}

func (c *UnsubscribeButton) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	usr := upd.SentFrom()
	db := database.GetDB()

	err := db.RemoveSource(ctx, usr.ID, args)
	if err != nil {
		return err
	}

	return nil
}

type NavigationButton struct{}

func (c *NavigationButton) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	message := upd.CallbackQuery.Message

	var (
		text   string
		markup tgbotapi.InlineKeyboardMarkup
	)

	switch args {
	case consts.NextText:
		text = consts.SecondMenu
		markup = consts.SecondMenuMarkup
	case consts.BackText:
		text = consts.FirstMenu
		markup = consts.FirstMenuMarkup
	}

	// Replace menu text and keyboard
	msg := tgbotapi.NewEditMessageTextAndMarkup(message.Chat.ID, message.MessageID, text, markup)
	msg.ParseMode = consts.ParseMode
	_, err := bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

//endregion
