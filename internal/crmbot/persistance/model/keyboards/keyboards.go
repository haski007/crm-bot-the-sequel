package keyboards

import (
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var MainMenuButton = tgbotapi.NewInlineKeyboardButtonData("Главное Меню "+emoji.House, "home")

var MainMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Настройки "+emoji.Gear, "configs"),
		tgbotapi.NewInlineKeyboardButtonData("Склад "+emoji.Package, "stock"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Ревизия "+emoji.Page, "revision"),
		tgbotapi.NewInlineKeyboardButtonData("Касса "+emoji.MoneyFace, "cashbox"),
	),
)
