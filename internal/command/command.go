package command

import (
	"github.com/Inoi-K/RSS-Feed-Bot/configs/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ICommand interface {
	Execute(bot *tgbotapi.BotAPI, upd tgbotapi.Update, cfg *util.Config) error
}

type Scream struct{}

func (c *Scream) Execute(bot *tgbotapi.BotAPI, upd tgbotapi.Update, cfg *util.Config) error {
	cfg.Screaming = true
	return nil
}

type Whisper struct{}

func (c *Whisper) Execute(bot *tgbotapi.BotAPI, upd tgbotapi.Update, cfg *util.Config) error {
	cfg.Screaming = false
	return nil
}

type Menu struct{}

func (c *Menu) Execute(bot *tgbotapi.BotAPI, upd tgbotapi.Update, cfg *util.Config) error {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, cfg.FirstMenu)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = cfg.FirstMenuMarkup
	_, err := bot.Send(msg)
	return err
}
