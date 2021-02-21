package resource

import (
	"strings"

	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) HandleRoutes(updates tgbotapi.UpdatesChannel) {
	for update := range updates {

		if update.EditedMessage != nil {
			continue
		}

		// ---> Callbacks of Inline Keyboard
		if update.CallbackQuery != nil {
			if !bot.AuthService.IsUser(update.CallbackQuery.From.ID) {
				bot.Reply(
					update.CallbackQuery.Message.Chat.ID,
					"Не зарегистрированный пользователь. Обратитесь к @pdemian за помощью")
				continue
			}

			switch update.CallbackQuery.Data {
			case "home":
				go bot.callHomeHandler(update)
			case "settings":
				go bot.callSettingsHandler(update)

			// ---> Categories
			case "category_settings":
				go bot.callCategorySettingsHandler(update)
			case "category_add":
				go bot.callCategoryAddHandler(update)
			case "category_get_all":
				go bot.callCategoryGetAllHandler(update)

			// ---> Suppliers
			case "supplier_settings":
				go bot.callSupplierSettingsHandler(update)
			case "supplier_add":
				go bot.callSupplierAddHandler(update)
			case "supplier_get_all":
				go bot.callSupplierGetAllHandler(update)

			// ---> Products
			case "product_settings":
				go bot.callProductSettingsHandler(update)
			case "product_add":
				go bot.callProductAddHandler(update)
			case "product_get_all":
				go bot.callProductGetAllHandler(update)

			// ---> Stock
			case "stock":
				go bot.callStockHandler(update)
			case "quantity_set":
				go bot.callQuantitySetHandler(update)
			case "quantity_add":
				go bot.callQuantityAddHandler(update)
			case "quantity_all":
				go bot.callQuantityAllHandler(update)

			// ---> Cash
			case "cash":
				go bot.callCashHandler(update)
			case "cash_add":
				go bot.callCashAddHandler(update)
			case "transactions":
				go bot.callTransactionsGetAllHandler(update)
			}

			continue
		}

		// ---> Commands
		if update.Message.IsCommand() {
			if !bot.AuthService.IsUser(update.Message.From.ID) {
				bot.Reply(
					update.Message.Chat.ID,
					"Не зарегистрированный пользователь. Обратитесь к @pdemian за помощью")
				continue
			}

			command := update.Message.CommandWithAt()
			switch {
			case command == "menu":
				go bot.commandMenuHandler(update)
			case strings.Contains(command, "remove_category_"):
				go bot.commandCategoryRemove(update)
			case strings.Contains(command, "edit_category_"):
				go bot.commandCategoryEditHandler(update)
			case strings.Contains(command, "remove_supplier_"):
				go bot.commandSupplierRemove(update)
			case strings.Contains(command, "edit_supplier_"):
				go bot.commandSupplierEditHandler(update)
			case strings.Contains(command, "remove_product_"):
				go bot.commandProductRemove(update)
			case strings.Contains(command, "show_product_"):
				go bot.commandProductShow(update)
			case strings.Contains(command, "edit_product_"):
				go bot.commandProductEditHandler(update)
			case strings.Contains(command, "remove_user_"):
				go bot.commandUserRemove(update)
			case strings.Contains(command, "remove_tx_"):
				go bot.commandTransactionRemoveHandler(update)

			// ---> Users
			case command == "register":
				go bot.commandRegisterHandler(update)
			case command == "get_users":
				go bot.commandGetUsersHandler(update)
			case command == "test":
				go bot.commandTestHandler(update)

			default:
				go bot.Reply(update.Message.Chat.ID, "Such command does not exist! "+emoji.NoEntry)
			}

			// ---> Hooks to process prompt
		} else {
			if !bot.AuthService.IsUser(update.Message.From.ID) {
				bot.Reply(
					update.Message.Chat.ID,
					"Не зарегистрированный пользователь. Обратитесь к @pdemian за помощью")
				continue
			}

			if op, ok := OpsQueue[update.Message.From.ID]; ok {
				switch op.Name {
				case OperationType_CategoryAdd:
					go bot.hookCategoryAdd(update)
				case OperationType_CategoryEdit:
					go bot.hookCategoryEdit(update)
				case OperationType_SupplierAdd:
					go bot.hookSupplierAdd(update)
				case OperationType_SupplierEdit:
					go bot.hookSupplierEdit(update)
				case OperationType_ProductAdd:
					go bot.hookProductAdd(update)
				case OperationType_ProductEdit:
					go bot.hookProductEdit(update)
				case OperationType_ProductGetByCategory:
					go bot.hookProductGetByCategory(update)

				case OperationType_QuantityAdd:
					go bot.hookQuantityAdd(update)
				case OperationType_QuantitySet:
					go bot.hookQuantitySet(update)
				case OperationType_QuantityAll:
					go bot.hookQuantityAll(update)

				case OperationType_CashAdd:
					go bot.hookCashAdd(update)
				}
			}
		}
	}
}
