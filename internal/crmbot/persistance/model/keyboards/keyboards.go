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
		tgbotapi.NewInlineKeyboardButtonData("Ревизия "+emoji.MagnifyingGlass, "revision"),
		tgbotapi.NewInlineKeyboardButtonData("Касса "+emoji.FaceMoney, "cash"),
	),
)

var Settings = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Продукты "+emoji.Products, "product_settings"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Категории "+emoji.Page, "category_settings"),
		tgbotapi.NewInlineKeyboardButtonData("Поставщики "+emoji.Lorry, "supplier_settings"),
	),
	tgbotapi.NewInlineKeyboardRow(
		MainMenuButton,
	),
)

// Categories

var CategorySettings = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить "+emoji.Plus, "category_add"),
		tgbotapi.NewInlineKeyboardButtonData("Посмотреть все "+emoji.Page, "category_get_all"),
	),
	tgbotapi.NewInlineKeyboardRow(
		MainMenuButton,
	),
)

// Suppliers

var SupplierSettings = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить "+emoji.Plus, "supplier_add"),
		tgbotapi.NewInlineKeyboardButtonData("Посмотреть все "+emoji.Page, "supplier_get_all"),
	),
	tgbotapi.NewInlineKeyboardRow(
		MainMenuButton,
	),
)

// Products

var ProductSettings = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить "+emoji.Plus, "product_add"),
		tgbotapi.NewInlineKeyboardButtonData("Посмотреть все "+emoji.Page, "product_get_all"),
	),
	tgbotapi.NewInlineKeyboardRow(
		MainMenuButton,
	),
)

// Stock
var Stock = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Посмотреть всё "+emoji.Page, "quantity_all"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Пополнить продукт "+emoji.Plus, "quantity_add"),
		tgbotapi.NewInlineKeyboardButtonData("Задать количество "+emoji.Pencil, "quantity_set"),
	),
	tgbotapi.NewInlineKeyboardRow(
		MainMenuButton,
	),
)

// Cash
var Cash = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Транзакции "+emoji.DollarBanknote, "transactions"),
		tgbotapi.NewInlineKeyboardButtonData("Добавить в кассу "+emoji.Plus, "cash_add"),
	),
	tgbotapi.NewInlineKeyboardRow(
		MainMenuButton,
	),
)

// Revision
var Revision = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("История "+emoji.Page, "revision_history"),
		tgbotapi.NewInlineKeyboardButtonData("Провести ревизию "+emoji.Plus, "revision_process"),
	),
	tgbotapi.NewInlineKeyboardRow(
		MainMenuButton,
	),
)

// Utils

func MarkupByArray(array []string) tgbotapi.ReplyKeyboardMarkup {
	countRows := len(array) / 3
	if len(array) > 3 || countRows == 0 {
		countRows++
	}
	rows := make([][]tgbotapi.KeyboardButton, countRows)
	var x int
	for i, c := range array {
		if i%3 == 0 && i != 0 {
			x++
		}
		rows[x] = append(rows[x], tgbotapi.NewKeyboardButton(c))
	}
	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	return keyboard
}
