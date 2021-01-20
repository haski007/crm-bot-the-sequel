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

		if update.Message.IsCommand() {
			switch update.Message.CommandWithAt() {
			case "menu":
				bot.commandMenuHandler(update)

			case "test":
				bot.commandTestHandler(update)
			default:
				bot.Reply(update.Message.Chat.ID, "Such command does not exist! "+emoji.NoEntry)
			}
		}
	}
}