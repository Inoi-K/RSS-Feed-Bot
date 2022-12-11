package consts

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

const (
	FirstMenu  = "*Menu 1*\n\nA beautiful menu with a shiny inline button"
	SecondMenu = "*Menu 2*\n\nA better menu with even more shiny inline buttons"

	NextText = "Next"
	BackText = "Back"
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
	UpdateKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(fmt.Sprintf("/%v", UpdateCommand)),
		),
	)
)
