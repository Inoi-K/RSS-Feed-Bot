package main

import (
	"context"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/consts"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

// receiveUpdates handles updates and context cancel
func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		// stop looping if ctx is cancelled
		case <-ctx.Done():
			return
		case update := <-updates:
			handleUpdate(ctx, update)
		}
	}
}

// handleUpdate distributes incoming update
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

// handleMessage defines the type of the message (command or other - replies as echo in the latter case)
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
	if message.IsCommand() {
		err = handleCommand(ctx, update)
	} else {
		// This is equivalent to forwarding, without the sender's name
		copyMsg := tgbotapi.NewCopyMessage(message.Chat.ID, message.Chat.ID, message.MessageID)
		_, err = bot.CopyMessage(copyMsg)
	}

	if err != nil {
		log.Printf("couldn't process the message: %s", err.Error())
	}
}

// handleCommand handles commands specifically
func handleCommand(ctx context.Context, update tgbotapi.Update) error {
	curCommand := update.Message.Command()
	return commands[curCommand].Execute(ctx, bot, update, update.Message.CommandArguments())
}

// handleButton handles buttons callback specifically
func handleButton(ctx context.Context, update tgbotapi.Update) {
	query := update.CallbackQuery
	command, args, _ := strings.Cut(query.Data, consts.ArgumentsSeparator)

	err := commands[command].Execute(ctx, bot, update, args)
	if err != nil {
		log.Printf("couldn't process button callback: %v", err)
	}

	// close the query
	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	_, err = bot.Request(callbackCfg)
	if err != nil {
		log.Printf("callback config error: %v", err)
	}
}
