package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	// `for {` means the loop is infinite until we manually stop it
	for {
		select {
		// stop looping if ctx is cancelled
		case <-ctx.Done():
			return
		// receive update from channel and then handle it
		case update := <-updates:
			handleUpdate(ctx, update)
		}
	}
}

func handleUpdate(ctx context.Context, update tgbotapi.Update) {
	switch {
	// Handle messages
	case update.Message != nil:
		handleMessage(ctx, update)

	// Handle button clicks
	case update.CallbackQuery != nil:
		handleButton(ctx, update)
	}
}

func handleMessage(ctx context.Context, update tgbotapi.Update) {
	message := update.Message
	user := message.From
	text := message.Text

	if user == nil {
		return
	}

	// Print to console
	log.Printf("%s wrote %s", user.FirstName, text)

	var err error
	if strings.HasPrefix(text, "/") {
		err = handleCommand(ctx, update, text)
	} else if cfg.Screaming && len(text) > 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, strings.ToUpper(text))
		// To preserve markdown, we attach entities (bold, italic..)
		msg.Entities = message.Entities
		_, err = bot.Send(msg)
	} else {
		// This is equivalent to forwarding, without the sender's name
		copyMsg := tgbotapi.NewCopyMessage(message.Chat.ID, message.Chat.ID, message.MessageID)
		_, err = bot.CopyMessage(copyMsg)
	}

	if err != nil {
		log.Printf("couldn't process the message: %s", err.Error())
	}
}

// When we get a command, we react accordingly
func handleCommand(ctx context.Context, update tgbotapi.Update, curCommand string) error {
	return commands[curCommand].Execute(ctx, bot, update, cfg)
}

func handleButton(ctx context.Context, update tgbotapi.Update) {
	query := update.CallbackQuery

	var text string

	markup := tgbotapi.NewInlineKeyboardMarkup()
	message := query.Message

	if query.Data == cfg.NextButton {
		text = cfg.SecondMenu
		markup = cfg.SecondMenuMarkup
	} else if query.Data == cfg.BackButton {
		text = cfg.FirstMenu
		markup = cfg.FirstMenuMarkup
	}

	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	_, err := bot.Send(callbackCfg)
	if err != nil {
		log.Printf("callback config error: %v", err)
	}

	// Replace menu text and keyboard
	msg := tgbotapi.NewEditMessageTextAndMarkup(message.Chat.ID, message.MessageID, text, markup)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err = bot.Send(msg)
	if err != nil {
		log.Printf("menu text and keyboard error: %v", err)
	}
}
