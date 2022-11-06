package command

import (
	"context"
	"fmt"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/consts"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/database"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/feed"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/structs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func reply(bot *tgbotapi.BotAPI, chat *tgbotapi.Chat, text string) error {
	msg := tgbotapi.NewMessage(chat.ID, text)
	msg.ParseMode = consts.ParseMode
	_, err := bot.Send(msg)
	return err
}

// ICommand provides an interface for all commands and buttons callbacks
type ICommand interface {
	Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error
}

// Menu command replies with first menu
type Menu struct{}

func (c *Menu) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, consts.FirstMenu)
	msg.ParseMode = consts.ParseMode
	msg.ReplyMarkup = consts.FirstMenuMarkup
	_, err := bot.Send(msg)
	return err
}

// Start command begins an interaction with the chat and creates the record in database
type Start struct{}

func (c *Start) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()
	usr := upd.SentFrom()

	db := database.GetDB()

	return db.AddChat(ctx, chat.ID, usr.LanguageCode)
}

// Subscribe command adds sources to database and associates it with the chat
type Subscribe struct{}

func (c *Subscribe) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()
	db := database.GetDB()

	urls := strings.Split(args, consts.ArgumentsSeparator)

	// TODO add url checker here

	for _, url := range urls {
		err := db.AddSource(ctx, chat.ID, url)
		if err != nil {
			return err
		}
	}

	return reply(bot, chat, "Successfully subscribed")
}

// Unsubscribe command removes provided sources from the chat
// or replies with menu with buttons as sources
type Unsubscribe struct{}

func (c *Unsubscribe) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()
	db := database.GetDB()

	// Remove urls if args are specified
	// Otherwise display inline buttons with sources
	if len(args) > 0 {
		urls := strings.Split(args, consts.ArgumentsSeparator)
		for _, url := range urls {
			err := db.RemoveSource(ctx, chat.ID, url)
			if err != nil {
				return err
			}
		}
	} else {
		infoText := "Please choose a subscription you'd like to unsubscribe from:"
		err := replyInlineChatSourceKeyboard(ctx, bot, upd, nil, infoText, consts.UnsubscribeButton)
		if err != nil {
			return err
		}
	}

	return nil
}

// Activate command switches source for the chat to active state
type Activate struct{}

func (c *Activate) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	return setIsActive(ctx, bot, upd, args, true)
}

// Deactivate command switches source for the chat to inactive state
type Deactivate struct{}

func (c *Deactivate) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	return setIsActive(ctx, bot, upd, args, false)
}

// setIsActive switches provided sources for the chat to provided state
// or replies with menu with inline buttons as corresponding sources
func setIsActive(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string, isActive bool) error {
	chat := upd.FromChat()
	db := database.GetDB()

	// Alter sources if args are specified
	// Otherwise display inline buttons with sources
	if len(args) > 0 {
		urls := strings.Split(args, consts.ArgumentsSeparator)
		for _, url := range urls {
			err := db.AlterChatSource(ctx, chat.ID, url, structs.ChatSource{IsActive: isActive})
			if err != nil {
				return err
			}
		}
	} else {
		var state string
		if isActive {
			state = consts.ActivateButton
		} else {
			state = consts.DeactivateButton
		}
		infoText := fmt.Sprintf("Please choose a subscription you'd like to %v:", state)

		err := replyInlineChatSourceKeyboard(ctx, bot, upd, &structs.ChatSource{IsActive: !isActive}, infoText, state)
		if err != nil {
			return err
		}
	}

	return nil
}

// replyInlineChatSourceKeyboard gets title and url of the sources associated with the chat
// and replies with inline buttons with commandButton as their beginning of the data
func replyInlineChatSourceKeyboard(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, cs *structs.ChatSource, infoText string, commandButton string) error {
	chat := upd.FromChat()
	db := database.GetDB()

	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, infoText)
	msg.ParseMode = consts.ParseMode

	sourcesTitleURL, err := db.GetChatSourceTitleURL(ctx, chat.ID, cs)
	if err != nil {
		return err
	}

	var keyboard [][]tgbotapi.InlineKeyboardButton
	for _, sourceTitleURL := range sourcesTitleURL {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(sourceTitleURL[0], strings.Join([]string{commandButton, sourceTitleURL[1]}, consts.ArgumentsSeparator)),
		)
		keyboard = append(keyboard, row)
	}
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboard...)

	_, err = bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// Ticker command starts a ticker
type Ticker struct{}

func (c *Ticker) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	feed.Begin(ctx, bot)

	return reply(bot, upd.FromChat(), "Ticker started")
}

// StopTicker command stops the ticker started in Ticker command
type StopTicker struct{}

func (c *StopTicker) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	feed.End()

	return reply(bot, upd.FromChat(), "Ticker stopped")
}

// Update command gets recent posts
type Update struct{}

func (c *Update) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()

	msg := tgbotapi.NewMessage(chat.ID, "Updates:")
	msg.ReplyMarkup = consts.UpdateKeyboard
	_, err := bot.Send(msg)
	if err != nil {
		return err
	}

	feed.ProcessNewPosts(ctx, bot)

	return nil
}

// List command shows all current subscriptions of the chat
type List struct{}

func (c *List) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()
	db := database.GetDB()

	sourcesTitleURL, err := db.GetChatSourceTitleURL(ctx, chat.ID, nil)
	if err != nil {
		return err
	}

	text := "*Subscription list*"
	for _, sourceTitleUrl := range sourcesTitleURL {
		text += fmt.Sprintf("\n[%v](%v)", sourceTitleUrl[0], sourceTitleUrl[1])
	}

	msg := tgbotapi.NewMessage(chat.ID, text)
	msg.ParseMode = consts.ParseMode
	_, err = bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
