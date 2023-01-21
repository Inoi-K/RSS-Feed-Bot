package command

import (
	"github.com/Inoi-K/RSS-Feed-Bot/configs/consts"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

// reply builds message and sends it to the chat
func reply(bot *tgbotapi.BotAPI, chat *tgbotapi.Chat, text string) error {
	msg := newMessage(chat, text, nil)
	_, err := bot.Send(msg)
	return err
}

// reply builds message with keyboard and sends it to the chat
func replyKeyboard(bot *tgbotapi.BotAPI, chat *tgbotapi.Chat, text string, keyboard interface{}) error {
	msg := newMessage(chat, text, keyboard)
	_, err := bot.Send(msg)
	return err
}

// newMessage builds message with all needed parameters
func newMessage(chat *tgbotapi.Chat, text string, keyboard interface{}) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chat.ID, text)
	msg.ParseMode = consts.ParseMode
	msg.ReplyMarkup = keyboard
	return msg
}

// makeInlineKeyboard builds inline keyboard from the content
func makeInlineKeyboard(content []model.Content, commandButton string) tgbotapi.InlineKeyboardMarkup {
	var keyboard [][]tgbotapi.InlineKeyboardButton

	for _, c := range content {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(c.Text, strings.Join([]string{commandButton, c.Data}, consts.ArgumentsSeparator)),
		)
		keyboard = append(keyboard, row)
	}

	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}
