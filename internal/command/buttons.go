package command

import (
	"context"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/consts"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/database"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/structs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

type ActivateButton struct{}

func (c *ActivateButton) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()
	db := database.GetDB()

	err := db.AlterChatSource(ctx, chat.ID, args, structs.ChatSource{IsActive: true})
	if err != nil {
		return err
	}

	return nil
}

type DeactivateButton struct{}

func (c *DeactivateButton) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()
	db := database.GetDB()

	err := db.AlterChatSource(ctx, chat.ID, args, structs.ChatSource{IsActive: false})
	if err != nil {
		return err
	}

	return nil
}
