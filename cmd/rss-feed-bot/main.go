package main

import (
	"bufio"
	"context"
	"flag"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/env"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/util"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/command"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

var (
	cfg      *util.Config
	bot      *tgbotapi.BotAPI
	commands map[string]command.ICommand
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

	commands = makeCommands()

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

func makeCommands() map[string]command.ICommand {
	return map[string]command.ICommand{
		"/scream":  &command.Scream{},
		"/whisper": &command.Whisper{},
		"/menu":    &command.Menu{},
	}
}