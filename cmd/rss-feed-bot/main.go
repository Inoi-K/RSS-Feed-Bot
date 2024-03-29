package main

import (
	"context"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/consts"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/flags"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/command"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/database"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/feed"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

var (
	bot      *tgbotapi.BotAPI
	db       *database.Database
	commands map[string]command.ICommand
)

func main() {
	var err error
	// Connect to the bot
	bot, err = tgbotapi.NewBotAPI(*flags.Token)
	if err != nil {
		log.Panic(err)
	}
	// Set this to true to log all interactions with telegram servers
	bot.Debug = true

	// Create a new cancellable background context. Calling `cancel()` leads to the cancellation of the context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// Create a new database connection
	db, err = database.ConnectDB(ctx)
	if err != nil {
		log.Fatalf("couldn't connect to DB: %v", err)
	}
	err = database.SetUp(ctx)
	if err != nil {
		log.Fatalf("couldn't create tables in an empty database %v", err)
	}

	// Generate structs for commands
	commands = makeCommands()

	// Set update rate
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	// `updates` is a golang channel which receives telegram updates
	updates := bot.GetUpdatesChan(u)
	// Pass cancellable context to goroutine
	go receiveUpdates(ctx, updates)

	// Begin feed updating
	feed.Begin(ctx, bot)

	// Tell the user the bot is online
	log.Println("Start listening for updates...")

	select {}

	// Wait for a newline symbol, then cancel handling updates
	//bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()
}

// makeCommands creates all bot commands and buttons
func makeCommands() map[string]command.ICommand {
	return map[string]command.ICommand{
		//consts.MenuCommand:  &command.Menu{},
		consts.StartCommand: &command.Start{},

		consts.SubscribeCommand:   &command.Subscribe{},
		consts.UnsubscribeCommand: &command.Unsubscribe{},
		consts.UnsubscribeButton:  &command.UnsubscribeButton{},

		consts.NavigationButton: &command.NavigationButton{},

		//"tick": &command.Ticker{},
		//"cancel": &command.StopTicker{},

		consts.UpdateCommand: &command.Update{},

		consts.ActivateCommand:   &command.Activate{},
		consts.DeactivateCommand: &command.Deactivate{},
		consts.ActivateButton:    &command.ActivateButton{},
		consts.DeactivateButton:  &command.DeactivateButton{},

		consts.ListCommand: &command.List{},
		consts.HelpCommand: &command.Help{},

		consts.LanguageCommand: &command.Language{},
		consts.LanguageButton:  &command.LanguageButton{},

		consts.PingCommand: &command.Ping{},
	}
}
