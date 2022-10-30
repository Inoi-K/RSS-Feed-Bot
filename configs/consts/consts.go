package consts

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

const (
	ParseMode = tgbotapi.ModeMarkdownV2

	FirstMenu  = "*Menu 1*\n\nA beautiful menu with a shiny inline button"
	SecondMenu = "*Menu 2*\n\nA better menu with even more shiny inline buttons"

	//region COMMANDS

	MenuCommand        = "menu"
	StartCommand       = "start"
	SubscribeCommand   = "sub"
	UnsubscribeCommand = "unsub"
	NavigationButton   = "navigation"
	TutorialButton     = "tutorial"
	UnsubscribeButton  = "unsubscribe"

	//endregion

	NextText = "Next"
	BackText = "Back"

	ArgumentsSeparator = " "

	// region ERRORS

	// DuplicationCode aka UniqueConstraintCode
	DuplicationCode = "23505"

	//endregion

	//region INTERVALS

	FeedUpdateIntervalSeconds = 5

	//endregion
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
