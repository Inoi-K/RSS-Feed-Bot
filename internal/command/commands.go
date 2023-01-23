package command

import (
	"context"
	"fmt"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/consts"
	loc "github.com/Inoi-K/RSS-Feed-Bot/configs/localization"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/builder"
	db "github.com/Inoi-K/RSS-Feed-Bot/internal/database"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/feed"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/model"
	"github.com/Inoi-K/RSS-Feed-Bot/pkg/parser"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

// ICommand provides an interface for all commands and buttons callbacks
type ICommand interface {
	Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error
}

// Menu command replies with first menu
type Menu struct{}

func (c *Menu) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()
	return builder.ReplyKeyboard(bot, chat, consts.FirstMenu, consts.FirstMenuMarkup)
}

// Start command begins an interaction with the chat and creates the record in database
type Start struct{}

func (c *Start) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()
	usr := upd.SentFrom()

	err := db.AddChat(ctx, chat.ID, usr.LanguageCode)
	if err != nil {
		return err
	}

	ok := loc.ChangeLanguage(usr.LanguageCode)
	// if user's language is not supported then set default language to english
	if !ok {
		loc.ChangeLanguage("en")
	}

	defer builder.Reply(bot, chat, loc.Message(loc.Help))
	return builder.Reply(bot, chat, loc.Message(loc.Start))
}

// Subscribe command adds sources to database and associates it with the chat
type Subscribe struct{}

func (c *Subscribe) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()

	urls := strings.Split(args, consts.ArgumentsSeparator)
	for _, url := range urls {
		ans := fmt.Sprintf(loc.Message(loc.SubFail), loc.Message(loc.NotValidLink))

		// SERVICE VALIDATION
		//res, err := client.Validate(url)
		//if err != nil {
		//	err = reply(bot, chat, ans)
		//	if err != nil {
		//		return err
		//	}
		//	continue
		//}
		//
		//if res.Valid {
		//	err := db.AddSource(ctx, chat.ID, res.Title, url)
		//	if err != nil {
		//		return err
		//	}
		//} else {
		//	err = reply(bot, chat, ans)
		//	if err != nil {
		//		return err
		//	}
		//	continue
		//}
		//ans = fmt.Sprintf(loc.Message(consts.SubscribeCommand], res.Title, url)

		// LIB VALIDATION
		source, err := parser.Parse(url)
		if err != nil {
			err = builder.Reply(bot, chat, ans)
			if err != nil {
				return err
			}
			continue
		}

		err = db.AddSource(ctx, chat.ID, source.Title, url)
		if err != nil {
			err = builder.Reply(bot, chat, ans)
			if err != nil {
				return err
			}
			continue
		}
		ans = fmt.Sprintf(loc.Message(loc.Sub), source.Title, url)

		err = builder.Reply(bot, chat, ans)
		if err != nil {
			return err
		}
	}

	return nil
}

// Unsubscribe command removes provided sources from the chat
// or replies with menu with buttons as sources
type Unsubscribe struct{}

func (c *Unsubscribe) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()

	// Remove urls if args are specified
	// Otherwise display inline buttons with sources
	if len(args) > 0 {
		urls := strings.Split(args, consts.ArgumentsSeparator)
		for _, url := range urls {
			err := db.RemoveSource(ctx, chat.ID, url)
			if err != nil {
				return err
			}
			ans := fmt.Sprintf(loc.Message(loc.UnsubSuccess), url)
			err = builder.Reply(bot, chat, ans)
		}
	} else {
		infoText := loc.Message(loc.Unsub)
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

	// Alter sources if args are specified
	// Otherwise display inline buttons with sources
	if len(args) > 0 {
		urls := strings.Split(args, consts.ArgumentsSeparator)
		for _, url := range urls {
			err := db.AlterChatSource(ctx, chat.ID, url, model.ChatSource{IsActive: isActive})
			if err != nil {
				return err
			}
		}
	} else {
		state := consts.ActivateButton
		infoText := loc.Message(loc.Activate)
		if !isActive {
			state = consts.DeactivateButton
			infoText = loc.Message(loc.Deactivate)
		}

		err := replyInlineChatSourceKeyboard(ctx, bot, upd, &model.ChatSource{IsActive: !isActive}, infoText, state)
		if err != nil {
			return err
		}
	}

	return nil
}

// replyInlineChatSourceKeyboard gets title and url of the sources associated with the chat
// and replies with inline buttons with commandButton as their beginning of the data
func replyInlineChatSourceKeyboard(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, cs *model.ChatSource, infoText string, commandButton string) error {
	chat := upd.FromChat()

	sourcesTitleURL, err := db.GetChatSourceTitleID(ctx, chat.ID, cs)
	if err != nil {
		return err
	}

	return builder.ReplyKeyboard(bot, chat, infoText, builder.MakeInlineKeyboard(sourcesTitleURL, commandButton))
}

// Ticker command starts a ticker
type Ticker struct{}

func (c *Ticker) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	feed.Begin(ctx, bot)

	return builder.Reply(bot, upd.FromChat(), "Ticker started")
}

// StopTicker command stops the ticker started in Ticker command
type StopTicker struct{}

func (c *StopTicker) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	feed.End()

	return builder.Reply(bot, upd.FromChat(), "Ticker stopped")
}

// Update command gets recent posts
type Update struct{}

func (c *Update) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()

	go feed.ProcessNewPosts(ctx, bot)

	return builder.ReplyKeyboard(bot, chat, loc.Message(loc.Upd), consts.UpdateKeyboard)
}

// List command shows all current subscriptions of the chat
type List struct{}

func (c *List) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()

	sourcesTitleURL, err := db.GetChatSourceTitleID(ctx, chat.ID, nil)
	if err != nil {
		return err
	}

	text := loc.Message(loc.List)
	for _, sourceTitleURL := range sourcesTitleURL {
		text += fmt.Sprintf("\n[%v](%v)", sourceTitleURL.Text, sourceTitleURL.Data)
	}

	return builder.Reply(bot, chat, text)
}

// Help command shows information about all commands
type Help struct{}

func (c *Help) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()

	return builder.Reply(bot, chat, loc.Message(loc.Help))
}

type Language struct{}

func (c *Language) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()

	return builder.ReplyKeyboard(bot, chat, loc.Message(loc.Lang), builder.MakeInlineKeyboard(loc.SupportedLanguages, consts.LanguageButton))
}

type Ping struct{}

func (c *Ping) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()
	return builder.Reply(bot, chat, loc.Message(loc.Pong))
}
