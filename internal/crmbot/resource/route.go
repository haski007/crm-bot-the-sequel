package resource

import (
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) HandleRoutes(updates tgbotapi.UpdatesChannel) {
	//botCreds, err := bot.Bot.GetMe()
	//if err != nil {
	//	bot.ReportToTheCreator(
	//		fmt.Sprintf("[bot GetMe] err: %s", err))
	//	return
	//}
	for update := range updates {
		if update.EditedMessage != nil {
			continue
		}

		if update.CallbackQuery != nil {
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
			}
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.CommandWithAt() {
			case "menu":
				go bot.commandMenuHandler(update)

			case "test":
				go bot.commandTestHandler(update)
			default:
				go bot.Reply(update.Message.Chat.ID, "Such command does not exist! "+emoji.NoEntry)
			}
		} else {
			if op, ok := OpsQueue[update.Message.From.ID]; ok {
				switch op.Name {
				case OperationType_CategoryAdd:
					go bot.hookCategoryAdd(update)
				}
			}
		}

	}
}
