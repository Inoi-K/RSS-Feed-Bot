// TODO find a better name (not util) and place for the package
package util

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Config struct {
	Screaming bool

	NextButton              string
	BackButton              string
	TutorialButton          string
	UnsubscribeCallbackData string

	// TODO organize the same fields in a struct
	FirstMenu       string
	FirstMenuMarkup tgbotapi.InlineKeyboardMarkup

	SecondMenu       string
	SecondMenuMarkup tgbotapi.InlineKeyboardMarkup
}

func NewConfig() *Config {
	cfg := &Config{}

	cfg.Screaming = false

	cfg.FirstMenu = "<b>Menu 1</b>\n\nA beautiful menu with a shiny inline button."
	// Keyboard layout for the first menu. One button, one row
	cfg.FirstMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(cfg.NextButton, cfg.NextButton),
		),
	)

	cfg.SecondMenu = "<b>Menu 2</b>\n\nA better menu with even more shiny inline buttons."
	// Keyboard layout for the second menu. Two buttons, one per row
	cfg.SecondMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(cfg.BackButton, cfg.BackButton),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(cfg.TutorialButton, "https://core.telegram.org/bots/api"),
		),
	)

	cfg.NextButton = "Next"
	cfg.BackButton = "Back"
	cfg.TutorialButton = "Tutorial"
	cfg.UnsubscribeCallbackData = "Unsub"

	return cfg
}
