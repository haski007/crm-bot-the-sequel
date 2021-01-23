package keyboards

import (
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var Empty = tgbotapi.NewInlineKeyboardMarkup()

var MainMenuButton = tgbotapi.NewInlineKeyboardButtonData("Главное Меню "+emoji.House, "home")

var MainMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Настройки "+emoji.Gear, "settings"),
		tgbotapi.NewInlineKeyboardButtonData("Склад "+emoji.Package, "stock"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Ревизия "+emoji.Page, "revision"),
		tgbotapi.NewInlineKeyboardButtonData("Касса "+emoji.MoneyFace, "cash"),
	),
)

var Settings = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Категории "+emoji.Pencil, "category_settings"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Ревизия "+emoji.Page, "revision"),
	),
	tgbotapi.NewInlineKeyboardRow(
		MainMenuButton,
	),
)

var CategorySettings = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить "+emoji.Plus, "category_add"),
		tgbotapi.NewInlineKeyboardButtonData("Изменить "+emoji.Pencil, "category_edit"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Удалить "+emoji.Basket, "category_remove"),
		tgbotapi.NewInlineKeyboardButtonData("Посмотреть все "+emoji.Page, "category_get_all"),
	),
	tgbotapi.NewInlineKeyboardRow(
		MainMenuButton,
	),
)
