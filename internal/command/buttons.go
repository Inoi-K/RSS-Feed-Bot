package command

import (
	"context"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/consts"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/database"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/structs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

// UnsubscribeButton gets called by button callback from 'Unsubscribe menu' and removes provided source from the chat
type UnsubscribeButton struct{}

func (c *UnsubscribeButton) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()
	db := database.GetDB()

	err := db.RemoveSource(ctx, chat.ID, args)
	if err != nil {
		return err
	}

	return editInlineChatSourceKeyboard(bot, upd, args)
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

// ActivateButton gets called by button callback from 'Activate menu' and switches state of the provided source
type ActivateButton struct{}

func (c *ActivateButton) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	return setIsActiveButton(ctx, bot, upd, args, true)
}

// DeactivateButton gets called by button callback from 'Deactivate menu' and switches state of the provided source
type DeactivateButton struct{}

func (c *DeactivateButton) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	return setIsActiveButton(ctx, bot, upd, args, false)
}

// setIsActiveButton switches provide sources for the chat to provided state
// and edits the keyboard from which the callback was sent
func setIsActiveButton(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string, isActive bool) error {
	chat := upd.FromChat()
	db := database.GetDB()

	err := db.AlterChatSource(ctx, chat.ID, args, structs.ChatSource{IsActive: isActive})
	if err != nil {
		return err
	}

	return editInlineChatSourceKeyboard(bot, upd, args)
}

// editInlineChatSourceKeyboard removes the provided button in the keyboard from which the callback was sent
func editInlineChatSourceKeyboard(bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	message := upd.CallbackQuery.Message

	var newKeyboard [][]tgbotapi.InlineKeyboardButton
	currentKeyboard := message.ReplyMarkup.InlineKeyboard
	for rowIndex, row := range currentKeyboard {
		button := row[0]
		if strings.HasSuffix(*button.CallbackData, args) {
			newKeyboard = append(currentKeyboard[:rowIndex], currentKeyboard[rowIndex+1:]...)
			break
		}
	}
	newMarkup := tgbotapi.NewInlineKeyboardMarkup(newKeyboard...)

	msg := tgbotapi.NewEditMessageReplyMarkup(message.Chat.ID, message.MessageID, newMarkup)
	_, err := bot.Send(msg)
	return err
}
