package main

import (
	"bufio"
	"context"
	"flag"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/env"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strings"
)

var (
	cfg *util.Config
	bot *tgbotapi.BotAPI
)

func main() {
	flag.Parse()

	// Create configuration
	cfg = util.NewConfig()

	var err error
	// Connect to the bot
	bot, err = tgbotapi.NewBotAPI(*env.Token)
	if err != nil {
		// Abort if something is wrong
		log.Panic(err)
	}
	// Set this to true to log all interactions with telegram servers
	bot.Debug = false

	// Set update rate
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Create a new cancellable background context. Calling `cancel()` leads to the cancellation of the context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// `updates` is a golang channel which receives telegram updates
	updates := bot.GetUpdatesChan(u)

	// Pass cancellable context to goroutine
	go receiveUpdates(ctx, updates)

	// Tell the user the bot is online
	log.Println("Start listening for updates. Press enter to stop")

	// Wait for a newline symbol, then cancel handling updates
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()
}

func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	// `for {` means the loop is infinite until we manually stop it
	for {
		select {
		// stop looping if ctx is cancelled
		case <-ctx.Done():
			return
		// receive update from channel and then handle it
		case update := <-updates:
			handleUpdate(update)
		}
	}
}

func handleUpdate(update tgbotapi.Update) {
	switch {
	// Handle messages
	case update.Message != nil:
		handleMessage(update.Message)
		break

	// Handle button clicks
	case update.CallbackQuery != nil:
		handleButton(update.CallbackQuery)
		break
	}
}

func handleMessage(message *tgbotapi.Message) {
	user := message.From
	text := message.Text

	if user == nil {
		return
	}

	// Print to console
	log.Printf("%s wrote %s", user.FirstName, text)

	var err error
	if strings.HasPrefix(text, "/") {
		err = handleCommand(message.Chat.ID, text)
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
		log.Printf("An error occured: %s", err.Error())
	}
}

// When we get a command, we react accordingly
func handleCommand(chatId int64, command string) error {
	var err error

	switch command {
	case "/scream":
		cfg.Screaming = true
		break

	case "/whisper":
		cfg.Screaming = false
		break

	case "/menu":
		err = sendMenu(chatId, cfg.FirstMenu, cfg.FirstMenuMarkup)
		break
	}

	return err
}

func handleButton(query *tgbotapi.CallbackQuery) {
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

func sendMenu(chatId int64, text string, markup tgbotapi.InlineKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = markup
	_, err := bot.Send(msg)
	return err
}
