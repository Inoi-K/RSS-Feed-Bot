package consts

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

const (
	ParseMode = tgbotapi.ModeMarkdownV2

	FirstMenu  = "*Menu 1*\n\nA beautiful menu with a shiny inline button"
	SecondMenu = "*Menu 2*\n\nA better menu with even more shiny inline buttons"

	MenuCommand        = "menu"
	StartCommand       = "start"
	SubscribeCommand   = "sub"
	UnsubscribeCommand = "unsub"
	ActivateCommand    = "act"
	DeactivateCommand  = "deact"

	NavigationButton  = "navigation"
	TutorialButton    = "tutorial"
	UnsubscribeButton = "unsubscribe"
	SetIsActiveButton = "setActive"

	ActivateText   = "activate"
	DeactivateText = "deactivate"

	NextText = "Next"
	BackText = "Back"

	ArgumentsSeparator = " "

	// ERRORS

	// DuplicationCode aka UniqueConstraintCode
	DuplicationCode = "23505"
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
