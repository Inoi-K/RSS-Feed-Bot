package command

import (
	"context"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/consts"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/database"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/feed"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

// ICommand provides an interface for all commands and buttons callbacks
type ICommand interface {
	Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error
}

// Menu command replies with first menu
type Menu struct{}

func (c *Menu) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, consts.FirstMenu)
	msg.ParseMode = consts.ParseMode
	msg.ReplyMarkup = consts.FirstMenuMarkup
	_, err := bot.Send(msg)
	return err
}

// Start command begins an interaction with the chat and creates the record in database
type Start struct{}

func (c *Start) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()
	usr := upd.SentFrom()

	db := database.GetDB()

	return db.AddChat(ctx, chat.ID, usr.LanguageCode)
}

// Subscribe command adds sources to database and associates it with the chat
type Subscribe struct{}

func (c *Subscribe) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()

	db := database.GetDB()

	urls := strings.Split(args, consts.ArgumentsSeparator)
	for _, url := range urls {
		err := db.AddSource(ctx, chat.ID, url)
		if err != nil {
			return err
		}
	}

	return nil
}

// Unsubscribe command removes provided sources from the chat or replies with menu with buttons as sources
type Unsubscribe struct{}

func (c *Unsubscribe) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()

	db := database.GetDB()

	// Remove urls if args are specified
	// Otherwise display inline buttons with sources
	if len(args) > 0 {
		urls := strings.Split(args, consts.ArgumentsSeparator)
		for _, url := range urls {
			err := db.RemoveSource(ctx, chat.ID, url)
			if err != nil {
				return err
			}
		}
	} else {
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Please choose a subscription you'd like to unsubscribe from:")
		msg.ParseMode = consts.ParseMode

		sourcesTitleURL, err := db.GetChatSourcesTitleURL(ctx, chat.ID)
		if err != nil {
			return err
		}

		var buttons [][]tgbotapi.InlineKeyboardButton
		for _, sourceTitleURL := range sourcesTitleURL {
			row := tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(sourceTitleURL[0], strings.Join([]string{consts.UnsubscribeButton, sourceTitleURL[1]}, consts.ArgumentsSeparator)),
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

// UnsubscribeButton command gets called by button callback from 'Unsubscribe menu' and then removes provided source from the chat
type UnsubscribeButton struct{}

func (c *UnsubscribeButton) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()
	db := database.GetDB()

	err := db.RemoveSource(ctx, chat.ID, args)
	if err != nil {
		return err
	}

	return nil
}

// NavigationButton gets called by button callback from 'Menu' and handles next/back navigation between menus/pages
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

type Ticker struct{}

func (c *Ticker) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	feed.Begin(ctx, bot)

	return nil
}

type StopTicker struct{}

func (c *StopTicker) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	feed.End()

	return nil
}
