package main

import (
	"bufio"
	"context"
	"flag"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/consts"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/flags"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/command"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

var (
	bot      *tgbotapi.BotAPI
	db       *database.Database
	commands map[string]command.ICommand
)

func main() {
	flag.Parse()

	var err error
	// Connect to the bot
	bot, err = tgbotapi.NewBotAPI(*flags.Token)
	if err != nil {
		// Abort if something is wrong
		log.Panic(err)
	}
	// Set this to true to log all interactions with telegram servers
	bot.Debug = true

	// Create a new cancellable background context. Calling `cancel()` leads to the cancellation of the context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// Create a new database connection
	db, err = database.ConnectDB(ctx)

	commands = makeCommands()

	// Set update rate
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

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

// makeCommands creates all bot commands and buttons
func makeCommands() map[string]command.ICommand {
	return map[string]command.ICommand{
		consts.MenuCommand:  &command.Menu{},
		consts.StartCommand: &command.Start{},

		consts.SubscribeCommand:   &command.Subscribe{},
		consts.UnsubscribeCommand: &command.Unsubscribe{},
		consts.UnsubscribeButton:  &command.UnsubscribeButton{},

		consts.NavigationButton: &command.NavigationButton{},

		"tick": &command.Ticker{},
		"stop": &command.StopTicker{},

		consts.UpdateCommand: &command.Update{},

		consts.ActivateCommand:   &command.Activate{},
		consts.DeactivateCommand: &command.Deactivate{},
		consts.ActivateButton:    &command.ActivateButton{},
		consts.DeactivateButton:  &command.DeactivateButton{},

		consts.ListCommand: &command.List{},
	}
}
