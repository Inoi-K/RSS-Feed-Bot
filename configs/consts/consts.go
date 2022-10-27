package consts

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

const (
	ParseMode = tgbotapi.ModeMarkdownV2

	FirstMenu  = "*Menu 1*\n\nA beautiful menu with a shiny inline button"
	SecondMenu = "*Menu 2*\n\nA better menu with even more shiny inline buttons"

	NavigationButton  = "Navigation"
	TutorialButton    = "Tutorial"
	UnsubscribeButton = "Unsubscribe"

	NextText = "Next"
	BackText = "Back"

	ArgumentsSeparator = " "
)

var (
	FirstMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(NextText, strings.Join([]string{NavigationButton, NextText}, ArgumentsSeparator)),
		),
	)
	SecondMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(BackText, strings.Join([]string{NavigationButton, BackText}, ArgumentsSeparator)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(TutorialButton, "https://core.telegram.org/bots/api"),
		),
	)
)
